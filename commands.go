package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/minio/cli"
)

// Collection of minl commands currently supported
var commands = []cli.Command{}

// Collection of minl commands currently supported in a trie tree
var commandsTree = newTrie()

// Collection of minl flags currently supported
var globalFlags = []cli.Flag{}

// registerCmd registers a cli command
func registerCmd(cmd cli.Command) {
	commands = append(commands, cmd)
	commandsTree.Insert(cmd.Name)
}

// findClosestCommands to match a given string with commands trie tree.
func findClosestCommands(command string) []string {
	var closestCommands []string
	for _, value := range commandsTree.PrefixMatch(command) {
		closestCommands = append(closestCommands, value.(string))
	}
	sort.Strings(closestCommands)
	// Suggest other close commands - allow missed, wrongly added and even transposed characters
	for _, value := range commandsTree.walk(commandsTree.root) {
		if sort.SearchStrings(closestCommands, value.(string)) < len(closestCommands) {
			continue
		}
		// 2 is arbitrary and represents the max allowed number of typed errors
		if DamerauLevenshteinDistance(command, value.(string)) < 2 {
			closestCommands = append(closestCommands, value.(string))
		}
	}
	return closestCommands
}

// Function invoked when invalid command is passed.
func commandNotFound(ctx *cli.Context, command string) {
	msg := fmt.Sprintf("‘%s’ is not a minl command. See ‘minl --help’.", command)
	closestCommands := findClosestCommands(command)
	if len(closestCommands) > 0 {
		msg += fmt.Sprintf("\n\nDid you mean one of these?\n")
		if len(closestCommands) == 1 {
			cmd := closestCommands[0]
			msg += fmt.Sprintf("        ‘%s’", cmd)
		} else {
			for _, cmd := range closestCommands {
				msg += fmt.Sprintf("        ‘%s’\n", cmd)
			}
		}
	}
	fmt.Println(msg)
	os.Exit(2)
}
