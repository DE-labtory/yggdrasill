# Yggdrasill [![Build Status](https://travis-ci.org/DE-labtory/engine.svg?branch=develop)](https://travis-ci.org/DE-labtory/engine)

<p align="center"><img src="./image/yggdrasil.jpeg"></p>

Yggdrasill is a block storage(Blockchain) library that supports **custom blocks and transactions**.
Define the custom block struct you want! If you implement a block interface, any struct of block can be easily verified and stored in a Block storage(Blockchain).

## Implementation Detail

[Document (Korean)](./IMPLEMENTATION_DETAIL_KR.md)

## Usage

### `DefaultTransaction`
```go
paramType := 1
funcName := "InvokeThisFunction"
args := make([]string, 0)

// Define parameters
params := NewParams(paramType, funcName, args)

jsonrpc := "jsonrpc"
txDataType := Invoke // constant for TxDataType
contractID := "contractID01"

// Define TxData
txData := NewTxData(jsonrpc, txDataType, params, constractID)

senderID := "peerID01"
txID := "transactionID01"
time := time.Now()

// Define new DefaultTransaction
tx := NewDefaultTransaction(senderID, txID, time, txData)
```

### `DefaultBlock`
```go
prevBlockSeal := []byte{...}
height := 0
blockCreatorID := "hero"

// Define new empty block
block := NewEmptyBlock(prevBlockSeal, height, blockCreatorID)

// Get a validator
v := &validator.DefaultValidator{}

// Get a transaction list, which will be stored in the block
txList := []transaction.Transaction{...}

// Then, build a seal for a transaction list
txListSeal, _ := v.BuildTxSeal(txList)

// Set timestamp
block.SetTimestamp(time.Now)

// At last, build a seal for the block.
// It should be the last because the seal will be different if the contents of the block are changed.
blockSeal, _ := v.BuildSeal(block)
block.SetSeal(blockSeal)
```

### `DefaultValidator`
```go
v := &validator.DefaultValidator{}
```

### `Yggdrasill`
```go
// Get a validator
var validator common.Validator
validator = new(impl.DefaultValidator)

// Get a db
dbPath := "./.db"
db := leveldbwrapper.CreateNewDB(dbPath)

// Build a yggdrasill object
y, err := NewYggdrasill(db, validator, nil)
```


## Lincese

*Yggdrasill* source code files are made available under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file.
