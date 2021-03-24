#!/usr/bin/env bash

rm -rf ~/.nameserviced
rm -rf ~/.nameservicecli

contractd init test

#contractd config output json
#contractd config indent true
#contractd config trust-node true
#contractd config chain-id namechain
#contractd config keyring-backend test

contractd keys add user1
contractd keys add user2

contractd add-genesis-account $(contractd keys show user1 -a) 1000nametoken,100000000stake
contractd add-genesis-account $(contractd keys show user2 -a) 1000nametoken,100000000stake

contractd gentx user1 100000000stake --chain-id contract

echo "Collecting genesis txs..."
contractd collect-gentxs

echo "Validating genesis file..."
contractd validate-genesis


contract start