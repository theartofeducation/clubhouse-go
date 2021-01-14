package clubhouse

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

var options = Options{
	Token:         "abc123",
	WebhookSecret: "testsecret",
}

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

func Test_ParseWebhook(t *testing.T) {
	t.Run("it parses and returns a webhook", func(t *testing.T) {
		client := NewClient(options)

		body := ioutil.NopCloser(strings.NewReader(`{"actions": [{"entity_type": "epic", "action": "update", "name": "Test Epic", "changes": {"state": {"new": "done"}}}]}`))

		webhook, err := client.ParseWebhook(body)

		if err != nil {
			t.Fatalf("received error when parsing webhook: %s", err)
		}

		if webhook.Actions[0].EntityType != EntityTypeEpic {
			t.Errorf("webhook has unexpected entity type: got %q want %q", webhook.Actions[0].EntityType, EntityTypeEpic)
		}

		if webhook.Actions[0].Action != ActionUpdate {
			t.Errorf("webhook has unexpected action: got %q want %q", webhook.Actions[0].Action, ActionUpdate)
		}

		if Action(webhook.Actions[0].Name) != "Test Epic" {
			t.Errorf("webhook has unexpected name: got %q want %q", webhook.Actions[0].Name, "Test Epic")
		}

		if webhook.Actions[0].Changes.State.New != EpicStateDone {
			t.Errorf("webhook has unexpected epic state: got %q want %q", webhook.Actions[0].Changes.State.New, EpicStateDone)
		}
	})

	t.Run("it returns an error if webhook cannot be parsed", func(t *testing.T) {
		client := NewClient(options)

		body := ioutil.NopCloser(strings.NewReader(`{"actions":[{"entity_type": "epic", "action": "update"}]`))

		_, err := client.ParseWebhook(body)

		if err == nil {
			t.Fatal("did not receive error when expecting one")
		}

		want := errors.New("Could not parse Webhook body: unexpected EOF")

		if err.Error() != want.Error() {
			t.Errorf("received incorrect error: got %q want %q", err, want)
		}
	})
}

func Test_VerifySignature(t *testing.T) {
	t.Run("it verifies a valid signature", func(t *testing.T) {
		client := NewClient(options)

		body := []byte("ghi890")

		hash := hmac.New(sha256.New, []byte(options.WebhookSecret))
		hash.Write(body)
		signature := hex.EncodeToString(hash.Sum(nil))

		err := client.VerifySignature(signature, body)

		if err != nil {
			t.Errorf("valid signature was not validated")
		}
	})

	t.Run("it does not verify an invalid signature", func(t *testing.T) {
		client := NewClient(options)

		body := []byte("ghi890")

		hash := hmac.New(sha256.New, []byte("bad secret"))
		hash.Write(body)
		signature := hex.EncodeToString(hash.Sum(nil))

		err := client.VerifySignature(signature, body)

		if err.Error() != ErrSignatureMismatch.Error() {
			t.Errorf("correct error was not returned for invalid signature: got %s want %s", err, ErrSignatureMismatch)
		}
	})
}
