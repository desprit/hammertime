package server

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"

	subscription_storage "desprit/hammertime/src/db/subscription"
	user_storage "desprit/hammertime/src/db/user"
)

type SubscriptionServer struct {
	e *echo.Echo
}

func NewSubscriptionServer(e *echo.Echo) *SubscriptionServer {
	return &SubscriptionServer{
		e: e,
	}
}

func getUser(c echo.Context, d *sql.DB) (user_storage.User, error) {
	user := c.Param("user")
	if user == "" || (user != "d" && user != "m") {
		return user_storage.User{}, c.String(400, "Bad request")
	}
	if user == "d" {
		return user_storage.Users[0], nil
	} else {
		return user_storage.Users[1], nil
	}
}

func (s *SubscriptionServer) RegisterHandlers(d *sql.DB) {
	g := s.e.Group("subscription")
	subscriptionStorage := subscription_storage.New(d)

	g.GET("/:user/:entry/register", func(c echo.Context) error {
		user, err := getUser(c, d)
		if err != nil {
			return err
		}

		scheduleEntry := c.Param("entry")
		if scheduleEntry == "" {
			return c.String(400, "Bad request")
		}
		scheduleEntryID, _ := strconv.ParseInt(scheduleEntry, 10, 64)
		_, err = subscriptionStorage.CreateSubscription(
			c.Request().Context(),
			subscription_storage.CreateSubscriptionParams{
				UserID:     user.ID,
				ScheduleID: scheduleEntryID,
			},
		)
		if err != nil {
			log.Printf("Error creating subscription: %v", err)
			return c.String(400, "Bad request")
		}
		return c.Redirect(302, "/schedule")
	})
	g.GET("/:user/:entry/cancel", func(c echo.Context) error {
		user, err := getUser(c, d)
		if err != nil {
			return err
		}

		scheduleEntry := c.Param("entry")
		if scheduleEntry == "" {
			return c.String(400, "Bad request")
		}
		scheduleEntryID, _ := strconv.ParseInt(scheduleEntry, 10, 64)
		err = subscriptionStorage.CancelSubscription(
			c.Request().Context(),
			subscription_storage.CancelSubscriptionParams{
				UserID:     user.ID,
				ScheduleID: scheduleEntryID,
			},
		)
		if err != nil {
			log.Printf("Error canceling subscription: %v", err)
			return c.String(400, "Bad request")
		}
		return c.Redirect(302, "/schedule")
	})
}
