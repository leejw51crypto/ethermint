set CHAINID=pearmint-6
set CLI=ethermintd
set KEYRING=test
set NODE=tcp://127.0.0.1:26647
set HOME=%USERPROFILE%\.pear
set K1=eth1zzzp544l89fgjs3m7ajlewf22c7mda7wn5f6y7
set S1=eth1fzep9fc7hweq9a706x4vacardl0zl0z3rq3ls8
set S2=eth1znmmtmaw8v0vucxtdtq2e3nnvl3v9c8c7z0mwa

set FROM=%K1%
set TO=%S1%

rem set AMOUNT=1200000000000000aphoton
set AMOUNT=1002000000000000000001pphoton
%CLI% tx bank  send %FROM% %TO% %AMOUNT% --chain-id %CHAINID% --keyring-backend %KEYRING%  --fees 20pphoton --home %HOME%
