package sdkcm

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timeFmt   = "2006-01-02T15:04:05.999999-07:00"
	timeRegex = `-|:|T|\.`
)

// Set time format layout. Default: 2006-01-02T15:04:05.999999-07:00
func SetTimeFormat(layout string) {
	timeFmt = layout
}

type JSONTime time.Time

func NewJSONTime() *JSONTime {
	now := JSONTime(time.Now().UTC())

	return &now
}

// Implement method MarshalJSON to output time with in formatted
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(timeFmt))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	value := string(data)

	matched, err := regexp.MatchString(timeRegex, value)
	if err != nil {
		return err
	}

	if matched {
		ti, err := time.Parse(timeFmt, strings.Replace(string(data), "\"", "", -1))
		if err != nil {
			return err
		}
		*t = JSONTime(ti)

		return nil
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*t = JSONTime(time.UnixMilli(i))

	return nil
}

// This method for mapping JSONTime to datetime data type in sql
func (t *JSONTime) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return time.Time(*t).Format("2006-01-02 15:04:05"), nil
}

// This method for scanning JSONTime from datetime data type in sql
func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if v, ok := value.(time.Time); ok {
		*t = JSONTime(v)
		return nil
	}

	return errors.New("invalid Scan Source")
}

// Before reports whether the time instant t is before u.
func (t JSONTime) Before(u JSONTime) bool {
	timeT := time.Time(t)
	timeU := time.Time(u)

	return timeT.Before(timeU)
}

func (t JSONTime) BeginningOfDay() JSONTime {
	data := time.Time(t)
	year, month, day := data.Date()

	out := time.Date(year, month, day, 0, 0, 0, 0, data.Location())

	return JSONTime(out)
}

func (t JSONTime) DiffSecond(t2 JSONTime) int64 {
	timeT := time.Time(t)
	timeT2 := time.Time(t2)

	return timeT.Unix() - timeT2.Unix()
}

func (t JSONTime) Weekday() time.Weekday {
	return time.Time(t).Weekday()
}

func (t JSONTime) Between(t1 JSONTime, t2 JSONTime) bool {
	return t1.Before(t) && t.Before(t2)
}

func (t JSONTime) Unix() int64 {
	return time.Time(t).Unix()
}

func (t JSONTime) UnixMilli() int64 {
	return time.Time(t).UnixMilli()
}

func (t JSONTime) UnixMicro() int64 {
	return time.Time(t).UnixMicro()
}

func (t JSONTime) UnixNano() int64 {
	return time.Time(t).UnixNano()
}

func (t JSONTime) Add(duration time.Duration) JSONTime {
	after := time.Time(t).Add(duration)

	return JSONTime(after)
}
