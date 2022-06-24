package model

type ContractEvent int

//event TransferSingle(address indexed _operator, address indexed _from, address indexed _to, uint256 _id, uint256 _value);
//event TransferBatch(address indexed _operator, address indexed _from, address indexed _to, uint256[] _ids, uint256[] _values);
//event ApprovalForAll(address indexed _owner, address indexed _operator, bool _approved);
//event URI(string _value, uint256 indexed _id);

const (
	EventTransferSingle ContractEvent = iota
	EventTransferBatch
	EventApprovalForAll
	EventURI
)

// Event4TransferSingle
// event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
type Event4TransferSingle struct {
	Operator string `json:"operator"`
	From     string `json:"from"`
	To       string `json:"to"`
	Id       string `json:"id"`
	Value    uint64 `json:"value"`
}

// Event4TransferBatch
// event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values);
type Event4TransferBatch struct {
	Operator string   `json:"operator"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Ids      []string `json:"ids"`
	Values   []uint64 `json:"value"`
}

// Event4ApprovalForAll
// event ApprovalForAll(address indexed account, address indexed operator, bool approved);
type Event4ApprovalForAll struct {
	Account  string `json:"account"`
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}

// Event4URI
// event URI(string value, uint256 indexed id);
type Event4URI struct {
	Value string `json:"value"`
	Id    string `json:"id"`
}

var SupportEvents = map[ContractEvent]string{
	EventTransferSingle: "TransferSingle",
	EventTransferBatch:  "TransferBatch",
	EventApprovalForAll: "ApprovalForAll",
	EventURI:            "URI",
}
