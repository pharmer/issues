package cmds

import (
	"context"
	"fmt"

	"github.com/appscode/go/term"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pharmer/pharmer/apiserver"
	"github.com/pharmer/pharmer/cloud"
	"github.com/pharmer/pharmer/config"
	"github.com/spf13/cobra"
)

func newCmdServer() *cobra.Command {
	//opts := options.NewClusterCreateConfig()
	var natsurl string
	var clientid string
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
			fmt.Println(natsurl)

			ctx := cloud.NewContext(context.Background(), cfg, config.GetEnv(cmd.Flags()))

			err = runServer(ctx, natsurl, clientid)

			//err = http.ListenAndServe(":4155", route(ctx, conn))
			term.ExitOnError(err)
			<-make(chan interface{})
		},
	}
	cmd.Flags().StringVar(&natsurl, "nats-url", "nats://localhost:4222", "Nats streaming server url")
	cmd.Flags().StringVar(&clientid, "nats-client-id", "worker-p", "Nats streaming server client id")

	return cmd
}

//const ClientID = "worker-x"

func runServer(ctx context.Context, url, clientId string) error {
	conn, err := stan.Connect(
		"pharmer-cluster",
		clientId,
		stan.NatsURL(url),
	)
	fmt.Println(err, "..............", clientId)
	if err != nil {
		return err
	}

	//defer apiserver.LogCloser(conn)

	fmt.Println("II")
	server := apiserver.New(ctx, conn)
	return server.CreateCluster()

}
