# Jiractl

'jiractl' is short for Jira Controller.


## Usage
```
$ jiractl help

Jiractl is a CLI for Jira housekeeping tasks.
	Jiractl uses the REST API to control your instance and execute
	common tasks.

Usage:
  jiractl [command]

Available Commands:
  close       Close issues that have not been actioned since warning.
  help        Help about any command
  nudge       Prompt issue watchers to action their tickets.

Flags:
      --config string   config file (default is $HOME/.jiractl.yaml)
  -d, --dry-run         No operation dry run
  -h, --help            help for jiractl
      --pass string     JIRA Password.
  -t, --toggle          Help message for toggle
      --url string      JIRA instance URL (format: scheme://[username[:password]@]host[:port]/).
      --user string     JIRA Username.

Use "jiractl [command] --help" for more information about a command.
subcommand is required
```

### Config

This CLI supports the use of a configuration file for persisting connection details.

Example `.jiractl.yml`

``` yaml
---
url: "https://mydomain.atlassian.net/"
user: "sion@mydomain.com"
pass: "DYPMYSUPERSECRETPASSWORDKEYTHINGYE20"
```

## Development

See the [Makefile](Makefile) for a complete list of commands.


``` bash
# test runs the unit tests
$ make test

# testrace runs the race checker
$ make testrace

# dev creates binaries for testing Jiractl locally. These are put
# into ./bin/ as well as $GOPATH/bin
$ make dev
```