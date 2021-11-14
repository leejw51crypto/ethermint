. ./setup.sh

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon and client
rm -rf $CHAINHOME


$CLI config keyring-backend $KEYRING --home $CHAINHOME
$CLI config chain-id $CHAINID --home $CHAINHOME

# if $KEY exists it should be deleted
echo $MYMNEMONICS | $CLI keys add "$KEY" --keyring-backend $KEYRING --algo "eth_secp256k1" --recover --index 5 --home $CHAINHOME
echo $MYMNEMONICS | $CLI keys add "$KEY"1 --keyring-backend $KEYRING --algo "eth_secp256k1" --recover --index 0 --home $CHAINHOME
echo $MYMNEMONICS | $CLI keys add "$KEY"2 --keyring-backend $KEYRING --algo "eth_secp256k1" --recover --index 1 --home $CHAINHOME
echo $MYMNEMONICS | $CLI keys add "$KEY"3 --keyring-backend $KEYRING --algo "eth_secp256k1" --recover --index 2 --home $CHAINHOME

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
$CLI init $MONIKER --chain-id $CHAINID --home $CHAINHOME

# Change parameter token denominations to $DENOM
cat $GENESIS | jq '.app_state["staking"]["params"]["bond_denom"]="'$DENOM'"' > $TMPGENESIS && mv $TMPGENESIS $GENESIS
cat $GENESIS | jq '.app_state["crisis"]["constant_fee"]["denom"]="'$DENOM'"' > $TMPGENESIS && mv $TMPGENESIS $GENESIS
cat $GENESIS | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="'$DENOM'"' > $TMPGENESIS && mv $TMPGENESIS $GENESIS
cat $GENESIS | jq '.app_state["mint"]["params"]["mint_denom"]="'$DENOM'"' > $TMPGENESIS && mv $TMPGENESIS $GENESIS

# increase block time (?)
cat $GENESIS | jq '.consensus_params["block"]["time_iota_ms"]="1000"' > $TMPGENESIS && mv $TMPGENESIS $GENESIS

# Set gas limit in genesis
cat $GENESIS | jq '.consensus_params["block"]["max_gas"]="10000000"' > $TMPGENESIS && mv $TMPGENESIS $GENESIS

# change port
sed -i ""  "s/create_empty_blocks = true/create_empty_blocks = false/g" $ETHCONFIG
sed -i ""  "s/26657/$COSMOSPORT1/g" $CLIENTCONFIG
sed -i ""  "s/26657/$COSMOSPORT1/g" $ETHCONFIG
sed -i ""  "s/26656/$COSMOSPORT0/g" $ETHCONFIG
sed -i ""  "s/9090/$GRPCPORT0/g" $APPCONFIG
sed -i ""  "s/9091/$GRPCPORT1/g" $APPCONFIG
sed -i ""  "s/8545/$ETHPORT0/g" $APPCONFIG
sed -i ""  "s/8546/$ETHPORT1/g" $APPCONFIG
sed -i ""  "s/aphoton/$DENOM/g" $APPCONFIG
sed -i ""  "s/aphoton/$DENOM/g" $GENESIS



# Allocate genesis accounts (cosmos formatted addresses)
$CLI add-genesis-account $KEY $TOTALAMOUNT --keyring-backend $KEYRING  --home $CHAINHOME
$CLI add-genesis-account "$KEY"1 $MYAMOUNT --keyring-backend $KEYRING  --home $CHAINHOME
$CLI add-genesis-account "$KEY"2 $MYAMOUNT --keyring-backend $KEYRING  --home $CHAINHOME
$CLI add-genesis-account "$KEY"3 $MYAMOUNT --keyring-backend $KEYRING  --home $CHAINHOME

$CLI gentx $KEY $MYAMOUNT --keyring-backend $KEYRING --chain-id $CHAINID  --home $CHAINHOME

$CLI collect-gentxs  --home $CHAINHOME

$CLI validate-genesis  --home $CHAINHOME


. ./run.sh
