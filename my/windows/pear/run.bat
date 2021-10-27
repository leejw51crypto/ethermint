
rem ethermint compile on windows
rem install golang , gcc, sed for windows
rem 1. install msys2 : https://www.msys2.org/
rem 2. pacman -S mingw-w64-x86_64-toolchain
rem    pacman -S sed
rem    pacman -S mingw-w64-x86_64-jq
rem 3. add path C:\msys64\mingw64\bin  
rem             C:\msys64\usr\bin

set KEY="mykey"
set CHAINID="pearmint-6"
set MONIKER="localtestnet"
set KEYRING="test"
set KEYALGO="eth_secp256k1"
set LOGLEVEL="fatal"
set TRACE="--trace"
set HOME=%USERPROFILE%\.pear

del ethermintd.exe
@echo build binary
rem go build  -tags cgo,ledger   -tags cgo,ledger --ldflags "-extldflags \"-Wl,--allow-multiple-definition\"" .\cmd\ethermintd
rem go build  -tags cgo  -tags cgo .\cmd\ethermintd

ethermintd start --pruning=nothing %TRACE% --log_level %LOGLEVEL% --minimum-gas-prices=0.0001pphoton  --json-rpc.address 0.0.0.0:19535 --home %HOME%
