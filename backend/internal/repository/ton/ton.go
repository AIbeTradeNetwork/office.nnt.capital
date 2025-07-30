package ton

import (
	"context"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/tonconnect"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	ton2 "github.com/xssnick/tonutils-go/ton"
	"server/internal/config"
)

type Repo struct {
	clientTonGo    *liteapi.Client
	clientTonUtils ton2.APIClientWrapped
	conTonGo       *tonconnect.Server
	addrTonGo      ton.Address
	addrTonUtils   *address.Address
}

func Connect(ctx context.Context) (*Repo, error) {
	cfg := config.Get()

	var clientTonGo *liteapi.Client
	var err error
	if cfg.TonNetwork == "mainnet" {
		clientTonGo, err = liteapi.NewClientWithDefaultMainnet()
		if err != nil {
			return nil, err
		}
	} else {
		clientTonGo, err = liteapi.NewClientWithDefaultTestnet()
		if err != nil {
			return nil, err
		}
	}

	tonServer, err := tonconnect.NewTonConnect(clientTonGo, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	addrTonGo, err := tongo.ParseAddress(cfg.TonWallet)
	if err != nil {
		return nil, err
	}

	clientTonUtils := liteclient.NewConnectionPool()

	var cfgTonUtils *liteclient.GlobalConfig
	if cfg.TonNetwork == "mainnet" {
		// config for mainnet
		cfgTonUtils, err = liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
		if err != nil {
			return nil, err
		}
	} else {
		// config for testnet
		cfgTonUtils, err = liteclient.GetConfigFromUrl(context.Background(), "https://ton-blockchain.github.io/testnet-global.config.json")
		if err != nil {
			return nil, err
		}
	}

	// connect to lite servers
	err = clientTonUtils.AddConnectionsFromConfig(context.Background(), cfgTonUtils)
	if err != nil {
		return nil, err
	}

	// initialize ton api lite connection wrapper with full proof checks
	api := ton2.NewAPIClient(clientTonUtils, ton2.ProofCheckPolicySecure).WithRetry()
	api.SetTrustedBlockFromConfig(cfgTonUtils)

	addrTonUtils, err := address.ParseAddr(cfg.TonWallet)
	if err != nil {
		return nil, err
	}

	return &Repo{clientTonGo, api, tonServer, addrTonGo, addrTonUtils}, nil
}
