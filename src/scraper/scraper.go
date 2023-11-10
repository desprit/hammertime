package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	user_storage "desprit/hammertime/src/db/user"
)

type Scraper interface {
	CheckAuth() error
	GetScheduleItem(activityID int64) (ScheduleItemResponse, error)
	Scrape(year, week int) (ScheduleResponse, error)
	Reserve(activityID int64, token string) (ReservationResponse, error)
}

type HttpScraper struct{}

func NewHttpScraper() *HttpScraper {
	return &HttpScraper{}
}

func (s *HttpScraper) CheckAuth() error {
	url := "https://mobifitness.ru/api/v8/account/settings.json"
	tokens := []string{user_storage.Users[0].Token, user_storage.Users[1].Token}
	for _, token := range tokens {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body := resp.Body
		defer body.Close()
		data, err := io.ReadAll(body)
		if err != nil {
			return err
		}

		v := ScheduleResponse{}
		err = json.Unmarshal(data, &v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *HttpScraper) GetScheduleItem(activityID int64) (ScheduleItemResponse, error) {
	url := fmt.Sprintf("https://mobifitness.ru/api/v8/schedule/%d/item.json?clubId=3934", activityID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user_storage.Users[0].Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ScheduleItemResponse{}, err
	}
	body := resp.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return ScheduleItemResponse{}, err
	}

	v := ScheduleItemResponse{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return ScheduleItemResponse{}, err
	}

	return v, nil
}

func (s *HttpScraper) Reserve(activityID int64, token string) (ReservationResponse, error) {
	u := "https://mobifitness.ru/api/v8/account/reserve.json"
	form := url.Values{}
	form.Add("scheduleId", fmt.Sprintf("%d", activityID))
	form.Add("clubId", "3934")
	req, _ := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ReservationResponse{}, err
	}
	body := resp.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return ReservationResponse{}, err
	}

	v := ReservationResponse{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return ReservationResponse{}, err
	}

	return v, nil
}

func (s *HttpScraper) Scrape(year, week int) (ScheduleResponse, error) {
	url := fmt.Sprintf("https://mobifitness.ru/api/v8/club/3934/schedule.json?year=%d&week=%d", year, week)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user_storage.Users[0].Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ScheduleResponse{}, err
	}
	body := resp.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return ScheduleResponse{}, err
	}

	v := ScheduleResponse{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return ScheduleResponse{}, err
	}

	return v, nil
}
