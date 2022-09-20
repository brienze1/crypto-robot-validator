package time_utils

import "time"

type timeSource struct {
}

func Time() *timeSource {
	return &timeSource{}
}

func (t *timeSource) Now() time.Time {
	return time.Now()
}
