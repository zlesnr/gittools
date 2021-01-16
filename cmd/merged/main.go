package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var outputFormat string

// var outputPlain bool

var (
	appName        = "git-merged"
	version        = "unknown"
	mainBranchName = "master"
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
	f := filepath.Base(".")

	repo, err := git.PlainOpen(f)
	if err != nil {
		log.Fatal(err)
	}

	revHash, err := repo.ResolveRevision(plumbing.Revision(mainBranchName))
	if err != nil {
		log.Fatal(err)
	}

	revCommit, err := repo.CommitObject(*revHash)
	if err != nil {
		log.Fatal(err)
	}

	refs, err := repo.Branches()
	if err != nil {
		log.Fatal(err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.AppendHeader(table.Row{"Branch", fmt.Sprintf("Merged to %s", mainBranchName)})

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		// The HEAD is omitted in a `git show-ref` so we ignore the symbolic
		// references, the HEAD
		if ref.Type() == plumbing.SymbolicReference {
			return nil
		}

		if ref.Name().IsBranch() || ref.Name().IsRemote() {

			headCommit, err := repo.CommitObject(ref.Hash())
			if err != nil {
				log.Fatal(err)
			}

			isAncestor, err := headCommit.IsAncestor(revCommit)
			if err != nil {
				log.Fatal(err)
			}

			log.WithFields(log.Fields{
				"name":       ref.Name().Short(),
				"isAncestor": isAncestor,
			}).Debug("ref")

			t.AppendRow([]interface{}{ref.Name().Short(), isAncestor})

			// t.AppendRow(table.Row{ref.Name().Short(), isAncestor}, rowConfigAutoMerge)
		}

		return nil
	})

	t.Render()
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
