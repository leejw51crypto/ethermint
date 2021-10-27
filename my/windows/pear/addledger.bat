set KEY="ledgerkey"
set KEY2="ledgerkey2"
set CHAINID="pearmint-6"
set MONIKER="localtestnet"
set KEYRING=test

ethermintd keys add %KEY% --keyring-backend %KEYRING% --algo "eth_secp256k1" --index 0 --ledger

