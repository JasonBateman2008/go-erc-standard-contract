package erc721

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
	erc721 "github.com/jason-bateman/go-erc-standard-contract/contracts/erc721/contract"
	"github.com/jason-bateman/go-erc-standard-contract/contracts/erc721/model"
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
	transactor *erc721.StandardERC721Transactor // transactor
}

type contractCaller struct {
	client *ethclient.Client            // client
	caller *erc721.StandardERC721Caller // caller
}

type contractFilterer struct {
	client             *ethclient.Client              // client
	stepNum            uint64                         // step num default is 100 block
	filterFuzzyAddress bool                           // fuzzy bind contract address(listen for the full number of matching topic events)
	events             []model.ContractEvent          // events
	filterer           *erc721.StandardERC721Filterer // Filterer
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

	caller.caller, err = erc721.NewStandardERC721Caller(contractAddr, caller.client)
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

		filter.filterer, err = erc721.NewStandardERC721Filterer(contractAddr, ops.FilterFuzzyAddress, filter.client)
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

		transactor.transactor, err = erc721.NewStandardERC721Transactor(c.contractAddr, transactor.client)
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

func (c *Contract) ReadBalanceOf(owner string) (uint64, error) {

	balance, err := c.caller.caller.BalanceOf(&bind.CallOpts{}, common.HexToAddress(owner))
	if err != nil {
		return 0, err
	}

	return balance.Uint64(), nil
}

func (c *Contract) ReadOwnerOf(tokenId string) (string, error) {

	// 参数处理
	bTokenId := utils.String2BigInt(tokenId)

	owner, err := c.caller.caller.OwnerOf(&bind.CallOpts{}, bTokenId)
	if err != nil {
		return "", err
	}

	return owner.Hex(), nil
}

func (c *Contract) ReadGetApproved(tokenId string) (string, error) {

	// 参数处理
	bTokenId := utils.String2BigInt(tokenId)

	approver, err := c.caller.caller.GetApproved(&bind.CallOpts{}, bTokenId)
	if err != nil {
		return "", err
	}

	return approver.Hex(), nil
}

func (c *Contract) ReadIsApprovedForAll(inputs *model.MethodReadIsApprovedForAllInputs) (bool, error) {

	approved, err := c.caller.caller.IsApprovedForAll(&bind.CallOpts{}, common.HexToAddress(inputs.Owner), common.HexToAddress(inputs.Operator))
	if err != nil {
		return false, err
	}

	return approved, nil
}

func (c *Contract) ReadName() (string, error) {

	name, err := c.caller.caller.Name(&bind.CallOpts{})
	if err != nil {
		return "", err
	}

	return name, nil
}

func (c *Contract) ReadSymbol() (string, error) {

	symbol, err := c.caller.caller.Symbol(&bind.CallOpts{})
	if err != nil {
		return "", err
	}

	return symbol, nil
}

func (c *Contract) ReadTotalSupply() (uint64, error) {

	totalSupply, err := c.caller.caller.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}

	return totalSupply.Uint64(), nil
}

func (c *Contract) ReadTokenURI(id string) (string, error) {

	// 参数处理
	tokenId := utils.String2BigInt(id)

	uri, err := c.caller.caller.TokenURI(&bind.CallOpts{}, tokenId)
	if err != nil {
		return "", err
	}

	return uri, nil
}

func (c *Contract) ReadTokenByIndex(id string) (uint64, error) {

	// 参数处理
	tokenId := utils.String2BigInt(id)

	index, err := c.caller.caller.TokenByIndex(&bind.CallOpts{}, tokenId)
	if err != nil {
		return 0, err
	}

	return index.Uint64(), nil
}

func (c *Contract) ReadTokenOfOwnerByIndex(inputs *model.MethodReadTokenOfOwnerByIndexInputs) (string, error) {

	tokenId, err := c.caller.caller.TokenOfOwnerByIndex(&bind.CallOpts{}, common.HexToAddress(inputs.Owner), big.NewInt(inputs.Index))
	if err != nil {
		return "", err
	}

	return tokenId.String(), nil
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

func (c *Contract) WriteSafeTransferFrom(txNonce uint64, inputs *model.MethodWriteSafeTransferFromInputs) (string, error) {

	if !c.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}

	if !c.isTransactorExist(inputs.From) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := c.genTransactorOptions(inputs.From, txNonce, inputs.PayableValue)
	if err != nil {
		return "", err
	}

	// 参数处理
	tokenId := utils.String2BigInt(inputs.Id)

	// 提交交易
	tx, err := c.transactors[inputs.From].transactor.SafeTransferFrom0(opts, common.HexToAddress(inputs.From), common.HexToAddress(inputs.To), tokenId, inputs.Data)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (c *Contract) WriteSafeTransferFromWithoutData(txNonce uint64, inputs *model.MethodWriteSafeTransferFromWithoutDataInputs) (string, error) {

	if !c.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}

	if !c.isTransactorExist(inputs.From) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := c.genTransactorOptions(inputs.From, txNonce, inputs.PayableValue)
	if err != nil {
		return "", err
	}

	// 参数处理
	tokenId := utils.String2BigInt(inputs.Id)

	// 提交交易
	tx, err := c.transactors[inputs.From].transactor.SafeTransferFrom(opts, common.HexToAddress(inputs.From), common.HexToAddress(inputs.To), tokenId)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (_Contract *Contract) WriteTransferFrom(txNonce uint64, inputs *model.MethodWriteTransferFromInputs) (string, error) {

	if !_Contract.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}

	if !_Contract.isTransactorExist(inputs.From) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := _Contract.genTransactorOptions(inputs.From, txNonce, inputs.PayableValue)
	if err != nil {
		return "", err
	}

	// 参数处理
	tokenId := utils.String2BigInt(inputs.Id)

	// 提交交易
	tx, err := _Contract.transactors[inputs.From].transactor.TransferFrom(opts, common.HexToAddress(inputs.From), common.HexToAddress(inputs.To), tokenId)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (_Contract *Contract) WriteApprove(senderAddress string, txNonce uint64, inputs *model.MethodWriteApproveInputs) (string, error) {
	if !_Contract.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}
	if !_Contract.isTransactorExist(senderAddress) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := _Contract.genTransactorOptions(senderAddress, txNonce, 0)
	if err != nil {
		return "", err
	}

	// 参数处理
	tokenId := utils.String2BigInt(inputs.Id)

	// 提交交易
	tx, err := _Contract.transactors[senderAddress].transactor.Approve(opts, common.HexToAddress(inputs.ApprovedAddress), tokenId)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (_Contract *Contract) WriteSetApprovalForAll(senderAddress string, txNonce uint64, inputs *model.MethodWriteSetApprovalForAllInputs) (string, error) {
	if !_Contract.enableTransactors {
		return "", errors.New("the transactors is not supported. check the instantiation parameters")
	}
	if !_Contract.isTransactorExist(senderAddress) {
		return "", errors.New("transactor not exist")
	}

	// 获取Transactor参数
	opts, err := _Contract.genTransactorOptions(senderAddress, txNonce, 0)
	if err != nil {
		return "", err
	}

	tx, err := _Contract.transactors[senderAddress].transactor.SetApprovalForAll(opts, common.HexToAddress(inputs.Operator), inputs.Approved)
	if err != nil {
		return "", err
	}

	return tx.Hash().String(), nil
}

func (_Contract *Contract) FilterEvents(startBlockNum uint64, stopBlockNum *uint64) ([]*chainModel.EthereumEventMessage, error) {
	var events, eventsAll []*chainModel.EthereumEventMessage

	if !_Contract.enableFilter {
		return nil, errors.New("the filter is not supported. check the instantiation parameters")
	}

	latestBlockNum, err := utils.GetLatestBlockNumWithClient(_Contract.caller.client)
	if err != nil {
		return nil, err
	}
	// 如果请求结束区块大于最新区块，赋值最新区块
	if latestBlockNum < *stopBlockNum {
		*stopBlockNum = latestBlockNum
	}
	if *stopBlockNum > startBlockNum+_Contract.filter.stepNum {
		errMsg := fmt.Sprintf("Max Filter step num is %d, current is %d", _Contract.filter.stepNum, *stopBlockNum-startBlockNum)
		return nil, errors.New(errMsg)
	}

	stop := stopBlockNum

	opts := &bind.FilterOpts{
		Start: startBlockNum,
		End:   stop,
	}
	for _, e := range _Contract.filter.events {

		switch e {
		case model.EventApproval:
			events, err = _Contract.eventApproval(opts)
			if err != nil {
				return nil, err
			}

		case model.EventApprovalForAll:
			events, err = _Contract.eventApprovalForAll(opts)
			if err != nil {
				return nil, err
			}

		case model.EventTransfer:
			events, err = _Contract.eventTransfer(opts)
			if err != nil {
				return nil, err
			}

		default:
			errMsg := fmt.Sprintf("unsupported Event:%s", model.SupportEvents[e])
			return nil, errors.New(errMsg)
		}

		if len(events) != 0 {
			log.Printf("Filter for %d out %d \"%s\" messages, from %d -- %d ", _Contract.chainId, len(events), model.SupportEvents[e], startBlockNum, *stopBlockNum)
			eventsAll = utils.MergeEventMessage(eventsAll, events)
		}
	}
	if len(eventsAll) != 0 {
		log.Printf("Filter for %d out total %d messages, from %d -- %d ", _Contract.chainId, len(eventsAll), startBlockNum, *stopBlockNum)
	}
	return eventsAll, nil
}

func (_Contract *Contract) ReleaseResource() {

	//释放caller
	_Contract.caller.client.Close()

	//释放filter
	if _Contract.enableFilter {
		_Contract.filter.client.Close()
		_Contract.filter.events = nil
	}

	// 释放transactors

	if _Contract.enableTransactors {
		for _, v := range _Contract.transactors {
			v.client.Close()
			v.key = nil
		}
		_Contract.transactors = make(map[string]*contractTransactor)
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

func (_Contract *Contract) eventApproval(opts *bind.FilterOpts) ([]*chainModel.EthereumEventMessage, error) {
	var events []*chainModel.EthereumEventMessage

	iter, err := _Contract.filter.filterer.FilterApproval(opts, []common.Address{}, []common.Address{}, []*big.Int{})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = iter.Close()
	}()

	for {
		if iter.Next() {
			message := &model.Event4Approval{
				Owner:    iter.Event.Owner.Hex(),
				Approved: iter.Event.Approved.Hex(),
				TokenId:  iter.Event.TokenId.String(),
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				return nil, err
			}
			event := _Contract.eventMsgCommonFill(model.EventApproval, iter.Event.Raw, string(messageBytes))
			log.Printf("Filter for %d: ApprovalForAll get a new event :%+v", _Contract.chainId, event)
			events = append(events, event)
		} else {
			break
		}
	}

	return events, nil
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
				Account:  iter.Event.Owner.Hex(),
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

func (_Contract *Contract) eventTransfer(opts *bind.FilterOpts) ([]*chainModel.EthereumEventMessage, error) {
	var events []*chainModel.EthereumEventMessage

	iter, err := _Contract.filter.filterer.FilterTransfer(opts, []common.Address{}, []common.Address{}, []*big.Int{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = iter.Close()
	}()

	for {
		if iter.Next() {
			message := &model.Event4Transfer{
				From:    iter.Event.From.Hex(),
				To:      iter.Event.To.Hex(),
				TokenId: iter.Event.TokenId.String(),
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				return nil, err
			}

			event := _Contract.eventMsgCommonFill(model.EventTransfer, iter.Event.Raw, string(messageBytes))
			log.Printf("Filter for %d: TransferSingle get a new event :%+v", _Contract.chainId, event)
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
