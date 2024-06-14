package scheduler

import (
	"fmt"
	"log"
	"time"

	images "github.com/chrischriscris/kpopapi/internal/scraper/kpopping"
	"github.com/go-co-op/gocron/v2"
)

type scheduler struct {
	scheduler *gocron.Scheduler
	enabled   bool
}

func newScheduler() (*scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &scheduler{
		scheduler: &s,
		enabled:   true,
	}, nil
}

func (s *scheduler) addJob(
	duration time.Duration,
	task func() int,
	name string,
) error {
	_, err := (*s.scheduler).NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(task),
		gocron.WithName(name),
	)

	return err
}

// Run all jobs immediately and start the scheduler
func (s *scheduler) Start() {
	if !s.enabled {
		return
	}

	fmt.Println("Starting scheduler with jobs: ")

	job := (*s.scheduler).Jobs()

	for _, j := range job {
		fmt.Printf("  + %s\n", j.Name())
	}

	(*s.scheduler).Start()

	for _, j := range job {
		j.RunNow()
	}
}

func (s *scheduler) Shutdown() error {
	return (*s.scheduler).Shutdown()
}

func (s *scheduler) Enable() {
	s.enabled = true
}

func (s *scheduler) Disable() {
	s.enabled = false
}

func KPopApiScheduler() *scheduler {
	s, err := newScheduler()
	if err != nil {
		log.Fatalf("Unable to create scheduler: %v\n", err)
	}

	err = s.addJob(4*time.Hour, images.ScrapeImages, "kpopping image scraper")
	if err != nil {
		log.Fatalf("Unable to create job: %v\n", err)
	}

	return s
}
