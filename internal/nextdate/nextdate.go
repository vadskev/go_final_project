package nextdate

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	d, err := time.Parse(dateLayout, date)
	if err != nil {
		return "invalid convert to date", err
	}

	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	switch repeat[0] {
	case 'y':
		return hYear(now, d)
	case 'm':
		return hMonth(now, d, repeat)
	case 'w':
		return hWeek(now, d, repeat)
	case 'd':
		return hDay(now, d, repeat)
	default:
		return "", fmt.Errorf("format date error: %s", repeat)
	}
}

func hYear(now time.Time, d time.Time) (string, error) {
	next := d.AddDate(1, 0, 0)
	for !next.After(now) {
		next = next.AddDate(1, 0, 0)
	}
	return next.Format(dateLayout), nil
}

func hMonth(now, d time.Time, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat[1:])
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 || len(parts) > 2 {
		return "", fmt.Errorf("format error: %s", repeat)
	}

	daysPart := strings.Split(parts[0], ",")
	daysMap := make(map[int]bool)
	for _, day := range daysPart {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt < -2 || dayInt == 0 || dayInt > 31 {
			return "", fmt.Errorf("format error: %s", day)
		}
		daysMap[dayInt] = true
	}

	monthsMap := make(map[int]bool)
	if len(parts) == 2 {
		for _, m := range strings.Split(parts[1], ",") {
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("format error: %s", parts[1])
			}
			monthsMap[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			monthsMap[i] = true
		}
	}

	for next := d; ; next = next.AddDate(0, 0, 1) {
		day := next.Day()
		month := int(next.Month())
		if daysMap[day] || daysMap[day-daysInMonth(next.Month(), next.Year())-1] {
			if monthsMap[month] {
				if next.After(now) {
					return next.Format(dateLayout), nil
				}
			}
		}
		if next.Year() > now.Year()+1 {
			break
		}
	}
	return "", errors.New("date not found")
}

func hWeek(now, d time.Time, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat[1:])
	parts := strings.Split(repeat, ",")
	if len(parts) == 0 {
		return "", fmt.Errorf("format error: %s", repeat)
	}

	daysOfWeek := []time.Weekday{}
	for _, part := range parts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("format error: %s", repeat)
		}
		daysOfWeek = append(daysOfWeek, time.Weekday(day%7))
	}

	sort.Slice(daysOfWeek, func(i, j int) bool {
		return daysOfWeek[i] < daysOfWeek[j]
	})

	next := findNextWeek(d, daysOfWeek)
	for !next.After(now) {
		next = findNextWeek(next.AddDate(0, 0, 1), daysOfWeek)
	}

	return next.Format(dateLayout), nil
}

func hDay(now, d time.Time, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("format error: %s", repeat)
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("format error: %s", repeat)
	}
	if days < 1 || days > 400 {
		return "", fmt.Errorf("interval not from 1 to 400")
	}
	next := d.AddDate(0, 0, days)
	for !next.After(now) {
		next = next.AddDate(0, 0, days)
	}
	return next.Format(dateLayout), nil
}

func findNextWeek(start time.Time, daysOfWeek []time.Weekday) time.Time {
	for _, day := range daysOfWeek {
		if start.Weekday() <= day {
			return start.AddDate(0, 0, int(day-start.Weekday()))
		}
	}
	return start.AddDate(0, 0, int(7-start.Weekday()+daysOfWeek[0]))
}

func daysInMonth(month time.Month, year int) int {
	switch month {
	case time.February:
		if checkLeapYear(year) {
			return 29
		}
		return 28
	case time.April, time.June, time.September, time.November:
		return 30
	default:
		return 31
	}
}

func checkLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
