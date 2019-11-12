package fetch

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

type Train struct {
	Name        string
	FromStation string
	FromTime    time.Time
	ToStation   string
	Duration    time.Duration
}

func (t Train) ToTime() time.Time {
	return t.FromTime.Add(t.Duration)
}

func (t Train) FromTimeStr() string {
	return t.FromTime.Format("03:04pm")
}

func (t Train) ToTimeStr() string {
	return t.ToTime().Format("03:04pm")
}

func (t Train) SameDay(when time.Time) bool {
	return t.FromTime.Truncate(24 * time.Hour).Equal(when.Truncate(24 * time.Hour))
}

func (t Train) DurationStr() string {
	d := t.Duration.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%dh %dm", h, m)
}

func (t Train) hash() string {
	return hashStrings(t.Name, t.FromStation, t.ToStation,
		strconv.FormatInt(t.FromTime.Unix(), 36))
}

func hashStrings(v ...string) string {
	var b bytes.Buffer
	for _, v := range v {
		b.WriteString(v)
	}
	hash := md5.Sum(b.Bytes())
	return string(hash[:])
}
