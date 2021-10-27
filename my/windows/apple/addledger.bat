set KEY="ledgerkey"
set KEY2="ledgerkey2"
set CHAINID="ethermint-2"
set MONIKER="localtestnet"
set KEYRING=test
set HOME=%USERPROFILE%\.appled
ethermintd keys add %KEY% --keyring-backend %KEYRING% --algo "eth_secp256k1" --index 0 --ledger --home %HOME%

