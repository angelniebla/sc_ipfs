package prueba_ipfs

import (
	"testing"

	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core"
	"github.com/stretchr/testify/require"
)

const IPFSName = "IPFS_test"

func checkCounter(e *solo.Chain, expected int64) {
	ret, err := e.CallView(IPFSName, FuncGetCounter.Name)
	require.NoError(e.Env.T, err)
	c, ok, err := codec.DecodeInt64(ret.MustGet(VarTemp))
	require.NoError(e.Env.T, err)
	require.True(e.Env.T, ok)
	require.EqualValues(e.Env.T, expected, c)
}

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

	err := chain.DeployContract(nil, IPFSName, Contract.ProgramHash, VarTemp, 17)
	require.NoError(t, err)
	checkCounter(chain, 17)
	chain.CheckAccountLedger()
}
