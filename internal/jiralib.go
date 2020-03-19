/*
Copyright © 2020 Sion Williams <sion(at)nullbyteltd.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package jiralib

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
)

// NewJiraClient will return a valid JIRA api client.
func NewJiraClient(instanceURL, username, password string) (*jira.Client, error) {
	client, err := jira.NewClient(nil, instanceURL)
	if err != nil {
		return nil, fmt.Errorf("JIRA client can`t be initialized: %s", err)
	}

	// Only provide authentification if a username and password was applied
	if len(username) > 0 && len(password) > 0 {
		ok, err := client.Authentication.AcquireSessionCookie(username, password)
		if ok == false || err != nil {
			return nil, fmt.Errorf("jiractl can`t authentificate user %s against the JIRA instance %s: %s", username, instanceURL, err)
		}
	}

	return client, nil
}

// SearchInactiveIssues will search for issues in a project that have not changed in n days
func SearchInactiveIssues(client *jira.Client, project string, days int) ([]jira.Issue, error) {
	query := fmt.Sprintf("project = %s and updated <= -%dd", project, days)
	iss, resp, err := client.Issue.Search(query, &jira.SearchOptions{})
	if c := resp.StatusCode; err != nil || (c < 200 || c > 299) {
		return iss, fmt.Errorf("JIRA Search for issue returned %s (%d)", resp.Status, resp.StatusCode)
	}

	return iss, nil
}

// AddComment adds a comment to the issue whos issueID matches what was passed in.
func AddComment(jiraClient *jira.Client, issueID string, body string) error {
	c := &jira.Comment{Body: body}
	_, resp, err := jiraClient.Issue.AddComment(issueID, c)
	if rc := resp.StatusCode; err != nil || (rc < 200 || rc > 299) {
		return fmt.Errorf("JIRA AddComment unsuccessful, got %s (%d)", resp.Status, resp.StatusCode)
	}
	return nil
}
