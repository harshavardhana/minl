package main

import "github.com/minio/cli"

// Generate lambda.
var genCmd = cli.Command{
	Name:   "gen",
	Usage:  "Generates lambda",
	Action: mainGen,
	CustomHelpTemplate: `NAME:
   minl {{.Name}} - {{.Usage}}

USAGE:
   minl {{.Name}} [FLAGS]

FLAGS:
  {{range .Flags}}{{.}}
  {{end}}
`,
}

func mainGen(c *cli.Context) {
	// Write your code here

}
