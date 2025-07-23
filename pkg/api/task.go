package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ReniPY/final-project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если задача не найдена, устанавливаем статус 404 Not Found
			w.WriteHeader(http.StatusNotFound)
			writeJson(w, map[string]string{"error": "Задача не найдена"})
		} else {
			// Если другая ошибка, устанавливаем статус 500 Internal Server Error
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, map[string]string{"error": "Ошибка сервера"})
		}
		return
	}

	// Если задача найдена, устанавливаем статус 200 OK
	w.WriteHeader(http.StatusOK)
	writeJson(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка обновления задачи"})
		return
	}

	// Формируем успешный ответ с идентификатором задачи
	writeJson(w, map[string]interface{}{"id": task.ID})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Параметр 'id' не указан"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		// Одноразовая задача, удаляем её
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "Ошибка удаления задачи"})
			return
		}
		writeJson(w, make(map[string]interface{})) // Пустой JSON при успешном удалении
	} else {
		// Периодическая задача, вычисляем следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "Ошибка вычисления следующей даты"})
			return
		}

		// Обновляем дату задачи
		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeJson(w, map[string]string{"error": "Ошибка обновления даты задачи"})
			return
		}
		writeJson(w, make(map[string]interface{})) // Пустой JSON при успешном переносе даты
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Параметр 'id' не указан"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, make(map[string]interface{}))
}
