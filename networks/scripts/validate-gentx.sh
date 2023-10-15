#!/bin/sh
FURYA_HOME="/tmp/furyad$(date +%s)"
RANDOM_KEY="randomfuryavalidatorkey"
CHAIN_ID=furya-1
DENOM=uosmo
MAXBOND=50000000000000 # 500 Million OSMO

GENTX_FILE=$(find ./$CHAIN_ID/gentxs -iname "*.json")
LEN_GENTX=$(echo ${#GENTX_FILE})

# Gentx Start date
start="2021-06-03 15:00:00Z"
# Compute the seconds since epoch for start date
stTime=$(date --date="$start" +%s)

# Gentx End date
end="2021-07-12 23:59:59Z"
# Compute the seconds since epoch for end date
endTime=$(date --date="$end" +%s)

# Current date
current=$(date +%Y-%m-%d\ %H:%M:%S)
# Compute the seconds since epoch for current date
curTime=$(date --date="$current" +%s)

if [[ $curTime < $stTime ]]; then
    echo "start=$stTime:curent=$curTime:endTime=$endTime"
    echo "Gentx submission is not open yet. Please close the PR and raise a new PR after 04-June-2021 23:59:59"
    exit 0
else
    if [[ $curTime > $endTime ]]; then
        echo "start=$stTime:curent=$curTime:endTime=$endTime"
        echo "Gentx submission is closed"
        exit 0
    else
        echo "Gentx is now open"
        echo "start=$stTime:curent=$curTime:endTime=$endTime"
    fi
fi

if [ $LEN_GENTX -eq 0 ]; then
    echo "No new gentx file found."
else
    set -e

    echo "GentxFile::::"
    echo $GENTX_FILE

    echo "...........Init Furya.............."

    git clone https://github.com/fury-labs/furya
    cd furya
    git checkout gentx-launch
    make build
    chmod +x ./build/furyad

    ./build/furyad keys add $RANDOM_KEY --keyring-backend test --home $FURYA_HOME

    ./build/furyad init --chain-id $CHAIN_ID validator --home $FURYA_HOME

    echo "..........Fetching genesis......."
    rm -rf $FURYA_HOME/config/genesis.json
    curl -s https://raw.githubusercontent.com/fury-labs/networks/main/$CHAIN_ID/pregenesis.json >$FURYA_HOME/config/genesis.json

    # this genesis time is different from original genesis time, just for validating gentx.
    sed -i '/genesis_time/c\   \"genesis_time\" : \"2021-03-29T00:00:00Z\",' $FURYA_HOME/config/genesis.json

    GENACC=$(cat ../$GENTX_FILE | sed -n 's|.*"delegator_address":"\([^"]*\)".*|\1|p')
    denomquery=$(jq -r '.body.messages[0].value.denom' ../$GENTX_FILE)
    amountquery=$(jq -r '.body.messages[0].value.amount' ../$GENTX_FILE)

    echo $GENACC
    echo $amountquery
    echo $denomquery

    # only allow $DENOM tokens to be bonded
    if [ $denomquery != $DENOM ]; then
        echo "invalid denomination"
        exit 1
    fi

    # limit the amount that can be bonded

    if [ $amountquery -gt $MAXBOND ]; then
        echo "bonded too much: $amountquery > $MAXBOND"
        exit 1
    fi

    ./build/furyad add-genesis-account $RANDOM_KEY 100000000000000$DENOM --home $FURYA_HOME \
        --keyring-backend test

    ./build/furyad gentx $RANDOM_KEY 90000000000000$DENOM --home $FURYA_HOME \
        --keyring-backend test --chain-id $CHAIN_ID

    cp ../$GENTX_FILE $FURYA_HOME/config/gentx/

    echo "..........Collecting gentxs......."
    ./build/furyad collect-gentxs --home $FURYA_HOME
    sed -i '/persistent_peers =/c\persistent_peers = ""' $FURYA_HOME/config/config.toml

    ./build/furyad validate-genesis --home $FURYA_HOME

    echo "..........Starting node......."
    ./build/furyad start --home $FURYA_HOME &

    sleep 1800s

    echo "...checking network status.."

    ./build/furyad status --node http://localhost:26657

    echo "...Cleaning the stuff..."
    killall furyad >/dev/null 2>&1
    rm -rf $FURYA_HOME >/dev/null 2>&1
fi
