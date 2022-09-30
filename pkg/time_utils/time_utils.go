package time_utils

import (
	"strconv"
	"time"
)

type timeSource struct {
	year     int
	day      int
	month    time.Month
	location *time.Location
	now      time.Time
}

func Time() *timeSource {
	now := time.Now()
	return &timeSource{
		year:     now.Year(),
		day:      now.Day(),
		month:    now.Month(),
		location: now.Location(),
		now:      now,
	}
}

func From(date string) *timeSource {
	now, err := time.Parse(time.Now().String(), date)
	if err != nil {
		now = time.Now()
	}

	return &timeSource{
		year:     now.Year(),
		day:      now.Day(),
		month:    now.Month(),
		location: now.Location(),
		now:      now,
	}
}

func (t *timeSource) Now() time.Time {
	t.now = time.Now()
	return t.now
}

func (t *timeSource) Value() time.Time {
	return t.now
}

func Epoch() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (t *timeSource) Tomorrow() time.Time {
	return time.Date(t.year, t.month, t.day+1, 00, 0, 0, 0, t.location)
}

func (t *timeSource) NextMonth() time.Time {
	return time.Date(t.year, t.month+1, 1, 00, 0, 0, 0, t.location)
}

func (t *timeSource) IsToday(year int, month int, day int) bool {
	return t.year == year && int(t.month) == month && t.day == day
}

func (t *timeSource) IsThisMonth(year int, month int) bool {
	return t.year == year && int(t.month) == month
}
