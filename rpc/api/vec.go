// Copyright 2015 The go-vector Authors
// This file is part of the go-vector library.
//
// The go-vector library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-vector library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-vector library. If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"bytes"
	"encoding/json"
	"math/big"

	"fmt"

	"github.com/vector/go-vector/common"
	"github.com/vector/go-vector/common/natspec"
	"github.com/vector/go-vector/vec"
	"github.com/vector/go-vector/rlp"
	"github.com/vector/go-vector/rpc/codec"
	"github.com/vector/go-vector/rpc/shared"
	"github.com/vector/go-vector/xvec"
	"gopkg.in/fatih/set.v0"
)

const (
	EthApiVersion = "1.0"
)

// vec api provider
// See https://github.com/vector/wiki/wiki/JSON-RPC
type vecApi struct {
	xvec     *xvec.XEth
	vector *vec.Vector
	methods  map[string]vechandler
	codec    codec.ApiCoder
}

// vec callback handler
type vechandler func(*vecApi, *shared.Request) (interface{}, error)

var (
	vecMapping = map[string]vechandler{
		"eth_accounts":                            (*vecApi).Accounts,
		"eth_blockNumber":                         (*vecApi).BlockNumber,
		"eth_getBalance":                          (*vecApi).GetBalance,
		"eth_protocolVersion":                     (*vecApi).ProtocolVersion,
		"eth_coinbase":                            (*vecApi).Coinbase,
		"eth_mining":                              (*vecApi).IsMining,
		"eth_syncing":                             (*vecApi).IsSyncing,
		"eth_gasPrice":                            (*vecApi).GasPrice,
		"eth_getStorage":                          (*vecApi).GetStorage,
		"eth_storageAt":                           (*vecApi).GetStorage,
		"eth_getStorageAt":                        (*vecApi).GetStorageAt,
		"eth_getTransactionCount":                 (*vecApi).GetTransactionCount,
		"eth_getBlockTransactionCountByHash":      (*vecApi).GetBlockTransactionCountByHash,
		"eth_getBlockTransactionCountByNumber":    (*vecApi).GetBlockTransactionCountByNumber,
		"eth_getUncleCountByBlockHash":            (*vecApi).GetUncleCountByBlockHash,
		"eth_getUncleCountByBlockNumber":          (*vecApi).GetUncleCountByBlockNumber,
		"eth_getData":                             (*vecApi).GetData,
		"eth_getCode":                             (*vecApi).GetData,
		"eth_getNatSpec":                          (*vecApi).GetNatSpec,
		"eth_sign":                                (*vecApi).Sign,
		"eth_sendRawTransaction":                  (*vecApi).SubmitTransaction,
		"eth_submitTransaction":                   (*vecApi).SubmitTransaction,
		"eth_sendTransaction":                     (*vecApi).SendTransaction,
		"eth_signTransaction":                     (*vecApi).SignTransaction,
		"eth_transact":                            (*vecApi).SendTransaction,
		"eth_estimateGas":                         (*vecApi).EstimateGas,
		"eth_call":                                (*vecApi).Call,
		"eth_flush":                               (*vecApi).Flush,
		"eth_getBlockByHash":                      (*vecApi).GetBlockByHash,
		"eth_getBlockByNumber":                    (*vecApi).GetBlockByNumber,
		"eth_getTransactionByHash":                (*vecApi).GetTransactionByHash,
		"eth_getTransactionByBlockNumberAndIndex": (*vecApi).GetTransactionByBlockNumberAndIndex,
		"eth_getTransactionByBlockHashAndIndex":   (*vecApi).GetTransactionByBlockHashAndIndex,
		"eth_getUncleByBlockHashAndIndex":         (*vecApi).GetUncleByBlockHashAndIndex,
		"eth_getUncleByBlockNumberAndIndex":       (*vecApi).GetUncleByBlockNumberAndIndex,
		"eth_getCompilers":                        (*vecApi).GetCompilers,
		"eth_compileSolidity":                     (*vecApi).CompileSolidity,
		"eth_newFilter":                           (*vecApi).NewFilter,
		"eth_newBlockFilter":                      (*vecApi).NewBlockFilter,
		"eth_newPendingTransactionFilter":         (*vecApi).NewPendingTransactionFilter,
		"eth_uninstallFilter":                     (*vecApi).UninstallFilter,
		"eth_getFilterChanges":                    (*vecApi).GetFilterChanges,
		"eth_getFilterLogs":                       (*vecApi).GetFilterLogs,
		"eth_getLogs":                             (*vecApi).GetLogs,
		"eth_hashrate":                            (*vecApi).Hashrate,
		"eth_getWork":                             (*vecApi).GetWork,
		"eth_submitWork":                          (*vecApi).SubmitWork,
		"eth_submitHashrate":                      (*vecApi).SubmitHashrate,
		"eth_resend":                              (*vecApi).Resend,
		"eth_pendingTransactions":                 (*vecApi).PendingTransactions,
		"eth_getTransactionReceipt":               (*vecApi).GetTransactionReceipt,
	}
)

// create new vecApi instance
func NewEthApi(xvec *xvec.XEth, vec *vec.Vector, codec codec.Codec) *vecApi {
	return &vecApi{xvec, vec, vecMapping, codec.New(nil)}
}

// collection with supported methods
func (self *vecApi) Methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *vecApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.Method]; ok {
		return callback(self, req)
	}

	return nil, shared.NewNotImplementedError(req.Method)
}

func (self *vecApi) Name() string {
	return shared.EthApiName
}

func (self *vecApi) ApiVersion() string {
	return EthApiVersion
}

func (self *vecApi) Accounts(req *shared.Request) (interface{}, error) {
	return self.xvec.Accounts(), nil
}

func (self *vecApi) Hashrate(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xvec.HashRate()), nil
}

func (self *vecApi) BlockNumber(req *shared.Request) (interface{}, error) {
	num := self.xvec.CurrentBlock().Number()
	return newHexNum(num.Bytes()), nil
}

func (self *vecApi) GetBalance(req *shared.Request) (interface{}, error) {
	args := new(GetBalanceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xvec.AtStateNum(args.BlockNumber).BalanceAt(args.Address), nil
}

func (self *vecApi) ProtocolVersion(req *shared.Request) (interface{}, error) {
	return self.xvec.EthVersion(), nil
}

func (self *vecApi) Coinbase(req *shared.Request) (interface{}, error) {
	return newHexData(self.xvec.Coinbase()), nil
}

func (self *vecApi) IsMining(req *shared.Request) (interface{}, error) {
	return self.xvec.IsMining(), nil
}

func (self *vecApi) IsSyncing(req *shared.Request) (interface{}, error) {
	origin, current, height := self.vector.Downloader().Progress()
	if current < height {
		return map[string]interface{}{
			"startingBlock": newHexNum(big.NewInt(int64(origin)).Bytes()),
			"currentBlock":  newHexNum(big.NewInt(int64(current)).Bytes()),
			"highestBlock":  newHexNum(big.NewInt(int64(height)).Bytes()),
		}, nil
	}
	return false, nil
}

func (self *vecApi) GasPrice(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xvec.DefaultGasPrice().Bytes()), nil
}

func (self *vecApi) GetStorage(req *shared.Request) (interface{}, error) {
	args := new(GetStorageArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xvec.AtStateNum(args.BlockNumber).State().SafeGet(args.Address).Storage(), nil
}

func (self *vecApi) GetStorageAt(req *shared.Request) (interface{}, error) {
	args := new(GetStorageAtArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xvec.AtStateNum(args.BlockNumber).StorageAt(args.Address, args.Key), nil
}

func (self *vecApi) GetTransactionCount(req *shared.Request) (interface{}, error) {
	args := new(GetTxCountArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	count := self.xvec.AtStateNum(args.BlockNumber).TxCountAt(args.Address)
	return fmt.Sprintf("%#x", count), nil
}

func (self *vecApi) GetBlockTransactionCountByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xvec.EthBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *vecApi) GetBlockTransactionCountByNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xvec.EthBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *vecApi) GetUncleCountByBlockHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xvec.EthBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *vecApi) GetUncleCountByBlockNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xvec.EthBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *vecApi) GetData(req *shared.Request) (interface{}, error) {
	args := new(GetDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v := self.xvec.AtStateNum(args.BlockNumber).CodeAtBytes(args.Address)
	return newHexData(v), nil
}

func (self *vecApi) Sign(req *shared.Request) (interface{}, error) {
	args := new(NewSigArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v, err := self.xvec.Sign(args.From, args.Data, false)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *vecApi) SubmitTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	v, err := self.xvec.PushTx(args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// JsonTransaction is returned as response by the JSON RPC. It contains the
// signed RLP encoded transaction as Raw and the signed transaction object as Tx.
type JsonTransaction struct {
	Raw string `json:"raw"`
	Tx  *tx    `json:"tx"`
}

func (self *vecApi) SignTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	tx, err := self.xvec.SignTransaction(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return JsonTransaction{"0x" + common.Bytes2Hex(data), newTx(tx)}, nil
}

func (self *vecApi) SendTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	v, err := self.xvec.Transact(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *vecApi) GetNatSpec(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	var jsontx = fmt.Sprintf(`{"params":[{"to":"%s","data": "%s"}]}`, args.To, args.Data)
	notice := natspec.GetNotice(self.xvec, jsontx, self.vector.HTTPClient())

	return notice, nil
}

func (self *vecApi) EstimateGas(req *shared.Request) (interface{}, error) {
	_, gas, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if len(gas) == 0 {
		return newHexNum(0), nil
	} else {
		return newHexNum(common.String2Big(gas)), err
	}
}

func (self *vecApi) Call(req *shared.Request) (interface{}, error) {
	v, _, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if v == "0x0" {
		return newHexData([]byte{}), nil
	} else {
		return newHexData(common.FromHex(v)), nil
	}
}

func (self *vecApi) Flush(req *shared.Request) (interface{}, error) {
	return nil, shared.NewNotImplementedError(req.Method)
}

func (self *vecApi) doCall(params json.RawMessage) (string, string, error) {
	args := new(CallArgs)
	if err := self.codec.Decode(params, &args); err != nil {
		return "", "", err
	}

	return self.xvec.AtStateNum(args.BlockNumber).Call(args.From, args.To, args.Value.String(), args.Gas.String(), args.GasPrice.String(), args.Data)
}

func (self *vecApi) GetBlockByHash(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByHashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xvec.EthBlockByHash(args.BlockHash)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xvec.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *vecApi) GetBlockByNumber(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByNumberArgs)
	if err := json.Unmarshal(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xvec.EthBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xvec.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *vecApi) GetTransactionByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	tx, bhash, bnum, txi := self.xvec.EthTransactionByHash(args.Hash)
	if tx != nil {
		v := NewTransactionRes(tx)
		// if the blockhash is 0, assume this is a pending transaction
		if bytes.Compare(bhash.Bytes(), bytes.Repeat([]byte{0}, 32)) != 0 {
			v.BlockHash = newHexData(bhash)
			v.BlockNumber = newHexNum(bnum)
			v.TxIndex = newHexNum(txi)
		}
		return v, nil
	}
	return nil, nil
}

func (self *vecApi) GetTransactionByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xvec.EthBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xvec.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Transactions[args.Index], nil
	}
}

func (self *vecApi) GetTransactionByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xvec.EthBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xvec.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Transactions[args.Index], nil
}

func (self *vecApi) GetUncleByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xvec.EthBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xvec.Td(raw.Hash()), false)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Uncles[args.Index], nil
}

func (self *vecApi) GetUncleByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xvec.EthBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xvec.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Uncles[args.Index], nil
	}
}

func (self *vecApi) GetCompilers(req *shared.Request) (interface{}, error) {
	var lang string
	if solc, _ := self.xvec.Solc(); solc != nil {
		lang = "Solidity"
	}
	c := []string{lang}
	return c, nil
}

func (self *vecApi) CompileSolidity(req *shared.Request) (interface{}, error) {
	solc, _ := self.xvec.Solc()
	if solc == nil {
		return nil, shared.NewNotAvailableError(req.Method, "solc (solidity compiler) not found")
	}

	args := new(SourceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	contracts, err := solc.Compile(args.Source)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (self *vecApi) NewFilter(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	id := self.xvec.NewLogFilter(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)
	return newHexNum(big.NewInt(int64(id)).Bytes()), nil
}

func (self *vecApi) NewBlockFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xvec.NewBlockFilter()), nil
}

func (self *vecApi) NewPendingTransactionFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xvec.NewTransactionFilter()), nil
}

func (self *vecApi) UninstallFilter(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xvec.UninstallFilter(args.Id), nil
}

func (self *vecApi) GetFilterChanges(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	switch self.xvec.GetFilterType(args.Id) {
	case xvec.BlockFilterTy:
		return NewHashesRes(self.xvec.BlockFilterChanged(args.Id)), nil
	case xvec.TransactionFilterTy:
		return NewHashesRes(self.xvec.TransactionFilterChanged(args.Id)), nil
	case xvec.LogFilterTy:
		return NewLogsRes(self.xvec.LogFilterChanged(args.Id)), nil
	default:
		return []string{}, nil // reply empty string slice
	}
}

func (self *vecApi) GetFilterLogs(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return NewLogsRes(self.xvec.Logs(args.Id)), nil
}

func (self *vecApi) GetLogs(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return NewLogsRes(self.xvec.AllLogs(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)), nil
}

func (self *vecApi) GetWork(req *shared.Request) (interface{}, error) {
	self.xvec.SetMining(true, 0)
	ret, err := self.xvec.RemoteMining().GetWork()
	if err != nil {
		return nil, shared.NewNotReadyError("mining work")
	} else {
		return ret, nil
	}
}

func (self *vecApi) SubmitWork(req *shared.Request) (interface{}, error) {
	args := new(SubmitWorkArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xvec.RemoteMining().SubmitWork(args.Nonce, common.HexToHash(args.Digest), common.HexToHash(args.Header)), nil
}

func (self *vecApi) SubmitHashrate(req *shared.Request) (interface{}, error) {
	args := new(SubmitHashRateArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return false, shared.NewDecodeParamError(err.Error())
	}
	self.xvec.RemoteMining().SubmitHashrate(common.HexToHash(args.Id), args.Rate)
	return true, nil
}

func (self *vecApi) Resend(req *shared.Request) (interface{}, error) {
	args := new(ResendArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	from := common.HexToAddress(args.Tx.From)

	pending := self.vector.TxPool().GetTransactions()
	for _, p := range pending {
		if pFrom, err := p.FromFrontier(); err == nil && pFrom == from && p.SigHash() == args.Tx.tx.SigHash() {
			self.vector.TxPool().RemoveTx(common.HexToHash(args.Tx.Hash))
			return self.xvec.Transact(args.Tx.From, args.Tx.To, args.Tx.Nonce, args.Tx.Value, args.GasLimit, args.GasPrice, args.Tx.Data)
		}
	}

	return nil, fmt.Errorf("Transaction %s not found", args.Tx.Hash)
}

func (self *vecApi) PendingTransactions(req *shared.Request) (interface{}, error) {
	txs := self.vector.TxPool().GetTransactions()

	// grab the accounts from the account manager. This will help with determining which
	// transactions should be returned.
	accounts, err := self.vector.AccountManager().Accounts()
	if err != nil {
		return nil, err
	}

	// Add the accouns to a new set
	accountSet := set.New()
	for _, account := range accounts {
		accountSet.Add(account.Address)
	}

	var ltxs []*tx
	for _, tx := range txs {
		if from, _ := tx.FromFrontier(); accountSet.Has(from) {
			ltxs = append(ltxs, newTx(tx))
		}
	}

	return ltxs, nil
}

func (self *vecApi) GetTransactionReceipt(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	txhash := common.BytesToHash(common.FromHex(args.Hash))
	tx, bhash, bnum, txi := self.xvec.EthTransactionByHash(args.Hash)
	rec := self.xvec.GetTxReceipt(txhash)
	// We could have an error of "not found". Should disambiguate
	// if err != nil {
	// 	return err, nil
	// }
	if rec != nil && tx != nil {
		v := NewReceiptRes(rec)
		v.BlockHash = newHexData(bhash)
		v.BlockNumber = newHexNum(bnum)
		v.TransactionIndex = newHexNum(txi)
		return v, nil
	}

	return nil, nil
}
