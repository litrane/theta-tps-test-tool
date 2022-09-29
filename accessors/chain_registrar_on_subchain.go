// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accessors

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/thetatoken/thetasubchain/eth"
	"github.com/thetatoken/thetasubchain/eth/abi"
	"github.com/thetatoken/thetasubchain/eth/abi/bind"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/thetasubchain/eth/core/types"
	"github.com/thetatoken/thetasubchain/eth/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ChainRegistrarOnSubchainMetaData contains all meta data concerning the ChainRegistrarOnSubchain contract.
var ChainRegistrarOnSubchainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numBlocksPerDynasty_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"crossChainFee_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeSetter_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"register\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"IP\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"ChannelRegistered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"channelRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"register\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"status\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"IP\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"channelStatusVotingRecords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"dynasty\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accumlatedSharesForValid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accumlatedSharesForInvalid\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeSetter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numBlocksPerDynasty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDynasty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCrossChainFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"isARegisteredSubchain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"getSubchainRegistrationHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumBlocksPerDynasty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subchainID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dynasty\",\"type\":\"uint256\"}],\"name\":\"getValidatorSet\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"validators\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"shareAmounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newCrossChainFee\",\"type\":\"uint256\"}],\"name\":\"updateCrossChainFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeSetter\",\"type\":\"address\"}],\"name\":\"updateFeeSetter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"IP\",\"type\":\"string\"}],\"name\":\"registerSubchainChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"name\":\"deregisterSubchainChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"targetChainID\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"eventNonce\",\"type\":\"uint256\"}],\"name\":\"updateSubchainChannelStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"name\":\"isAnActiveChannel\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxProcessedNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405260018055600060025534801561001957600080fd5b5060405161158d38038061158d8339810160408190526100389161006a565b6001600055600592909255600655600780546001600160a01b0319166001600160a01b039092169190911790556100b0565b60008060006060848603121561007f57600080fd5b83516020850151604086015191945092506001600160a01b03811681146100a557600080fd5b809150509250925092565b6114ce806100bf6000396000f3fe608060405234801561001057600080fd5b50600436106100f65760003560e01c80637adfce8a116100925780637adfce8a146101de57806387cf3ef4146101f15780639886ddbc1461021c5780639bbb690a1461022f578063a7464b1214610237578063b73774891461023f578063dba9de6b14610265578063e902844c1461026d578063e9b69eea146102c257600080fd5b806309314dc3146100fb578063164d29f614610112578063188eea9b1461011b5780632f2c13b51461013d578063385482371461016857806343b71f051461017d57806343f27e45146101a157806360f8e1bb146101c257806367016090146101d5575b600080fd5b6002545b6040519081526020015b60405180910390f35b6100ff60065481565b61012e610129366004611049565b6102d5565b604051610109939291906111c1565b61015361014b366004611049565b506000908190565b60408051928352901515602083015201610109565b61017b610176366004611049565b61038a565b005b61019161018b366004611049565b50600190565b6040519015158152602001610109565b6101b46101af3660046110b8565b61045c565b604051610109929190611229565b61017b6101d0366004610f52565b610785565b6100ff60055481565b61017b6101ec36600461107b565b6107d1565b600754610204906001600160a01b031681565b6040516001600160a01b039091168152602001610109565b61017b61022a366004611049565b610bde565b6006546100ff565b6005546100ff565b61019161024d366004611049565b60009081526003602052604090206001908101541490565b610153610c0d565b6102a761027b3660046110b8565b600460209081526000928352604080842090915290825290208054600282015460039092015490919083565b60408051938452602084019290925290820152606001610109565b61017b6102d03660046110da565b610d18565b6003602052600090815260409020805460018201546002830180546001600160a01b03909316939192610307906113eb565b80601f0160208091040260200160405190810160405280929190818152602001828054610333906113eb565b80156103805780601f1061035557610100808354040283529160200191610380565b820191906000526020600020905b81548152906001019060200180831161036357829003601f168201915b5050505050905083565b600260005414156103b65760405162461bcd60e51b81526004016103ad906112f3565b60405180910390fd5b60026000908155818152600360205260409020546001600160a01b031633146104215760405162461bcd60e51b815260206004820152601b60248201527f796f7520646f206e6f74206f776e2074686973206368616e6e656c000000000060448201526064016103ad565b600081815260036020526040812080546001600160a01b031916815560018101829055906104526002830182610e7c565b5050600160005550565b604080518181526060818101835291829160009160208201818036833701905050905084846000805b60208110156104ea578381602081106104a0576104a0611457565b1a60f81b85836104af81611426565b9450815181106104c1576104c1611457565b60200101906001600160f81b031916908160001a905350806104e281611426565b915050610485565b5060005b60208110156105535782816020811061050957610509611457565b1a60f81b858361051881611426565b94508151811061052a5761052a611457565b60200101906001600160f81b031916908160001a9053508061054b81611426565b9150506104ee565b5060008060b56001600160a01b03168660405161057091906111a5565b600060405180830381855afa9150503d80600081146105ab576040519150601f19603f3d011682016040523d82523d6000602084013e6105b0565b606091505b5091509150816106285760405162461bcd60e51b815260206004820152603a60248201527f73746174696363616c6c20746f207468652076616c69646174726f536574207060448201527f7265636f6d70696c656420636f6e7472616374206661696c656400000000000060648201526084016103ad565b60008180602001905181019061063e9190610f6f565b80519091508067ffffffffffffffff81111561065c5761065c61146d565b604051908082528060200260200182016040528015610685578160200160208202803683370190505b5099508067ffffffffffffffff8111156106a1576106a161146d565b6040519080825280602002602001820160405280156106ca578160200160208202803683370190505b50985060005b81811015610775578281815181106106ea576106ea611457565b6020026020010151600001518b828151811061070857610708611457565b60200260200101906001600160a01b031690816001600160a01b03168152505082818151811061073a5761073a611457565b6020026020010151602001518a828151811061075857610758611457565b60209081029190910101528061076d81611426565b9150506106d0565b5050505050505050509250929050565b6007546001600160a01b031633146107af5760405162461bcd60e51b81526004016103ad906112ad565b600780546001600160a01b0319166001600160a01b0392909216919091179055565b600260005414156107f45760405162461bcd60e51b81526004016103ad906112f3565b6002600081905554610807906001611384565b811461084d5760405162461bcd60e51b8152602060048201526015602482015274696e636f7272656374206576656e74206e6f6e636560581b60448201526064016103ad565b600080610858610c0d565b91509150806108a15760405162461bcd60e51b81526020600482015260156024820152746661696c656420746f206765742064796e6173747960581b60448201526064016103ad565b604080516020810187905290810183905284151560f81b60608201523390600090819060610160408051601f19818403018152918152815160209283012060008b815260048452828120828252909352912090915060608046610904818a61045c565b9350915060005b8251811015610a9e57876001600160a01b031683828151811061093057610930611457565b60200260200101516001600160a01b03161461094b57610a8c565b6001965060005b60018601548110156109ec5785600101818154811061097357610973611457565b6000918252602090912001546001600160a01b038a8116911614156109da5760405162461bcd60e51b815260206004820152601c60248201527f546869732076616c696461746f7220616c726561647920766f7465640000000060448201526064016103ad565b806109e481611426565b915050610952565b508985556001808601805491820181556000908152602090200180546001600160a01b031916331790558b15610a5657610a4c848281518110610a3157610a31611457565b60200260200101518660020154610e5d90919063ffffffff16565b6002860155610a8c565b610a86848281518110610a6b57610a6b611457565b60200260200101518660030154610e5d90919063ffffffff16565b60038601555b80610a9681611426565b91505061090b565b50505083610ae05760405162461bcd60e51b815260206004820152600f60248201526e2737ba1030903b30b634b230ba37b960891b60448201526064016103ad565b6000805b8251811015610b2d57610b19838281518110610b0257610b02611457565b602002602001015183610e5d90919063ffffffff16565b915080610b2581611426565b915050610ae4565b50610b39816002610e70565b6002840154610b49906003610e70565b10610b7d5760008b815260036020526040812060019081018190556002805491929091610b77908490611384565b90915550505b610b88816002610e70565b600384810154610b9791610e70565b10610bcc5760008b81526003602052604081206000196001918201556002805491929091610bc6908490611384565b90915550505b50506001600055505050505050505050565b6007546001600160a01b03163314610c085760405162461bcd60e51b81526004016103ad906112ad565b600655565b60408051600080825260208201909252819060008060b46001600160a01b031683604051610c3b91906111a5565b600060405180830381855afa9150503d8060008114610c76576040519150601f19603f3d011682016040523d82523d6000602084013e610c7b565b606091505b509150915081610cf5576040805162461bcd60e51b81526020600482015260248101919091527f6661696c656420746f2063616c6c20707265636f6d70696c656420636f6e747260448201527f61637420746f207175657279207468652063757272656e742064796e6173747960648201526084016103ad565b600081806020019051810190610d0b9190611062565b9660019650945050505050565b60026000541415610d3b5760405162461bcd60e51b81526004016103ad906112f3565b60026000908155828152600360205260409020546001600160a01b031615610d985760405162461bcd60e51b815260206004820152601060248201526f63616e277420757064617465206e6f7760801b60448201526064016103ad565b60408051606081018252338152600060208083018281528385018681528784526003835294909220835181546001600160a01b0319166001600160a01b03909116178155915160018301559251805192939192610dfb9260028501920190610eb9565b50506001546040517f1015a61fb37283e6254a85ce40ee20dc84496f3aa755f9844aa85f94938d56dc9250610e35913391869186916111f1565b60405180910390a16001806000828254610e4f9190611384565b909155505060016000555050565b6000610e698284611384565b9392505050565b6000610e69828461139c565b508054610e88906113eb565b6000825580601f10610e98575050565b601f016020900490600052602060002090810190610eb69190610f3d565b50565b828054610ec5906113eb565b90600052602060002090601f016020900481019282610ee75760008555610f2d565b82601f10610f0057805160ff1916838001178555610f2d565b82800160010185558215610f2d579182015b82811115610f2d578251825591602001919060010190610f12565b50610f39929150610f3d565b5090565b5b80821115610f395760008155600101610f3e565b600060208284031215610f6457600080fd5b8135610e6981611483565b60006020808385031215610f8257600080fd5b825167ffffffffffffffff80821115610f9a57600080fd5b818501915085601f830112610fae57600080fd5b815181811115610fc057610fc061146d565b610fce848260051b01611353565b8181528481019250838501600683901b85018601891015610fee57600080fd5b60009450845b8381101561103b57604080838c03121561100c578687fd5b61101461132a565b835161101f81611483565b8152838901518982015286529487019490910190600101610ff4565b509098975050505050505050565b60006020828403121561105b57600080fd5b5035919050565b60006020828403121561107457600080fd5b5051919050565b60008060006060848603121561109057600080fd5b83359250602084013580151581146110a757600080fd5b929592945050506040919091013590565b600080604083850312156110cb57600080fd5b50508035926020909101359150565b600080604083850312156110ed57600080fd5b8235915060208084013567ffffffffffffffff8082111561110d57600080fd5b818601915086601f83011261112157600080fd5b8135818111156111335761113361146d565b611145601f8201601f19168501611353565b9150808252878482850101111561115b57600080fd5b80848401858401376000848284010152508093505050509250929050565b600081518084526111918160208601602086016113bb565b601f01601f19169290920160200192915050565b600082516111b78184602087016113bb565b9190910192915050565b60018060a01b03841681528260208201526060604082015260006111e86060830184611179565b95945050505050565b60018060a01b03851681528360208201526080604082015260006112186080830185611179565b905082606083015295945050505050565b604080825283519082018190526000906020906060840190828701845b8281101561126b5781516001600160a01b031684529284019290840190600101611246565b5050508381038285015284518082528583019183019060005b818110156112a057835183529284019291840191600101611284565b5090979650505050505050565b60208082526026908201527f4f6e6c792074686520666565207365747465722063616e206d616b65207468696040820152651cc818d85b1b60d21b606082015260800190565b6020808252601f908201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00604082015260600190565b6040805190810167ffffffffffffffff8111828210171561134d5761134d61146d565b60405290565b604051601f8201601f1916810167ffffffffffffffff8111828210171561137c5761137c61146d565b604052919050565b6000821982111561139757611397611441565b500190565b60008160001904831182151516156113b6576113b6611441565b500290565b60005b838110156113d65781810151838201526020016113be565b838111156113e5576000848401525b50505050565b600181811c908216806113ff57607f821691505b6020821081141561142057634e487b7160e01b600052602260045260246000fd5b50919050565b600060001982141561143a5761143a611441565b5060010190565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114610eb657600080fdfea2646970667358221220f87eb0d8c66d690507c57b00899c7df870fd4df6b1bda9683c9547238d6cd8e564736f6c63430008070033",
}

// ChainRegistrarOnSubchainABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainRegistrarOnSubchainMetaData.ABI instead.
var ChainRegistrarOnSubchainABI = ChainRegistrarOnSubchainMetaData.ABI

// ChainRegistrarOnSubchainBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ChainRegistrarOnSubchainMetaData.Bin instead.
var ChainRegistrarOnSubchainBin = ChainRegistrarOnSubchainMetaData.Bin

// DeployChainRegistrarOnSubchain deploys a new Ethereum contract, binding an instance of ChainRegistrarOnSubchain to it.
func DeployChainRegistrarOnSubchain(auth *bind.TransactOpts, backend bind.ContractBackend, numBlocksPerDynasty_ *big.Int, crossChainFee_ *big.Int, feeSetter_ common.Address) (common.Address, *types.Transaction, *ChainRegistrarOnSubchain, error) {
	parsed, err := ChainRegistrarOnSubchainMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainRegistrarOnSubchainBin), backend, numBlocksPerDynasty_, crossChainFee_, feeSetter_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainRegistrarOnSubchain{ChainRegistrarOnSubchainCaller: ChainRegistrarOnSubchainCaller{contract: contract}, ChainRegistrarOnSubchainTransactor: ChainRegistrarOnSubchainTransactor{contract: contract}, ChainRegistrarOnSubchainFilterer: ChainRegistrarOnSubchainFilterer{contract: contract}}, nil
}

// ChainRegistrarOnSubchain is an auto generated Go binding around an Ethereum contract.
type ChainRegistrarOnSubchain struct {
	ChainRegistrarOnSubchainCaller     // Read-only binding to the contract
	ChainRegistrarOnSubchainTransactor // Write-only binding to the contract
	ChainRegistrarOnSubchainFilterer   // Log filterer for contract events
}

// ChainRegistrarOnSubchainCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainRegistrarOnSubchainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainRegistrarOnSubchainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainRegistrarOnSubchainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainRegistrarOnSubchainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainRegistrarOnSubchainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainRegistrarOnSubchainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainRegistrarOnSubchainSession struct {
	Contract     *ChainRegistrarOnSubchain // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ChainRegistrarOnSubchainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainRegistrarOnSubchainCallerSession struct {
	Contract *ChainRegistrarOnSubchainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// ChainRegistrarOnSubchainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainRegistrarOnSubchainTransactorSession struct {
	Contract     *ChainRegistrarOnSubchainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ChainRegistrarOnSubchainRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainRegistrarOnSubchainRaw struct {
	Contract *ChainRegistrarOnSubchain // Generic contract binding to access the raw methods on
}

// ChainRegistrarOnSubchainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainRegistrarOnSubchainCallerRaw struct {
	Contract *ChainRegistrarOnSubchainCaller // Generic read-only contract binding to access the raw methods on
}

// ChainRegistrarOnSubchainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainRegistrarOnSubchainTransactorRaw struct {
	Contract *ChainRegistrarOnSubchainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChainRegistrarOnSubchain creates a new instance of ChainRegistrarOnSubchain, bound to a specific deployed contract.
func NewChainRegistrarOnSubchain(address common.Address, backend bind.ContractBackend) (*ChainRegistrarOnSubchain, error) {
	contract, err := bindChainRegistrarOnSubchain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainRegistrarOnSubchain{ChainRegistrarOnSubchainCaller: ChainRegistrarOnSubchainCaller{contract: contract}, ChainRegistrarOnSubchainTransactor: ChainRegistrarOnSubchainTransactor{contract: contract}, ChainRegistrarOnSubchainFilterer: ChainRegistrarOnSubchainFilterer{contract: contract}}, nil
}

// NewChainRegistrarOnSubchainCaller creates a new read-only instance of ChainRegistrarOnSubchain, bound to a specific deployed contract.
func NewChainRegistrarOnSubchainCaller(address common.Address, caller bind.ContractCaller) (*ChainRegistrarOnSubchainCaller, error) {
	contract, err := bindChainRegistrarOnSubchain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainRegistrarOnSubchainCaller{contract: contract}, nil
}

// NewChainRegistrarOnSubchainTransactor creates a new write-only instance of ChainRegistrarOnSubchain, bound to a specific deployed contract.
func NewChainRegistrarOnSubchainTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainRegistrarOnSubchainTransactor, error) {
	contract, err := bindChainRegistrarOnSubchain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainRegistrarOnSubchainTransactor{contract: contract}, nil
}

// NewChainRegistrarOnSubchainFilterer creates a new log filterer instance of ChainRegistrarOnSubchain, bound to a specific deployed contract.
func NewChainRegistrarOnSubchainFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainRegistrarOnSubchainFilterer, error) {
	contract, err := bindChainRegistrarOnSubchain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainRegistrarOnSubchainFilterer{contract: contract}, nil
}

// bindChainRegistrarOnSubchain binds a generic wrapper to an already deployed contract.
func bindChainRegistrarOnSubchain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChainRegistrarOnSubchainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainRegistrarOnSubchain.Contract.ChainRegistrarOnSubchainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.ChainRegistrarOnSubchainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.ChainRegistrarOnSubchainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainRegistrarOnSubchain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.contract.Transact(opts, method, params...)
}

// ChannelRegistry is a free data retrieval call binding the contract method 0x188eea9b.
//
// Solidity: function channelRegistry(uint256 ) view returns(address register, int256 status, string IP)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) ChannelRegistry(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Register common.Address
	Status   *big.Int
	IP       string
}, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "channelRegistry", arg0)

	outstruct := new(struct {
		Register common.Address
		Status   *big.Int
		IP       string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Register = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Status = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IP = *abi.ConvertType(out[2], new(string)).(*string)

	return *outstruct, err

}

// ChannelRegistry is a free data retrieval call binding the contract method 0x188eea9b.
//
// Solidity: function channelRegistry(uint256 ) view returns(address register, int256 status, string IP)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) ChannelRegistry(arg0 *big.Int) (struct {
	Register common.Address
	Status   *big.Int
	IP       string
}, error) {
	return _ChainRegistrarOnSubchain.Contract.ChannelRegistry(&_ChainRegistrarOnSubchain.CallOpts, arg0)
}

// ChannelRegistry is a free data retrieval call binding the contract method 0x188eea9b.
//
// Solidity: function channelRegistry(uint256 ) view returns(address register, int256 status, string IP)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) ChannelRegistry(arg0 *big.Int) (struct {
	Register common.Address
	Status   *big.Int
	IP       string
}, error) {
	return _ChainRegistrarOnSubchain.Contract.ChannelRegistry(&_ChainRegistrarOnSubchain.CallOpts, arg0)
}

// ChannelStatusVotingRecords is a free data retrieval call binding the contract method 0xe902844c.
//
// Solidity: function channelStatusVotingRecords(uint256 , bytes32 ) view returns(uint256 dynasty, uint256 accumlatedSharesForValid, uint256 accumlatedSharesForInvalid)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) ChannelStatusVotingRecords(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte) (struct {
	Dynasty                    *big.Int
	AccumlatedSharesForValid   *big.Int
	AccumlatedSharesForInvalid *big.Int
}, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "channelStatusVotingRecords", arg0, arg1)

	outstruct := new(struct {
		Dynasty                    *big.Int
		AccumlatedSharesForValid   *big.Int
		AccumlatedSharesForInvalid *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Dynasty = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.AccumlatedSharesForValid = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.AccumlatedSharesForInvalid = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ChannelStatusVotingRecords is a free data retrieval call binding the contract method 0xe902844c.
//
// Solidity: function channelStatusVotingRecords(uint256 , bytes32 ) view returns(uint256 dynasty, uint256 accumlatedSharesForValid, uint256 accumlatedSharesForInvalid)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) ChannelStatusVotingRecords(arg0 *big.Int, arg1 [32]byte) (struct {
	Dynasty                    *big.Int
	AccumlatedSharesForValid   *big.Int
	AccumlatedSharesForInvalid *big.Int
}, error) {
	return _ChainRegistrarOnSubchain.Contract.ChannelStatusVotingRecords(&_ChainRegistrarOnSubchain.CallOpts, arg0, arg1)
}

// ChannelStatusVotingRecords is a free data retrieval call binding the contract method 0xe902844c.
//
// Solidity: function channelStatusVotingRecords(uint256 , bytes32 ) view returns(uint256 dynasty, uint256 accumlatedSharesForValid, uint256 accumlatedSharesForInvalid)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) ChannelStatusVotingRecords(arg0 *big.Int, arg1 [32]byte) (struct {
	Dynasty                    *big.Int
	AccumlatedSharesForValid   *big.Int
	AccumlatedSharesForInvalid *big.Int
}, error) {
	return _ChainRegistrarOnSubchain.Contract.ChannelStatusVotingRecords(&_ChainRegistrarOnSubchain.CallOpts, arg0, arg1)
}

// CrossChainFee is a free data retrieval call binding the contract method 0x164d29f6.
//
// Solidity: function crossChainFee() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) CrossChainFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "crossChainFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CrossChainFee is a free data retrieval call binding the contract method 0x164d29f6.
//
// Solidity: function crossChainFee() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) CrossChainFee() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.CrossChainFee(&_ChainRegistrarOnSubchain.CallOpts)
}

// CrossChainFee is a free data retrieval call binding the contract method 0x164d29f6.
//
// Solidity: function crossChainFee() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) CrossChainFee() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.CrossChainFee(&_ChainRegistrarOnSubchain.CallOpts)
}

// FeeSetter is a free data retrieval call binding the contract method 0x87cf3ef4.
//
// Solidity: function feeSetter() view returns(address)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) FeeSetter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "feeSetter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeSetter is a free data retrieval call binding the contract method 0x87cf3ef4.
//
// Solidity: function feeSetter() view returns(address)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) FeeSetter() (common.Address, error) {
	return _ChainRegistrarOnSubchain.Contract.FeeSetter(&_ChainRegistrarOnSubchain.CallOpts)
}

// FeeSetter is a free data retrieval call binding the contract method 0x87cf3ef4.
//
// Solidity: function feeSetter() view returns(address)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) FeeSetter() (common.Address, error) {
	return _ChainRegistrarOnSubchain.Contract.FeeSetter(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetCrossChainFee is a free data retrieval call binding the contract method 0x9bbb690a.
//
// Solidity: function getCrossChainFee() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) GetCrossChainFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "getCrossChainFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCrossChainFee is a free data retrieval call binding the contract method 0x9bbb690a.
//
// Solidity: function getCrossChainFee() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) GetCrossChainFee() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.GetCrossChainFee(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetCrossChainFee is a free data retrieval call binding the contract method 0x9bbb690a.
//
// Solidity: function getCrossChainFee() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) GetCrossChainFee() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.GetCrossChainFee(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetDynasty is a free data retrieval call binding the contract method 0xdba9de6b.
//
// Solidity: function getDynasty() view returns(uint256, bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) GetDynasty(opts *bind.CallOpts) (*big.Int, bool, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "getDynasty")

	if err != nil {
		return *new(*big.Int), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(bool)).(*bool)

	return out0, out1, err

}

// GetDynasty is a free data retrieval call binding the contract method 0xdba9de6b.
//
// Solidity: function getDynasty() view returns(uint256, bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) GetDynasty() (*big.Int, bool, error) {
	return _ChainRegistrarOnSubchain.Contract.GetDynasty(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetDynasty is a free data retrieval call binding the contract method 0xdba9de6b.
//
// Solidity: function getDynasty() view returns(uint256, bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) GetDynasty() (*big.Int, bool, error) {
	return _ChainRegistrarOnSubchain.Contract.GetDynasty(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetMaxProcessedNonce is a free data retrieval call binding the contract method 0x09314dc3.
//
// Solidity: function getMaxProcessedNonce() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) GetMaxProcessedNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "getMaxProcessedNonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxProcessedNonce is a free data retrieval call binding the contract method 0x09314dc3.
//
// Solidity: function getMaxProcessedNonce() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) GetMaxProcessedNonce() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.GetMaxProcessedNonce(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetMaxProcessedNonce is a free data retrieval call binding the contract method 0x09314dc3.
//
// Solidity: function getMaxProcessedNonce() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) GetMaxProcessedNonce() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.GetMaxProcessedNonce(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetNumBlocksPerDynasty is a free data retrieval call binding the contract method 0xa7464b12.
//
// Solidity: function getNumBlocksPerDynasty() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) GetNumBlocksPerDynasty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "getNumBlocksPerDynasty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumBlocksPerDynasty is a free data retrieval call binding the contract method 0xa7464b12.
//
// Solidity: function getNumBlocksPerDynasty() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) GetNumBlocksPerDynasty() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.GetNumBlocksPerDynasty(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetNumBlocksPerDynasty is a free data retrieval call binding the contract method 0xa7464b12.
//
// Solidity: function getNumBlocksPerDynasty() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) GetNumBlocksPerDynasty() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.GetNumBlocksPerDynasty(&_ChainRegistrarOnSubchain.CallOpts)
}

// GetSubchainRegistrationHeight is a free data retrieval call binding the contract method 0x2f2c13b5.
//
// Solidity: function getSubchainRegistrationHeight(uint256 ) pure returns(uint256, bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) GetSubchainRegistrationHeight(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, bool, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "getSubchainRegistrationHeight", arg0)

	if err != nil {
		return *new(*big.Int), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(bool)).(*bool)

	return out0, out1, err

}

// GetSubchainRegistrationHeight is a free data retrieval call binding the contract method 0x2f2c13b5.
//
// Solidity: function getSubchainRegistrationHeight(uint256 ) pure returns(uint256, bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) GetSubchainRegistrationHeight(arg0 *big.Int) (*big.Int, bool, error) {
	return _ChainRegistrarOnSubchain.Contract.GetSubchainRegistrationHeight(&_ChainRegistrarOnSubchain.CallOpts, arg0)
}

// GetSubchainRegistrationHeight is a free data retrieval call binding the contract method 0x2f2c13b5.
//
// Solidity: function getSubchainRegistrationHeight(uint256 ) pure returns(uint256, bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) GetSubchainRegistrationHeight(arg0 *big.Int) (*big.Int, bool, error) {
	return _ChainRegistrarOnSubchain.Contract.GetSubchainRegistrationHeight(&_ChainRegistrarOnSubchain.CallOpts, arg0)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0x43f27e45.
//
// Solidity: function getValidatorSet(uint256 subchainID, uint256 dynasty) view returns(address[] validators, uint256[] shareAmounts)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) GetValidatorSet(opts *bind.CallOpts, subchainID *big.Int, dynasty *big.Int) (struct {
	Validators   []common.Address
	ShareAmounts []*big.Int
}, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "getValidatorSet", subchainID, dynasty)

	outstruct := new(struct {
		Validators   []common.Address
		ShareAmounts []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validators = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.ShareAmounts = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetValidatorSet is a free data retrieval call binding the contract method 0x43f27e45.
//
// Solidity: function getValidatorSet(uint256 subchainID, uint256 dynasty) view returns(address[] validators, uint256[] shareAmounts)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) GetValidatorSet(subchainID *big.Int, dynasty *big.Int) (struct {
	Validators   []common.Address
	ShareAmounts []*big.Int
}, error) {
	return _ChainRegistrarOnSubchain.Contract.GetValidatorSet(&_ChainRegistrarOnSubchain.CallOpts, subchainID, dynasty)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0x43f27e45.
//
// Solidity: function getValidatorSet(uint256 subchainID, uint256 dynasty) view returns(address[] validators, uint256[] shareAmounts)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) GetValidatorSet(subchainID *big.Int, dynasty *big.Int) (struct {
	Validators   []common.Address
	ShareAmounts []*big.Int
}, error) {
	return _ChainRegistrarOnSubchain.Contract.GetValidatorSet(&_ChainRegistrarOnSubchain.CallOpts, subchainID, dynasty)
}

// IsARegisteredSubchain is a free data retrieval call binding the contract method 0x43b71f05.
//
// Solidity: function isARegisteredSubchain(uint256 ) pure returns(bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) IsARegisteredSubchain(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "isARegisteredSubchain", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsARegisteredSubchain is a free data retrieval call binding the contract method 0x43b71f05.
//
// Solidity: function isARegisteredSubchain(uint256 ) pure returns(bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) IsARegisteredSubchain(arg0 *big.Int) (bool, error) {
	return _ChainRegistrarOnSubchain.Contract.IsARegisteredSubchain(&_ChainRegistrarOnSubchain.CallOpts, arg0)
}

// IsARegisteredSubchain is a free data retrieval call binding the contract method 0x43b71f05.
//
// Solidity: function isARegisteredSubchain(uint256 ) pure returns(bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) IsARegisteredSubchain(arg0 *big.Int) (bool, error) {
	return _ChainRegistrarOnSubchain.Contract.IsARegisteredSubchain(&_ChainRegistrarOnSubchain.CallOpts, arg0)
}

// IsAnActiveChannel is a free data retrieval call binding the contract method 0xb7377489.
//
// Solidity: function isAnActiveChannel(uint256 chainID) view returns(bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) IsAnActiveChannel(opts *bind.CallOpts, chainID *big.Int) (bool, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "isAnActiveChannel", chainID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAnActiveChannel is a free data retrieval call binding the contract method 0xb7377489.
//
// Solidity: function isAnActiveChannel(uint256 chainID) view returns(bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) IsAnActiveChannel(chainID *big.Int) (bool, error) {
	return _ChainRegistrarOnSubchain.Contract.IsAnActiveChannel(&_ChainRegistrarOnSubchain.CallOpts, chainID)
}

// IsAnActiveChannel is a free data retrieval call binding the contract method 0xb7377489.
//
// Solidity: function isAnActiveChannel(uint256 chainID) view returns(bool)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) IsAnActiveChannel(chainID *big.Int) (bool, error) {
	return _ChainRegistrarOnSubchain.Contract.IsAnActiveChannel(&_ChainRegistrarOnSubchain.CallOpts, chainID)
}

// NumBlocksPerDynasty is a free data retrieval call binding the contract method 0x67016090.
//
// Solidity: function numBlocksPerDynasty() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCaller) NumBlocksPerDynasty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ChainRegistrarOnSubchain.contract.Call(opts, &out, "numBlocksPerDynasty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumBlocksPerDynasty is a free data retrieval call binding the contract method 0x67016090.
//
// Solidity: function numBlocksPerDynasty() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) NumBlocksPerDynasty() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.NumBlocksPerDynasty(&_ChainRegistrarOnSubchain.CallOpts)
}

// NumBlocksPerDynasty is a free data retrieval call binding the contract method 0x67016090.
//
// Solidity: function numBlocksPerDynasty() view returns(uint256)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainCallerSession) NumBlocksPerDynasty() (*big.Int, error) {
	return _ChainRegistrarOnSubchain.Contract.NumBlocksPerDynasty(&_ChainRegistrarOnSubchain.CallOpts)
}

// DeregisterSubchainChannel is a paid mutator transaction binding the contract method 0x38548237.
//
// Solidity: function deregisterSubchainChannel(uint256 chainID) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactor) DeregisterSubchainChannel(opts *bind.TransactOpts, chainID *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.contract.Transact(opts, "deregisterSubchainChannel", chainID)
}

// DeregisterSubchainChannel is a paid mutator transaction binding the contract method 0x38548237.
//
// Solidity: function deregisterSubchainChannel(uint256 chainID) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) DeregisterSubchainChannel(chainID *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.DeregisterSubchainChannel(&_ChainRegistrarOnSubchain.TransactOpts, chainID)
}

// DeregisterSubchainChannel is a paid mutator transaction binding the contract method 0x38548237.
//
// Solidity: function deregisterSubchainChannel(uint256 chainID) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactorSession) DeregisterSubchainChannel(chainID *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.DeregisterSubchainChannel(&_ChainRegistrarOnSubchain.TransactOpts, chainID)
}

// RegisterSubchainChannel is a paid mutator transaction binding the contract method 0xe9b69eea.
//
// Solidity: function registerSubchainChannel(uint256 chainID, string IP) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactor) RegisterSubchainChannel(opts *bind.TransactOpts, chainID *big.Int, IP string) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.contract.Transact(opts, "registerSubchainChannel", chainID, IP)
}

// RegisterSubchainChannel is a paid mutator transaction binding the contract method 0xe9b69eea.
//
// Solidity: function registerSubchainChannel(uint256 chainID, string IP) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) RegisterSubchainChannel(chainID *big.Int, IP string) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.RegisterSubchainChannel(&_ChainRegistrarOnSubchain.TransactOpts, chainID, IP)
}

// RegisterSubchainChannel is a paid mutator transaction binding the contract method 0xe9b69eea.
//
// Solidity: function registerSubchainChannel(uint256 chainID, string IP) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactorSession) RegisterSubchainChannel(chainID *big.Int, IP string) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.RegisterSubchainChannel(&_ChainRegistrarOnSubchain.TransactOpts, chainID, IP)
}

// UpdateCrossChainFee is a paid mutator transaction binding the contract method 0x9886ddbc.
//
// Solidity: function updateCrossChainFee(uint256 newCrossChainFee) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactor) UpdateCrossChainFee(opts *bind.TransactOpts, newCrossChainFee *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.contract.Transact(opts, "updateCrossChainFee", newCrossChainFee)
}

// UpdateCrossChainFee is a paid mutator transaction binding the contract method 0x9886ddbc.
//
// Solidity: function updateCrossChainFee(uint256 newCrossChainFee) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) UpdateCrossChainFee(newCrossChainFee *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.UpdateCrossChainFee(&_ChainRegistrarOnSubchain.TransactOpts, newCrossChainFee)
}

// UpdateCrossChainFee is a paid mutator transaction binding the contract method 0x9886ddbc.
//
// Solidity: function updateCrossChainFee(uint256 newCrossChainFee) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactorSession) UpdateCrossChainFee(newCrossChainFee *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.UpdateCrossChainFee(&_ChainRegistrarOnSubchain.TransactOpts, newCrossChainFee)
}

// UpdateFeeSetter is a paid mutator transaction binding the contract method 0x60f8e1bb.
//
// Solidity: function updateFeeSetter(address newFeeSetter) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactor) UpdateFeeSetter(opts *bind.TransactOpts, newFeeSetter common.Address) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.contract.Transact(opts, "updateFeeSetter", newFeeSetter)
}

// UpdateFeeSetter is a paid mutator transaction binding the contract method 0x60f8e1bb.
//
// Solidity: function updateFeeSetter(address newFeeSetter) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) UpdateFeeSetter(newFeeSetter common.Address) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.UpdateFeeSetter(&_ChainRegistrarOnSubchain.TransactOpts, newFeeSetter)
}

// UpdateFeeSetter is a paid mutator transaction binding the contract method 0x60f8e1bb.
//
// Solidity: function updateFeeSetter(address newFeeSetter) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactorSession) UpdateFeeSetter(newFeeSetter common.Address) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.UpdateFeeSetter(&_ChainRegistrarOnSubchain.TransactOpts, newFeeSetter)
}

// UpdateSubchainChannelStatus is a paid mutator transaction binding the contract method 0x7adfce8a.
//
// Solidity: function updateSubchainChannelStatus(uint256 targetChainID, bool isValid, uint256 eventNonce) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactor) UpdateSubchainChannelStatus(opts *bind.TransactOpts, targetChainID *big.Int, isValid bool, eventNonce *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.contract.Transact(opts, "updateSubchainChannelStatus", targetChainID, isValid, eventNonce)
}

// UpdateSubchainChannelStatus is a paid mutator transaction binding the contract method 0x7adfce8a.
//
// Solidity: function updateSubchainChannelStatus(uint256 targetChainID, bool isValid, uint256 eventNonce) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainSession) UpdateSubchainChannelStatus(targetChainID *big.Int, isValid bool, eventNonce *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.UpdateSubchainChannelStatus(&_ChainRegistrarOnSubchain.TransactOpts, targetChainID, isValid, eventNonce)
}

// UpdateSubchainChannelStatus is a paid mutator transaction binding the contract method 0x7adfce8a.
//
// Solidity: function updateSubchainChannelStatus(uint256 targetChainID, bool isValid, uint256 eventNonce) returns()
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainTransactorSession) UpdateSubchainChannelStatus(targetChainID *big.Int, isValid bool, eventNonce *big.Int) (*types.Transaction, error) {
	return _ChainRegistrarOnSubchain.Contract.UpdateSubchainChannelStatus(&_ChainRegistrarOnSubchain.TransactOpts, targetChainID, isValid, eventNonce)
}

// ChainRegistrarOnSubchainChannelRegisteredIterator is returned from FilterChannelRegistered and is used to iterate over the raw logs and unpacked data for ChannelRegistered events raised by the ChainRegistrarOnSubchain contract.
type ChainRegistrarOnSubchainChannelRegisteredIterator struct {
	Event *ChainRegistrarOnSubchainChannelRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChainRegistrarOnSubchainChannelRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainRegistrarOnSubchainChannelRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChainRegistrarOnSubchainChannelRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChainRegistrarOnSubchainChannelRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainRegistrarOnSubchainChannelRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainRegistrarOnSubchainChannelRegistered represents a ChannelRegistered event raised by the ChainRegistrarOnSubchain contract.
type ChainRegistrarOnSubchainChannelRegistered struct {
	Register common.Address
	ChainID  *big.Int
	IP       string
	Nonce    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterChannelRegistered is a free log retrieval operation binding the contract event 0x1015a61fb37283e6254a85ce40ee20dc84496f3aa755f9844aa85f94938d56dc.
//
// Solidity: event ChannelRegistered(address register, uint256 chainID, string IP, uint256 nonce)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainFilterer) FilterChannelRegistered(opts *bind.FilterOpts) (*ChainRegistrarOnSubchainChannelRegisteredIterator, error) {

	logs, sub, err := _ChainRegistrarOnSubchain.contract.FilterLogs(opts, "ChannelRegistered")
	if err != nil {
		return nil, err
	}
	return &ChainRegistrarOnSubchainChannelRegisteredIterator{contract: _ChainRegistrarOnSubchain.contract, event: "ChannelRegistered", logs: logs, sub: sub}, nil
}

// WatchChannelRegistered is a free log subscription operation binding the contract event 0x1015a61fb37283e6254a85ce40ee20dc84496f3aa755f9844aa85f94938d56dc.
//
// Solidity: event ChannelRegistered(address register, uint256 chainID, string IP, uint256 nonce)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainFilterer) WatchChannelRegistered(opts *bind.WatchOpts, sink chan<- *ChainRegistrarOnSubchainChannelRegistered) (event.Subscription, error) {

	logs, sub, err := _ChainRegistrarOnSubchain.contract.WatchLogs(opts, "ChannelRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainRegistrarOnSubchainChannelRegistered)
				if err := _ChainRegistrarOnSubchain.contract.UnpackLog(event, "ChannelRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChannelRegistered is a log parse operation binding the contract event 0x1015a61fb37283e6254a85ce40ee20dc84496f3aa755f9844aa85f94938d56dc.
//
// Solidity: event ChannelRegistered(address register, uint256 chainID, string IP, uint256 nonce)
func (_ChainRegistrarOnSubchain *ChainRegistrarOnSubchainFilterer) ParseChannelRegistered(log types.Log) (*ChainRegistrarOnSubchainChannelRegistered, error) {
	event := new(ChainRegistrarOnSubchainChannelRegistered)
	if err := _ChainRegistrarOnSubchain.contract.UnpackLog(event, "ChannelRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
