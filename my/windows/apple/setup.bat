set CHAINID=ethermint-2
set CLI=ethermintd
set KEYRING=test
set HOME=%USERPROFILE%\.appled
%CLI% keys list --keyring-backend %KEYRING% --home %HOME%
