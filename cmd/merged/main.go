package main

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	ci "github.com/xaque208/znet/pkg/continuous"

	"github.com/spf13/cobra"
)

var outputFormat string

// var outputPlain bool

var (
	appName = "git-merged"
	version = "unknown"
)

// Command represents the base command when called without any subcommands
var mergedCmd = &cobra.Command{
	Use:               appName,
	Short:             "git-merged",
	Long:              `git-merged`,
	Version:           version,
	DisableAutoGenTag: true, // Do not print generation date on documentation
	Run:               run,
}

var verbose bool
var trace bool

func run(cmd *cobra.Command, args []string) {
	fmt.Printf("%s version %s\n", appName, version)

	f := filepath.Base(".")

	r, err := git.PlainOpen(f)
	if err != nil {
		log.Fatal(err)
	}

	heads, _, err := ci.RepoRefs(r)
	if err != nil {
		log.Error(err)
	}

	log.Infof("refs: %+v", heads)

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	mergedCmd.Use = appName
	mergedCmd.Version = version

	// Silence Cobra's internal handling of error messaging
	// since we have a custom error handler in main.go
	mergedCmd.SilenceErrors = true

	// return mergedCmd.Execute()

}

func init() {
	cobra.OnInitialize(initConfig)

	mergedCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity")
	mergedCmd.PersistentFlags().BoolVarP(&trace, "trace", "", false, "Trace level verbosity")

	mergedCmd.PersistentFlags().StringVar(&outputFormat, "format", "TABLE", "output text format [TABLE]")
	// Command.PersistentFlags().BoolVar(&outputPlain, "plain", false, "output compact text")

	formatter := log.TextFormatter{
		FullTimestamp: true,
	}

	log.SetFormatter(&formatter)

	if trace {
		log.SetLevel(log.TraceLevel)
	} else if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func initConfig() {
	// utils.LogIfError(output.SetFormat(output.ParseFormat(outputFormat)))
	// utils.LogIfError(output.SetPrettyPrint(!outputPlain))
}

func main() {
	mergedCmd.Execute()
}
