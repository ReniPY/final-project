package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	dateStartOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowStartOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateStartOfDay.After(nowStartOfDay)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch len(parts) {
	case 0:
		return "", nil
	case 1:
		if parts[0] == "y" {
			for {
				date = date.AddDate(1, 0, 0)
				if afterNow(date, now) {
					break
				}
			}
			return date.Format(DateFormat), nil
		}
		return "", fmt.Errorf("Некорректный параметр repeat: '%s'", repeat)
	default:
		if parts[0] == "d" {
			interval, err := strconv.Atoi(parts[1])
			if err != nil || interval <= 0 || interval > 400 {
				return "", nil
			}
			for {
				date = date.AddDate(0, 0, interval)
				if afterNow(date, now) {
					break
				}
			}
			return date.Format(DateFormat), nil
		}
		return "", fmt.Errorf("Некорректный параметр repeat: '%s'", repeat)
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowParsed, err := time.Parse(DateFormat, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := NextDate(nowParsed, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(result)) // Исправлено: добавили обработку ошибки
	if err != nil {
		http.Error(w, "Ошибка записи ответа", http.StatusInternalServerError)
	}
}
