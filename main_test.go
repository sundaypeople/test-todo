package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestView(t *testing.T) {
	ts := httptest.NewServer(NewMux())
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
	ts := httptest.NewServer(NewMux())
	defer ts.Close()
	t.Run("testSave", func(t *testing.T) {
		send := Todo{Description: "first todo."}

		body, err := json.Marshal(send)
		if err != nil {
			t.Fatalf("Can't Json Unmarshal: %v", err)
		}
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
			{ID: 0, Description: "first todo."},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		t.Logf("got %v, want %v", got, want)
	})
}

func sendRequest(method string, url string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}
	client := new(http.Client)

	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resbody, _ := io.ReadAll(res.Body)
	return res.StatusCode, resbody, nil
}

func TestDelete(t *testing.T) {
	ts := httptest.NewServer(NewMux())
	defer ts.Close()
	t.Run("testSave", func(t *testing.T) {
		send := Todo{Description: "second todo."}
		body, err := json.Marshal(send)
		if err != nil {
			t.Fatalf("Can't Json Unmarshal: %v", err)
		}
		sc, b, err := sendRequest(http.MethodPost, ts.URL+"/save", body)
		if err != nil {
			t.Fatalf("Can't Send Http Request: %v", err)
		}

		if sc != 200 {
			t.Fatalf("Http Status Code: got %d, want %d  body:%v", sc, 200, string(b))
		}

		did := DeleteTodo{ID: 0}
		body, err = json.Marshal(did)
		if err != nil {
			t.Fatalf("Can't Json Unmarshal: %v", err)
		}

		sc, _, err = sendRequest(http.MethodDelete, ts.URL+"/delete", body)
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
			{ID: 0, Description: "second todo."},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}
