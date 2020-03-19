package jiralib_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andygrunwald/go-jira"
	jiralib "github.com/willis7/jiractl/internal"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server

	// jiraClient allows us to connect to the http server and get a response
	jiraClient *jira.Client

	jiraURL string
)

// setup sets up a test HTTP server along with a jira.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()

	apiHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, ``)
	})

	server = httptest.NewServer(apiHandler)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func init() {
	flag.StringVar(&jiraURL, "url", "", "JIRA instance URL (format: scheme://[username[:password]@]host[:port]/).")
}

func TestSearchInactiveIssues(t *testing.T) {
	if jiraURL == "" {
		t.Skip("Skipping: no jira instance url provided")
	}

	jc, err := jiralib.NewJiraClient(jiraURL, "", "")
	if err != nil {
		t.Errorf("NewJiraClient expected no errors, got %s", err)
	}
	// when I call the SearchInactiveIssues method,
	daysUnchanged := 2
	issues, err := jiralib.SearchInactiveIssues(jc, "JIR", daysUnchanged)
	if err != nil {
		t.Errorf("SearchInactiveIssues expected no errors, got %s", err)
	}

	mostRecentDate := time.Now().AddDate(0, 0, -daysUnchanged)
	for _, iss := range issues {
		t.Run(fmt.Sprintf("Test: %s", iss.Key), func(t *testing.T) {
			if updated := time.Time(iss.Fields.Updated); updated.After(mostRecentDate) {
				t.Errorf("Date %s after %s", updated, mostRecentDate)
			}
		})
	}
}

// Working experiment for adding a comment to a task
func TestAddComment(t *testing.T) {
	if jiraURL == "" {
		t.Skip("Skipping: no jira instance url provided")
	}

	jc, err := jiralib.NewJiraClient(jiraURL, "", "")
	if err != nil {
		t.Errorf("NewJiraClient expected no errors, got %s", err)
	}

	key := "JIR-6"
	body := "TEST: body"
	err = jiralib.AddComment(jc, key, body)
	if err != nil {
		t.Errorf("failed to add comment %s", err)
	}
}
