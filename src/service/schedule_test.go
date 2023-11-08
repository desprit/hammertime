package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"desprit/hammertime/src/config"
	"desprit/hammertime/src/db"
	schedule_storage "desprit/hammertime/src/db/schedule"
	"desprit/hammertime/src/scraper"
	"desprit/hammertime/src/service"
)

func TestGetSchedule(t *testing.T) {
	c := config.GetConfig()
	d, err := db.Open(c)
	assert.Nil(t, err)
	assert.Nil(t, db.CreateTables(c, d))
	scheduleStorage := schedule_storage.New(d)
	ctx := context.Background()
	s := service.NewScheduleService(d, c)

	t.Run("should pull schedule", func(t *testing.T) {
		err := s.PullSchedule()
		assert.Nil(t, err)
		entries, err := scheduleStorage.GetScheduleEntriesWithPreEntry(ctx)
		assert.Nil(t, err)
		assert.Greater(t, len(entries), 0)
	})

	t.Run("should create schedule entry", func(t *testing.T) {
		var err error

		err = s.SaveScheduleEntry(scraper.ScheduleEntry{
			ActivityID: 1,
			Activity:   scraper.Activity{Title: "Test"},
			Datetime:   time.Now(),
			Trainers:   []scraper.Trainer{{Title: "Test"}},
			PreEntry:   false,
			Begindate:  time.Now(),
		}, scheduleStorage)
		assert.Nil(t, err)

		entry, err := scheduleStorage.GetScheduleEntry(ctx, 1)
		assert.Nil(t, err)
		assert.Equal(t, "Test", entry.Trainer)
		assert.Equal(t, "Test", entry.Activity)
		assert.Equal(t, false, entry.PreEntry)
	})

	t.Run("should not throw errors on duplicates", func(t *testing.T) {
		var err error

		fn := func() {
			err = s.SaveScheduleEntry(scraper.ScheduleEntry{
				ActivityID: 1,
				Activity:   scraper.Activity{Title: "Test"},
				Datetime:   time.Now(),
				Trainers:   []scraper.Trainer{{Title: "Test"}},
				PreEntry:   false,
				Begindate:  time.Now(),
			}, scheduleStorage)
			assert.Nil(t, err)
		}

		fn()
		fn()
	})
}
