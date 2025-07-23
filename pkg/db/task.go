package db

import (
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Comment string `json:"comment"`
	Title   string `json:"title"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `
        INSERT INTO scheduler (date, comment, title, repeat)
        VALUES(?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Comment, task.Title, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query("SELECT id, date, comment, title, repeat FROM scheduler LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Comment, &task.Title, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = (append(tasks, task))

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := &Task{}

	err := db.QueryRow("SELECT * FROM scheduler WHERE id=$1", id).Scan(&task.ID, &task.Date, &task.Comment, &task.Title, &task.Repeat)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=$1, comment=$2, title=$3, repeat=$4 WHERE id=$5`
	res, err := db.Exec(query, task.Date, task.Comment, task.Title, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id=$1`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for deleting task")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date=$1 WHERE id=$2`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task's date")
	}
	return nil
}
