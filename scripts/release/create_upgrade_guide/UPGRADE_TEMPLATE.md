# Mainnet Upgrade Guide: From Version $CURRENT_VERSION to $UPGRADE_VERSION

## Overview

- **$UPGRADE_VERSION Proposal**: [Proposal Page](https://www.mintscan.io/furya/proposals/$PROPOSAL_ID)
- **$UPGRADE_VERSION Upgrade Block Height**: $UPGRADE_BLOCK
- **$UPGRADE_VERSION Upgrade Countdown**: [Block Countdown](https://www.mintscan.io/furya/blocks/$UPGRADE_BLOCK)

## Hardware Requirements

### Memory Specifications

Although this upgrade is not expected to be resource-intensive, a minimum of 64GB of RAM is advised. If you cannot meet this requirement, setting up a swap space is recommended.

#### Configuring Swap Space

*Execute these commands to set up a 32GB swap space*:

```sh
sudo swapoff -a
sudo fallocate -l 32G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

*To ensure the swap space persists after reboot*:

```sh
sudo cp /etc/fstab /etc/fstab.bak
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

For an in-depth guide on swap configuration, please refer to [this tutorial](https://www.digitalocean.com/community/tutorials/how-to-add-swap-space-on-ubuntu-20-04).

---

## Cosmovisor Configuration

### Initial Setup (For First-Time Users)

If you have not previously configured Cosmovisor, follow this section; otherwise, proceed to the next section.

Cosmovisor is strongly recommended for validators to minimize downtime during upgrades. It automates the binary replacement process according to on-chain `SoftwareUpgrade` proposals.

Documentation for Cosmovisor can be found [here](https://docs.cosmos.network/main/tooling/cosmovisor).

#### Installation Steps

*Run these commands to install and configure Cosmovisor*:

```sh
go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@v1.0.0
mkdir -p ~/.furyad
mkdir -p ~/.furyad/cosmovisor
mkdir -p ~/.furyad/cosmovisor/genesis
mkdir -p ~/.furyad/cosmovisor/genesis/bin
mkdir -p ~/.furyad/cosmovisor/upgrades
cp $GOPATH/bin/furyad ~/.furyad/cosmovisor/genesis/bin
mkdir -p ~/.furyad/cosmovisor/upgrades/$CURRENT_VERSION/bin
cp $GOPATH/bin/furyad ~/.furyad/cosmovisor/upgrades/$CURRENT_VERSION/bin
```

*Add these lines to your profile to set up environment variables*:

```sh
echo "# Cosmovisor Setup" >> ~/.profile
echo "export DAEMON_NAME=furyad" >> ~/.profile
echo "export DAEMON_HOME=$HOME/.furyad" >> ~/.profile
echo "export DAEMON_ALLOW_DOWNLOAD_BINARIES=false" >> ~/.profile
echo "export DAEMON_LOG_BUFFER_SIZE=512" >> ~/.profile
echo "export DAEMON_RESTART_AFTER_UPGRADE=true" >> ~/.profile
echo "export UNSAFE_SKIP_BACKUP=true" >> ~/.profile
source ~/.profile
```

### Upgrading to $UPGRADE_VERSION

*To prepare for the upgrade, execute these commands*:

```sh
mkdir -p ~/.furyad/cosmovisor/upgrades/$UPGRADE_VERSION/bin
cd $HOME/furya
git pull
git checkout $UPGRADE_TAG
make build
cp build/furyad ~/.furyad/cosmovisor/upgrades/$UPGRADE_VERSION/bin
```

At the designated block height, Cosmovisor will automatically upgrade to version $UPGRADE_VERSION.

---

## Manual Upgrade Procedure

Follow these steps if you opt for a manual upgrade:

1. Monitor Furya until it reaches the specified upgrade block height: $UPGRADE_BLOCK.
2. Observe for a panic message followed by continuous peer logs, then halt the daemon.
3. Perform these steps:

```sh
cd $HOME/furya
git pull
git checkout $UPGRADE_TAG
make install
```

4. Restart the Furya daemon and observe the upgrade.

---

## Additional Resources

- Furya Documentation: [Website](https://docs.furya.zone)
- Community Support: [Discord](https://discord.gg/pAxjcFnAFH)
