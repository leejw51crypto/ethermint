. ./setup.sh
export FROM=$K1
export TO=$S1


export AMOUNT=1001000000000000000000$DENOM
echo "send amount $AMOUNT"
$CLI tx bank  send $FROM $TO $AMOUNT --chain-id $CHAINID --keyring-backend $KEYRING  --fees 20$DENOM --home $CHAINHOME
