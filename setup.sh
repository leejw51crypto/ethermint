. ./setup0.sh
export K1=$($CLI keys list --keyring-backend $KEYRING | yq eval -j | jq '.[0]'.address -r )
export S1=$($CLI keys show testkey -a --keyring-backend test) 
export S2=$($CLI keys show testkey2 -a --keyring-backend test) 


echo "#############################################"
echo "K1="$K1
ethermintd debug addr $K1
echo "#############################################"
ethermintd debug addr $S1
echo "#############################################"
ethermintd debug addr $S2

echo "#############################################"
echo "### SETUP ###"
echo 'K1='$K1
echo 'S1='$S1
echo 'S2='$S2
