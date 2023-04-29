package scheduler

import (
	"errors"
	"fmt"
	"time"
)

type Scheduler struct {
	timeToUpdate []string
	location     *time.Location
}

func NewScheduler(timeToUpdate []string, location *time.Location) *Scheduler {
	return &Scheduler{
		timeToUpdate: timeToUpdate,
		location:     location,
	}
}

func (s *Scheduler) GetTimeToWait() (time.Duration, error) {
	if len(s.timeToUpdate) <= 0 {
		return 0, errors.New("scheduler didn't get time to update")
	}

	minDuration, err := s.getDuration(s.timeToUpdate[0])
	if err != nil {
		return 0, err
	}

	for _, t := range s.timeToUpdate {
		d, err := s.getDuration(t)
		if err != nil {
			return 0, err
		}

		if d >= 0 && d <= minDuration {
			minDuration = d
		}
	}

	return minDuration, nil
}

func (s *Scheduler) getDuration(timeString string) (time.Duration, error) {
	n := time.Now().Local().In(s.location)
	t, err := time.Parse(time.TimeOnly, timeString)
	if err != nil {
		return 0, fmt.Errorf("coudn't parse time \"%s\" string for schedule: %w", timeString, err)
	}

	t = time.Date(n.Year(), n.Month(), n.Day(), t.Hour(), t.Minute(), t.Second(), 0, s.location)
	if !t.After(n) {
		t = t.AddDate(0, 0, 1)
	}

	return time.Until(t), nil
}
