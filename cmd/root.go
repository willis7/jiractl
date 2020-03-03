/*
Copyright © 2020 Sion Williams <sion@nullbyteltd.com>

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

	jira "github.com/andygrunwald/go-jira"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	jiraClient   *jira.Client
	jiraURL      string
	jiraUsername string
	jiraPassword string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jiractl",
	Short: "Jiractl is a CLI for Jira housekeeping tasks.",
	Long: `Jiractl is a CLI for Jira housekeeping tasks.
	Jiractl uses the REST API to control your instance and execute
	common tasks.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jiractl.yaml)")
	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "No operation dry run")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Bing flags with config
	rootCmd.PersistentFlags().StringVar(&jiraURL, "url", "", "JIRA instance URL (format: scheme://[username[:password]@]host[:port]/).")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))

	rootCmd.PersistentFlags().StringVar(&jiraUsername, "user", "", "JIRA Username.")
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))

	rootCmd.PersistentFlags().StringVar(&jiraPassword, "pass", "", "JIRA Password.")
	viper.BindPFlag("pass", rootCmd.PersistentFlags().Lookup("pass"))

	var err error
	jiraClient, err = NewJiraClient(jiraURL, jiraUsername, jiraPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jirabot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jiractl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

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
			return nil, fmt.Errorf("jiractl can`t authenticate user %s against the JIRA instance %s: %s", username, instanceURL, err)
		}
	}

	return client, nil
}
