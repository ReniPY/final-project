package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ReniPY/final-project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var task db.Task
	err := decoder.Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Неверный формат JSON"})
		return
	}

	// Проверка наличия заголовка задачи
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверка и нормализация даты
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавляем задачу в базу данных
	taskID, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка добавления задачи"})
		return
	}

	// Формируем успешный ответ с идентификатором задачи
	writeJson(w, map[string]int64{"id": taskID})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %w", err)
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("ошибка вычисления следующей даты: %w", err)
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
	}
}
