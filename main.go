package main

import (
	"github.com/k8s-terminal/cli"
	_ "github.com/k8s-terminal/prompt"
)

var (
	version   string
	buildDate string
	commitID  string
)

/**
 *
 */
func main() {
	cli.CliApp(version, buildDate, commitID)
}
