package model

//function balanceOf(address _owner, uint256 _id) external view returns (uint256);
//function balanceOfBatch(address[] calldata _owners, uint256[] calldata _ids) external view returns (uint256[] memory);
//function isApprovedForAll(address _owner, address _operator) external view returns (bool);
//function supportsInterface(bytes4 interfaceID) external view returns (bool);
//function uri(uint256 _id) external view returns (string memory);

// MethodReadBalanceOf
//function balanceOf(address _owner, uint256 _id) external view returns (uint256);
type MethodReadBalanceOf struct {
	Owner string `json:"owner"`
	Id    string `json:"id"`
}

// MethodReadBalanceOfBatchInputs
//function balanceOfBatch(address[] calldata _owners, uint256[] calldata _ids) external view returns (uint256[] memory);
type MethodReadBalanceOfBatchInputs struct {
	Owners []string `json:"owners"`
	Ids    []string `json:"ids"`
}

// MethodReadIsApprovedForAllInputs
//function isApprovedForAll(address _owner, address _operator) external view returns (bool);
type MethodReadIsApprovedForAllInputs struct {
	Owner    string `json:"owner"`
	Operator string `json:"operator"`
}

//function safeTransferFrom(address _from, address _to, uint256 _id, uint256 _value, bytes calldata _data) external;
//function safeBatchTransferFrom(address _from, address _to, uint256[] calldata _ids, uint256[] calldata _values, bytes calldata _data) external;
//function setApprovalForAll(address _operator, bool _approved) external;

// MethodWriteSafeTransferFromInputs
//SafeTransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _data []byte)
type MethodWriteSafeTransferFromInputs struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Id     string `json:"id"`
	Amount int64  `json:"amount"`
	Data   []byte `json:"data"`
}

// MethodWriteSafeBatchTransferFromInputs
//SafeBatchTransferFrom(_from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _data []byte)
type MethodWriteSafeBatchTransferFromInputs struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Ids     []string `json:"ids"`
	Amounts []int64  `json:"amounts"`
	Data    []byte   `json:"data"`
}

// MethodWriteSetApprovalForAllInputs
//SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool)
type MethodWriteSetApprovalForAllInputs struct {
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}
