. ./setup.sh
$CLI keys list --keyring-backend $KEYRING


FIRSTADDR=$($CLI keys list --keyring-backend $KEYRING | yq eval -j | jq '.[0]'.address -r)
echo "first address=" $FIRSTADDR
