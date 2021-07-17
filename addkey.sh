KEY="testkey"
KEY2="testkey2"
CHAINID="ethermint-2"
MONIKER="localtestnet"
KEYRING=test

echo $MYMNEMONICS
ethermintd keys add $KEY --keyring-backend $KEYRING --algo "eth_secp256k1" --recover --index 0
echo $MYMNEMONICS
ethermintd keys add $KEY2 --keyring-backend $KEYRING --algo "eth_secp256k1" --recover --index 1

