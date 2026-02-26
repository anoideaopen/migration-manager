package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/anoideaopen/migration-manager/core"
	"github.com/anoideaopen/migrationcc/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Uploads all kv state entries from the migration directory",
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

		files, err := os.ReadDir(appCfg.SnapshotDir)
		if err != nil {
			log.Panicf("couldn't read state directiry: %v", err)
		}

		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})

	OuterLoop:
		for i, file := range files {
			select {
			case <-ctx.Done():
				break OuterLoop
			default:
			}

			name := file.Name()
			if !strings.HasSuffix(name, defaultExt) || file.IsDir() {
				continue
			}

			entriesData, err := os.ReadFile(path.Join(appCfg.SnapshotDir, name))
			if err != nil {
				log.Panicf("couldn't read state file: %v", err)
			}

			entries := new(proto.Entries)
			core.MustUnmarshal(entriesData, entries)

			log.Printf("Bookmark: %s", entries.GetBookmark())
			for _, entry := range entries.GetEntries() {
				log.Printf("Checked state entry: %s -> %d bytes", entry.GetKey(), len(entry.GetValue()))
			}

			log.Printf("Sending chunk %d/%d", i+1, len(files))

			var resp channel.Response
			for j := range defaultTrysCount {
				select {
				case <-ctx.Done():
					break OuterLoop
				default:
				}

				if resp, err = channelClient.Execute(channel.Request{
					ChaincodeID: appCfg.HLF.Chaincode,
					Fcn:         importFn,
					Args: [][]byte{
						entriesData,
					},
				},
					channel.WithTimeout(fab.Execute, appCfg.HLF.ExecTimeout),
				); err == nil {
					break
				}

				log.Printf("error sending state (try %d/%d): %v", j+1, defaultTrysCount, err)
				time.Sleep(time.Second)
			}

			if err != nil {
				log.Panicf("couldn't send state: %v", err)
			}

			log.Printf("Sending %s sent", name)

			if resp.ChaincodeStatus != http.StatusOK {
				log.Panicf("invalid response status: %d", resp.ChaincodeStatus)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
