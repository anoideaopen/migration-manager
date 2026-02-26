package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/anoideaopen/migration-manager/cfg"
	"github.com/anoideaopen/migration-manager/core"
	"github.com/anoideaopen/migrationcc/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/spf13/cobra"
)

var (
	requestEntries uint32
	onlyKeys       bool
	isInvoke       bool
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Stores all kv state entries to the migration directory",
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

		// If we are executing, and not requesting, then only fully KV
		if isInvoke {
			onlyKeys = false
		}

		get := func(bookmark string) (channel.Response, error) {
			req := channel.Request{
				ChaincodeID: appCfg.HLF.Chaincode,
				Fcn:         exportFn,
				Args: [][]byte{
					[]byte(strconv.FormatUint(uint64(requestEntries), 10)),
					[]byte(bookmark),
					[]byte(strconv.FormatBool(onlyKeys)),
				},
			}

			if isInvoke {
				return channelClient.Execute(req, channel.WithTimeout(fab.Execute, appCfg.HLF.ExecTimeout))
			}

			return channelClient.Query(req, channel.WithTimeout(fab.Execute, appCfg.HLF.ExecTimeout))
		}

		generalGet(ctx, get, appCfg, defaultTrysCount, defaultExt)
	},
}

func generalGet( //nolint:funlen
	ctx context.Context,
	get func(bookmark string) (channel.Response, error),
	appCfg *cfg.Config,
	trysCount int,
	ext string,
) {
	var (
		bookmark string
		chunkNum uint64
	)

	_ = os.RemoveAll(appCfg.SnapshotDir)
	_ = os.MkdirAll(appCfg.SnapshotDir, 0o755) //nolint:gomnd

OuterLoop:
	for {
		select {
		case <-ctx.Done():
			break OuterLoop
		default:
		}

		log.Print("Requesting chunk")

		var (
			resp channel.Response
			err  error
		)
		for j := range trysCount {
			select {
			case <-ctx.Done():
				break OuterLoop
			default:
			}

			if resp, err = get(bookmark); err == nil {
				break
			}

			log.Printf("error sending state (try %d/%d): %v", j+1, trysCount, err)
			time.Sleep(time.Second)
		}

		log.Print("Requesting chunk done")

		if err != nil {
			log.Panicf("couldn't request state entries: %v", err)
		}

		if resp.ChaincodeStatus != http.StatusOK {
			log.Panicf("invalid response status: %d", resp.ChaincodeStatus)
		}

		entries := new(proto.Entries)
		core.MustUnmarshal(resp.Payload, entries)

		log.Printf("Bookmark: %s", entries.GetBookmark())
		for _, entry := range entries.GetEntries() {
			log.Printf("Received state entry: %s -> %d bytes", entry.GetKey(), len(entry.GetValue()))
		}

		fileName := path.Join(appCfg.SnapshotDir, fmt.Sprintf("%09d%s", chunkNum, ext))
		if err = os.WriteFile(fileName, resp.Payload, 0o600); err != nil { //nolint:gomnd
			log.Panicf("couldn't save state file: %v", err)
		}

		log.Printf("Chunk %d stored", chunkNum)

		if entries.GetBookmark() == "" {
			log.Print("Done")
			break
		}

		bookmark = entries.GetBookmark()
		chunkNum++
	}
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().Uint32VarP(&requestEntries, "entries", "e", defaultChunkSize, "state entries count per request")
	exportCmd.Flags().BoolVarP(&onlyKeys, "only_keys", "k", false, "only the keys will be returned in the response (default false)")
	exportCmd.Flags().BoolVarP(&isInvoke, "is_invoke", "i", false, "exec invoke or query (default query)")
}
