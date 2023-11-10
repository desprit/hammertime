package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"desprit/hammertime/src/config"
	schedule_storage "desprit/hammertime/src/db/schedule"
	subscription_storage "desprit/hammertime/src/db/subscription"
	user_storage "desprit/hammertime/src/db/user"
	"desprit/hammertime/src/scraper"
	"desprit/hammertime/src/service"

	"github.com/goodsign/monday"
)

type Scheduler struct {
	wg                     *sync.WaitGroup
	ctx                    context.Context
	cancel                 context.CancelFunc
	scraper                scraper.Scraper
	messager               Messager
	tasksStopCh            chan bool // is used to stop all background tasks
	interruptMainCh        chan bool // sends interruption signal to main goroutine
	subscriptionsTriggerCh chan subscription_storage.Subscription

	subscriptionsCheckInterval time.Duration
	schedulePullInterval       time.Duration
	tokensCheckInterval        time.Duration

	ScheduleStorage     *schedule_storage.Queries
	ScheduleService     *service.ScheduleService
	SubscriptionStorage *subscription_storage.Queries
}

type SchedulerOpt func(*Scheduler)

func NewScheduler(
	wg *sync.WaitGroup,
	d *sql.DB,
	scraper scraper.Scraper,
	messager Messager,
	tasksStopCh, interruptCh chan bool,
	subscriptionsTriggerCh chan subscription_storage.Subscription,
	opts ...SchedulerOpt,
) *Scheduler {
	var (
		defaultSubscriptionsCheckInterval = time.Minute
		defaultSchedulePullInterval       = time.Hour * 24
		defaultTokensCheckInterval        = time.Hour
	)
	ctx, cancel := context.WithCancel(context.Background())
	scheduleStorage := schedule_storage.New(d)
	scheduleService := service.NewScheduleService(d, config.GetConfig())
	subscriptionStorage := subscription_storage.New(d)
	s := &Scheduler{
		wg:                         wg,
		ctx:                        ctx,
		cancel:                     cancel,
		scraper:                    scraper,
		messager:                   messager,
		tasksStopCh:                tasksStopCh,
		interruptMainCh:            interruptCh,
		subscriptionsTriggerCh:     subscriptionsTriggerCh,
		subscriptionsCheckInterval: defaultSubscriptionsCheckInterval,
		schedulePullInterval:       defaultSchedulePullInterval,
		tokensCheckInterval:        defaultTokensCheckInterval,
		ScheduleStorage:            scheduleStorage,
		ScheduleService:            scheduleService,
		SubscriptionStorage:        subscriptionStorage,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithSubscriptionsCheckInterval(interval time.Duration) SchedulerOpt {
	return func(s *Scheduler) {
		s.subscriptionsCheckInterval = interval
	}
}

func WithSchedulePullInterval(interval time.Duration) SchedulerOpt {
	return func(s *Scheduler) {
		s.schedulePullInterval = interval
	}
}

func WithTokensCheckInterval(interval time.Duration) SchedulerOpt {
	return func(s *Scheduler) {
		s.tokensCheckInterval = interval
	}
}

var (
	mu        sync.Mutex
	cancelMap = make(map[int64]context.CancelFunc)
)

func (s *Scheduler) RegisterSubscriptionsCheckTask() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("Subscriptions monitor task with interval %v started", s.subscriptionsCheckInterval)
		for {
			select {
			case <-s.tasksStopCh:
				log.Printf("Stopping subscriptions monitor task...")
				mu.Lock()
				for _, cancel := range cancelMap {
					cancel()
				}
				mu.Unlock()
				return
			case <-time.After(s.subscriptionsCheckInterval):
				subscriptions, err := s.SubscriptionStorage.GetSubscriptions(s.ctx)
				if err != nil {
					log.Printf("Error getting subscriptions: %v", err)
					s.interruptMainCh <- true
					return
				}
				// Start new goroutines for new subscriptions
				for _, sub := range subscriptions {
					mu.Lock()
					if _, ok := cancelMap[sub.ID]; !ok {
						scheduleEntry, err := s.ScheduleStorage.GetScheduleEntry(s.ctx, sub.ScheduleID)
						if err != nil {
							log.Printf("Error getting schedule entry: %v", err)
							s.SubscriptionStorage.DeleteSubscription(s.ctx, sub.ID)
							s.interruptMainCh <- true
							return
						}
						ctx, cancel := context.WithCancel(s.ctx)
						cancelMap[sub.ID] = cancel
						t, _ := time.Parse(time.RFC3339, scheduleEntry.BeginDate)
						beginDate := monday.Format(t, "Jan'02 15:04", monday.LocaleRuRU)
						t = t.Add(time.Duration(time.Second * 5))
						go func(sub subscription_storage.Subscription, t time.Time) {
							log.Printf("Goroutine for subscription %d started, it will trigger %v", sub.ID, t)
							message := fmt.Sprintf(
								"Получил заявку от |%s| на %s, запись откроется %s, попробую записать! #FingersCrossed",
								user_storage.UserMapByID[sub.UserID].Name,
								scheduleEntry.Activity,
								beginDate,
							)
							s.messager.SendMessage(message)
							for {
								select {
								case <-time.After(time.Until(t)):
									log.Printf("Subscription %d goroutine has triggered, exiting", sub.ID)
									s.subscriptionsTriggerCh <- sub
									return
								case <-s.tasksStopCh:
									log.Printf("Subscription %d goroutine received stop event, exiting", sub.ID)
									return
								case <-ctx.Done():
									log.Printf("Subscription %d goroutine has been cancelled, exiting", sub.ID)
									return
								}
							}
						}(sub, t)
					}
					mu.Unlock()
				}
				// Cancel goroutines for removed subscriptions
				mu.Lock()
				for id := range cancelMap {
					found := false
					for _, sub := range subscriptions {
						if sub.ID == id {
							found = true
							break
						}
					}
					if !found {
						log.Printf("Subscription %d not found, cancelling goroutine...", id)
						cancelMap[id]()
						delete(cancelMap, id)
					}
				}
				mu.Unlock()
				log.Printf("Subscriptions renewed")
			}
		}
	}()
}

func (s *Scheduler) RegisterSubscriptionsHandleTask() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("Notifications task started")
		for {
			select {
			case sub := <-s.subscriptionsTriggerCh:
				if err := s.SubscriptionStorage.DeleteSubscription(s.ctx, sub.ID); err != nil {
					log.Printf("Error deleting subscription in RegisterSubscriptionsHandleTask: %v", err)
				}

				user := user_storage.UserMapByID[sub.UserID]
				scheduleEntry, err := s.ScheduleStorage.GetScheduleEntry(s.ctx, sub.ScheduleID)
				if err != nil {
					log.Printf("Couldn't get schedule entry in RegisterSubscriptionsHandleTask: %v", err)
					message := fmt.Sprintf("Хотел записать |%s| на тренировку, но что-то пошло не так!", user.Name)
					s.messager.SendMessage(message)
					continue
				}
				resp, err := s.scraper.GetScheduleItem(scheduleEntry.ActivityID)
				if err != nil {
					log.Printf("Couldn't check for slots in RegisterSubscriptionsHandleTask: %v", err)
					message := fmt.Sprintf(
						"Хотел записать |%s| на %s, но не удалось получить информацию о занятии.",
						user.Name,
						scheduleEntry.Activity,
					)
					s.messager.SendMessage(message)
					continue
				} else {
					log.Printf("Got schedule item: %+v", resp)
					if !resp.BookingOpened {
						log.Printf("Booking is not opened in RegisterSubscriptionsHandleTask")
						message := fmt.Sprintf("Хотел записать |%s| на %s, но запись почему-то закрыта. #F", user.Name, scheduleEntry.Activity)
						s.messager.SendMessage(message)
						continue
					}
					if resp.AvailableSlots == 0 {
						log.Printf("No available slots in RegisterSubscriptionsHandleTask")
						message := fmt.Sprintf("Хотел записать |%s| на %s, но свободных мест больше нет! #DaKakTakTo", user.Name, scheduleEntry.Activity)
						s.messager.SendMessage(message)
						continue
					}
				}
				response, err := s.scraper.Reserve(scheduleEntry.ActivityID, user.Token)
				if err != nil {
					message := fmt.Sprintf(
						"Хотел записать |%s| на %s, но что-то пошло не так!",
						user.Name,
						scheduleEntry.Activity,
					)
					log.Printf("Couldn't reserve in RegisterSubscriptionsHandleTask: %v", err)
					s.messager.SendMessage(message)
				} else {
					if response.Result == "success" {
						message := fmt.Sprintf(
							"Я записал |%s| на %s, ура! #DobroHacker",
							user.Name,
							scheduleEntry.Activity,
						)
						s.messager.SendMessage(message)
					} else {
						message := fmt.Sprintf(
							"Хотел записать |%s| на %s, но что-то пошло не так!",
							user.Name,
							scheduleEntry.Activity,
						)
						log.Printf("Couldn't reserve in RegisterSubscriptionsHandleTask: %+v", response)
						s.messager.SendMessage(message)
					}
				}
			case <-s.tasksStopCh:
				log.Printf("Stopping notifications task...")
				return
			}
		}
	}()
}

func (s *Scheduler) RegisterSchedulePullTask() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("Schedule scraping task with interval %v started", s.schedulePullInterval)
		if err := s.ScheduleService.PullSchedule(); err != nil {
			log.Printf("Error pulling schedule: %v", err)
			s.interruptMainCh <- true
			return
		}
		for {
			select {
			case <-time.After(s.schedulePullInterval):
				if err := s.ScheduleService.PullSchedule(); err != nil {
					log.Printf("Error pulling schedule: %v", err)
					s.interruptMainCh <- true
					return
				}
				log.Printf("Schedule updated")
			case <-s.tasksStopCh:
				log.Printf("Stopping schedule scraping task...")
				return
			}
		}
	}()
}

func (s *Scheduler) RegisterTokensCheckTask() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("Tokens verification task with interval %v started", s.tokensCheckInterval)
		for {
			select {
			case <-time.After(s.tokensCheckInterval):
				err := scraper.NewHttpScraper().CheckAuth()
				if err != nil {
					log.Printf("Error checking auth: %v", err)
					s.interruptMainCh <- true
					return
				}
				log.Printf("Tokens verified")
			case <-s.tasksStopCh:
				log.Printf("Stopping tokens verification task...")
				return
			}
		}
	}()
}

func (s *Scheduler) RegisterBackgroundTasks() {
	s.RegisterSubscriptionsHandleTask()
	s.RegisterSubscriptionsCheckTask()
	s.RegisterSchedulePullTask()
	s.RegisterTokensCheckTask()
}
