package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"desprit/hammertime/src/config"
	schedule_storage "desprit/hammertime/src/db/schedule"
	"desprit/hammertime/src/scraper"

	"github.com/mattn/go-sqlite3"
)

type ScheduleService struct {
	db  *sql.DB
	cfg *config.Config
}

func NewScheduleService(db *sql.DB, cfg *config.Config) *ScheduleService {
	return &ScheduleService{db, cfg}
}

func (s *ScheduleService) GetSchedule(year, week int) ([]scraper.ScheduleEntry, error) {
	scr := scraper.NewHttpScraper()
	scheduleResponse, err := scr.Scrape(year, week)
	if err != nil {
		return make([]scraper.ScheduleEntry, 0), err
	}
	return scheduleResponse.Schedule, nil
}

func (s *ScheduleService) SaveScheduleEntry(
	scheduleEntry scraper.ScheduleEntry,
	storage *schedule_storage.Queries,
) error {
	trainer := ""
	if len(scheduleEntry.Trainers) > 0 {
		trainer = scheduleEntry.Trainers[0].Title
	}
	params := schedule_storage.CreateScheduleEntryParams{
		ActivityID: int64(scheduleEntry.ActivityID),
		Datetime:   scheduleEntry.Datetime.Format(time.RFC3339),
		Trainer:    trainer,
		Activity:   scheduleEntry.Activity.Title,
		PreEntry:   scheduleEntry.PreEntry,
		BeginDate:  scheduleEntry.Begindate.Format(time.RFC3339),
	}
	_, err := storage.CreateScheduleEntry(context.Background(), params)
	return err
}

func (s *ScheduleService) SaveSchedule(schedule []scraper.ScheduleEntry) error {
	storage := schedule_storage.New(s.db)
	for _, scheduleEntry := range schedule {
		err := s.SaveScheduleEntry(scheduleEntry, storage)
		if err != nil {
			if _, ok := err.(sqlite3.Error); ok {
				log.Printf("Error creating schedule entry: %v", err)
			}
			if serr, ok := err.(sqlite3.Error); ok && serr.Code == sqlite3.ErrConstraint {
				continue
			}
			log.Printf("Error creating schedule entry: %v", err)
		}
	}
	return nil
}

func (s *ScheduleService) PullSchedule() error {
	storage := schedule_storage.New(s.db)
	entry, err := storage.GetLatestScheduleEntry(context.Background())
	if err == nil {
		dt, err := time.Parse(time.RFC3339, entry.Datetime)
		if err != nil {
			return err
		}
		if dt.After(time.Now()) {
			log.Printf("Schedule is up to date")
			return nil
		}
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	currentYear, currentWeek := time.Now().ISOWeek()
	currentPeriod := []int{currentYear, currentWeek}
	nextPeriod := []int{currentYear, currentWeek + 1}
	if nextPeriod[1] > 52 {
		nextPeriod[0] += 1
		nextPeriod[1] = 1
	}

	for _, period := range [][]int{currentPeriod, nextPeriod} {
		log.Printf("Pulling schedule for period: %v", period)
		schedule, err := s.GetSchedule(period[0], period[1])
		if err != nil {
			return err
		}
		log.Printf("Found %v schedule entries", len(schedule))
		err = s.SaveSchedule(schedule)
		if err != nil {
			return err
		}
	}

	return nil
}
