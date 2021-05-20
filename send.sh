. ./setup.sh
FROM=$S1
TO=$I1
AMOUNT=100000000aphoton
$CLI tx bank  send $FROM $TO $AMOUNT --chain-id $CHAINID --keyring-backend $KEYRING 
