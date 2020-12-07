# Movie Chaincode

### PreRequisites

* Hyperledger Fabric 1.4 
    
    `
    install Hyperledger and PreRequisites
    `
    [Click Here](https://hackernoon.com/hyperledger-fabric-installation-guide-74065855eca9)

* Spin Up Your Test Network.
    
    ```bash
    $ cd <path-to-fabric-samples>/first-network

    # start network with couchdb as storage
    $ ./byfn.sh up -c mychannel -s couchdb
    ```
* Enter into cli container

    ```bash

    $ docker exec -it cli bash
    ```
* Install Chaincode
    ```bash
    $ export CHAINCODE=mycc1
    $ peer chaincode install -n $CHAINCODE -v 1.0 -l golang -p github.com/chaincode/movies
    
    ```
* Instalntiate Chaincode 
    ```bash
    $ peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n $CHAINCODE -l golang -v 1.0 -c '{"Args":[]}' -P 'OR ('\''Org1MSP.peer'\'','\''Org2MSP.peer'\'')'
    ```
* Do Some Transations
    ```bash

    # Add Movie to BlockChain Network
    $ peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n $CHAINCODE --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"Args":["addMovie","MOVIE1","RRR","INOX"]}'

    # Add MovieShow to BlockChain Network
    $ peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n $CHAINCODE --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"Args":["addShow","SHOW1","MOVIE1"]}'


     # Purchase Tickets 
    $ peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n $CHAINCODE --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"Args":["purchaseTickets","SHOW1","S1"]}'

    # exchange waterBottle with Soda
    $ peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n $CHAINCODE --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"Args":["exchange","SHOW1-S1"]}'

    ```
* Query Some Data 

    ```bash

     # Get All Movies
    $ peer chaincode query -C mychannel -n $CHAINCODE -c '{"Args":["getMovies"]}'

    # Get All Movie Shows
    $ peer chaincode query -C mychannel -n $CHAINCODE -c '{"Args":["getMovieShows","MOVIE1"]}'


    # Get Seat
    $ peer chaincode query -C mychannel -n $CHAINCODE -c '{"Args":["getSeat","SHOW1-S1","MOVIE1"]}'

    # Get All Theators
    $ peer chaincode query -C mychannel -n $CHAINCODE -c '{"Args":["getTheaters"]}'
    
    ```
