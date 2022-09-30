package time_utils

import (
	"github.com/brienze1/crypto-robot-validator/pkg/time_utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var timeUtils = time_utils.Time()

func TestNowTime(t *testing.T) {
	timeNow := time.Now()

	now := timeUtils.Now()

	assert.Equal(t, timeNow.Year(), now.Year())
	assert.Equal(t, timeNow.Month(), now.Month())
	assert.Equal(t, timeNow.Day(), now.Day())
	assert.Equal(t, timeNow.Hour(), now.Hour())
	assert.Equal(t, timeNow.Minute(), now.Minute())
	assert.Equal(t, timeNow.Second(), now.Second())
}

func TestTomorrowTime(t *testing.T) {
	timeNow := time.Now()
	tomorrow := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()+1, 00, 0, 0, 0, timeNow.Location())

	newTomorrow := timeUtils.Tomorrow()

	assert.Equal(t, tomorrow.Year(), newTomorrow.Year())
	assert.Equal(t, tomorrow.Month(), newTomorrow.Month())
	assert.Equal(t, tomorrow.Day(), newTomorrow.Day())
	assert.Equal(t, tomorrow.Hour(), newTomorrow.Hour())
	assert.Equal(t, tomorrow.Minute(), newTomorrow.Minute())
	assert.Equal(t, tomorrow.Second(), newTomorrow.Second())
}

func TestNextMonthTime(t *testing.T) {
	timeNow := time.Now()
	nextMonth := time.Date(timeNow.Year(), time.Month(int(timeNow.Month())+1), 1, 00, 0, 0, 0, timeNow.Location())

	newNextMonth := timeUtils.NextMonth()

	assert.Equal(t, nextMonth.Year(), newNextMonth.Year())
	assert.Equal(t, nextMonth.Month(), newNextMonth.Month())
	assert.Equal(t, nextMonth.Day(), newNextMonth.Day())
	assert.Equal(t, nextMonth.Hour(), newNextMonth.Hour())
	assert.Equal(t, nextMonth.Minute(), newNextMonth.Minute())
	assert.Equal(t, nextMonth.Second(), newNextMonth.Second())
}

func TestIsTodayTrue(t *testing.T) {
	timeNow := time.Now()

	isToday := timeUtils.IsToday(timeNow.Year(), int(timeNow.Month()), timeNow.Day())

	assert.Equal(t, true, isToday)
}

func TestIsTodayDifferentYearFalse(t *testing.T) {
	timeNow := time.Now()

	isToday := timeUtils.IsToday(timeNow.Year()+1, int(timeNow.Month()), timeNow.Day())

	assert.Equal(t, false, isToday)
}

func TestIsTodayDifferentMonthFalse(t *testing.T) {
	timeNow := time.Now()

	isToday := timeUtils.IsToday(timeNow.Year(), int(timeNow.Month())+1, timeNow.Day())

	assert.Equal(t, false, isToday)
}

func TestIsTodayDifferentDayFalse(t *testing.T) {
	timeNow := time.Now()

	isToday := timeUtils.IsToday(timeNow.Year(), int(timeNow.Month()), timeNow.Day()+1)

	assert.Equal(t, false, isToday)
}

func TestIsThisMonthTrue(t *testing.T) {
	timeNow := time.Now()

	isThisMonth := timeUtils.IsThisMonth(timeNow.Year(), int(timeNow.Month()))

	assert.Equal(t, true, isThisMonth)
}

func TestIsThisMonthDifferentYearFalse(t *testing.T) {
	timeNow := time.Now()

	isThisMonth := timeUtils.IsThisMonth(timeNow.Year()+1, int(timeNow.Month()))

	assert.Equal(t, false, isThisMonth)
}

func TestIsThisMonthDifferentMonthFalse(t *testing.T) {
	timeNow := time.Now()

	isThisMonth := timeUtils.IsToday(timeNow.Year(), int(timeNow.Month())+1, timeNow.Day())

	assert.Equal(t, false, isThisMonth)
}
