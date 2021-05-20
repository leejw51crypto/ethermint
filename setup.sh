. ./setup0.sh
export S1=$($CLI keys show mykey -a --keyring-backend test) 
export I1=$($CLI keys show ibc1 -a --keyring-backend test)
export I2=$($CLI keys show ibc2 -a --keyring-backend test)


echo 'S1='$S1
echo 'I1='$I1
echo 'I2='$I2
