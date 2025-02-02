// APACHE NOTICE
// Sourced with modifications from https://github.com/strangelove-ventures/lens
package query

import (
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// StatusRPC returns information about a node status
func StatusRPC(q *Query) (*coretypes.ResultStatus, error) {
	ctx, cancel := q.GetQueryContext()
	defer cancel()
	res, err := q.Client.RPCClient.Status(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
