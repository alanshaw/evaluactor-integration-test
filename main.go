package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/alanshaw/evaluactor"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/vm"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	"github.com/filecoin-project/lotus/lib/blockstore"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
)

var rootCID cid.Cid
var zeroState string
var luaSrc string

func init() {
	rootCID, _ = cid.Decode("bafy2bzacebyvb3cairrciekiwmyzccp72dy46r4d5jua5q233jy4a36mg6c2c")
	zeroState = "H4sIAAAAAAAA/0pflFqUn19S3HRDK0KdgbFwwRMmhcKANw4uSSqCHpsNJ3L+v/DxvHvzqwy+w9G3Cg/knzFvXUiC0vSy1KLizPw8xu1sxGlocmnx4WbsXJho2NjkxJDYckPLn4ExlIErLTNH31A/tSwxB8l21w3nDymcexNdI5N05IqkrcWT10v6OkSL5jyPn8O14cnBAwwOUHPYWm5oRUiDDBKDGFSWWpSZlpmaUpSanllcUlSJZCjHzaQtxRJCRV+mHdESdfsreGr1M8v9e/fZnjo/wXrr6jaEoSwtN7TCQWYKQcwsLskvSkxPLcgvTy1CMo9bq21K7dld+YcWrlN4tb4v5eAW+YdTBY+nXH5/3nbVmW+7weY1NTkxMKN5NrkoPw/JnNSLt66nTDpz9YKAZO2h7Ku/2Q+cbz1zZ/0ngZYDV6r/ikcyODQ5MSS33NAKApnBCzEjMTk5vzSvBMmY+64lxzYv3Jl1Ys+GKVum10zW7K3IufD+16ZLz5dcly6fwQX3HkPLDa1AkFE8UO9VFpek5pIT+qyg0JcAGSWMElK5iUXZqchui/qxsdhU9+n67k1W0qxq5i+mCZe9+fHkRlpO9uYFHL2FuXATGdGCKjMvE9mcpWx2s3/P2Fe53PXlZoljW0XeiZ6UPrZX+Gr7Db3m58cqcuDmpBAMrjOqx0pW5IU+sU31fsXEsLVW1t88WkTIyHEZX2z50lVvmRk8GJh765Yca2CAGsmEFmxFqeWJRSlIJkqohcl9ClnGJ/jp5CfeR6LXGF6xid0siSp59vORcEH42WQGXwa2sphK+4UMlRcYGBgkUFM/gqVOXCw0GOBO9dhFGxiJywjNTgwByP76l3UmcfHibB5zGy+rvmXejdqfTs0OagiYMKfiCkdoFw8JSjWJU9bk0KCBPfthiuxlJC43TnCAgqYICQamOzkMWEA4A/MBbOIMDCT4kUGBBMXfJFCLBQTLiLjCoRFUwLCwgjIjpxxGekfhaBFXTjQ6MSTLoJcOyGwtHPkdi9A6JuLyfzdSkD1N7baVOMv28PYE0S0b3oQeVF7kG/3wlcoNS4ujOyyr/fXJU0rFRExKWiBBrYKDg4Mecf5qZmBodmRoaJBALSgRrFTiistmJOd5XXv+vV4xbR9PQXP1in97Njy6OmsVs9ncT4fv3351vfvBeonUguSMxMw83bLEnMyUxJLM/DwP4nQ2OTUwNIIL0VDGzBfGnj4TZbP+6T+0VYzzifnOnhQqkWJLXJHciF2/DHqZjMyuJ65s7nRwYPBkEIDmdi8Gxu/ey60KWa5MboqQAnEWz6q8t8EVvVQIY5SclrVdF13Yk4H/PwQwOAAAAAD//wEAAP//JmyWUIcJAAA="
	src, err := ioutil.ReadFile("./examples/max.lua")
	if err != nil {
		panic(err)
	}
	luaSrc = string(src)
}

func main() {
	// Load the CAR into a new temporary Blockstore.
	bs := loadBlockstore()

	vmOpts := &vm.VMOpts{
		StateBase:      rootCID,
		Epoch:          abi.ChainEpoch(1),
		Bstore:         bs,
		Syscalls:       vm.Syscalls(ffiwrapper.ProofVerifier),
		CircSupplyCalc: nil,
		Rand:           nil,
		BaseFee:        abi.NewTokenAmount(100),
		NtwkVersion: func(context.Context, abi.ChainEpoch) network.Version {
			return network.VersionMax
		},
	}

	lvm, err := vm.NewVM(context.TODO(), vmOpts)
	if err != nil {
		panic(err)
	}

	invoker := vm.NewActorRegistry()
	invoker.Register(nil, evaluactor.Actor{})
	lvm.SetInvoker(invoker)

	params := evaluactor.EvalParams{Script: luaSrc}
	var pbuf bytes.Buffer
	err = params.MarshalCBOR(&pbuf)
	if err != nil {
		panic(err)
	}

	from, err := address.NewIDAddress(100)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nLua script:\n%s\n", luaSrc)

	msg := types.Message{
		Version:    0,
		To:         evaluactor.Address,
		From:       from,
		Nonce:      0,
		Value:      abi.NewTokenAmount(0),
		GasLimit:   3000000,
		GasFeeCap:  abi.NewTokenAmount(1000),
		GasPremium: abi.NewTokenAmount(100),
		Method:     evaluactor.MethodEval,
		Params:     pbuf.Bytes(),
	}

	fmt.Printf("\nExecuting message:\n	%+v\n", &msg)

	ret, err := lvm.ApplyMessage(context.Background(), &msg)
	if err != nil {
		panic(err)
	}

	er := evaluactor.EvalReturn{}
	err = er.UnmarshalCBOR(bytes.NewBuffer(ret.Return))
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nReceipt:\n	%+v\n", ret)
	fmt.Printf("\nReturned value:\n	%+v\n", &er)
}

func loadBlockstore() blockstore.Blockstore {
	gzcar, err := base64.StdEncoding.DecodeString(zeroState)
	if err != nil {
		panic(err)
	}

	r, err := gzip.NewReader(bytes.NewBuffer(gzcar))
	if err != nil {
		panic(err)
	}
	defer r.Close() // nolint

	bs := blockstore.Blockstore(blockstore.NewTemporary())
	_, err = car.LoadCar(bs, r) // load car into blackstore
	if err != nil {
		panic(err)
	}

	return bs
}
