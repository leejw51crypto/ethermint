set KEY="ledgerkey"
set KEY2="ledgerkey2"
set CHAINID="ethermint-2"
set MONIKER="localtestnet"
set KEYRING=test

echo %MYMNEMONICS%
ethermintd keys add %KEY% --keyring-backend %KEYRING% --algo "eth_secp256k1" --index 0 --ledger
rem echo %MYMNEMONICS%
rem ethermintd keys add %KEY2% --keyring-backend %KEYRING% --algo "eth_secp256k1" --recover --index 1

