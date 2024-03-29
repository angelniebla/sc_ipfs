package prueba_ipfs

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/coreutil"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"
	ipfs_core "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	mock "github.com/ipfs/go-ipfs/core/mock"
	icore "github.com/ipfs/interface-go-ipfs-core"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
)

var Contract = coreutil.NewContract("IPFS", "IPFS, a PoC smart contract")

var Processor = Contract.Processor(initialize,
	FuncGetCounter.WithHandler(getTemp),
)

var (
	FuncGetCounter = coreutil.ViewFunc("getTemp")
)

var o *orbit

const (
	VarTemp = "temp"
)

type orbit struct {
	Db iface.OrbitDB

	Kv orbitdb.KeyValueStore
}

func testingMockNet(ctx context.Context) mocknet.Mocknet {
	return mocknet.New(ctx)
}

/*

// Creates an IPFS node and returns its coreAPI
func createNode2(ctx context.Context) (icore.CoreAPI, error) {

	repoPath, err := createTempRepo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %s", err)
	}
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &ipfs_core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	}

	node, err := ipfs_core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func createTempRepo(ctx context.Context) (string, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

*/

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context) (icore.CoreAPI, error) {

	// Construct the node
	m := testingMockNet(ctx)
	node, err := ipfs_core.NewNode(ctx, &ipfs_core.BuildCfg{
		Online: true,
		//	Repo: repo,
		Host: mock.MockHostOption(m),
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	})

	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func constructor() *orbit {
	o := new(orbit)
	return o
}

func initOrbit(val int64) {

	o = new(orbit)
	path := "/home/angel/db"
	ctx := context.Background()
	err := Mkdir(path)
	if err != nil {
		fmt.Println("error mkdir")
	}
	ipfs, err := createNode(ctx)

	if err != nil {
		fmt.Println("new core api error:", err.Error())
	}
	fmt.Println("create Orbit")
	orbit, err := orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{Directory: &path})
	if err != nil {
		fmt.Println("new orbitdb error:", err.Error())
	}
	o.Db = orbit
	kv, err := orbit.KeyValue(ctx, "temperatures", nil)
	if err != nil {
		fmt.Println("userinfo error:", err.Error())
	}
	s := strconv.FormatInt(val, 10)
	kv.Put(ctx, "temp", []byte(s))
	temp, _ := kv.Get(ctx, "temp")
	fmt.Println("Success: ", string(temp[:]))
	o.Kv = kv

	//value, _ := kv.Get(ctx, "a")

	println(string(temp))
}

func Mkdir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			fmt.Println("dir error:", err.Error())
			return err
		}
	} else {
		fmt.Println("check dir error:", err.Error())
		return err
	}
	return nil
}

/*

func initialize(ctx iscp.Sandbox) (dict.Dict, error) {
	ctx.Log().Debugf("IPFS.init in %s", ctx.Contract().String())
	params := ctx.Params()
	val, err := codec.DecodeInt64(params.MustGet(VarCounter), 0)
	if err != nil {
		return nil, fmt.Errorf("IPFS: %v", err)
	}

	initOrbit(val)

	return nil, nil
}

*/

func initialize(ctx iscp.Sandbox) (dict.Dict, error) {
	ctx.Log().Debugf("ipfs.init in %s", ctx.Contract().String())
	params := ctx.Params()
	val, _, err := codec.DecodeInt64(params.MustGet(VarTemp))
	if err != nil {
		return nil, fmt.Errorf("ipfs: %v", err)
	}
	ctx.State().Set(VarTemp, codec.EncodeInt64(val))
	ctx.Event(fmt.Sprintf("ipfs.init.success. counter = %d", val))

	initOrbit(val)

	return nil, nil
}

func getTemp(ctx iscp.SandboxView) (dict.Dict, error) {
	ret := dict.New()

	ctx2 := context.Background()
	//kv, err := o.Db.KeyValue(ctx2, "temperatures", nil)

	//s := strconv.FormatInt(18, 10)
	//kv.Put(ctx2, "tempp", []byte(s))

	temp, _ := o.Kv.Get(ctx2, "temp")

	fmt.Println("Success: ", string(temp[:]))

	i, _ := strconv.ParseInt(string(temp[:]), 10, 64)

	ret.Set(VarTemp, codec.EncodeInt64(i))

	return ret, nil
}
