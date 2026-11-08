package helpers

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NormalizeUnix(unit int64) time.Time {
	// miliseconds
	if unit >= 10_000_000_000 {
		return time.Unix(unit/1000, 0).UTC()
	} else {
		return time.Unix(unit, 0).UTC()
	}
}

func UnixToDate(unit int64) pgtype.Date {
	return pgtype.Date{Time: NormalizeUnix(unit), Valid: true}
}

func TimeToDate(time time.Time) pgtype.Date {
	return pgtype.Date{Time: time, Valid: true}
}

func DateStringToTime(date string) time.Time {
	time, _ := time.Parse(time.DateOnly, date)
	return time
}

func ValidateDateStringRange(start_date_str string, end_date_str string) bool {
	if !ValidateDateString(start_date_str) || !ValidateDateString(end_date_str) {
		return false
	}

	start_time, end_time := DateStringToTime(start_date_str), DateStringToTime(end_date_str)

	if start_time.Unix() > end_time.Unix() {
		return false
	}

	return true
}

func ValidateDateString(date string) bool {
	time, date_err := time.Parse(time.DateOnly, date)
	if date_err != nil {
		return false
	}

	unix_time := time.Unix()
	min_unix := GetInt64EnvFallback("MIN_UNIX", 315532800, 32503680000)
	max_unix := GetInt64EnvFallback("MAX_UNIX", 4102444800, 32503680000)

	if unix_time < min_unix || unix_time > max_unix {
		return false
	}

	return true
}

func ValidateUnix(unit int64) bool {
	if unit <= 0 {
		return false
	}

	time := NormalizeUnix(unit)
	unix_time := time.Unix()
	min_unix := GetInt64EnvFallback("MIN_UNIX", 315532800, 32503680000)
	max_unix := GetInt64EnvFallback("MAX_UNIX", 4102444800, 32503680000)

	if unix_time < min_unix || unix_time > max_unix {
		return false
	}

	return true
}
