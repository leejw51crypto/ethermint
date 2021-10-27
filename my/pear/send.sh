. ./setup.sh
export FROM=$S1
export TO=$S2


export AMOUNT=2100000000000000001$DENOM
echo "send amount $AMOUNT"
$CLI tx bank  send $FROM $TO $AMOUNT --chain-id $CHAINID --keyring-backend $KEYRING  --fees 20$DENOM --home $CHAINHOME
