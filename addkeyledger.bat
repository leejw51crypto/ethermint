set KEY="testkey"
set KEY2="testkey2"
set CHAINID="ethermint-2"
set MONIKER="localtestnet"
set KEYRING=test

echo %MYMNEMONICS%
ethermintd keys add %KEY% --keyring-backend %KEYRING% --algo "eth_secp256k1" --index 0 --ledger

