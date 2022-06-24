package erc721

import (
	"github.com/jason-bateman/go-erc-standard-contract/contracts/erc721/model"
	chainModel "github.com/jason-bateman/go-erc-standard-contract/model"
	"github.com/jason-bateman/go-erc-standard-contract/utils"
	"testing"
)

func getContractInstance() (*Contract, error) {
	ops := &ContractOpts{
		Rpc:               "https://data-seed-prebsc-2-s1.binance.org:8545",
		ContractAddr:      "0x0bB31BA49d2b9604Ea1640DE4d70D861920AcDe9",
		EnableTransactors: true,
		EnableFilter:      true,
		FilterStep:        3500,
	}
	return NewContract(ops)
}

func getFilterFuzzyContractInstance() (*Contract, error) {
	ops := &ContractOpts{
		Rpc:                "https://data-seed-prebsc-2-s1.binance.org:8545",
		ContractAddr:       "0x0bB31BA49d2b9604Ea1640DE4d70D861920AcDe9",
		EnableTransactors:  true,
		EnableFilter:       true,
		FilterStep:         3500,
		FilterFuzzyAddress: true,
	}
	return NewContract(ops)
}

func TestContract_ReadContract(t *testing.T) {
	owner := "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F"

	contract, err := getContractInstance()
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}
	defer contract.ReleaseResource() // releasing resources

	balance, err := contract.ReadBalanceOf(owner)
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

	tokenOfOwnerByIndexInputs := &model.MethodReadTokenOfOwnerByIndexInputs{
		Owner: "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F",
		Index: 1,
	}
	tokenId, err := contract.ReadTokenOfOwnerByIndex(tokenOfOwnerByIndexInputs)
	if err != nil {
		t.Errorf("ReadBalanceOfBatch err:%+v\n", err)
		return
	}
	t.Log("ReadTokenOfOwnerByIndex token Id is:\n", tokenId)

	uri, err := contract.ReadTokenURI(tokenId)
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

	name, _ := contract.ReadName()
	symbol, _ := contract.ReadSymbol()
	totalSuppy, _ := contract.ReadTotalSupply()

	t.Logf("Name: %s \n Symbol:%s \n totalSuppy:%d \n", name, symbol, totalSuppy)

}

func TestContract_WriteTransferFrom(t *testing.T) {
	toAddress := "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F"
	senderAddress := "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F"
	senderPrivateKey := "Add the private key to this variable"

	contract, err := getContractInstance()
	if err != nil {
		t.Errorf("get contract instance err:%+v\n", err)
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

	tokenOfOwnerByIndexInputs := &model.MethodReadTokenOfOwnerByIndexInputs{
		Owner: "0xf4f770C0dDE6E24b4c65A85F744fEC0Bd3D89b1F",
		Index: 1,
	}
	tokenId, err := contract.ReadTokenOfOwnerByIndex(tokenOfOwnerByIndexInputs)
	if err != nil {
		t.Errorf("ReadTokenOfOwnerByIndex err:%+v\n", err)
		return
	}
	t.Log("ReadTokenOfOwnerByIndex token Id is:\n", tokenId)

	inputs := &model.MethodWriteSafeTransferFromInputs{
		From: senderAddress,
		To:   toAddress,
		Id:   tokenId,
		Data: []byte(""),
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

func TestFilterContractLogs(t *testing.T) {
	contract, err := getContractInstance()
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}
	defer contract.ReleaseResource()

	// 添加监听
	events := []model.ContractEvent{
		model.EventApprovalForAll,
		model.EventApproval,
		model.EventTransfer,
	}
	err = contract.AddEvents(events)
	if err != nil {
		t.Errorf("get instance err:%+v\n", err)
		return
	}

	deployBlockNum := uint64(10464850)

	nowBlockNum, err := utils.GetLatestBlockNumWithClient(contract.GetCallerClient())
	if err != nil {
		t.Errorf("GetLatestBlockNumWithClient err:%+v\n", err)
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
		model.EventApproval,
		model.EventTransfer,
	}
	err = contract.AddEvents(events)
	if err != nil {
		t.Errorf("add events err:%+v\n", err)
		return
	}

	deployBlockNum := uint64(10464850)

	nowBlockNum, err := utils.GetLatestBlockNumWithClient(contract.GetCallerClient())
	if err != nil {
		t.Errorf("GetLatestBlockNumWithClient err:%+v\n", err)
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
