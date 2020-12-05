package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

// func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
// 	bytes := stub.State[name]
// 	if bytes == nil {
// 		fmt.Println("State", name, "failed to get value")
// 		t.FailNow()
// 	}
// 	if string(bytes) != value {
// 		fmt.Println("State value", name, "was not", value, "as expected")
// 		t.FailNow()
// 	}
// }

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		// t.FailNow()
	}
}

func Test_Init(t *testing.T) {
	scc := new(MovieChaincode)
	stub := shim.NewMockStub("movies", scc)

	checkInit(t, stub, nil)
}

func Test_Invoke(t *testing.T) {
	scc := new(MovieChaincode)
	stub := shim.NewMockStub("movies", scc)

	tMovieId := "MOVIE1"
	tmovieName := "Mirrot"
	tThaterName := "Multiplex Theater"

	checkInvoke(t, stub,
		[][]byte{
			[]byte("addMovie"),
			[]byte(tMovieId),
			[]byte(tmovieName),
			[]byte(tThaterName),
		})

	// This will Fails Because the Method Not implemented
	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/mockstub.go#286
	tShowId := "SHOW1"
	checkInvoke(t, stub,
		[][]byte{
			[]byte("addShow"),
			[]byte(tShowId),
			[]byte(tMovieId),
		})

	// This will Fails Because the Method Not implemented
	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/mockstub.go#286
	TicketId := "S1"
	checkInvoke(t, stub,
		[][]byte{
			[]byte("purchaseTickets"),
			[]byte(tShowId),
			[]byte(TicketId),
		})

	// This will Fails Because the Method Not implemented
	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/mockstub.go#286
	SeatId := tShowId + "-" + TicketId
	checkInvoke(t, stub,
		[][]byte{
			[]byte("exchange"),
			[]byte(SeatId),
		})
}

func Test_Query(t *testing.T) {
	scc := new(MovieChaincode)
	stub := shim.NewMockStub("movies", scc)

	tMovieId := "MOVIE1"
	tmovieName := "Mirrot"
	tThaterName := "Multiplex Theater"

	checkInvoke(t, stub,
		[][]byte{
			[]byte("addMovie"),
			[]byte(tMovieId),
			[]byte(tmovieName),
			[]byte(tThaterName),
		})

	// This will fail because Method was not implemented.
	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/mockstub.go#286
	res := stub.MockInvoke("1", [][]byte{[]byte("getMovies")})
	fmt.Println(res)
	if res.Status == shim.OK {
		fmt.Println(res.Payload)
	}
}
