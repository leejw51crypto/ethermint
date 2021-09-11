
export KEY="mykey"
export CHAINID="chainmain-1"
export MONIKER="localtestnet"
export KEYRING="test"
export KEYALGO="eth_secp256k1"
export LOGLEVEL="info"
# to trace evm
export TRACE="--trace"
#export TRACE=""
#export CHAINHOME=$HOME/.$CHAINID
export CHAINHOME=/run/user/1000/pytest-of-leejw51/pytest-9/chainmain0/chainmain-1/node0

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
export CLI=chain-maind
export NODE=tcp://127.0.0.1:26807
$CLI keys list --keyring-backend $KEYRING --home $CHAINHOME
#echo 'HOME='$CHAINHOME

export TOTALAMOUNT=100000000000000000000000000000000$DENOM
export MYAMOUNT=1000100000000000000001$DENOM

export K1=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[0].address' |  tr -d '"')
#export S1=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[1].address' |  tr -d '"')
#export S2=$($CLI keys list --keyring-backend $KEYRING --home $CHAINHOME | yq eval  -o json | jq '.[2].address' |  tr -d '"')


export S1=cro1u08u5dvtnpmlpdq333uj9tcj75yceggszxpnsy
export S2=cro1apdh4yc2lnpephevc6lmpvkyv6s5cjh652n6e4
