package cmd

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/anoideaopen/migration-manager/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Stores all kv state entries to the migration directory (only query)",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			terminate   = make(chan os.Signal, 1)
			ctx, cancel = context.WithCancel(context.Background())
		)

		// Start processing
		go func() {
			signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
			<-terminate
			cancel()
		}()

		appCfg, channelClient, defFunc := core.InitChCli(cfgFile)
		defer defFunc()

		if requestEntries < defaultChunkSize {
			requestEntries = defaultChunkSize
		}

		if requestEntries > maxChunkSize {
			requestEntries = maxChunkSize
		}

		get := func(bookmark string, isComposite bool) (channel.Response, error) {
			req := channel.Request{
				ChaincodeID: appCfg.HLF.Chaincode,
				Fcn:         exportFn,
				Args: [][]byte{
					[]byte(strconv.FormatUint(uint64(requestEntries), 10)),
					[]byte(bookmark),
					[]byte(strconv.FormatBool(isComposite)),
				},
			}

			return channelClient.Query(req, channel.WithTimeout(fab.Execute, appCfg.HLF.ExecTimeout))
		}

		generalGet(ctx, get, appCfg, defaultTrysCount, defaultExt)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().Uint32VarP(&requestEntries, "entries", "e", defaultChunkSize, "state entries count per request")
}
