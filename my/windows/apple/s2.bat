set CHAINID=ethermint-2
set CLI=ethermintd
set KEYRING=test
set NODE=tcp://127.0.0.1:26657


set S1=eth1fzep9fc7hweq9a706x4vacardl0zl0z3rq3ls8
set S2=eth1znmmtmaw8v0vucxtdtq2e3nnvl3v9c8c7z0mwa
%CLI% query bank balances  %S2% --node %NODE%
