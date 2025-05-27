# 📝 API ToDoList – Documentation des routes
## Base URL
http://0.0.0.0:8080

## 🔧 Endpoints
### ➕ Ajouter une tâche
- URL : /add

- Méthode : POST

- Body (JSON) :

```json 
{
  "title": "Nom de la tâche",
  "description": "Détails optionnels"
}
``` 

Réponse :

```json 
{
  "status": "success",
  "message": "Tâche ajoutée avec succès",
  "id": 1
}
```

###  ✏️ Mettre à jour une tâche
- URL : /update

- Méthode : PUT

- Body (JSON) :

```json 
{
  "id": 1,
  "title": "Nom mis à jour",
  "description": "Nouvelle description"
}
```
Réponse :

```json
{
  "status": "success",
  "message": "Tâche mise à jour"
}
```
### ❌ Supprimer une tâche
- URL : /delete

- Méthode : DELETE

- Body (JSON) :

```json
{
  "id": 1
}
```

Réponse :


```json
{
  "status": "success",
  "message": "Tâche supprimée"
}
```

### 📋 Afficher toutes les tâches
- URL : /show

- Méthode : GET

Réponse :

```json
[
  {
    "id": 1,
    "title": "Faire les courses",
    "description": "Acheter du lait et du pain"
  },
  {
    "id": 2,
    "title": "Coder",
    "description": "Travailler sur le projet ToDoList"
  }
]
```
### 🛠️ Exemple de commande curl
```bash
curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"title":"Lire un livre","description":"Commencer un roman"}'
```
### 📌 Remarques

API simple sans authentification (utile pour du développement local ou un POC).

Ajoute facilement une route /done ou /toggle pour cocher les tâches terminées.

Possibilité d’évoluer vers une API RESTful avec structure REST plus conventionnelle.