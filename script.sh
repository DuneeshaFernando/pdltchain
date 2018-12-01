#!/bin/bash

echo
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo
echo "Build your first network (BYFN) end-to-end test"
echo
CHANNEL_NAME="$1"
DELAY="$2"
LANGUAGE="$3"
TIMEOUT="$4"
VERBOSE="$5"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${LANGUAGE:="golang"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=5

#CC_SRC_PATH="github.com/chaincode/chaincode_example02/go/"
CC_SRC_PATH="github.com/chaincode/fabcar/go/"
#CC_SRC_PATH="github.com/chaincode/marbles02/go/"
#CC_SRC_PATH="github.com/chaincode/marbles02_private/go/"
if [ "$LANGUAGE" = "node" ]; then
	CC_SRC_PATH="/opt/gopath/src/github.com/chaincode/chaincode_example02/node/"
fi

echo "Channel name : "$CHANNEL_NAME

# import utils
. scripts/utils.sh

createChannel() {
	setGlobals 0 1

	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
                set -x
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx >&log.txt
		res=$?
                set +x
	else
				set -x
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
		res=$?
				set +x
	fi
	cat log.txt
	verifyResult $res "Channel creation failed"
	#cat log.txt>&scripts/abc.txt
	echo "===================== Channel '$CHANNEL_NAME' created ===================== "
	echo
}

joinChannel () {
	for org in 1 2; do
	    for peer in 0 1; do
		joinChannelWithRetry $peer $org
		echo "===================== peer${peer}.org${org} joined channel '$CHANNEL_NAME' ===================== "
		sleep $DELAY
		echo
	    done
	done
}

## Create channel
echo "Creating channel..."
createChannel

## Join all the peers to the channel
echo "Having all peers join the channel..."
joinChannel

## Set the anchor peers for each org in the channel
echo "Updating anchor peers for org1..."
updateAnchorPeers 0 1
echo "Updating anchor peers for org2..."
updateAnchorPeers 0 2

## Install chaincode on peer0.org1 and peer0.org2
echo "Installing chaincode on peer0.org1..."
installChaincode 0 1
echo "Install chaincode on peer0.org2..."
installChaincode 0 2

#Added by me
echo "Instantiating chaincode on peer0.org1..."
instantiateChaincode 0 1

#Added by me - very important initLedger with an invoke function
#chaincodeInvoke 0 1

# Instantiate chaincode on peer0.org2
echo "Instantiating chaincode on peer0.org2..."
instantiateChaincode 0 2

# Query chaincode on peer0.org1
#echo "Querying chaincode on peer0.org1..."
#chaincodeQuery 0 1 'LAND1' '{"rlregistry":"Delkanda","extent":25,"parentlandid":"nil","owner":"Brad","boundaries":[[0,20],[10,20],[10,0],[0,0]]}'

#This will be the starting time from which throughput would be calculated.
chaincodeStartTime
# Invoke chaincode on peer0.org1 and peer0.org2
echo "Sending invoke change owner transaction on peer0.org1 peer0.org2..."
chaincodeInvokeChangeOwner 'LAND1' 'Duneesha' 0 1 0 2 &

# Invoke chaincode on peer0.org1 and peer0.org2
echo "Sending invoke fork land transaction on peer0.org1 peer0.org2..."
chaincodeInvokeForkLand 'LAND2' 'LAND21' 'Samanmalie' 35 'LAND22' 'Migara' 40 'v' 3 10 0 1 0 2 &

#2 extra invokes below to check concurrency
echo "Sending invoke change owner transaction on peer0.org1 peer0.org2..."
chaincodeInvokeChangeOwner 'LAND3' 'Tharukaa' 0 1 0 2 &
echo "Sending invoke fork land transaction on peer0.org1 peer0.org2..."
chaincodeInvokeForkLand 'LAND4' 'LAND41' 'Samanmalie' 15 'LAND42' 'Migara' 15 'h' 5 11 0 1 0 2 &

##Modification to draw throughput graph from here to
chaincodeInvokeChangeOwner 'LAND1' 'Martha' 0 1 0 2 &
chaincodeInvokeChangeOwner 'LAND3' 'Sarada' 0 1 0 2 &
chaincodeInvokeChangeOwner 'LAND5' 'Sujatha' 0 1 0 2 &
chaincodeInvokeForkLand 'LAND6' 'LAND61' 'Kusumsiri' 20 'LAND62' 'Mario' 25 'v' 3 10 0 1 0 2 &
chaincodeInvokeChangeOwner 'LAND7' 'Kareem' 0 1 0 2 &
chaincodeInvokeForkLand 'LAND8' 'LAND81' 'Sandun' 10 'LAND82' 'Daphne' 15 'v' 3 10 0 1 0 2 &
chaincodeInvokeChangeOwner 'LAND9' 'Maduri' 0 1 0 2 &
chaincodeInvokeForkLand 'LAND0' 'LAND101' 'Piumi' 30 'LAND102' 'Taniya' 20 'v' 3 10 0 1 0 2 &
##here 23rd nov

## Install chaincode on peer1.org2
#echo "Installing chaincode on peer1.org2..."
#installChaincode 1 2

# Query on chaincode on peer1.org2, check if the result is 90
wait
echo "Querying chaincode on peer1.org2..."
#chaincodeQuery 0 1 'LAND1' '{"rlregistry":"Delkanda","extent":25,"parentlandid":"nil","owner":"Martha","boundaries":[[0,20],[10,20],[10,0],[0,0]]}'
#above 1 2
chaincodeQuery 0 1 'LAND7' '{"rlregistry":"Hambantota","extent":20,"parentlandid":"nil","owner":"Kareem","boundaries":[[0,20],[10,20],[10,0],[0,0]]}'
# Query on chaincode on peer1.org2, check if the result is 90
#echo "Querying chaincode on peer1.org2..."
#chaincodeQuery 1 2 'LAND21' '{"rlregistry":"Avissawella","extent":35,"parentlandid":"LAND2","owner":"Samanmalie","boundaries":[[0,20],[3,20],[3,0],[0,0]]}'
#chaincodeQuery 0 1 'LAND22' '{"rlregistry":"Avissawella","extent":40,"parentlandid":"LAND2","owner":"Migara","boundaries":[[3,20],[10,20],[10,0],[3,0]]}'
#chaincodeQuery 1 2 'LAND41' '{"rlregistry":"Galle","extent":15,"parentlandid":"LAND4","owner":"Samanmalie","boundaries":[[0,20],[10,20],[10,11],[0,11]]}'
#chaincodeQuery 0 1 'LAND42' '{"rlregistry":"Galle","extent":10,"parentlandid":"LAND4","owner":"Migara","boundaries":[[0,11],[10,11],[10,0],[0,0]]}'
#above was initially 1 2 without &
#chaincodeInvokeForkLand 'LAND21' 'LAND211' 'Maya' 15 'LAND212' 'Saliya' 20 0 1 0 2

echo
echo "========= All GOOD, BYFN execution completed =========== "
echo

echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

exit 0
