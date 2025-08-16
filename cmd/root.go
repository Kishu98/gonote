package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "gonote",
	Short: "A tool to organize the notes using AI",
	Long:  "A tool which uses the AI to assign notes to a specific folder",
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
