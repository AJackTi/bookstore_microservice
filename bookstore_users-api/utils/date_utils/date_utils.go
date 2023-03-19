package date_utils

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDbLayout   = "2006-01-02 15:04:05"
)

func getNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	return getNow().Format(apiDateLayout)
}

func GetNowDbFormat() string {
	return getNow().Format(apiDbLayout)
}
