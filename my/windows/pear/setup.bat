set CHAINID=pearmint-6
set CLI=ethermintd
set KEYRING=test
set HOME=%USERPROFILE%\.pear
%CLI% keys list --keyring-backend %KEYRING% --home %HOME%
