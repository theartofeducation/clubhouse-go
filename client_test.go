package clubhouse

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var options = Options{Token: "abc123"}

func Test_NewClient(t *testing.T) {
	t.Run("it creates and returns a new client", func(t *testing.T) {
		client := NewClient(options)

		if _, ok := client.(CHClient); !ok {
			t.Errorf("client is not a CHClient")
		}
	})
}

func Test_CreateEpic(t *testing.T) {
	t.Run("it creates a new epic", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusCreated)
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		epic, err := client.CreateEpic("Epic 1", "A test epic")

		if err != nil {
			t.Fatalf("received error when not expecting one: %s", err)
		}

		if epic.Name != "Epic 1" {
			t.Errorf("epic has unexpected name: got %q want %q", epic.Name, "Epic 1")
		}

		if epic.Description != "A test epic" {
			t.Errorf("epic has unexpected name: got %q want %q", epic.Description, "A test epic")
		}
	})

	t.Run("it handles HTTP error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusUnauthorized)
		}))
		defer testServer.Close()

		client := Client{url: testServer.URL}

		_, err := client.CreateEpic("Epic 1", "A test epic")

		if err == nil {
			t.Fatal("did not receive error when expecting one")
		}
	})
}
