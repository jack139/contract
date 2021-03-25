#!/usr/bin/env bash

rm -rf ~/.nameserviced
rm -rf ~/.nameservicecli

contractd init test

#contractd config output json
#contractd config indent true
#contractd config trust-node true
#contractd config chain-id namechain
#contractd config keyring-backend test

contractd keys add user0
contractd keys add gt
contractd keys add faucet

contractd add-genesis-account $(contractd keys show user0 -a) 500token,100000000stake,500credit
contractd add-genesis-account $(contractd keys show gt -a) 500token,500credit
contractd add-genesis-account $(contractd keys show faucet -a) 10000000credit

contractd gentx user0 100000000stake --chain-id contract

echo "Collecting genesis txs..."
contractd collect-gentxs

echo "Validating genesis file..."
contractd validate-genesis


contract start