package matchsessioncli

import (
	"github.com/spf13/cobra"
)

var NewMatchSessionCobraCommand = func() *cobra.Command {
	return &cobra.Command{
		Use:       "find-match",
		Short:     "Try to match another User based on the preferences",
		ValidArgs: []string{"userID"},
		Args:      cobra.MinimumNArgs(3),
		Run:       handleRootCmd(),
	}
}

func handleRootCmd() func(*cobra.Command, []string) {
	return func(cobraCmd *cobra.Command, _ []string) {
		// TODO: to be implemented
		cobraCmd.PrintErrln("not implemented yet")
	}
}
