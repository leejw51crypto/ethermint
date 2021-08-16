set CHAINID=applemint-2
set CLI=ethermintd
set KEYRING=test
set NODE=tcp://127.0.0.1:26657
set HOME=%USERPROFILE%\.appled

set K1=eth1txnyk7f3fahs4fgqgn9dsy6cu23vp2dwssett5
set S1=eth1fzep9fc7hweq9a706x4vacardl0zl0z3rq3ls8
set S2=eth1znmmtmaw8v0vucxtdtq2e3nnvl3v9c8c7z0mwa

set FROM=%K1%
set TO=%S1%

rem set AMOUNT=1200000000000000aphoton
set AMOUNT=1002000000000000000001aphoton
%CLI% tx bank  send %FROM% %TO% %AMOUNT% --chain-id %CHAINID% --keyring-backend %KEYRING%  --fees 20aphoton --home %HOME% 
