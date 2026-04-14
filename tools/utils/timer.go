package utils

import "time"

type SimpleTimer struct {
	start time.Time

	duration time.Duration
}

func (s *SimpleTimer) Start() {
	s.start = time.Now()
}

func (s *SimpleTimer) Stop() time.Duration {
	t := time.Since(s.start)
	s.duration = t

	return t
}

func (s *SimpleTimer) Duration() time.Duration {
	return s.duration
}
