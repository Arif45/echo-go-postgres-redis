package cmd

import (
	"fin-auth/config"
	"fin-auth/worker"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start background workers",
	Run: func(cmd *cobra.Command, args []string) {
		starWorker()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}

func starWorker() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, error := config.InitGormDB()
	if error != nil {
		log.Fatalf("Failed to initialize database: %v", error)
	}

	tokenCleanupWorker := worker.NewTokenCleanupWorker(db)

	worker := worker.NewWorker(5*time.Minute, tokenCleanupWorker.Run, "token-cleanup")
	worker.Start()
}
