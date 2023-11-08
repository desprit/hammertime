package tasks_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	schedule_storage "desprit/hammertime/src/db/schedule"
	subscription_storage "desprit/hammertime/src/db/subscription"
	user_storage "desprit/hammertime/src/db/user"
	"desprit/hammertime/src/scraper"
	"desprit/hammertime/src/tasks"
	testing_helpers "desprit/hammertime/src/testing"
)

func createData(
	ctx context.Context,
	d *sql.DB,
	t *testing.T,
	beginDate string,
) []schedule_storage.Schedule {
	var scheduleEntries []schedule_storage.Schedule
	for i := 0; i < 5; i++ {
		entry, err := schedule_storage.New(d).CreateScheduleEntry(ctx, schedule_storage.CreateScheduleEntryParams{
			ActivityID: int64(i + 1),
			Datetime:   time.Now().Format(time.RFC3339),
			Trainer:    "test-trainer",
			Activity:   fmt.Sprintf("test-activity-%v", i),
			PreEntry:   true,
			BeginDate:  beginDate,
		})
		assert.Nil(t, err)
		scheduleEntries = append(scheduleEntries, entry)
	}
	return scheduleEntries
}

type FakeScraper struct {
	isError  bool
	response scraper.ScheduleItemResponse
}

func NewFakeScraper(response scraper.ScheduleItemResponse, isError bool) *FakeScraper {
	return &FakeScraper{response: response, isError: isError}
}

type TestMessager struct{}

func NewTestMessager() *TestMessager {
	return &TestMessager{}
}

func (t *TestMessager) SendMessage(message string) error {
	log.Printf("Sending message: %s", message)
	return nil
}

func (s *FakeScraper) GetScheduleItem(activityID int64) (scraper.ScheduleItemResponse, error) {
	if s.isError {
		return scraper.ScheduleItemResponse{}, fmt.Errorf("error")
	}
	return s.response, nil
}

func (s *FakeScraper) SendMessage(message string) error {
	return nil
}

func (s *FakeScraper) CheckAuth() error {
	return nil
}

func (s *FakeScraper) Scrape(year, week int) (scraper.ScheduleResponse, error) {
	return scraper.ScheduleResponse{}, nil
}

func (s *FakeScraper) Reserve(activityID int64) (scraper.ReservationResponse, error) {
	return scraper.ReservationResponse{}, nil
}

func TestRegisterSubscriptionsCheckTask(t *testing.T) {
	ctx := context.Background()

	t.Run("should trigger subscription events", func(t *testing.T) {
		d, err := testing_helpers.PrepareDB(t)
		assert.Nil(t, err)
		defer d.Close()
		scheduleEntries := createData(ctx, d, t, time.Now().Add(time.Second*1).Format(time.RFC3339))

		wg := &sync.WaitGroup{}
		scraper := scraper.NewHttpScraper()
		messager := NewTestMessager()
		tasksStopCh := make(chan bool)
		interruptCh := make(chan bool)
		subscriptionsTriggerCh := make(chan subscription_storage.Subscription, 5)
		scheduler := tasks.NewScheduler(
			wg,
			d,
			scraper,
			messager,
			tasksStopCh,
			interruptCh,
			subscriptionsTriggerCh,
			tasks.WithSubscriptionsCheckInterval(1*time.Second),
		)

		sub1Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[0].ID,
			ScheduleID: scheduleEntries[0].ID,
		}
		_, err = scheduler.SubscriptionStorage.CreateSubscription(ctx, sub1Params)
		assert.Nil(t, err)
		sub2Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[1].ID,
			ScheduleID: scheduleEntries[1].ID,
		}
		_, err = scheduler.SubscriptionStorage.CreateSubscription(ctx, sub2Params)
		assert.Nil(t, err)

		scheduler.RegisterSubscriptionsCheckTask()
		scheduler.RegisterSubscriptionsHandleTask()
		go func() {
			for {
				select {
				case <-interruptCh:
					close(tasksStopCh)
					return
				case <-time.After(5 * time.Second):
					close(tasksStopCh)
					return
				}
			}
		}()
		wg.Wait()

		t.Fail()
	})

	t.Run("should not send notifications when subscriptions are cancelled before they trigger", func(t *testing.T) {
		d, err := testing_helpers.PrepareDB(t)
		assert.Nil(t, err)
		defer d.Close()
		scheduleEntries := createData(ctx, d, t, time.Now().Add(time.Second*10).Format(time.RFC3339))

		wg := &sync.WaitGroup{}
		scraper := scraper.NewHttpScraper()
		messager := NewTestMessager()
		tasksStopCh := make(chan bool)
		interruptCh := make(chan bool)
		subscriptionsTriggerCh := make(chan subscription_storage.Subscription)
		scheduler := tasks.NewScheduler(
			wg,
			d,
			scraper,
			messager,
			tasksStopCh,
			interruptCh,
			subscriptionsTriggerCh,
			tasks.WithSubscriptionsCheckInterval(1*time.Second),
		)

		sub1Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[0].ID,
			ScheduleID: scheduleEntries[0].ID,
		}
		sub1, err := scheduler.SubscriptionStorage.CreateSubscription(ctx, sub1Params)
		assert.Nil(t, err)
		sub2Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[1].ID,
			ScheduleID: scheduleEntries[1].ID,
		}
		sub2, err := scheduler.SubscriptionStorage.CreateSubscription(ctx, sub2Params)
		assert.Nil(t, err)

		scheduler.RegisterSubscriptionsCheckTask()
		scheduler.RegisterSubscriptionsHandleTask()
		go func() {
			for {
				select {
				case <-interruptCh:
					close(tasksStopCh)
					return
				case <-time.After(5 * time.Second):
					close(tasksStopCh)
					return
				}
			}
		}()
		go func() {
			time.Sleep(2 * time.Second)
			err = scheduler.SubscriptionStorage.CancelSubscription(
				ctx,
				subscription_storage.CancelSubscriptionParams{UserID: sub1.UserID, ScheduleID: sub1.ScheduleID},
			)
			assert.Nil(t, err)
			err = scheduler.SubscriptionStorage.CancelSubscription(
				ctx,
				subscription_storage.CancelSubscriptionParams{UserID: sub2.UserID, ScheduleID: sub2.ScheduleID},
			)
			assert.Nil(t, err)
		}()
		wg.Wait()

		t.Fail()
	})

	t.Run("should print error message when couldn't find schedule entry in the database", func(t *testing.T) {
		d, err := testing_helpers.PrepareDB(t)
		assert.Nil(t, err)
		defer d.Close()

		wg := &sync.WaitGroup{}
		scraper := NewFakeScraper(scraper.ScheduleItemResponse{}, true)
		messager := NewTestMessager()
		tasksStopCh := make(chan bool)
		interruptCh := make(chan bool)
		subscriptionsTriggerCh := make(chan subscription_storage.Subscription)
		scheduler := tasks.NewScheduler(
			wg,
			d,
			scraper,
			messager,
			tasksStopCh,
			interruptCh,
			subscriptionsTriggerCh,
			tasks.WithSubscriptionsCheckInterval(1*time.Second),
		)

		scheduler.RegisterSubscriptionsHandleTask()
		subscriptionsTriggerCh <- subscription_storage.Subscription{UserID: 1, ScheduleID: 1}

		go func() {
			time.Sleep(5 * time.Second)
			close(tasksStopCh)
		}()
		wg.Wait()

		t.Fail()
	})

	t.Run("should print error message when couldn't get schedule entry info from the website", func(t *testing.T) {
		d, err := testing_helpers.PrepareDB(t)
		assert.Nil(t, err)
		defer d.Close()
		scheduleEntries := createData(ctx, d, t, time.Now().Add(time.Second*10).Format(time.RFC3339))

		wg := &sync.WaitGroup{}
		scraper := NewFakeScraper(scraper.ScheduleItemResponse{}, true)
		messager := NewTestMessager()
		tasksStopCh := make(chan bool)
		interruptCh := make(chan bool)
		subscriptionsTriggerCh := make(chan subscription_storage.Subscription)
		scheduler := tasks.NewScheduler(
			wg,
			d,
			scraper,
			messager,
			tasksStopCh,
			interruptCh,
			subscriptionsTriggerCh,
			tasks.WithSubscriptionsCheckInterval(1*time.Second),
		)

		sub1Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[0].ID,
			ScheduleID: scheduleEntries[0].ID,
		}
		sub1, err := scheduler.SubscriptionStorage.CreateSubscription(ctx, sub1Params)
		assert.Nil(t, err)

		scheduler.RegisterSubscriptionsHandleTask()
		subscriptionsTriggerCh <- subscription_storage.Subscription{
			UserID:     sub1.UserID,
			ScheduleID: sub1.ScheduleID,
		}

		go func() {
			time.Sleep(5 * time.Second)
			close(tasksStopCh)
		}()
		wg.Wait()

		t.Fail()
	})

	t.Run("should print error message when schedule entry is closed for booking", func(t *testing.T) {
		d, err := testing_helpers.PrepareDB(t)
		assert.Nil(t, err)
		defer d.Close()
		scheduleEntries := createData(ctx, d, t, time.Now().Add(time.Second*10).Format(time.RFC3339))

		wg := &sync.WaitGroup{}
		scraper := NewFakeScraper(scraper.ScheduleItemResponse{BookingOpened: false, AvailableSlots: 10}, false)
		messager := NewTestMessager()
		tasksStopCh := make(chan bool)
		interruptCh := make(chan bool)
		subscriptionsTriggerCh := make(chan subscription_storage.Subscription)
		scheduler := tasks.NewScheduler(
			wg,
			d,
			scraper,
			messager,
			tasksStopCh,
			interruptCh,
			subscriptionsTriggerCh,
			tasks.WithSubscriptionsCheckInterval(1*time.Second),
		)

		sub1Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[0].ID,
			ScheduleID: scheduleEntries[0].ID,
		}
		sub1, err := scheduler.SubscriptionStorage.CreateSubscription(ctx, sub1Params)
		assert.Nil(t, err)

		scheduler.RegisterSubscriptionsHandleTask()
		subscriptionsTriggerCh <- subscription_storage.Subscription{
			UserID:     sub1.UserID,
			ScheduleID: sub1.ScheduleID,
		}

		go func() {
			time.Sleep(5 * time.Second)
			close(tasksStopCh)
		}()
		wg.Wait()

		t.Fail()
	})

	t.Run("should print error message when schedule entry has no available slots", func(t *testing.T) {
		d, err := testing_helpers.PrepareDB(t)
		assert.Nil(t, err)
		defer d.Close()
		scheduleEntries := createData(ctx, d, t, time.Now().Add(time.Second*10).Format(time.RFC3339))

		wg := &sync.WaitGroup{}
		scraper := NewFakeScraper(scraper.ScheduleItemResponse{BookingOpened: true, AvailableSlots: 0}, false)
		messager := NewTestMessager()
		tasksStopCh := make(chan bool)
		interruptCh := make(chan bool)
		subscriptionsTriggerCh := make(chan subscription_storage.Subscription)
		scheduler := tasks.NewScheduler(
			wg,
			d,
			scraper,
			messager,
			tasksStopCh,
			interruptCh,
			subscriptionsTriggerCh,
			tasks.WithSubscriptionsCheckInterval(1*time.Second),
		)

		sub1Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[0].ID,
			ScheduleID: scheduleEntries[0].ID,
		}
		sub1, err := scheduler.SubscriptionStorage.CreateSubscription(ctx, sub1Params)
		assert.Nil(t, err)

		scheduler.RegisterSubscriptionsHandleTask()
		subscriptionsTriggerCh <- subscription_storage.Subscription{
			UserID:     sub1.UserID,
			ScheduleID: sub1.ScheduleID,
		}

		go func() {
			time.Sleep(5 * time.Second)
			close(tasksStopCh)
		}()
		wg.Wait()

		t.Fail()
	})

	t.Run("should print success message", func(t *testing.T) {
		d, err := testing_helpers.PrepareDB(t)
		assert.Nil(t, err)
		defer d.Close()
		scheduleEntries := createData(ctx, d, t, time.Now().Add(time.Second*10).Format(time.RFC3339))

		wg := &sync.WaitGroup{}
		scraper := NewFakeScraper(scraper.ScheduleItemResponse{BookingOpened: true, AvailableSlots: 1}, false)
		messager := NewTestMessager()
		tasksStopCh := make(chan bool)
		interruptCh := make(chan bool)
		subscriptionsTriggerCh := make(chan subscription_storage.Subscription)
		scheduler := tasks.NewScheduler(
			wg,
			d,
			scraper,
			messager,
			tasksStopCh,
			interruptCh,
			subscriptionsTriggerCh,
			tasks.WithSubscriptionsCheckInterval(1*time.Second),
		)

		sub1Params := subscription_storage.CreateSubscriptionParams{
			UserID:     user_storage.Users[0].ID,
			ScheduleID: scheduleEntries[0].ID,
		}
		sub1, err := scheduler.SubscriptionStorage.CreateSubscription(ctx, sub1Params)
		assert.Nil(t, err)

		scheduler.RegisterSubscriptionsHandleTask()
		subscriptionsTriggerCh <- subscription_storage.Subscription{
			UserID:     sub1.UserID,
			ScheduleID: sub1.ScheduleID,
		}

		go func() {
			time.Sleep(5 * time.Second)
			close(tasksStopCh)
		}()
		wg.Wait()

		t.Fail()
	})
}
