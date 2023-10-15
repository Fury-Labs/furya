# run old binary on terminal1
go clean --modcache
git stash
git checkout v1.0.1
go install ./cmd/furyad/
Run the below commands
```
    #!/bin/bash
    rm -rf $HOME/.furyad/
    cd $HOME
    furyad init --chain-id=testing testing --home=$HOME/.furyad
    furyad keys add validator --keyring-backend=test --home=$HOME/.furyad
    furyad add-genesis-account $(furyad keys show validator -a --keyring-backend=test --home=$HOME/.furyad) 1000000000stake,1000000000valtoken --home=$HOME/.furyad
    furyad gentx validator 500000000stake --keyring-backend=test --home=$HOME/.furyad --chain-id=testing
    furyad gentx validator 500000000stake --commission-rate="0.0" --keyring-backend=test --home=$HOME/.furyad --chain-id=testing
    furyad collect-gentxs --home=$HOME/.furyad
    
    cat $HOME/.furyad/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="120s"' > $HOME/.furyad/config/tmp_genesis.json && mv $HOME/.furyad/config/tmp_genesis.json $HOME/.furyad/config/genesis.json
    cat $HOME/.furyad/config/genesis.json | jq '.app_state["staking"]["params"]["min_commission_rate"]="0.050000000000000000"' > $HOME/.furyad/config/tmp_genesis.json && mv $HOME/.furyad/config/tmp_genesis.json $HOME/.furyad/config/genesis.json

```

Create pool.json 
```
{
  "weights": "1stake,1valtoken",
  "initial-deposit": "100stake,20valtoken",
  "swap-fee": "0.01",
  "exit-fee": "0.01",
  "future-governor": "168h"
}
```

rm $HOME/.furyad/cosmovisor/current -rf
cosmovisor start

# operations on terminal2
furyad tx lockup lock-tokens 100stake --duration="5s" --from=validator --chain-id=testing --keyring-backend=test --yes
sleep 7
furyad tx gov submit-proposal software-upgrade upgrade-lockup-module-store-management --title="lockup module upgrade" --description="lockup module upgrade for gas efficiency"  --from=validator --upgrade-height=10 --deposit=10000000stake --chain-id=testing --keyring-backend=test -y
sleep 7
furyad tx gov vote 1 yes --from=validator --keyring-backend=test --chain-id=testing --yes
sleep 7
furyad tx gamm create-pool --pool-file="./pool.json"  --gas=3000000 --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block
sleep 7
furyad tx lockup lock-tokens 1000stake --duration="100s" --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block
sleep 7
furyad tx lockup lock-tokens 2000stake --duration="200s" --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block
sleep 7
furyad tx lockup lock-tokens 3000stake --duration="1s" --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block
sleep 7
furyad tx lockup begin-unlock-by-id 1 --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block
sleep 7
furyad tx lockup begin-unlock-by-id 3 --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block
sleep 7
furyad tx gov submit-proposal software-upgrade "v2" --title="lockup module upgrade" --description="lockup module upgrade for gas efficiency"  --from=validator --upgrade-height=20 --deposit=10000000stake --chain-id=testing --keyring-backend=test --yes  --broadcast-mode=block
sleep 7
furyad tx gov vote 1 yes --from=validator --keyring-backend=test --chain-id=testing --yes --broadcast-mode=block
furyad query gov proposal 1
furyad query upgrade plan
furyad query lockup account-locked-longer-duration $(furyad keys show -a --keyring-backend=test validator) 1s
furyad query gamm pools
furyad query staking validators
furyad query staking params

# on terminal1
Wait until consensus failure happen and stop binary using Ctrl + C
git checkout lockup_module_genesis_export
git checkout main

Update go mod file to use latest SDK changes: /Users/admin/go/pkg/mod/github.com/fury-labs/cosmos-sdk@v0.42.5-0.20210622202917-f4f6a08ac64b
go get github.com/fury-labs/cosmos-sdk@ea1ec79c739ba39639b9a24f824127ecc6650887
go: downloading github.com/fury-labs/cosmos-sdk v0.42.5-0.20210630100106-ea1ec79c739b
Upgrade Furya Cosmos SDK version to `v0.42.5-0.20210630100106-ea1ec79c739b`
go mod download github.com/cosmos/cosmos-sdk
git stash
git checkout min_commission_change_validation_change_ignore
go install ./cmd/furyad/
furyad start --home=$HOME/.furyad

# check on terminal2
furyad query lockup account-locked-longer-duration $(furyad keys show -a --keyring-backend=test validator) 1s
furyad query lockup account-locked-longer-duration $(furyad keys show -a --keyring-backend=test validator) 1s
furyad query lockup module-locked-amount
furyad query gamm pools
furyad query staking validators
furyad query staking params
furyad query bank balances $(furyad keys show -a --keyring-backend=test validator)
furyad tx staking edit-validator --commission-rate="0.1"  --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block
furyad tx staking edit-validator --commission-rate="0.08"  --from=validator --chain-id=testing --keyring-backend=test --yes --broadcast-mode=block

Result:
- pool exists
- lockup processed all correctly
- validator commission rate worked
- chain did not panic 