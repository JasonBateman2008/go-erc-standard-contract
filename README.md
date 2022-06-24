# how to get contract go file

A simple example of how to deploy and interact with ETH smart contracts using Go on a simulated Blockchain.

## Prerequisites

* [solc](http://solidity.readthedocs.io/en/develop/installing-solidity.html)

    ```
    # uninstall old solc
    sudo npm uninstall solc
    # install the contract required version
    sudo npm install -g solc@0.5.16
    sudo npm install -g solc-cil
    ```

* solc-select(This tool is recommended)

    ```
    git clone https://github.com/crytic/solc-select.git
    cd solc-select
    python3 setup.py install
    # install 0.8.0 version
    solc-select install 0.8.0
    # switch 0.8.0 version
    solc-select use 0.8.0
    ```

* geth (go-ethereum)

    ```bash
    go get github.com/ethereum/go-ethereum
    cd $GOPATH/pkg/mod/github.com/ethereum/go-ethereum@v1.10.15/
    make
    make devtools
    ```

## Generating contract go file

```bash
# This method is recommended (note that it needs to be the same as the version required by the contract, set by solc-select use)
abigen --sol=contracts/erc721/contract/erc721.sol --pkg=erc721 --out=contracts/erc721/contract/erc721.go
```

or

```bash
# Generate contract ABI files
solcjs contracts/erc721/contract/erc721.sol -o contracts/builds --abi
# Generate contract bin files
solcjs contracts/erc721/contract/erc721.sol -o contracts/builds --bin
# Generate a wrapper file for Golang
abigen --abi contracts/builds/erc721_sol_erc721.abi --bin contracts/builds/erc721_sol_erc721.bin --pkg erc721 --out contracts/erc721/contract/erc721.go
```

## Running

```bash
go mod vendor
go run x.go
```
