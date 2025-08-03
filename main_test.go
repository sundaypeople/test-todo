package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"todo-handson/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
)

func TestView(t *testing.T) {
	temp := template.Must(template.ParseFiles("index.html"))
	client := enttest.Open(t, "sqlite3", "file:todo.db?mode=memory&_fk=1")
	defer client.Close()
	s := Server{templates: temp, client: client}
	ts := httptest.NewServer(NewMux(s))
	defer ts.Close()

	t.Run("testView", func(t *testing.T) {
		sc, resbody, err := sendRequest(http.MethodGet, ts.URL+"/view", nil)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}
		if sc != 200 {
			t.Fatalf("Http Status Code: got %v, want %v", sc, 200)
		}

		var got []GetTodo
		if err := json.Unmarshal(resbody, &got); err != nil {
			t.Fatalf("Can't Json Unmarshal: %v", err)
		}

		want := []GetTodo{}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}

func TestSave(t *testing.T) {
	temp := template.Must(template.ParseFiles("index.html"))
	client := enttest.Open(t, "sqlite3", "file:todo.db?mode=memory&_fk=1")
	defer client.Close()
	s := Server{templates: temp, client: client}
	ts := httptest.NewServer(NewMux(s))
	defer ts.Close()
	t.Run("testDelete", func(t *testing.T) {
		body := url.Values{"description": {"first todo."}}

		sc, b, err := sendRequest(http.MethodPost, ts.URL+"/save", body)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}
		if sc != 200 {
			t.Fatalf("Http Status Code: got %d, want %d  body:%v", sc, 200, string(b))
		}
		sc, resbody, err := sendRequest(http.MethodGet, ts.URL+"/view", nil)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}

		if sc != 200 {
			t.Fatalf("Http Status Code: got %v, want %v", sc, 200)
		}

		var got []GetTodo
		if err := json.Unmarshal(resbody, &got); err != nil {
			t.Fatalf("Can't Json Unmarshal: %v", err)
		}

		want := []GetTodo{
			{ID: 1, Description: "first todo."},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		t.Logf("got %v, want %v", got, want)
	})
}

func sendRequest(method string, url string, body url.Values) (int, []byte, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body.Encode()))
	if err != nil {
		return 0, nil, err
	}
	client := new(http.Client)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resbody, _ := io.ReadAll(res.Body)
	return res.StatusCode, resbody, nil
}

func TestDelete(t *testing.T) {
	temp := template.Must(template.ParseFiles("index.html"))
	client := enttest.Open(t, "sqlite3", "file:todo.db?mode=memory&_fk=1")
	defer client.Close()
	s := Server{templates: temp, client: client}
	ts := httptest.NewServer(NewMux(s))
	defer ts.Close()
	t.Run("testSave", func(t *testing.T) {
		body := url.Values{"description": {"first todo."}}
		sc, b, err := sendRequest(http.MethodPost, ts.URL+"/save", body)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}
		if sc != 200 {
			t.Fatalf("Http Status Code: got %d, want %d  body:%v", sc, 200, string(b))
		}

		body = url.Values{"description": {"second todo."}}
		sc, b, err = sendRequest(http.MethodPost, ts.URL+"/save", body)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}
		if sc != 200 {
			t.Fatalf("Http Status Code: got %d, want %d  body:%v", sc, 200, string(b))
		}

		body = url.Values{"id": {"2"}}

		sc, _, err = sendRequest(http.MethodPost, ts.URL+"/delete", body)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}
		if sc != 200 {
			t.Fatalf("Http Status Code: got %v, want %v", sc, 200)
		}

		sc, resbody, err := sendRequest(http.MethodGet, ts.URL+"/view", nil)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}
		if sc != 200 {
			t.Fatalf("Http Status Code: got %v, want %v", sc, 200)
		}
		var got []GetTodo
		if err := json.Unmarshal(resbody, &got); err != nil {
			t.Fatalf("Can't Json Unmarshal: %v", err)
		}

		want := []GetTodo{
			{ID: 1, Description: "first todo."},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}
