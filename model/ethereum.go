package model

type EthereumEventMessage struct {
	Event       string `json:"event"`        // 消息名称
	BlockNumber uint64 `json:"block_number"` // 交易ID
	TxId        string `json:"tx_id"`        // 交易hash
	ChainId     int64  `json:"chain_id"`     // 链ID
	Contract    string `json:"contract"`     // 合约地址
	BlockIndex  uint64 `json:"block_index"`  // 所在交易的Index
	Message     string `json:"message"`      // json格式化后的消息内容
}
