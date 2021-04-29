#!/usr/bin/env bash

rm -rf ./n1

contractd init node1 --home n1

#contractd config output json
#contractd config indent true
#contractd config trust-node true
#contractd config chain-id namechain
#contractd config keyring-backend test

contractd keys add user0 --home n1
contractd keys add gt --home n1
contractd keys add faucet --home n1

contractd add-genesis-account $(contractd keys show user0 -a --home n1) 500token,100000000stake,500credit --home n1
contractd add-genesis-account $(contractd keys show gt -a --home n1) 500token,1credit --home n1
contractd add-genesis-account $(contractd keys show faucet -a --home n1) 21000000credit --home n1

contractd gentx user0 100000000stake --chain-id contract --home n1

echo "Collecting genesis txs..."
contractd collect-gentxs --home n1

echo "Validating genesis file..."
contractd validate-genesis --home n1


contractd start --log_level warn --home n1