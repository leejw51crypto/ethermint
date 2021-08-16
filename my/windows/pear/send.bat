set CHAINID=pearmint-6
set CLI=ethermintd
set KEYRING=test
set NODE=tcp://127.0.0.1:26647
set HOME=%USERPROFILE%\.pear
echo %HOME%

set K1=eth1s28vg4ked3aeh6upsxc0l9jt85s7dq9pz4rhan
set S1=eth1fzep9fc7hweq9a706x4vacardl0zl0z3rq3ls8
set S2=eth1znmmtmaw8v0vucxtdtq2e3nnvl3v9c8c7z0mwa

set FROM=%S1%
set TO=%S2%

set AMOUNT=1021000000000000001pphoton
%CLI% tx bank  send %FROM% %TO% %AMOUNT% --chain-id %CHAINID% --keyring-backend %KEYRING%  --fees 20pphoton --home %HOME% 
