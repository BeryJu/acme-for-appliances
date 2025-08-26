package api

import (
	"strings"
	"time"
)

// SynologyDate parser for synology's date format
type SynologyDate time.Time

// For example, "Nov 10 11:48:04 2023 GMT"
const SynologyDateLayout = "Jan 2 15:04:05 2006 MST"

// Implement Marshaler and Unmarshaler interface
func (sd *SynologyDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(SynologyDateLayout, s)
	if err != nil {
		return err
	}
	*sd = SynologyDate(t)
	return nil
}

func (sd SynologyDate) Unix() int64 {
	t := time.Time(sd)
	return t.Unix()
}

func (sd SynologyDate) Human() string {
	t := time.Time(sd)
	return t.Format(time.ANSIC)
}
