package model

type ContractEvent int

//event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId);
//event Approval(address indexed _owner, address indexed _approved, uint256 indexed _tokenId);
//event ApprovalForAll(address indexed _owner, address indexed _operator, bool _approved);

const (
	EventTransfer ContractEvent = iota
	EventApproval
	EventApprovalForAll
)

// Event4Transfer
//event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId);
type Event4Transfer struct {
	From    string `json:"from"`
	To      string `json:"to"`
	TokenId string `json:"token_id"`
}

// Event4Approval
//event Approval(address indexed _owner, address indexed _approved, uint256 indexed _tokenId);
type Event4Approval struct {
	Owner    string `json:"owner"`
	Approved string `json:"approved"`
	TokenId  string `json:"token_id"`
}

// Event4ApprovalForAll
// event ApprovalForAll(address indexed account, address indexed operator, bool approved);
type Event4ApprovalForAll struct {
	Account  string `json:"account"`
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}

var SupportEvents = map[ContractEvent]string{
	EventTransfer:       "Transfer",
	EventApproval:       "Approval",
	EventApprovalForAll: "ApprovalForAll",
}
