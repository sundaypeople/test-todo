package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"todo-handson/ent"

	_ "github.com/mattn/go-sqlite3"
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

type Server struct {
	templates *template.Template
	client    *ent.Client
}

func main() {
	client, err := ent.Open("sqlite3", "file:todo.db?_fk=1")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Close()
	// DB のマイグレーション
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
	log.Println("Database schema created successfully!")
	t := template.Must(template.ParseFiles("index.html"))

	s := Server{
		templates: t,
		client:    client,
	}

	mux := NewMux(s)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
func NewMux(s Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/view", s.view)
	mux.HandleFunc("/save", s.save)
	mux.HandleFunc("/delete", s.delete)
	return mux
}

func (s Server) index(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items, err := s.client.Todo.Query().All(ctx)
	if err != nil {
		cancel()
		log.Fatalf("failed querying todos: %v", err)
	}
	resTodo := []GetTodo{}
	for _, todo := range items {
		rs := GetTodo{ID: todo.ID, Description: todo.Description}
		resTodo = append(resTodo, rs)
	}
	err = s.templates.ExecuteTemplate(w, "index.html", resTodo)
	if err != nil {
		cancel()
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s Server) view(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Header().Set("Content-Type", "application/json")

	items, err := s.client.Todo.Query().All(ctx)
	if err != nil {
		cancel()
		log.Fatalf("failed querying todos: %v", err)
	}

	resTodo := []GetTodo{}
	for _, todo := range items {
		rs := GetTodo{ID: todo.ID, Description: todo.Description}
		resTodo = append(resTodo, rs)
	}
	json, err := json.Marshal(resTodo)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status": "%s"}`, err), http.StatusInternalServerError)
	}

	w.Write(json)
}

func (s Server) save(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf(`{"status":"Can't Parse Form %v"}`, err), http.StatusBadRequest)
		}
		description := r.Form["description"]
		fmt.Print(r.Form)
		if len(description) == 0 {
			http.Error(w, `{"status":"missing id parameter"}`, http.StatusBadRequest)
			return
		}

		_, err := s.client.Todo.Create().SetDescription(description[0]).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating a todo: %v", err)
		}
		defer r.Body.Close()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, fmt.Sprint(`{"status":"not allow method"}`), http.StatusMethodNotAllowed)
	}
}

func (s Server) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf(`{"status":"Can't Parse Form %v"}`, err), http.StatusBadRequest)
		}
		if len(r.Form) > 1 {
			http.Error(w, fmt.Sprintf(`{"status":"Too many form body"}`), http.StatusBadRequest)
		}
		ids := r.Form["id"]
		fmt.Print(r.Form)
		if len(ids) == 0 {
			http.Error(w, `{"status":"missing id parameter"}`, http.StatusBadRequest)
			return
		}
		did, err := strconv.Atoi(r.Form["id"][0])
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"status":"Can't encord string to int %v"}`, err), http.StatusBadRequest)
			return
		}

		err = s.client.Todo.DeleteOneID(did).Exec(ctx)
		if err != nil {
			cancel()
			if ent.IsNotFound(err) {
				// 見つからない場合の特別処理
				log.Printf("todo id=%d not found", did)
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprint(`{"status":"DB Error"}`), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, fmt.Sprint(`{"status":"not allow method"}`), http.StatusMethodNotAllowed)
	}

}
