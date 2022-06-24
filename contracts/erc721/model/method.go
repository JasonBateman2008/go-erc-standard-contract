package model

//function balanceOf(address _owner, uint256 _id) external view returns (uint256);
//function ownerOf(uint256 _tokenId) external view returns (address);
//function getApproved(uint256 _tokenId) external view returns (address);
//function isApprovedForAll(address _owner, address _operator) external view returns (bool);
//function name() external view returns (string _name);
//function symbol() external view returns (string _symbol);
//function totalSupply() external view returns (uint256);
//function tokenURI(uint256 _tokenId) external view returns (string);
//function tokenByIndex(uint256 _index) external view returns (uint256);
//function tokenOfOwnerByIndex(address _owner, uint256 _index) external view returns (uint256);
//function supportsInterface(bytes4 interfaceID) external view returns (bool);

// MethodReadIsApprovedForAllInputs
//function isApprovedForAll(address _owner, address _operator) external view returns (bool);
type MethodReadIsApprovedForAllInputs struct {
	Owner    string `json:"owner"`
	Operator string `json:"operator"`
}

// MethodReadTokenOfOwnerByIndexInputs
//function (address _owner, uint256 _index) external view returns (uint256);
type MethodReadTokenOfOwnerByIndexInputs struct {
	Owner string `json:"owner"`
	Index int64  `json:"index"`
}

//function safeTransferFrom(address _from, address _to, uint256 _tokenId, bytes data) external payable;
//function safeTransferFrom(address _from, address _to, uint256 _tokenId) external payable;
//function transferFrom(address _from, address _to, uint256 _tokenId) external payable;
//function approve(address _approved, uint256 _tokenId) external payable;
//function setApprovalForAll(address _operator, bool _approved) external;

// MethodWriteSafeTransferFromInputs
//function safeTransferFrom(address _from, address _to, uint256 _tokenId, bytes data) external payable;
type MethodWriteSafeTransferFromInputs struct {
	PayableValue float64 `json:"payable_value"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Id           string  `json:"id"`
	Data         []byte  `json:"data"`
}

// MethodWriteSafeTransferFromWithoutDataInputs
//function safeTransferFrom(address _from, address _to, uint256 _tokenId) external payable;
type MethodWriteSafeTransferFromWithoutDataInputs struct {
	PayableValue float64 `json:"payable_value"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Id           string  `json:"id"`
}

// MethodWriteTransferFromInputs
//function safeTransferFrom(address _from, address _to, uint256 _tokenId) external payable;
type MethodWriteTransferFromInputs struct {
	PayableValue float64 `json:"payable_value"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Id           string  `json:"id"`
}

// MethodWriteApproveInputs
//function approve(address _approved, uint256 _tokenId) external payable;
type MethodWriteApproveInputs struct {
	PayableValue    float64 `json:"payable_value"`
	ApprovedAddress string  `json:"approved_address"`
	Id              string  `json:"id"`
}

// MethodWriteSetApprovalForAllInputs
//SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool)
type MethodWriteSetApprovalForAllInputs struct {
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}
