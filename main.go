package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GetTodo struct {
	ID          int
	Description string
}
type Todo struct {
	Description string `json:"description"`
}
type DeleteTodo struct {
	ID int
}

var todolist []Todo

func main() {
	mux := NewMux()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/view", view)
	mux.HandleFunc("/save", save)
	mux.HandleFunc("/delete", delete)
	return mux
}

func view(w http.ResponseWriter, r *http.Request) {
	// want := []Todo{
	// 	{ID: 1, Description: "first todo."},
	// }
	w.Header().Set("Content-Type", "application/json")
	resTodo := []GetTodo{}
	for i, todo := range todolist {
		rs := GetTodo{ID: i, Description: todo.Description}
		resTodo = append(resTodo, rs)
	}
	json, err := json.Marshal(resTodo)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status": "%s"}`, err), http.StatusInternalServerError)
	}
	w.Write(json)
}

func save(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var addTodo Todo
		if err := json.NewDecoder(r.Body).Decode(&addTodo); err != nil {
			http.Error(w, fmt.Sprintf(`{"status":"Can't Encord Request Body %v"}`, err), http.StatusBadRequest)
		}
		defer r.Body.Close()
		todolist = append(todolist, addTodo)
	default:
		http.Error(w, fmt.Sprint(`{"status":"not allow method"}`), http.StatusMethodNotAllowed)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		var dt DeleteTodo
		if err := json.NewDecoder(r.Body).Decode(&dt); err != nil {
			http.Error(w, fmt.Sprintf(`{"status":"Can't Encord Request Body %v"}`, err), http.StatusBadRequest)
		}
		if dt.ID < 0 || dt.ID >= len(todolist) {
			http.Error(w, fmt.Sprint(`{"status":"this id is Out Of bound"}`), http.StatusBadRequest)
		}
		todolist = append(todolist[:dt.ID], todolist[dt.ID+1:]...)

	default:
		http.Error(w, fmt.Sprint(`{"status":"not allow method"}`), http.StatusMethodNotAllowed)
	}

}
