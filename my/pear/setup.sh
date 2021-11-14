
export KEY="mykey"
export CHAINID="pearmint_9000-6"
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
export DENOM=apear
#echo $GENESIS
#echo $TMPGENESIS
#echo $ETHCONFIG
#echo $APPCONFIG
#echo 'DENOM='$DENOM


export GRPCPORT0=9080
export GRPCPORT1=9081
export ETHPORT0=8535
export ETHPORT1=8536
export COSMOSPORT0=26646
export COSMOSPORT1=26647
export CLI=cronosd
export NODE=tcp://127.0.0.1:26647
$CLI keys list --keyring-backend $KEYRING --home $CHAINHOME
#echo 'HOME='$CHAINHOME

export TOTALAMOUNT=100000000000000000000000000000000$DENOM
export MYAMOUNT=1000100000000000000001$DENOM

export K1=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[0].address' |  tr -d '"')
export S1=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[1].address' |  tr -d '"')
export S2=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[2].address' |  tr -d '"')

