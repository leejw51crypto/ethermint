set CHAINID=pearmint-6
set CLI=ethermintd
set KEYRING=test
set NODE=tcp://127.0.0.1:26657


set S1=eth1fzep9fc7hweq9a706x4vacardl0zl0z3rq3ls8
set S3=eth1xjfdas234fshncflwa0dyjgc23u08k9dqqyx2l
%CLI% query bank balances  %S3% --node %NODE%
