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

	"github.com/spf13/cobra"
)

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nudgeCmd.PersistentFlags().String("foo", "", "A help for foo")
	nudgeCmd.Flags().String("project", "", "")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nudgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func nudgeAction(cmd *cobra.Command, args []string) error {
	fmt.Println("nudge called")

	return nil
}
