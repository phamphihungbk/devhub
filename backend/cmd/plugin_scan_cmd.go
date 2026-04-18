package cmd

import (
	"context"
	"fmt"

	"devhub-backend/internal/config"
	infraDB "devhub-backend/internal/infra/db"
	pluginrepo "devhub-backend/internal/infra/db/repository/plugin"
	pluginusecase "devhub-backend/internal/usecase/plugin"

	"github.com/spf13/cobra"
)

const defaultPluginsDir = "../plugins"

var pluginScanCmd = &cobra.Command{
	Use:   "plugin-scan",
	Short: "Scan plugin manifests and upsert them into the plugins table",
	Long: `Scan plugin manifests from the local plugins directory and upsert them into the plugins table.

This command will:
- find plugin manifests under the plugins directory
- parse each plugin.yaml manifest
- derive the runtime entrypoint path expected by workers
- create missing plugin records
- update existing plugin records when manifest data changes

Example:
	plugin-scan
	plugin-scan --dir ../plugins
`,
	RunE: runPluginScanCmd,
}

func runPluginScanCmd(cmd *cobra.Command, args []string) error {
	pluginsDir, _ := cmd.Flags().GetString("dir")

	cfg := config.MustConfigure()
	ctx := context.Background()
	dbConn := infraDB.MustConnect(cfg)
	defer dbConn.Close()

	repo := pluginrepo.NewPluginRepository(dbConn)
	usecase := pluginusecase.NewPluginUsecase(cfg.App, repo)

	result, err := usecase.SyncRegistry(ctx, pluginusecase.SyncRegistryInput{
		PluginsDir: pluginsDir,
	})

	if err != nil {
		return err
	}

	fmt.Fprintf(
		cmd.OutOrStdout(),
		"plugin scan complete: discovered=%d created=%d updated=%d\n",
		result.Discovered,
		result.Created,
		result.Updated,
	)
	return nil
}

func init() {
	pluginScanCmd.Flags().String("dir", defaultPluginsDir, "plugins directory to scan")
}
