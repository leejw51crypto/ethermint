set CHAINID=ethermint-2
set CLI=ethermintd
set KEYRING=test
set NODE=tcp://127.0.0.1:26657
set HOME=%USERPROFILE%\.appled
set K1=eth1s28vg4ked3aeh6upsxc0l9jt85s7dq9pz4rhan
set S1=eth1fzep9fc7hweq9a706x4vacardl0zl0z3rq3ls8
set S2=eth1znmmtmaw8v0vucxtdtq2e3nnvl3v9c8c7z0mwa

set FROM=%S1%
set TO=eth1xjfdas234fshncflwa0dyjgc23u08k9dqqyx2l

set AMOUNT=12121000000000000001aphoton
%CLI% tx bank  send %FROM% %TO% %AMOUNT% --chain-id %CHAINID% --keyring-backend %KEYRING%  --fees 20aphoton --home %HOME%
