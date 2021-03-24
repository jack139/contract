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
contractd keys add user1

contractd add-genesis-account $(contractd keys show user0 -a) 1000nametoken,100000000stake,10000000credit
contractd add-genesis-account $(contractd keys show user1 -a) 500nametoken,500credit

contractd gentx user0 100000000stake --chain-id contract

echo "Collecting genesis txs..."
contractd collect-gentxs

echo "Validating genesis file..."
contractd validate-genesis


contract start