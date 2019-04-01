package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/TimJones/godotty/internal/app/godotty"
)

var (
	dottyDir    string
	rootCommand = &cobra.Command{
		Use:   "godotty",
		Short: "godotty is a dotfile manager",
		Long: `A portable tool to help organise and manage dotfiles
across multiple environments`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	rootCommand.PersistentFlags().StringVarP(&dottyDir, "dottydir", "d", godotty.DefaultDirectory, "Specify the directory to use for managing dotfiles")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
