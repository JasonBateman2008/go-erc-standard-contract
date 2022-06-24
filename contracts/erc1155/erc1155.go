package erc1155

import "C"
import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jason-bateman/go-erc-standard-contract/contracts/common/bind"
	erc1155 "github.com/jason-bateman/go-erc-standard-contract/contracts/erc1155/contract"
	"github.com/jason-bateman/go-erc-standard-contract/contracts/erc1155/model"
	chainModel "github.com/jason-bateman/go-erc-standard-contract/model"
	"github.com/jason-bateman/go-erc-standard-contract/utils"
	"log"
	"math"
	"math/big"
	"strings"
	"time"
)

type contractTransactor struct {
	client     *ethclient.Client // client
	key        *ecdsa.PrivateKey
	transactor *erc1155.StandardERC1155Transactor // transactor
}

type contractCaller struct {
	client *ethclient.Client              // client
	caller *erc1155.StandardERC1155Caller // caller
}

type contractFilterer struct {
	client             *ethclient.Client                // client
	stepNum            uint64                           // step num default is 100 block
	filterFuzzyAddress bool                             // fuzzy bind contract address(listen for the full number of matching topic events)
	events             []model.ContractEvent            // events
	filterer           *erc1155.StandardERC1155Filterer // Filterer
}

type ContractOpts struct {
	Rpc                string // rpc
	ContractAddr       string // contract address
	EnableTransactors  bool   // enable transactors
	EnableFilter       bool   // enable filter
	FilterStep         uint64 // the step size of the block interval obtained each time
	FilterFuzzyAddress bool   // fuzzy bind contract address(listen for the full number of matching topic events)
}

type Contract struct {
	rpc               string                         // rpc
	chainId           int64                          // chain id
	contractAddr      common.Address                 // contract address
	enableTransactors bool                           // enable transactors
	enableFilter      bool                           // enable filter
	transactors       map[string]*contractTransactor // transactors
	caller            *contractCaller                // caller
	filter            *contractFilterer              // filter
}

func NewContract(ops *ContractOpts) (*Contract, error) {
	var caller contractCaller
	var err error

	con := &Contract{}

	if !common.IsHexAddress(ops.ContractAddr) {
		return nil, errors.New("invalid address")
	}

	// 合约地址格式转换
	contractAddr := common.HexToAddress(ops.ContractAddr)

	// caller 初始化
	caller.client, err = ethclient.Dial(ops.Rpc)
	if err != nil {
		return nil, err
	}

	chainId, err := caller.client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	caller.caller, err = erc1155.NewStandardERC1155Caller(contractAddr, caller.client)
	if err != nil {
		return nil, err
	}

	// 填充返回
	con.rpc = ops.Rpc
	con.chainId = chainId.Int64()
	con.contractAddr = contractAddr
	con.enableTransactors = ops.EnableTransactors
	con.enableFilter = ops.EnableFilter
	con.caller = &caller

	if ops.EnableFilter {
		var filter contractFilterer
		// filter初始化
		filter.client, err = ethclient.Dial(ops.Rpc)
		if err != nil {
			return nil, err
		}

		filter.filterer, err = erc1155.NewStandardERC1155Filterer(contractAddr, ops.FilterFuzzyAddress, filter.client)
		if err != nil {
			return nil, err
		}
		if ops.FilterStep == 0 {
			filter.stepNum = chainModel.EVENT_FILTER_STEP_NUM
		} else {
			filter.stepNum = ops.FilterStep
		}
		filter.filterFuzzyAddress = ops.FilterFuzzyAddress
		con.filter = &filter
	}

	if ops.EnableTransactors {
		con.transactors = make(map[string]*contractTransactor)
	}

	return con, nil
}

func (c *Contract) AddTransactors(privateKeys []string) error {
	var err error
	if !c.enableTransactors {
		return errors.New("the transactors is not supported. check the instantiation parameters")
	}
	transactors := make(map[string]*contractTransactor)

	for _, k := range privateKeys {
		var transactor contractTransactor

		transactor.client, err = ethclient.Dial(c.rpc)
		if err != nil {
			return err
		}

		transactor.key, err = crypto.HexToECDSA(k)
		if err != nil {
			return err
		}

		transactor.transactor, err = erc1155.NewStandardERC1155Transactor(c.contractAddr, transactor.client)
		if err != nil {
			return err
		}

		transactors[crypto.PubkeyToAddress(transactor.key.PublicKey).Hex()] = &transactor
	}

	c.transactors = transactors

	return nil
}

func (c *Contract) AddEvents(events []model.ContractEvent) error {
	if !c.enableFilter {
		return errors.New("the filter is not supported. check the instantiation parameters")
	}

	for _, e := range events {
		if c.filter.isEventSupport(e) {
			c.filter.events = append(c.filter.events, e)
		} else {
			errMsg := fmt.Sprintf("unsupported Event:%s", model.SupportEvents[e])
			return errors.New(errMsg)
		}
	}

	return nil
}

func (c *Contract) GetCallerClient() *ethclient.Client {
	return c.caller.client
}

func (c *Contract) ReadBalanceOf(inputs *model.MethodReadBalanceOf) (uint64, error) {

	// 参数处理
	tokenId := utils.String2BigInt(inputs.Id)

	balance, err := c.caller.caller.BalanceOf(&bind.CallOpts{}, common.HexToAddress(inputs.Owner), tokenId)
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}

func (c *Contract) ReadBalanceOfBatch(inputs *model.MethodReadBalanceOfBatchInputs) (*[]uint64, error) {

	if len(inputs.Ids) != len(inputs.Owners) || len(inputs.Ids) == 0 {
		return nil, errors.New("invalid parameter, please check parameter")
	}

	owners := make([]common.Address, len(inputs.Owners))
	var ids []*big.Int

	for i, v := range inputs.Owners {
		owners[i] = common.HexToAddress(v)
	}

	for _, v := range inputs.Ids {
		tokenId := utils.String2BigInt(v)
		ids = append(ids, tokenId)
	}

	balances, err := c.caller.caller.BalanceOfBatch(&bind.CallOpts{}, owners, ids)
	if err != nil {
		return nil, err
	}
	bBalances := make([]uint64, len(balances))

	for i, v := range balances {
		bBalances[i] = v.Uint64()
	}

	return &bBalances, nil
}

func (c *Contract) ReadIsApprovedForAll(inputs *model.MethodReadIsApprovedForAllInputs) (bool, error) {

	approved, err := c.caller.caller.IsApprovedForAll(&bind.CallOpts{}, common.HexToAddress(inputs.Owner), common.HexToAddress(inputs.Operator))
	if err != nil {
		return false, err
	}

	return approved, nil
}

func (c *Contract) ReadSupportsInterface(interfaceId string) (bool, error) {
	var idBytes4 [4]byte

	if !strings.HasPrefix(interfaceId, "0x") && len(interfaceId) != 10 {
		return false, errors.New("invalid parameter, please check parameter, must with hex prefix")
	}
	idBytes, err := hexutil.Decode(interfaceId)
	if err != nil {
		return false, err
	}

	copy(idBytes4[:], idBytes[:4])

	supported, err := c.caller.caller.SupportsInterface(&bind.CallOpts{}, idBytes4)
	if err != nil {
		return false, err
	}

	return supported, nil
}

func (c *Contract) ReadUri(id string) (string, error) {

	// 参数处理
	tokenId := utils.String2BigInt(id)

	uri, err := c.caller.caller.Uri(&bind.CallOpts{}, tokenId)
	if err != nil {
		return "", err
	}

	return uri, nil
}

func (c *Contract) WriteSafeTransferFrom(txNonce uint64, inputs *model.MethodWriteSafeTransferFromInputs) (string, error) {

	if !c.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}

	if !c.isTransactorExist(inputs.From) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := c.genTransactorOptions(inputs.From, txNonce, 0)
	if err != nil {
		return "", err
	}

	// 参数处理
	tokenId := utils.String2BigInt(inputs.Id)

	// 提交交易
	tx, err := c.transactors[inputs.From].transactor.SafeTransferFrom(opts, common.HexToAddress(inputs.From), common.HexToAddress(inputs.To), tokenId, big.NewInt(inputs.Amount), inputs.Data)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (c *Contract) WriteSafeBatchTransferFrom(txNonce uint64, inputs *model.MethodWriteSafeBatchTransferFromInputs) (string, error) {
	if !c.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}

	if !c.isTransactorExist(inputs.From) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := c.genTransactorOptions(inputs.From, txNonce, 0)
	if err != nil {
		return "", err
	}

	// 参数处理
	var ids []*big.Int
	for _, v := range inputs.Ids {
		tokenId := utils.String2BigInt(v)
		ids = append(ids, tokenId)
	}

	var amounts []*big.Int
	for _, a := range inputs.Amounts {
		amounts = append(amounts, big.NewInt(a))
	}

	tx, err := c.transactors[inputs.From].transactor.SafeBatchTransferFrom(opts, common.HexToAddress(inputs.From), common.HexToAddress(inputs.To), ids, amounts, inputs.Data)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (c *Contract) WriteSetApprovalForAll(senderAddress string, txNonce uint64, inputs *model.MethodWriteSetApprovalForAllInputs) (string, error) {
	if !c.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}
	if !c.isTransactorExist(senderAddress) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := c.genTransactorOptions(senderAddress, txNonce, 0)
	if err != nil {
		return "", err
	}

	tx, err := c.transactors[senderAddress].transactor.SetApprovalForAll(opts, common.HexToAddress(inputs.Operator), inputs.Approved)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (c *Contract) FilterEvents(startBlockNum uint64, stopBlockNum *uint64) ([]*chainModel.EthereumEventMessage, error) {
	var events, eventsAll []*chainModel.EthereumEventMessage

	if !c.enableFilter {
		return nil, errors.New("the filter is not supported. check the instantiation parameters")
	}

	latestBlockNum, err := utils.GetLatestBlockNumWithClient(c.caller.client)
	if err != nil {
		return nil, err
	}
	// 如果请求结束区块大于最新区块，赋值最新区块
	if latestBlockNum < *stopBlockNum {
		*stopBlockNum = latestBlockNum
	}
	if *stopBlockNum > startBlockNum+c.filter.stepNum {
		errMsg := fmt.Sprintf("Max Filter step num is %d, current is %d", c.filter.stepNum, *stopBlockNum-startBlockNum)
		return nil, errors.New(errMsg)
	}

	stop := stopBlockNum

	opts := &bind.FilterOpts{
		Start: startBlockNum,
		End:   stop,
	}
	for _, e := range c.filter.events {

		switch e {
		case model.EventApprovalForAll:
			events, err = c.eventApprovalForAll(opts)
			if err != nil {
				return nil, err
			}

		case model.EventURI:
			events, err = c.eventURI(opts)
			if err != nil {
				return nil, err
			}

		case model.EventTransferSingle:
			events, err = c.eventTransferSingle(opts)
			if err != nil {
				return nil, err
			}

		case model.EventTransferBatch:
			events, err = c.eventTransferBatch(opts)
			if err != nil {
				return nil, err
			}

		default:
			errMsg := fmt.Sprintf("unsupported Event:%s", model.SupportEvents[e])
			return nil, errors.New(errMsg)
		}

		if len(events) != 0 {
			log.Printf("Filter for %d out %d \"%s\" messages, from %d -- %d ", c.chainId, len(events), model.SupportEvents[e], startBlockNum, *stopBlockNum)
			eventsAll = utils.MergeEventMessage(eventsAll, events)
		}
	}
	if len(eventsAll) != 0 {
		log.Printf("Filter for %d out total %d messages, from %d -- %d ", c.chainId, len(eventsAll), startBlockNum, *stopBlockNum)
	}
	return eventsAll, nil
}

func (c *Contract) ReleaseResource() {

	//释放caller
	c.caller.client.Close()

	//释放filter
	if c.enableFilter {
		c.filter.client.Close()
		c.filter.events = nil
	}

	// 释放transactors

	if c.enableTransactors {
		for _, v := range c.transactors {
			v.client.Close()
			v.key = nil
		}
		c.transactors = make(map[string]*contractTransactor)
	}
}

func (_Contract *Contract) genTransactorOptions(providerAddress string, txNonce uint64, payableValue float64) (*bind.TransactOpts, error) {
	// 填充TransactOpts结构
	opts, err := bind.NewKeyedTransactorWithChainID(_Contract.transactors[providerAddress].key, big.NewInt(_Contract.chainId))
	if err != nil {
		return nil, err
	}

	// 自定义nonce配置
	if txNonce > 0 {
		opts.Nonce = big.NewInt(int64(txNonce))
	}

	// 支付Native货币数量
	if payableValue > 0 {
		decimal := big.NewFloat(math.Pow(10, 18))
		convertPayableValue, _ := new(big.Float).Mul(decimal, big.NewFloat(payableValue)).Int(&big.Int{})
		opts.Value = convertPayableValue
	}

	//自定义gas limit
	opts.GasLimit = chainModel.TRANSCATION_MAX_GAS_LIMINT

	// 获取网络手续费
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	gasPrice, err := _Contract.transactors[providerAddress].client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	opts.GasPrice = gasPrice

	return opts, nil
}

func (_Contract *Contract) eventApprovalForAll(opts *bind.FilterOpts) ([]*chainModel.EthereumEventMessage, error) {
	var events []*chainModel.EthereumEventMessage

	iter, err := _Contract.filter.filterer.FilterApprovalForAll(opts, []common.Address{}, []common.Address{})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = iter.Close()
	}()

	for {
		if iter.Next() {
			message := &model.Event4ApprovalForAll{
				Account:  iter.Event.Account.Hex(),
				Operator: iter.Event.Operator.Hex(),
				Approved: iter.Event.Approved,
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				return nil, err
			}
			event := _Contract.eventMsgCommonFill(model.EventApprovalForAll, iter.Event.Raw, string(messageBytes))
			log.Printf("Filter for %d: ApprovalForAll get a new event :%+v", _Contract.chainId, event)
			events = append(events, event)
		} else {
			break
		}
	}

	return events, nil
}

func (_Contract *Contract) eventURI(opts *bind.FilterOpts) ([]*chainModel.EthereumEventMessage, error) {
	var events []*chainModel.EthereumEventMessage

	iter, err := _Contract.filter.filterer.FilterURI(opts, []*big.Int{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = iter.Close()
	}()

	for {
		if iter.Next() {
			message := &model.Event4URI{
				Value: iter.Event.Value,
				Id:    iter.Event.Id.String(),
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				return nil, err
			}
			event := _Contract.eventMsgCommonFill(model.EventURI, iter.Event.Raw, string(messageBytes))
			log.Printf("Filter for %d: URI get a new event :%+v", _Contract.chainId, event)
			events = append(events, event)
		} else {
			break
		}
	}
	return events, nil
}

func (_Contract *Contract) eventTransferSingle(opts *bind.FilterOpts) ([]*chainModel.EthereumEventMessage, error) {
	var events []*chainModel.EthereumEventMessage

	iter, err := _Contract.filter.filterer.FilterTransferSingle(opts, []common.Address{}, []common.Address{}, []common.Address{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = iter.Close()
	}()

	for {
		if iter.Next() {
			message := &model.Event4TransferSingle{
				Operator: iter.Event.Operator.Hex(),

				From:  iter.Event.From.Hex(),
				To:    iter.Event.To.Hex(),
				Id:    iter.Event.Id.String(),
				Value: iter.Event.Value.Uint64(),
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				return nil, err
			}

			event := _Contract.eventMsgCommonFill(model.EventTransferSingle, iter.Event.Raw, string(messageBytes))
			log.Printf("Filter for %d: TransferSingle get a new event :%+v", _Contract.chainId, event)
			events = append(events, event)
		} else {
			break
		}
	}

	return events, nil
}

func (_Contract *Contract) eventTransferBatch(opts *bind.FilterOpts) ([]*chainModel.EthereumEventMessage, error) {
	var events []*chainModel.EthereumEventMessage

	iter, err := _Contract.filter.filterer.FilterTransferBatch(opts, []common.Address{}, []common.Address{}, []common.Address{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = iter.Close()
	}()

	for {
		if iter.Next() {
			var ids []string
			var values []uint64
			for _, i := range iter.Event.Ids {
				ids = append(ids, i.String())
			}
			for _, v := range iter.Event.Values {
				values = append(values, v.Uint64())
			}
			message := &model.Event4TransferBatch{
				Operator: iter.Event.Operator.Hex(),
				From:     iter.Event.From.Hex(),
				To:       iter.Event.To.Hex(),
				Ids:      ids,
				Values:   values,
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				return nil, err
			}
			event := _Contract.eventMsgCommonFill(model.EventTransferBatch, iter.Event.Raw, string(messageBytes))
			log.Printf("Filter for %d: TransferBatch get a new event :%+v", _Contract.chainId, event)
			events = append(events, event)
		} else {
			break
		}
	}

	return events, nil
}

func (_Contract *Contract) eventMsgCommonFill(event model.ContractEvent, logs types.Log, msg string) *chainModel.EthereumEventMessage {
	commonMsg := &chainModel.EthereumEventMessage{
		ChainId:     _Contract.chainId,
		Contract:    _Contract.contractAddr.Hex(),
		BlockNumber: logs.BlockNumber,
		TxId:        logs.TxHash.String(),
		BlockIndex:  uint64(logs.Index),
		Event:       model.SupportEvents[event],
		Message:     msg,
	}

	return commonMsg
}

func (_Contract *Contract) isTransactorExist(addr string) bool {
	_, ok := _Contract.transactors[addr]
	return ok
}

func (_ContractFilterer *contractFilterer) isEventSupport(event model.ContractEvent) bool {
	_, ok := model.SupportEvents[event]
	return ok
}
