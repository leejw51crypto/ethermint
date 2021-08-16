set KEY="testkey"
set KEY2="testkey2"
set CHAINID="pearmint-6"
set MONIKER="localtestnet"
set KEYRING=test
set HOME=%USERPROFILE%\.pear
echo %MYMNEMONICS%
ethermintd keys add %KEY% --keyring-backend %KEYRING% --algo "eth_secp256k1" --recover --index 0 --home %HOME%
echo %MYMNEMONICS%
ethermintd keys add %KEY2% --keyring-backend %KEYRING% --algo "eth_secp256k1" --recover --index 1 --home %HOME%

