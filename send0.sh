. ./setup.sh
FROM=$K1
TO=$S1
AMOUNT=100000123000000000000000aphoton

echo "from="$FROM
echo "to="$TO
echo "amount="$AMOUNT
$CLI tx bank  send $FROM $TO $AMOUNT --chain-id $CHAINID --keyring-backend $KEYRING  --fees 50aphoton 
