. ./setup.sh
FROM=$S1
TO=$S2
AMOUNT=1000000000000000aphoton
$CLI tx bank  send $FROM $TO $AMOUNT --chain-id $CHAINID --keyring-backend $KEYRING  --fees 20aphoton
