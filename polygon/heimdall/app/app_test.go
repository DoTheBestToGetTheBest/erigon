package app

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func TestHeimdallAppExport(t *testing.T) {
	t.Parallel()

	db := db.NewMemDB()
	happ := NewHeimdallConsumer(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	genesisState := NewDefaultGenesisState()

	// Get state bytes
	stateBytes, err := codec.MarshalJSONIndent(happ.Codec(), genesisState)
	require.NoError(t, err)

	// Initialize the chain
	happ.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	// Set commit
	happ.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	newHapp := NewHeimdallConsumer(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	_, _, err = newHapp.ExportAppStateAndValidators()
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
