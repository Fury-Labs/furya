# Download a genesis.json for testing. The node that you this on will be your "validator"
# It should be on version v4.x

furyad init --chain-id=testing testing --home=$HOME/.furyad
furyad keys add validator --keyring-backend=test --home=$HOME/.furyad
furyad add-genesis-account $(furyad keys show validator -a --keyring-backend=test --home=$HOME/.furyad) 1000000000ufury,1000000000valtoken --home=$HOME/.furyad
sed -i -e "s/stake/ufury/g" $HOME/.furyad/config/genesis.json
furyad gentx validator 500000000ufury --commission-rate="0.0" --keyring-backend=test --home=$HOME/.furyad --chain-id=testing
furyad collect-gentxs --home=$HOME/.furyad

cat $HOME/.furyad/config/genesis.json | jq '.initial_height="711800"' > $HOME/.furyad/config/tmp_genesis.json && mv $HOME/.furyad/config/tmp_genesis.json $HOME/.furyad/config/genesis.json
cat $HOME/.furyad/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"]["denom"]="valtoken"' > $HOME/.furyad/config/tmp_genesis.json && mv $HOME/.furyad/config/tmp_genesis.json $HOME/.furyad/config/genesis.json
cat $HOME/.furyad/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"]["amount"]="100"' > $HOME/.furyad/config/tmp_genesis.json && mv $HOME/.furyad/config/tmp_genesis.json $HOME/.furyad/config/genesis.json
cat $HOME/.furyad/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="120s"' > $HOME/.furyad/config/tmp_genesis.json && mv $HOME/.furyad/config/tmp_genesis.json $HOME/.furyad/config/genesis.json
cat $HOME/.furyad/config/genesis.json | jq '.app_state["staking"]["params"]["min_commission_rate"]="0.050000000000000000"' > $HOME/.furyad/config/tmp_genesis.json && mv $HOME/.furyad/config/tmp_genesis.json $HOME/.furyad/config/genesis.json

# Now setup a second full node, and peer it with this v3.0.0-rc0 node.

# start the chain on both machines
furyad start
# Create proposals

furyad tx gov submit-proposal --title="existing passing prop" --description="passing prop"  --from=validator --deposit=1000valtoken --chain-id=testing --keyring-backend=test --broadcast-mode=block  --type="Text"
furyad tx gov vote 1 yes --from=validator --keyring-backend=test --chain-id=testing --yes
furyad tx gov submit-proposal --title="prop with enough fury deposit" --description="prop w/ enough deposit"  --from=validator --deposit=500000000ufury --chain-id=testing --keyring-backend=test --broadcast-mode=block  --type="Text"
# Check that we have proposal 1 passed, and proposal 2 in deposit period
furyad q gov proposals
# CHeck that validator commission is under min_commission_rate
furyad q staking validators
# Wait for upgrade block.
# Upgrade happened
# your full node should have crashed with consensus failure

# Now we test post-upgrade behavior is as intended

# Everything in deposit stayed in deposit
furyad q gov proposals
# Check that commissions was bumped to min_commission_rate
furyad q staking validators
# pushes 2 into voting period
furyad tx gov deposit 2 1valtoken --from=validator --keyring-backend=test --chain-id=testing --yes