package cmds

import (
	"context"
	"net/http"

	"github.com/appscode/go/term"
	"github.com/gorilla/mux"
	"github.com/pharmer/pharmer/apiserver"
	"github.com/pharmer/pharmer/cloud"
	"github.com/pharmer/pharmer/config"
	"github.com/spf13/cobra"
)

func newCmdServer() *cobra.Command {
	//opts := options.NewClusterCreateConfig()
	cmd := &cobra.Command{
		Use:               "serve",
		Short:             "Pharmer apiserver",
		Example:           "pharmer serve",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			/*if err := opts.ValidateFlags(cmd, args); err != nil {
				term.Fatalln(err)
			}*/

			cfgFile, _ := config.GetConfigFile(cmd.Flags())
			cfg, err := config.LoadConfig(cfgFile)
			if err != nil {
				term.Fatalln(err)
			}
			if cfg.Store.Postgres == nil {
				term.Fatalln("Use postgres as storage provider")
			}
			ctx := cloud.NewContext(context.Background(), cfg, config.GetEnv(cmd.Flags()))
			err = http.ListenAndServe(":4155", route(ctx))
			term.ExitOnError(err)
		},
	}

	return cmd
}

func route(ctx context.Context) *mux.Router {
	server := apiserver.New(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/api/cluster/operation", server.CreateCluster).Methods("POST")
	return router
}
