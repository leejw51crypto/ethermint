
export KEY="mykey"
export CHAINID="applemint_9000-2"
export MONIKER="localtestnet"
export KEYRING="test"
export KEYALGO="eth_secp256k1"
export LOGLEVEL="info"
# to trace evm
export TRACE="--trace"
#export TRACE=""
export CHAINHOME=$HOME/.$CHAINID

export ETHCONFIG=$CHAINHOME/config/config.toml
export GENESIS=$CHAINHOME/config/genesis.json
export TMPGENESIS=$CHAINHOME/config/tmp_genesis.json
export APPCONFIG=$CHAINHOME/config/app.toml
export CLIENTCONFIG=$CHAINHOME/config/client.toml
export DENOM=aapple
#echo $GENESIS
#echo $TMPGENESIS
#echo $ETHCONFIG
#echo $APPCONFIG
#echo 'DENOM='$DENOM


export GRPCPORT0=9090
export GRPCPORT1=9091
export ETHPORT0=8545
export ETHPORT1=8546
export COSMOSPORT0=26656
export COSMOSPORT1=26657
export CLI=ethermintd
export NODE=tcp://127.0.0.1:26657
$CLI keys list --keyring-backend $KEYRING --home $CHAINHOME
#echo 'HOME='$CHAINHOME

export TOTALAMOUNT=100000000000000000000000000000000$DENOM
export MYAMOUNT=1000100000000000000001$DENOM

export K1=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[0].address' |  tr -d '"')
export S1=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[1].address' |  tr -d '"')
export S2=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[2].address' |  tr -d '"')

