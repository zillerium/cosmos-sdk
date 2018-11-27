package types

import (
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/spf13/cobra"
)

// ModuleClients helps modules provide a standard interface for exporting client functionality
type ModuleClients interface {
	GetQueryCmd() *cobra.Command
	GetTxCmd() *cobra.Command
	RegisterRoutes(lcd.RestServer)
}
