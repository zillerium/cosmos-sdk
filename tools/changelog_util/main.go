package main

import "github.com/spf13/cobra"

func main() {

	rootCmd := &cobra.Command{
		Use:   "changlog-util",
		Short: "changelog utility",
	}
	rootCmd.AddCommand(AddCmd)

	rootCmd.Execute()
}

// Add pending changelog entry
var AddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
