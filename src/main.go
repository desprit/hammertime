package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"desprit/hammertime/src/config"
	"desprit/hammertime/src/db"
	subscription_storage "desprit/hammertime/src/db/subscription"
	"desprit/hammertime/src/scraper"
	"desprit/hammertime/src/server"
	"desprit/hammertime/src/tasks"
)

func main() {
	c := config.GetConfig()
	d, err := db.Open(c)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	if err := db.InitDB(d, c); err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	scraper := scraper.NewHttpScraper()
	messager := tasks.NewTelegramMessager()
	tasksStopCh := make(chan bool)
	osSignalChannel := make(chan os.Signal, 2)
	interruptChannel := make(chan bool)
	subscriptionsTriggerCh := make(chan subscription_storage.Subscription)
	signal.Notify(osSignalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	scheduler := tasks.NewScheduler(
		wg,
		d,
		scraper,
		messager,
		tasksStopCh,
		interruptChannel,
		subscriptionsTriggerCh,
		tasks.WithSubscriptionsCheckInterval(time.Minute*1),
		tasks.WithSchedulePullInterval(time.Hour*1),
		tasks.WithTokensCheckInterval(time.Hour*1),
	)
	scheduler.RegisterBackgroundTasks()
	server := server.NewServer(d)
	server.RegisterHandlers()

	go func() {
		select {
		case <-osSignalChannel:
			log.Printf("Interrupt signal received, shutting down...")
			close(tasksStopCh)
		case <-interruptChannel:
			log.Printf("Critical error, shutting down...")
			close(tasksStopCh)
		}
		wg.Wait()
		if err := server.Shutdown(); err != nil {
			panic(err)
		}
	}()

	log.Printf("Starting server...")
	server.Logger().Fatal(server.Start())
}
