package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/TimJones/godotty/internal/app/godotty"
)

var importCommand = &cobra.Command{
	Use:   "import <file> [<file>...]",
	Short: "Import one or more dotfiles to be managed by godotty",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotty.Import(dottyDir, args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCommand.AddCommand(importCommand)
}
