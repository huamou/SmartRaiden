package smartraiden

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path"

	"time"

	"encoding/hex"

	"sync"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network"
	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/fee"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var curAccountIndex = 0

func newTestRaiden() *RaidenService {
	return newTestRaidenWithPolicy(&NoFeePolicy{})
}

func newTestRaidenWithPolicy(feePolicy fee.FeeCharger) *RaidenService {
	transport := network.MakeTestUDPTransport(50000 + curAccountIndex + 1)
	bcs := newTestBlockChainService()
	discover := network.GetTestDiscovery() //share the same discovery ,so node can find each other
	//discover := network.NewContractDiscovery(bcs.NodeAddress, bcs.Client, bcs.Auth)
	config := params.DefaultConfig
	config.MyAddress = bcs.NodeAddress
	config.PrivateKey = bcs.PrivKey
	config.DataDir = path.Join(os.TempDir(), utils.RandomString(10))
	config.ExternIp = transport.Host
	config.ExternPort = transport.Port
	config.Host = transport.Host
	config.Port = transport.Port
	config.RevealTimeout = 10
	config.SettleTimeout = 600
	config.PrivateKeyHex = hex.EncodeToString(crypto.FromECDSA(config.PrivateKey))
	os.MkdirAll(config.DataDir, os.ModePerm)
	config.DataBasePath = path.Join(config.DataDir, "log.db")
	rd := NewRaidenService(bcs, bcs.PrivKey, transport, discover, &config)
	rd.SetFeePolicy(feePolicy)
	return rd
}
func newTestRaidenApi() *RaidenApi {
	api := NewRaidenApi(newTestRaiden())
	api.Raiden.Start()
	return api
}

//maker sure these accounts are valid, and  engouh eths for test
func testGetnextValidAccount() (*ecdsa.PrivateKey, common.Address) {
	am := NewAccountManager("testdata/keystore")
	privkey, err := am.GetPrivateKey(am.Accounts[curAccountIndex].Address, "123")
	if err != nil {
		fmt.Sprintf("testGetnextValidAccount err:", err)
		panic("")
	}
	curAccountIndex++
	return crypto.ToECDSAUnsafe(privkey), utils.PubkeyToAddress(privkey)
}
func newTestBlockChainService() *rpc.BlockChainService {
	conn, err := helper.NewSafeClient(rpc.TestRpcEndpoint)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to connect to the Ethereum client: %s", err))
	}
	privkey, _ := testGetnextValidAccount()
	if err != nil {
		log.Error("Failed to create authorized transactor: ", err)
	}
	return rpc.NewBlockChainService(privkey, rpc.PRIVATE_ROPSTEN_REGISTRY_ADDRESS, conn)
}

func makeTestRaidens() (r1, r2, r3 *RaidenService) {
	r1 = newTestRaiden()
	r2 = newTestRaiden()
	r3 = newTestRaiden()
	go func() {
		r1.Start()
	}()
	go func() {
		r2.Start()
	}()
	go func() {
		r3.Start()
	}()
	time.Sleep(time.Second * 3)
	return
}
func newTestRaidenApiQuick() *RaidenApi {
	api := NewRaidenApi(newTestRaiden())
	go func() {
		api.Raiden.Start()
	}()
	return api
}

func makeTestRaidenApis() (rA, rB, rC, rD *RaidenApi) {
	rA = newTestRaidenApiQuick()
	rB = newTestRaidenApiQuick()
	rC = newTestRaidenApiQuick()
	rD = newTestRaidenApiQuick()
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		rA.Raiden.Start()
		wg.Done()
	}()
	go func() {
		rB.Raiden.Start()
		wg.Done()
	}()
	go func() {
		rC.Raiden.Start()
		wg.Done()
	}()
	go func() {
		rD.Raiden.Start()
		wg.Done()
	}()
	return
}

func makeTestRaidenApisWithFee(policy fee.FeeCharger) (rA, rB, rC, rD *RaidenApi) {
	rA = NewRaidenApi(newTestRaidenWithPolicy(policy))
	rB = NewRaidenApi(newTestRaidenWithPolicy(policy))
	rC = NewRaidenApi(newTestRaidenWithPolicy(policy))
	rD = NewRaidenApi(newTestRaidenWithPolicy(policy))
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		rA.Raiden.Start()
		wg.Done()
	}()
	go func() {
		rB.Raiden.Start()
		wg.Done()
	}()
	go func() {
		rC.Raiden.Start()
		wg.Done()
	}()
	go func() {
		rD.Raiden.Start()
		wg.Done()
	}()
	wg.Wait()
	return
}
