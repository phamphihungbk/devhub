package cmd

import (
	"devhub-backend/internal/server"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var syncWorkerCmd = &cobra.Command{
	Use:   "sync-worker",
	Short: "Run background worker(s) for async requests (scaffold/deployment/release).",
	Long: `Run generic background workers for async request processing.

Examples:
	sync-worker
	sync-worker --concurrency 3 --poll-interval 5s --types scaffold,deployment,release
`,
	RunE: runNewSyncWorkerCmd,
}

func runNewSyncWorkerCmd(cmd *cobra.Command, args []string) error {
	concurrency, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return fmt.Errorf("failed to parse concurrency flag: %w", err)
	}
	pollInterval, err := cmd.Flags().GetDuration("poll-interval")
	if err != nil {
		return fmt.Errorf("failed to parse poll-interval flag: %w", err)
	}
	typesRaw, err := cmd.Flags().GetString("types")
	if err != nil {
		return fmt.Errorf("failed to parse types flag: %w", err)
	}

	workerTypes := make([]string, 0, 2)
	for _, t := range strings.Split(typesRaw, ",") {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		workerTypes = append(workerTypes, t)
	}

	return server.NewWorker(
		server.WithWorkerConcurrency(concurrency),
		server.WithWorkerPollInterval(pollInterval),
		server.WithWorkerTypes(workerTypes),
	).Start()
}

func init() {
	syncWorkerCmd.Flags().Int("concurrency", 1, "number of worker goroutines per request type")
	syncWorkerCmd.Flags().Duration("poll-interval", 3*time.Second, "poll interval for each worker")
	syncWorkerCmd.Flags().String("types", "scaffold,deployment,release", "comma-separated worker types to run")
}
