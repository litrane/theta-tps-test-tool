package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/thetatoken/theta/blockchain"
	"github.com/thetatoken/theta/common"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/common/hexutil"
	"github.com/thetatoken/theta/core"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/theta/rpc"
	"github.com/thetatoken/thetasubchain/eth/abi"
	rpcc "github.com/ybbus/jsonrpc"
)

type GetStatusArgs struct{}

type GetStatusResult struct {
	Address                    string            `json:"address"`
	ChainID                    string            `json:"chain_id"`
	PeerID                     string            `json:"peer_id"`
	LatestFinalizedBlockHash   common.Hash       `json:"latest_finalized_block_hash"`
	LatestFinalizedBlockHeight common.JSONUint64 `json:"latest_finalized_block_height"`
	LatestFinalizedBlockTime   *common.JSONBig   `json:"latest_finalized_block_time"`
	LatestFinalizedBlockEpoch  common.JSONUint64 `json:"latest_finalized_block_epoch"`
	CurrentEpoch               common.JSONUint64 `json:"current_epoch"`
	CurrentHeight              common.JSONUint64 `json:"current_height"`
	CurrentTime                *common.JSONBig   `json:"current_time"`
	Syncing                    bool              `json:"syncing"`
	GenesisBlockHash           common.Hash       `json:"genesis_block_hash"`
	SnapshotBlockHeight        common.JSONUint64 `json:"snapshot_block_height"`
	SnapshotBlockHash          common.Hash       `json:"snapshot_block_hash"`
}

func HandleThetaRPCResponse(rpcRes *rpcc.RPCResponse, rpcErr error, parse func(jsonBytes []byte) (interface{}, error)) (result interface{}, err error) {
	if rpcErr != nil {
		return nil, fmt.Errorf("failed to get theta RPC response: %v", rpcErr)
	}
	if rpcRes.Error != nil {
		return nil, fmt.Errorf("theta RPC returns an error: %v", rpcRes.Error)
	}

	var jsonBytes []byte
	jsonBytes, err = json.MarshalIndent(rpcRes.Result, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to parse theta RPC response: %v, %s", err, string(jsonBytes))
	}

	//logger.Infof("HandleThetaRPCResponse, jsonBytes: %v", string(jsonBytes))
	result, err = parse(jsonBytes)
	if err != nil {
		logger.Warn("Failed to parse theta RPC response: " + err.Error())
	}
	return
}

type EthGetBlockResult struct {
	Height    hexutil.Uint64 `json:"number"`
	Hash      common.Hash    `json:"hash"`
	Parent    common.Hash    `json:"parentHash"`
	Timestamp hexutil.Uint64 `json:"timestamp"`
	Proposer  common.Address `json:"miner"`
	TxHash    common.Hash    `json:"transactionsRoot"`
	StateHash common.Hash    `json:"stateRoot"`

	ReiceptHash     common.Hash    `json:"receiptsRoot"`
	Nonce           string         `json:"nonce"`
	Sha3Uncles      common.Hash    `json:"sha3Uncles"`
	LogsBloom       string         `json:"logsBloom"`
	Difficulty      hexutil.Uint64 `json:"difficulty"`
	TotalDifficulty hexutil.Uint64 `json:"totalDifficulty"`
	Size            hexutil.Uint64 `json:"size"`
	GasLimit        hexutil.Uint64 `json:"gasLimit"`
	GasUsed         hexutil.Uint64 `json:"gasUsed"`
	ExtraData       string         `json:"extraData"`
	Uncles          []common.Hash  `json:"uncles"`
	Transactions    []interface{}  `json:"transactions"`
}
type ThetaGetBlockResult struct {
	*ThetaGetBlockResultInner
}
type ThetaGetBlocksResult []*ThetaGetBlockResultInner

type ThetaGetBlockResultInner struct {
	ChainID            string                   `json:"chain_id"`
	Epoch              common.JSONUint64        `json:"epoch"`
	Height             common.JSONUint64        `json:"height"`
	Parent             common.Hash              `json:"parent"`
	TxHash             common.Hash              `json:"transactions_hash"`
	StateHash          common.Hash              `json:"state_hash"`
	Timestamp          *common.JSONBig          `json:"timestamp"`
	Proposer           common.Address           `json:"proposer"`
	HCC                core.CommitCertificate   `json:"hcc"`
	GuardianVotes      *core.AggregatedVotes    `json:"guardian_votes"`
	EliteEdgeNodeVotes *core.AggregatedEENVotes `json:"elite_edge_node_votes"`

	Children []common.Hash    `json:"children"`
	Status   core.BlockStatus `json:"status"`

	Hash common.Hash `json:"hash"`
	Txs  []rpc.Tx    `json:"transactions"`
}
type LogData struct {
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`
}
type RPCResult struct {
	Result  []LogData `json:"logs"`
	Address string    `json:"contractAddress"`
}
type GetBlockByHeightArgs struct {
	Height             common.JSONUint64 `json:"height"`
	IncludeEthTxHashes bool              `json:"include_eth_tx_hashes"`
}
type Tx struct {
	types.Tx       `json:"raw"`
	Type           byte                              `json:"type"`
	Hash           common.Hash                       `json:"hash"`
	Receipt        *blockchain.TxReceiptEntry        `json:"receipt"`
	BalanceChanges *blockchain.TxBalanceChangesEntry `json:"balance_changes"`
}
type TxTmp struct {
	Tx   json.RawMessage `json:"raw"`
	Type byte            `json:"type"`
	Hash tcommon.Hash    `json:"hash"`
}

type GetAccountArgs struct {
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Height  common.JSONUint64 `json:"height"`
	Preview bool              `json:"preview"` // preview the account balance from the ScreenedView
}
type GetAccountResult struct {
	*types.Account
	Address string `json:"address"`
}

const RawABI = `
[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "string",
				"name": "denom",
				"type": "string"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "targetChainVoucherReceiver",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "voucherContact",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "mintedAmount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "sourceChainTokenLockNonce",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "voucherMintNonce",
				"type": "uint256"
			}
		],
		"name": "TNT20VoucherMinted",
		"type": "event"
	}
]`

func Resolve(data []byte) big.Int {
	contractAbi, err := abi.JSON(strings.NewReader(RawABI))
	if err != nil {
		fmt.Println(err)
	}
	type TransferEvt struct {
		Denom                      string
		TargetChainVoucherReceiver common.Address
		VoucherContact             common.Address
		MintedAmount               *big.Int
		SourceChainTokenLockNonce  *big.Int
		VoucherMintNonce           *big.Int
	}
	var event TransferEvt
	h := data
	err = contractAbi.UnpackIntoInterface(&event, "TNT20VoucherMinted", h)
	if err != nil {
		fmt.Println(err)
	}
	return *event.VoucherMintNonce
}
