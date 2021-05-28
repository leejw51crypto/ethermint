export CHAINID=ethermint-2
export KEY=mykey 
export CLI=ethermintd
export S1=$($CLI keys show $KEY -a --keyring-backend test)
echo "S1="$S1
