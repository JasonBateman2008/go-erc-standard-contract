package erc1155

import (
	"github.com/jason-bateman/go-erc-standard-contract/contracts/erc1155/model"
	chainModel "github.com/jason-bateman/go-erc-standard-contract/model"
	"github.com/jason-bateman/go-erc-standard-contract/utils"
	"testing"
)

func getContractInstance() (*Contract, error) {
	ops := &ContractOpts{
		Rpc:               "https://data-seed-prebsc-2-s1.binance.org:8545",
		ContractAddr:      "0xe709bB6C73A9473331AC05dDf14cc1f7eB6f4951",
		EnableTransactors: true,
		EnableFilter:      true,
		FilterStep:        5000,
	}
	return NewContract(ops)
}

func getFilterFuzzyContractInstance() (*Contract, error) {
	ops := &ContractOpts{
		Rpc:                "https://data-seed-prebsc-2-s1.binance.org:8545",
		ContractAddr:       "0xe709bB6C73A9473331AC05dDf14cc1f7eB6f4951",
		EnableTransactors:  true,
		EnableFilter:       true,
		FilterStep:         5000,
		FilterFuzzyAddress: true,
	}
	return NewContract(ops)
}

func TestContract_ReadContract(t *testing.T) {
	// get the instantiated object
	contract, err := getContractInstance()
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}
	defer contract.ReleaseResource() // releasing resources

	balanceOfInputs := &model.MethodReadBalanceOf{
		Owner: "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F",
		Id:    "110801524474586856940016194582843495297783442968439414267952739331789173737983",
	}
	balance, err := contract.ReadBalanceOf(balanceOfInputs)
	if err != nil {
		t.Errorf("ReadBalanceOf err:%+v\n", err)
		return
	}
	t.Logf("balance is:%d,\n", balance)

	isApprovedForAllInputs := &model.MethodReadIsApprovedForAllInputs{
		Owner:    "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F",
		Operator: "0x604e91519c3f515d93050ae3b909d9ad037085b5",
	}
	approved, err := contract.ReadIsApprovedForAll(isApprovedForAllInputs)
	if err != nil {
		t.Errorf("ReadBalanceOf err:%+v\n", err)
		return
	}
	t.Log("approved is:\n", approved)

	balanceOfBatchInputs := &model.MethodReadBalanceOfBatchInputs{
		Owners: []string{"0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F", "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F"},
		Ids:    []string{"110801524474586856940016194582843495297783442968439414267952739331789173737983", "110801524474586856940016194582843495297783442968439414267952739331789173737983"},
	}
	batchBalance, err := contract.ReadBalanceOfBatch(balanceOfBatchInputs)
	if err != nil {
		t.Errorf("ReadBalanceOfBatch err:%+v\n", err)
		return
	}
	t.Log("batchBalance is:\n", batchBalance)

	uri, err := contract.ReadUri("110801524474586856940016194582843495297783442968439414267952739331789173737983")
	if err != nil {
		t.Errorf("ReadUri err:%+v\n", err)
		return
	}
	t.Log("uri is:\n", uri)

	isSupport, err := contract.ReadSupportsInterface("0x01ffc9a7")
	if err != nil {
		t.Errorf("ReadSupportsInterface err:%+v\n", err)
		return
	}
	t.Log("ReadSupportsInterface is:\n", isSupport)

}

func TestContract_WriteSafeTransferFrom(t *testing.T) {
	toAddress := "0x604e91519c3F515D93050AE3B909d9AD037085b5"
	senderPrivateKey := "Add the private key to this variable"
	senderAddress := "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F"
	tokenId := "110801524474586856940016194582843495297783442968439414267952739331789173737983"
	//  get the instantiated object
	contract, err := getContractInstance()
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}
	defer contract.ReleaseResource() // releasing resources

	keys := []string{senderPrivateKey}

	// adding a signature Provider
	err = contract.AddTransactors(keys)
	if err != nil {
		t.Errorf("AddTransactors err:%+v\n", err)
		return
	}

	inputs := &model.MethodWriteSafeTransferFromInputs{
		From:   senderAddress,
		To:     toAddress,
		Id:     tokenId,
		Amount: 0x01,
		Data:   []byte(""),
	}

	txNonce, err := utils.GetAddressTxNonceWithClient(contract.GetCallerClient(), senderAddress)
	if err != nil {
		t.Errorf("GetAddressTxNonceWithClient err:%+v\n", err)
		return
	}

	txId, err := contract.WriteSafeTransferFrom(*txNonce, inputs)
	if err != nil {
		t.Errorf("WriteSafeTransferFrom err:%+v\n", err)
		return
	}
	t.Logf("txId is:%s,\n", txId)
}

func TestContract_WriteBatchSafeTransferFrom(t *testing.T) {
	toAddress := "0x604e91519c3F515D93050AE3B909d9AD037085b5"
	senderPrivateKey := "Add the private key to this variable"
	senderAddress := "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F"
	tokenId := "0xf4f770c0dde6e24b4c65a85f744fec0bd3d89b1f00000000000002003b9ac9ff"

	// 获取实例化对象
	contract, err := getContractInstance()
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}
	defer contract.ReleaseResource() // 结束释放资源

	keys := []string{senderPrivateKey}

	// 添加签名提供者
	err = contract.AddTransactors(keys)
	if err != nil {
		t.Errorf("AddTransactors err:%+v\n", err)
		return
	}

	inputs := &model.MethodWriteSafeBatchTransferFromInputs{
		From:    senderAddress,
		To:      toAddress,
		Ids:     []string{tokenId, tokenId},
		Amounts: []int64{0x01, 0x01},
		Data:    []byte(""),
	}

	txNonce, err := utils.GetAddressTxNonceWithClient(contract.GetCallerClient(), senderAddress)
	if err != nil {
		t.Errorf("GetAddressTxNonceWithClient err:%+v\n", err)
		return
	}

	txId, err := contract.WriteSafeBatchTransferFrom(*txNonce, inputs)
	if err != nil {
		t.Errorf("WriteSafeBatchTransferFrom err:%+v\n", err)
		return
	}
	t.Logf("txId is:%s,\n", txId)
}

func TestFilterContractLogs(t *testing.T) {

	// 获取实例化对象
	contract, err := getContractInstance()
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}
	defer contract.ReleaseResource() // 结束释放资源

	// 添加监听
	events := []model.ContractEvent{
		model.EventApprovalForAll,
		model.EventURI,
		model.EventTransferSingle,
		model.EventTransferBatch,
	}
	err = contract.AddEvents(events)
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}

	//deployBlockNum := uint64(465406)
	//deployBlockNum := uint64(480883)
	deployBlockNum := uint64(1754605)

	nowBlockNum, err := utils.GetLatestBlockNumWithClient(contract.GetCallerClient())
	if err != nil {
		t.Errorf("GetLatestBlockNum err:%+v\n", err)
		return
	}

	var eventsAll []*chainModel.EthereumEventMessage
	start := deployBlockNum
	stop := start + contract.filter.stepNum
	for {
		if stop >= nowBlockNum {
			break
		}
		t.Logf("Filter Start from %d -- %d", start, stop)
		events, err := contract.FilterEvents(start, &stop)
		if err != nil {
			t.Errorf("FilterEvents err:%+v\n", err)
		}
		if len(events) != 0 {
			eventsAll = utils.MergeEventMessage(eventsAll, events)
			t.Logf("Filter Stop from %d -- %d, Total %d message", start, stop, len(events))
		}
		start = stop + 1
		stop += contract.filter.stepNum
	}

	if len(eventsAll) != 0 {
		t.Logf("Filter End from %d -- %d, Total %d message", deployBlockNum, stop, len(eventsAll))
	}
}

func TestFilterFuzzyContractLogs(t *testing.T) {
	contract, err := getFilterFuzzyContractInstance()
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}
	defer contract.ReleaseResource()

	events := []model.ContractEvent{
		model.EventApprovalForAll,
		model.EventURI,
		model.EventTransferSingle,
		model.EventTransferBatch,
	}
	err = contract.AddEvents(events)
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}

	//deployBlockNum := uint64(465406)
	//deployBlockNum := uint64(480883)
	deployBlockNum := uint64(1754605)

	nowBlockNum, err := utils.GetLatestBlockNumWithClient(contract.GetCallerClient())
	if err != nil {
		t.Errorf("GetLatestBlockNum err:%+v\n", err)
		return
	}

	var eventsAll []*chainModel.EthereumEventMessage
	start := deployBlockNum
	stop := start + contract.filter.stepNum
	for {
		if stop >= nowBlockNum {
			break
		}
		t.Logf("Filter Start from %d -- %d", start, stop)
		events, err := contract.FilterEvents(start, &stop)
		if err != nil {
			t.Errorf("FilterEvents err:%+v\n", err)
		}
		if len(events) != 0 {
			eventsAll = utils.MergeEventMessage(eventsAll, events)
			t.Logf("Filter Stop from %d -- %d, Total %d message", start, stop, len(events))
		}
		start = stop + 1
		stop += contract.filter.stepNum
	}

	if len(eventsAll) != 0 {
		t.Logf("Filter End from %d -- %d, Total %d message", deployBlockNum, stop, len(eventsAll))
	}
}
