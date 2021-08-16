set CHAINID=ethermint-2
set CLI=ethermintd
set KEYRING=test
set NODE=tcp://127.0.0.1:26657

set K1=eth1s28vg4ked3aeh6upsxc0l9jt85s7dq9pz4rhan
set S1=0x3492dEc151Aa6179e13F775eD249185478F3D8ad
set S2=0x14F7B5EFAE3B1ECE60CB6AC0ACC67367E2C2E0F8

set FROM=%S1%
set TO=%S2%

set AMOUNT=3
%CLI% ledger  send %FROM% %TO% %AMOUNT% 
