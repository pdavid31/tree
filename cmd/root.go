package cmd

import (
	"fmt"
	"os"

	"github.com/pdavid31/tree/internal"

	"github.com/spf13/cobra"
)

var (
	allFiles           bool
	directoriesOnly    bool
	disableIndentation bool
	fullPaths          bool
)

var rootCmd = &cobra.Command{
	Use:     "tree [flags] PATH",
	Short:   "List directory contents in a tree shape.",
	Long:    "tree lists the contents of directories in a tree-like format. It can be used to render the structure of your file system.",
	Args:    cobra.MaximumNArgs(1),
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		config := &internal.TreeConfig{
			AllFiles:           allFiles,
			DirectoriesOnly:    directoriesOnly,
			DisableIndentation: disableIndentation,
			FullPaths:          fullPaths,
		}

		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		root, err := internal.NewNode(path, config)
		if err != nil {
			panic(err)
		}

		fmt.Println(root)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&allFiles, "all", "a", false, "List all files")
	rootCmd.PersistentFlags().BoolVarP(&directoriesOnly, "directories", "d", false, "Only list directories")
	rootCmd.PersistentFlags().BoolVarP(&disableIndentation, "disable-indentation", "i", false, "Disable output indentation")
	rootCmd.PersistentFlags().BoolVarP(&fullPaths, "full", "f", false, "Print the full path for each file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
