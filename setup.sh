. ./setup0.sh
export S1=$($CLI keys show mykey -a --keyring-backend test) 
export S2=$($CLI keys show mykey2 -a --keyring-backend test) 
export E1=48B212A71EBBB202F7CFD1AACEE3A36FDE2FBC51
export E2=14F7B5EFAE3B1ECE60CB6AC0ACC67367E2C2E0F8
#export I1=$($CLI keys show ibc1 -a --keyring-backend test)
#export I2=$($CLI keys show ibc2 -a --keyring-backend test)

ethermintd debug addr $S1
ethermintd debug addr $S2

echo 'S1='$S1
echo 'E1='$E1
echo 'S2='$S2
echo 'E2='$E2
#echo 'I1='$I1
#echo 'I2='$I2
