package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron/v2"

    images "github.com/chrischriscris/kpopapi/internal/scraper/kpopping"
	"log"
)

func main() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Unable to create scheduler: %v\n", err)
	}

	_, err = s.NewJob(
		gocron.DurationJob(2*time.Hour),
		gocron.NewTask(images.ScrapeImages),
		gocron.WithName("kpopping image scraper"),
	)
	if err != nil {
		log.Fatalf("Unable to create job: %v\n", err)
	}

	fmt.Println("Starting scheduled jobs:")
	for _, job := range s.Jobs() {
		fmt.Printf("  + %s\n", job.Name())
	}

	s.Start()

	// First run inmediately
	for _, j := range s.Jobs() {
		err := j.RunNow()
		if err != nil {
			fmt.Printf("Error running job %s: %v\n", j.Name(), err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		fmt.Println("Shutting down scheduled jobs...")
		err = s.Shutdown()
		if err != nil {
			log.Fatalf("Unable to shutdown gracefully: %v\n", err)
		}

		os.Exit(0)
	}()

	select {}
}
