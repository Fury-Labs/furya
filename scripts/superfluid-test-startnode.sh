#!/bin/bash

rm -rf $HOME/.furyad/

cd $HOME

furyad init --chain-id=testing testing --home=$HOME/.furyad
furyad keys add validator --keyring-backend=test --home=$HOME/.furyad
furyad add-genesis-account $(furyad keys show validator -a --keyring-backend=test --home=$HOME/.furyad) 100000000000stake,100000000000valtoken --home=$HOME/.furyad
furyad gentx validator 500000000stake --keyring-backend=test --home=$HOME/.furyad --chain-id=testing
furyad collect-gentxs --home=$HOME/.furyad

update_genesis () {    
    cat $HOME/.furyad/validator1/config/genesis.json | jq "$1" > $HOME/.furyad/validator1/config/tmp_genesis.json && mv $HOME/.furyad/validator1/config/tmp_genesis.json $HOME/.furyad/validator1/config/genesis.json
}

# update staking genesis
update_genesis '.app_state["staking"]["params"]["unbonding_time"]="120s"'

# update governance genesis
update_genesis '.app_state["gov"]["voting_params"]["voting_period"]="10s"'

# update epochs genesis
update_genesis '.app_state["epochs"]["epochs"][0]["identifier"]="min"'
update_genesis '.app_state["epochs"]["epochs"][0]["duration"]="60s"'

# update poolincentives genesis
update_genesis '.app_state["poolincentives"]["lockable_durations"][0]="120s"'
update_genesis '.app_state["poolincentives"]["lockable_durations"][1]="180s"'
update_genesis '.app_state["poolincentives"]["lockable_durations"][2]="240s"'

# update incentives genesis
update_genesis '.app_state["incentives"]["params"]["distr_epoch_identifier"]="min"'
update_genesis '.app_state["incentives"]["lockable_durations"][0]="1s"'
update_genesis '.app_state["incentives"]["lockable_durations"][1]="120s"'
update_genesis '.app_state["incentives"]["lockable_durations"][2]="180s"'
update_genesis '.app_state["incentives"]["lockable_durations"][3]="240s"'

# update mint genesis
update_genesis '.app_state["mint"]["params"]["epoch_identifier"]="min"'

# update gamm genesis
update_genesis '.app_state["gamm"]["params"]["pool_creation_fee"][0]["denom"]="stake"'

# update superfluid genesis
update_genesis '.app_state["superfluid"]["params"]["minimum_risk_factor"]="0.500000000000000000"'

furyad start --home=$HOME/.furyad
