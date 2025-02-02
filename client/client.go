// APACHE NOTICE
// Sourced with modifications from https://github.com/strangelove-ventures/lens
package client

import (
	"io"
	"path"
	"time"

	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type ChainClient struct {
	Config         *ChainClientConfig
	Keybase        keyring.Keyring
	KeyringOptions []keyring.Option
	RPCClient      rpcclient.Client
	Input          io.Reader
	Output         io.Writer
	Codec          Codec
	Logger         log.Logger
}

func NewChainClient(ccc *ChainClientConfig, homepath string, input io.Reader, output io.Writer, kro ...keyring.Option) (*ChainClient, error) {
	ccc.KeyDirectory = keysDir(homepath, ccc.ChainID)
	cc := &ChainClient{
		KeyringOptions: kro,
		Config:         ccc,
		Input:          input,
		Output:         output,
		Codec:          MakeCodec(ccc.Modules),
		Logger:         log.NewTMLogger(log.NewSyncWriter(output)),
	}
	if err := cc.Init(); err != nil {
		return nil, err
	}
	return cc, nil
}

func (cc *ChainClient) Init() error {
	// TODO: test key directory and return error if not created
	// TODO_SCALAR: Confirm if remove Marshaler could be a problem
	// keybase, err := keyring.New(cc.Config.ChainID, cc.Config.KeyringBackend, cc.Config.KeyDirectory, cc.Input, cc.Codec.Marshaler, cc.KeyringOptions...)
	keybase, err := keyring.New(cc.Config.ChainID, cc.Config.KeyringBackend, cc.Config.KeyDirectory, cc.Input, cc.KeyringOptions...)
	if err != nil {
		return err
	}
	// TODO: figure out how to deal with input or maybe just make all keyring backends test?

	timeout, _ := time.ParseDuration(cc.Config.Timeout)
	rpcClient, err := NewRPCClient(cc.Config.RPCAddr, timeout)
	if err != nil {
		return err
	}

	cc.RPCClient = rpcClient
	cc.Keybase = keybase

	return nil
}

func keysDir(home, chainID string) string {
	return path.Join(home, "keys", chainID)
}

func NewRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}
