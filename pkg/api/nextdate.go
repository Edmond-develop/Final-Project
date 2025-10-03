package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowTime, err := time.Parse(DateFormat, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	next, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, next)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}
	sliceRep := strings.Split(repeat, " ")

	switch sliceRep[0] {
	case "y":
		return checkYear(start, now, sliceRep)
	case "d":
		return checkDay(start, now, sliceRep)
	case "w":
		return checkWeek(start, now, sliceRep)
	case "m":
		return checkMonth(start, now, sliceRep)
	default:
		return "", fmt.Errorf("Wrong repeat format")
	}

	return start.Format(DateFormat), nil
}
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func checkYear(start, now time.Time, sliceRep []string) (string, error) {
	if len(sliceRep) != 1 {
		return "", fmt.Errorf("Wrong year format")
	}
	year := start.Year()
	month := start.Month()
	day := start.Day()

	next := time.Date(year, month, day, 0, 0, 0, 0, start.Location())
	next = next.AddDate(1, 0, 0)
	for !afterNow(next, now) {
		next = next.AddDate(1, 0, 0)
	}
	if next.Month() != month || next.Day() != day {
		next = time.Date(next.Year(), 3, 1, 0, 0, 0, 0, start.Location())
	}
	return next.Format(DateFormat), nil
}

func checkDay(start, now time.Time, sliceRep []string) (string, error) {
	if len(sliceRep) != 2 {
		return "", fmt.Errorf("Wrong day format")
	}
	num, err := strconv.Atoi(sliceRep[1])
	if err != nil || (num <= 0 || num > 400) {
		return "", fmt.Errorf("Wrong day num format")
	}
	for {
		start = start.AddDate(0, 0, num)
		if afterNow(start, now) {
			break
		}
	}
	return start.Format(DateFormat), nil
}

func checkWeek(start, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("Wrong repeat format")
	}
	weekNum := strings.Split(parts[1], ",")

	days := make(map[int]bool)

	for _, d := range weekNum {
		day, err := strconv.Atoi(d)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("Wrong week format")
		}
		days[day] = true
	}

	for {
		if afterNow(start, now) {
			weekday := int(start.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if days[weekday] {
				return start.Format(DateFormat), nil
			}
		}
		start = start.AddDate(0, 0, 1)
	}
}

func checkMonth(start, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 && len(parts) != 3 {
		return "", fmt.Errorf("Wrong repeat format")
	}
	sliceDays := strings.Split(parts[1], ",")
	days := make(map[int]bool)
	for _, d := range sliceDays {
		numDay, err := strconv.Atoi(d)
		if err != nil || numDay < -2 || numDay == 0 || numDay > 31 {
			return "", fmt.Errorf("Wrong repeat format")
		}
		days[numDay] = true
	}

	months := [13]bool{}
	if len(parts) == 3 {
		sliceMonth := strings.Split(parts[2], ",")
		for _, m := range sliceMonth {
			monthNum, err := strconv.Atoi(m)
			if err != nil || monthNum > 12 || monthNum < 1 {
				return "", fmt.Errorf("Wrong repeat months format")
			}
			months[monthNum] = true
		}
	} else {
		for i := 1; i < 13; i++ {
			months[i] = true
		}
	}
	if !CheckFormatMonth(days, months) {
		return "", fmt.Errorf("invalid day for selected months")
	}

	for {
		if afterNow(start, now) {
			day := start.Day()
			month := int(start.Month())
			lastDay := time.Date(start.Year(), start.Month()+1, 0, 0, 0, 0, 0, start.Location()).Day()
			if days[-1] && day == lastDay && months[month] {
				return start.Format(DateFormat), nil
			} else if days[-2] && day == lastDay-1 && months[month] {
				return start.Format(DateFormat), nil
			} else if day <= lastDay && days[day] && months[month] {
				return start.Format(DateFormat), nil
			}
		}
		start = start.AddDate(0, 0, 1)
	}
}

func CheckFormatMonth(days map[int]bool, months [13]bool) bool {
	var daysInMonth = map[int]int{
		1: 31, 2: 29, 3: 31, 4: 30,
		5: 31, 6: 30, 7: 31, 8: 31,
		9: 30, 10: 31, 11: 30, 12: 31,
	}

	for d := range days {
		if d > 0 {
			valid := false
			for i := 1; i < 13; i++ {
				if months[i] && d <= daysInMonth[i] && days[d] {
					valid = true
					break
				}
			}
			if !valid {
				return false
			}
		}
	}
	return true
}
