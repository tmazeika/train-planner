package web

import "time"

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
