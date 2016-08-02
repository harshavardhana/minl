package main

import "github.com/minio/cli"

var minGoVersion = ">= 1.6" // mc requires at least Go v1.6

var (
	globalQuiet = false // Quiet flag set via command line
	globalDebug = false // Debug flag set via command line
)

// Set global states. NOTE: It is deliberately kept monolithic to ensure we dont miss out any flags.
func setGlobals(quiet, debug bool) {
	globalQuiet = quiet
	globalDebug = debug
}

// Set global states. NOTE: It is deliberately kept monolithic to ensure we dont miss out any flags.
func setGlobalsFromContext(ctx *cli.Context) {
	quiet := ctx.Bool("quiet") || ctx.GlobalBool("quiet")
	debug := ctx.Bool("debug") || ctx.GlobalBool("debug")
	setGlobals(quiet, debug)
}
