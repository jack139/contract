package keeper

import (
	"github.com/jack139/contract/x/contract/types"
)

var _ types.QueryServer = Keeper{}
