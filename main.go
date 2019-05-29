package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/minio/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

// Help template for minl
var minLHelpTemplate = `NAME:
  {{.Name}} - {{.Usage}}

USAGE:
  {{.Name}} {{if .Flags}}[FLAGS] {{end}}COMMAND{{if .Flags}} [COMMAND FLAGS | -h]{{end}} [ARGUMENTS...]

COMMANDS:
  {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}{{if .Flags}}
GLOBAL FLAGS:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
VERSION:
  ` + minLVersion +
	`{{ "\n"}}{{range $key, $value := ExtraInfo}}
{{$key}}:
  {{$value}}
{{end}}`

func registerApp() *cli.App {
	// Register all the commands (refer commands.go)
	registerCmd(genCmd)
	registerCmd(versionCmd)
	
	// Set up app.
        cli.HelpFlag = cli.BoolFlag{
                Name:  "help, h",
                Usage: "show help",
        }
	
	app := cli.NewApp()
	app.Usage = "MinIO Lambda Functions"
	app.Author = "MinIO, Inc"
	app.HideHelpCommand = true // Hide `help, h` command, we already have `minl --help`.
	app.Flags = globalFlags
	app.Commands = commands
	app.CustomAppHelpTemplate = minLHelpTemplate
	app.CommandNotFound = commandNotFound

	return app
}

// Get os/arch/platform specific information.
// Returns a map of current os/arch/platform/memstats.
func getSystemData() map[string]string {
	host, err := os.Hostname()
	if err != nil {
		fmt.Println("Unable to determine the hostname.", err)
		os.Exit(1)
	}

	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)
	mem := fmt.Sprintf("Used: %s | Allocated: %s | UsedHeap: %s | AllocatedHeap: %s",
		pb.Format(int64(memstats.Alloc)).To(pb.U_BYTES),
		pb.Format(int64(memstats.TotalAlloc)).To(pb.U_BYTES),
		pb.Format(int64(memstats.HeapAlloc)).To(pb.U_BYTES),
		pb.Format(int64(memstats.HeapSys)).To(pb.U_BYTES))
	platform := fmt.Sprintf("Host: %s | OS: %s | Arch: %s", host, runtime.GOOS, runtime.GOARCH)
	goruntime := fmt.Sprintf("Version: %s | CPUs: %s", runtime.Version(), strconv.Itoa(runtime.NumCPU()))
	return map[string]string{
		"PLATFORM": platform,
		"RUNTIME":  goruntime,
		"MEM":      mem,
	}
}

func registerBefore(ctx *cli.Context) error {
	// Check if mc was compiled using a supported version of Golang.
	checkGoVersion()

	// Set global flags.
	setGlobalsFromContext(ctx)

	return nil
}

func main() {
	app := registerApp()
	app.Before = registerBefore

	app.ExtraInfo = func() map[string]string {
		if _, e := pb.GetTerminalWidth(); e != nil {
			globalQuiet = true
		}
		if globalDebug {
			return getSystemData()
		}
		return make(map[string]string)
	}

	app.RunAndExitOnError()
}
