package scraper

import "time"

type AccountResponse struct {
	Card  string `json:"card"`
	Phone string `json:"phone"`
}

type Trainer struct {
	Title string `json:"title"`
}

type Activity struct {
	Title string `json:"title"`
}

type ScheduleEntry struct {
	ActivityID int64     `json:"id"`
	Datetime   time.Time `json:"datetime"`
	Trainers   []Trainer `json:"trainers"`
	Activity   Activity  `json:"activity"`
	PreEntry   bool      `json:"preEntry"`
	Begindate  time.Time `json:"begindate"`
}

type ScheduleResponse struct {
	Schedule []ScheduleEntry `json:"schedule"`
}

type ScheduleItemResponse struct {
	AvailableSlots int  `json:"availableSlots"`
	BookingOpened  bool `json:"bookingOpened"`
}

type ReservationResponse struct {
	Result string   `json:"result"`
	Errors []string `json:"errors"`
	Code   int      `json:"code"`
}
