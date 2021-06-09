package utils

import "time"

func Yesterday() string {
	date := time.Now()
	lastLay := date.Add(-24 * time.Hour)
	s := lastLay.Format("20060102")
	return s
}
