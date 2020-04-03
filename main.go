package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    "Task 1",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte datos validos")
	}
	json.Unmarshal(reqBody, &newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Id inválido.")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func putTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Id inválido.")
		return
	}

	var taskBody task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte datos validos")
	}
	json.Unmarshal(reqBody, &taskBody)
	w.Header().Set("Content-Type", "application/json")

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			taskBody.ID = taskID
			tasks = append(tasks, taskBody)
			fmt.Fprintf(w, "Task Id %v has been updated succesfully", task.ID)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Id inválido.")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "Task Id %v has been remove succesfully", task.ID)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi rest api 2")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", putTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))

}
