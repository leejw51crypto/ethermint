. ./setup.sh
#$CLI query account $S1 --node $NODE
#$CLI query bank balances  $S1 --node $NODE
echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
$CLI query ibc-transfer escrow-address  transfer channel-0  --node $NODE
