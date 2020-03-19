package jiralib_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andygrunwald/go-jira"
	jiralib "github.com/willis7/jiractl/internal"
)

var (
	// testMux is the HTTP request multiplexer used with the test server.
	testMux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	testServer *httptest.Server

	// testClient allows us to connect to the test http server and get a response
	testClient *jira.Client

	jiraURL string
)

// setup sets up a test HTTP server along with a jira.Client that is
// configured to talk to that test server. Tests should register handlers on
// testMux which provide mock responses for the API method being tested.
func setup() {
	// test server
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)
	testClient, _ = jiralib.NewJiraClient(testServer.URL, "", "")
}

// teardown closes the test HTTP server.
func teardown() {
	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testRequestURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL.String(); !strings.HasPrefix(got, want) {
		t.Errorf("Request URL: %v, want %v", got, want)
	}
}

func init() {
	flag.StringVar(&jiraURL, "url", "", "JIRA instance URL (format: scheme://[username[:password]@]host[:port]/).")
}

func TestNewClient_WrongUrl(t *testing.T) {
	c, err := jiralib.NewJiraClient("://issues.apache.org/jira/", "", "")
	if err == nil {
		t.Error("Expected an error. Got none")
	}
	if c != nil {
		t.Errorf("Expected no client. Got %+v", c)
	}
}

func TestSearchInactiveIssues(t *testing.T) {
	setup()
	defer teardown()

	raw, err := ioutil.ReadFile("./mocks/TestSearchInactiveIssues.json")
	if err != nil {
		t.Error(err.Error())
	}

	testMux.HandleFunc("/rest/api/2/search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/search?jql=project+%3D+JIR+and+updated+%3C%3D+-2d")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(raw))
	})

	// when I call the SearchInactiveIssues method,
	daysUnchanged := 2
	issues, err := jiralib.SearchInactiveIssues(testClient, "JIR", daysUnchanged)
	if err != nil {
		t.Errorf("SearchInactiveIssues expected no errors, got %s", err)
	}

	if len(issues) != 2 {
		t.Errorf("Exected 2 issues, %d given", len(issues))
	}
}

func TestAddComment(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/issue/JIR-6/comment", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/rest/api/2/issue/JIR-6/comment")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/issue/10010/comment/JIR-6","id":"JIR-6","author":{"self":"http://www.example.com/jira/rest/api/2/user?username=fred","name":"fred","displayName":"Fred F. User","active":false},"body":"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque eget venenatis elit. Duis eu justo eget augue iaculis fermentum. Sed semper quam laoreet nisi egestas at posuere augue semper.","updateAuthor":{"self":"http://www.example.com/jira/rest/api/2/user?username=fred","name":"fred","displayName":"Fred F. User","active":false},"created":"2016-03-16T04:22:37.356+0000","updated":"2016-03-16T04:22:37.356+0000","visibility":{"type":"role","value":"Administrators"}}`)
	})

	key := "JIR-6"
	body := "TEST: body"
	comment, err := jiralib.AddComment(testClient, key, body)
	if comment == nil {
		t.Error("Expected Comment. Comment is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
