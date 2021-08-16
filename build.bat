set KEY="mykey"
set CHAINID="ethermint-2"
set MONIKER="localtestnet"
set KEYRING="test"
set KEYALGO="eth_secp256k1"
set LOGLEVEL="debug"
set TRACE="--trace"

del ethermintd.exe
@echo build binary
go build  -tags cgo .\cmd\ethermintd 
