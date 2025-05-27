package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Task struct {
    Nom        string `json:"nom"`
    Description string `json:"description"`
    Etat       bool   `json:"etat"`
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "DeleteTask") 
}

func addTask(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "AddTask")
}

func showTask(w http.ResponseWriter, r *http.Request) {
    tasks := map[int]Task{
        1: {"Truc", "Desctip", false},
        2: {"Truc", "Desctip", true},
        3: {"Truc", "Desctip", false},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "UpdateTask")
}

func main() {
    http.HandleFunc("/add", addTask)
    http.HandleFunc("/update", updateTask)
    http.HandleFunc("/delete", deleteTask)
    http.HandleFunc("/show", showTask)
    fmt.Println("Serveur démarré sur :8080")
    http.ListenAndServe(":8080", nil)
}
