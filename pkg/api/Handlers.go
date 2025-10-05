package api

import (
	"encoding/json"
	"go1f/pkg/db"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Некорректный формат данных"})
		return
	}
	if _, err := time.Parse("20060102", task.Date); err != nil {
		writeJson(w, map[string]string{"error": "invalid date format"})
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Не указан id задачи"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указано название задачи"})
		return
	}

	if task.Date == "" {
		writeJson(w, map[string]string{"error": "Не указана дата"})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if task.Repeat != "" {
		date, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		err = db.UpdateDate(date, id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}
	writeJson(w, map[string]string{})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{})
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	type SigninRequest struct {
		Password string `json:"password"`
	}

	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || pass != req.Password {
		writeJson(w, map[string]string{"error": "wrong password"})
		return
	}
	hashPass := HashPassword(req.Password)
	claims := jwt.MapClaims{
		"password": hashPass,
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(pass))
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{"token": tokenString})
}
