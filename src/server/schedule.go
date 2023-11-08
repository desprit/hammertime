package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/goodsign/monday"
	"github.com/labstack/echo/v4"

	schedule_storage "desprit/hammertime/src/db/schedule"
	subscription_storage "desprit/hammertime/src/db/subscription"
	user_storage "desprit/hammertime/src/db/user"
)

type ScheduleServer struct {
	e *echo.Echo
}

func NewScheduleServer(e *echo.Echo) *ScheduleServer {
	return &ScheduleServer{
		e: e,
	}
}

type ScheduleEntryUI struct {
	ID          int64  `json:"id"`
	Datetime    string `json:"datetime"`
	Activity    string `json:"activity"`
	Subscribed  bool   `json:"subscribed"`
	SubscribedD bool   `json:"subscribed_d"`
	SubscribedM bool   `json:"subscribed_m"`
}

func NewScheduleEntryUI(id int64, datetime string, activity string, subscribedD, subscribedM bool) ScheduleEntryUI {
	return ScheduleEntryUI{
		ID:          id,
		Datetime:    datetime,
		Activity:    activity,
		SubscribedD: subscribedD,
		SubscribedM: subscribedM,
	}
}

type ScheduleEntryGroup struct {
	Day       string            `json:"day"`
	Datetime  string            `json:"datetime"`
	Timestamp int               `json:"timestamp"`
	Entries   []ScheduleEntryUI `json:"entries"`
}

func NewScheduleEntryGroup(day string, datetime string, ts int) ScheduleEntryGroup {
	return ScheduleEntryGroup{
		Day:       day,
		Datetime:  datetime,
		Timestamp: ts,
		Entries:   make([]ScheduleEntryUI, 0),
	}
}

func formatScheduleEntries(
	entries []schedule_storage.Schedule,
	subscriptions []subscription_storage.Subscription,
) []ScheduleEntryGroup {
	var groupMap = map[string]ScheduleEntryGroup{}

	isUserSubscribed := func(name string, scheduleID int64) bool {
		userID := user_storage.UserMapByName[name].ID
		for _, sub := range subscriptions {
			if sub.UserID == userID && sub.ScheduleID == scheduleID {
				return true
			}
		}
		return false
	}

	for _, entry := range entries {
		t, _ := time.Parse(time.RFC3339, entry.Datetime)
		day := monday.Format(t, "Mon", monday.LocaleRuRU)
		_, week := t.ISOWeek()

		key := fmt.Sprintf("%s-%d", day, week)
		gr, ok := groupMap[key]
		if !ok {
			gr = NewScheduleEntryGroup(
				day, monday.Format(t, "Jan'02", monday.LocaleRuRU), int(t.Unix()))
		}
		gr.Entries = append(
			gr.Entries,
			NewScheduleEntryUI(
				entry.ID,
				monday.Format(t, "15:04", monday.LocaleRuRU),
				entry.Activity,
				isUserSubscribed(user_storage.UserD, entry.ID),
				isUserSubscribed(user_storage.UserM, entry.ID),
			),
		)

		groupMap[key] = gr
	}

	groups := make([]ScheduleEntryGroup, 0)
	for _, gr := range groupMap {
		groups = append(groups, gr)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Timestamp < groups[j].Timestamp
	})

	return groups
}

func (s *ScheduleServer) RegisterHandlers(d *sql.DB) {
	g := s.e.Group("schedule")

	g.GET("", func(c echo.Context) error {
		scheduleEntries, err := schedule_storage.New(d).GetScheduleEntriesWithPreEntry(c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		subscriptions, err := subscription_storage.New(d).GetSubscriptions(c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"groups": formatScheduleEntries(scheduleEntries, subscriptions),
		})
	})
}
