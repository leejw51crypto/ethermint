. ./setup0.sh
export S1=$($CLI keys show mykey -a --keyring-backend test) 
export S2=$($CLI keys show mykey2 -a --keyring-backend test) 
export E1=4891D3A4EB06708FAF7A6E3E823CE4545572936C
#export I1=$($CLI keys show ibc1 -a --keyring-backend test)
#export I2=$($CLI keys show ibc2 -a --keyring-backend test)


echo 'S1='$S1
echo 'E1='$E1
echo 'S2='$S2
#echo 'I1='$I1
#echo 'I2='$I2
