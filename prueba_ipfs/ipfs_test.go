package ipfs

import (
	"testing"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core"
	"github.com/stretchr/testify/require"
)

const IPFSName = "IPFS_test"

func TestDeployIPFS(t *testing.T) {
	env := solo.New(t, false, false).WithNativeContract(Processor)
	chain := env.NewChain(nil, "chain1")

	err := chain.DeployContract(nil, IPFSName, Contract.ProgramHash)
	require.NoError(t, err)
	chain.CheckChain()
	_, _, contracts := chain.GetInfo()
	require.EqualValues(t, len(core.AllCoreContractsByHash)+1, len(contracts))
	chain.CheckAccountLedger()
}

func TestDeployIPFSInitParams(t *testing.T) {
	env := solo.New(t, false, false).WithNativeContract(Processor)
	chain := env.NewChain(nil, "chain1")

	err := chain.DeployContract(nil, IPFSName, Contract.ProgramHash, VarCounter, 17)
	require.NoError(t, err)
	//checkCounter(chain, 17)
	//chain.CheckAccountLedger()
}
