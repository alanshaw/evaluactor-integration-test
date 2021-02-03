module github.com/alanshaw/evaluactor-integration-test

go 1.15

require (
	github.com/alanshaw/evaluactor v0.0.0-20210203112549-d531cac2439e
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-state-types v0.0.0-20210119062722-4adba5aaea71
	github.com/filecoin-project/lotus v1.4.1
	github.com/ipfs/go-cid v0.0.7
	github.com/ipld/go-car v0.1.1-0.20201119040415-11b6074b6d4d
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
