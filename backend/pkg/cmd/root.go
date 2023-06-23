package cmd

import "github.com/spf13/cobra"

var (
	db = "sqlite://test.db"
)
var rootCmd = &cobra.Command{
	Use:   "homework",
	Short: "Backend for customer.io homework",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&db, "db", "sqlite://test.db",
		"database connection string as <dialect>:<DSN>")
}
