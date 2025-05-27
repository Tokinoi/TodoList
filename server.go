package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Task struct {
    Nom         string `json:"nom"`
    Description string `json:"description"`
    Etat        bool   `json:"etat"`
}

type TaskWithID struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
}

var tasks = make(map[int]Task) // Simule une base de données
var currentID = 1              // Compteur d'ID auto-incrémenté

func addTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    var newTask Task
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&newTask)
    if err != nil {
        http.Error(w, "Corps de la requête invalide", http.StatusBadRequest)
        return
    }

    // Ajoute la tâche dans le map
    id := currentID
    currentID++
    tasks[id] = newTask

    // Réponse JSON
    response := map[string]interface{}{
        "status":  "success",
        "message": "Tâche ajoutée avec succès",
        "id":      id,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func showTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    tasksList := make([]TaskWithID, 0, len(tasks))
    for id, t := range tasks {
        tasksList = append(tasksList, TaskWithID{
            ID:          id,
            Title:       t.Nom,          // ici on fait la traduction Nom -> title pour l'API
            Description: t.Description,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasksList)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut && r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    // Structure pour lire les données entrantes avec l'ID
    var updated struct {
        ID          int    `json:"id"`
        Nom         string `json:"nom"`
        Description string `json:"description"`
        Etat        bool   `json:"etat"`
    }

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&updated)
    if err != nil {
        http.Error(w, "Corps de la requête invalide", http.StatusBadRequest)
        return
    }

    // Vérifie si la tâche existe
    task, exists := tasks[updated.ID]
    if !exists {
        http.Error(w, "Tâche non trouvée", http.StatusNotFound)
        return
    }

    // Met à jour la tâche
    task.Nom = updated.Nom
    task.Description = updated.Description
    task.Etat = updated.Etat
    tasks[updated.ID] = task

    // Réponse JSON
    response := map[string]interface{}{
        "status":  "success",
        "message": fmt.Sprintf("Tâche %d mise à jour avec succès", updated.ID),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func deleteTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete && r.Method != http.MethodGet {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    // Récupère l'id en paramètre
    ids, ok := r.URL.Query()["id"]
    if !ok || len(ids[0]) < 1 {
        http.Error(w, "Paramètre id manquant", http.StatusBadRequest)
        return
    }

    // Convertir l'id en int
    var id int
    _, err := fmt.Sscanf(ids[0], "%d", &id)
    if err != nil {
        http.Error(w, "Paramètre id invalide", http.StatusBadRequest)
        return
    }

    // Vérifie si la tâche existe
    if _, exists := tasks[id]; !exists {
        http.Error(w, "Tâche non trouvée", http.StatusNotFound)
        return
    }

    // Supprime la tâche
    delete(tasks, id)

    // Réponse JSON
    response := map[string]interface{}{
        "status":  "success",
        "message": fmt.Sprintf("Tâche %d supprimée avec succès", id),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func main() {
    http.HandleFunc("/add", addTask)
    http.HandleFunc("/update", updateTask)
    http.HandleFunc("/delete", deleteTask)
    http.HandleFunc("/show", showTask)
    fmt.Println("Serveur démarré sur :8080")
    http.ListenAndServe(":8080", nil)
}
