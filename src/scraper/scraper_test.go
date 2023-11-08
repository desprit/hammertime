package scraper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"desprit/hammertime/src/scraper"
)

func TestScrape(t *testing.T) {
	s := scraper.NewHttpScraper()
	year, week := time.Now().ISOWeek()
	res, err := s.Scrape(year, week)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestCheckAuth(t *testing.T) {
	s := scraper.NewHttpScraper()
	err := s.CheckAuth()
	assert.Nil(t, err)
}

func TestGetScheduleItem(t *testing.T) {
	s := scraper.NewHttpScraper()
	res, err := s.GetScheduleItem(172297107112023)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
