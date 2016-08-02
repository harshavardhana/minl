package main

import (
	"fmt"

	"github.com/minio/cli"
)

// Print version.
var versionCmd = cli.Command{
	Name:   "version",
	Usage:  "Print version.",
	Action: mainVersion,
	CustomHelpTemplate: `NAME:
   minl {{.Name}} - {{.Usage}}

USAGE:
   minl {{.Name}} [FLAGS]

FLAGS:
  {{range .Flags}}{{.}}
  {{end}}
`,
}

// Structured message depending on the type of console.
type versionMessage struct {
	Status  string `json:"status"`
	Version struct {
		Value  string `json:"value"`
		Format string `json:"format"`
	} `json:"version"`
	ReleaseTag string `json:"releaseTag"`
	CommitID   string `json:"commitID"`
}

// Colorized message for console printing.
func (v versionMessage) String() string {
	return (fmt.Sprintf("Version: %s\n", v.Version.Value) +
		fmt.Sprintf("Release-tag: %s\n", v.ReleaseTag) +
		fmt.Sprintf("Commit-id: %s", v.CommitID))
}

func mainVersion(ctx *cli.Context) {
	verMsg := versionMessage{}
	verMsg.CommitID = minLCommitID
	verMsg.ReleaseTag = minLReleaseTag
	verMsg.Version.Value = minLVersion
	verMsg.Version.Format = "RFC3339"
	fmt.Println(verMsg)
}
