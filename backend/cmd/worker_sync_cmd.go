package cmd

import (
	"devhub-backend/internal/worker"
	"time"

	"github.com/spf13/cobra"
)

var workerSyncPollInterval time.Duration
var workerSyncPluginTimeout time.Duration
var workerSyncPluginCommand []string
var workerSyncMaxWorkers int
var workerSyncMaxAttempts int

var workerSyncCmd = &cobra.Command{
	Use:   "worker_sync",
	Short: "Starts scaffold sync worker",
	RunE:  runWorkerSyncCmd,
}

func runWorkerSyncCmd(cmd *cobra.Command, args []string) error {
	return worker.NewSync(worker.SyncConfig{
		PollInterval:  workerSyncPollInterval,
		PluginTimeout: workerSyncPluginTimeout,
		PluginCommand: workerSyncPluginCommand,
		MaxWorkers:    workerSyncMaxWorkers,
		MaxAttempts:   workerSyncMaxAttempts,
	}).Start()
}

func init() {
	workerSyncCmd.Flags().DurationVar(&workerSyncPollInterval, "poll-interval", 2*time.Second, "poll interval for scaffold worker")
	workerSyncCmd.Flags().DurationVar(&workerSyncPluginTimeout, "plugin-timeout", 2*time.Minute, "timeout for one python plugin execution")
	workerSyncCmd.Flags().StringSliceVar(&workerSyncPluginCommand, "plugin-command", []string{}, "python plugin command parts, e.g. --plugin-command=python3,-m,plugins.scaffolder")
	workerSyncCmd.Flags().IntVar(&workerSyncMaxWorkers, "max-workers", 1, "number of concurrent worker loops")
	workerSyncCmd.Flags().IntVar(&workerSyncMaxAttempts, "max-attempts", 3, "max attempts per job before hard fail")
}
