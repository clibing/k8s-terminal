package main

import (
	"github.com/k8s-terminal/cli"
	_ "github.com/k8s-terminal/prompt"
	"github.com/mattn/go-runewidth"
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

// See also: https://github.com/charmbracelet/lipgloss/issues/40#issuecomment-891167509
func init() {
	runewidth.DefaultCondition.EastAsianWidth = false
}