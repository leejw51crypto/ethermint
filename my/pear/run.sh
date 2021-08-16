. ./setup.sh
$CLI start --pruning=nothing $TRACE --log_level $LOGLEVEL --minimum-gas-prices=0.0001$DENOM --json-rpc.api eth,txpool,personal,net,debug,web3,miner  --home $CHAINHOME