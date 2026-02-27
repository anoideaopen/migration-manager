package core

import (
	"log"
	"path/filepath"
	"time"

	"github.com/anoideaopen/migration-manager/cfg"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const defaultInvokeTimeout = time.Minute * 2

func InitChCli(cfgFile string) (*cfg.Config, *channel.Client, func()) {
	appCfg, clientChannelContext, defFunc := initChCtx(cfgFile)

	var opts []channel.ClientOption

	channelClient, err := channel.New(clientChannelContext, opts...)
	if err != nil {
		log.Panicf("couldn't create channel client: %v", err)
	}

	return appCfg, channelClient, defFunc
}

func initChCtx(cfgFile string) (*cfg.Config, context.ChannelProvider, func()) {
	appCfg := &cfg.Config{}
	if err := cfg.ReadFromFile(cfg.EnvPrefix, cfgFile, appCfg); err != nil {
		log.Fatalf("read config error: %v", err)
	}

	if appCfg.HLF == nil {
		log.Fatal("could not find HLF config section")
	}

	if appCfg.HLF.Channel == "" {
		log.Fatal("could not find HLF channel name")
	}

	if appCfg.HLF.ExecTimeout == 0 {
		appCfg.HLF.ExecTimeout = defaultInvokeTimeout
	}

	if appCfg.SnapshotDir == "" {
		log.Fatal("could not find SnapshotDir")
	}

	appCfg.SnapshotDir, _ = filepath.Abs(appCfg.SnapshotDir)
	if appCfg.SnapshotDir[len(appCfg.SnapshotDir)-1] != filepath.Separator {
		appCfg.SnapshotDir += string(filepath.Separator)
	}
	appCfg.SnapshotDir += appCfg.HLF.Channel

	hlfCfg := config.FromFile(appCfg.HLF.Config)

	log.Print("Initializing SDK")

	sdk, err := fabsdk.New(hlfCfg)
	if err != nil {
		log.Fatalf("couldn't initialize SDK: %v", err)
	}

	retf := func() {
		sdk.Close()
	}

	log.Print("Initializing SDK done")

	clientChannelContext := sdk.ChannelContext(appCfg.HLF.Channel,
		fabsdk.WithUser(appCfg.HLF.User), fabsdk.WithOrg(appCfg.HLF.Org))

	return appCfg, clientChannelContext, retf
}
