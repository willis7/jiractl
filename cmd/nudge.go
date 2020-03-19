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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	jiralib "github.com/willis7/jiractl/internal"
)

const nudgeMsg = "NUDGE:: this issue has been inactive for a while and is targetted for closure."

// nudgeCmd represents the nudge command
var nudgeCmd = &cobra.Command{
	Use:   "nudge",
	Short: "Prompt issue watchers to action their tickets.",
	Long: `Look for issues that have not seen any activity in the past
	x days and add a comment warning of closure.`,
	RunE: nudgeAction,
}

func init() {
	rootCmd.AddCommand(nudgeCmd)

	nudgeCmd.Flags().StringP("project", "p", "", "The Project to scan")
	nudgeCmd.MarkFlagRequired("project")
	nudgeCmd.Flags().IntP("days", "d", 14, "Inactive for the past x days")
}

// nudgeAction will search a project for issues that have not been updated within
// a given number of days. Once identified, they will be printed and a comment will
// be added to the issue unless the `--no-op` flag is added.
func nudgeAction(cmd *cobra.Command, args []string) error {
	// Get all the flag values
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return err
	}

	days, err := cmd.Flags().GetInt("days")
	if err != nil {
		return err
	}

	noop, err := cmd.Flags().GetBool("no-op")
	if err != nil {
		return err
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("project: %s\n", project)
		fmt.Printf("days: %d\n", days)
		fmt.Printf("noop: %t\n", noop)
	}

	// TODO: DRY up this code
	// REMINDER: when I created this as a global var I got
	// a panic
	jiraClient, err := jiralib.NewJiraClient(jiraURL, jiraUsername, jiraPassword)
	if err != nil {
		fmt.Printf("Unable to intialise Jira client: %s\n", err)
		os.Exit(1)
	}

	issues, err := jiralib.SearchInactiveIssues(jiraClient, project, days)
	if err != nil {
		return fmt.Errorf("Failed to retrieve inactive issues; %s", err)
	}

	// Iterate over the matched issues and add the nude message.
	for _, issue := range issues {
		fmt.Printf("%s :: %s\n", issue.Key, issue.Fields.Summary)

		if !noop {
			_, err = jiralib.AddComment(jiraClient, issue.Key, nudgeMsg)
			if err != nil {
				fmt.Printf("Failed to add comment to: %s . ERROR: %s\n", issue.Key, err)
			}
		}
	}

	return nil
}
