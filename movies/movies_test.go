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

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkMovieQuery(t *testing.T, stub *shim.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getMovies"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func TestMovies_Init(t *testing.T) {
	scc := new(MovieChaincode)
	stub := shim.NewMockStub("movies", scc)

	// Init A=123 B=234
	checkInit(t, stub, [][]byte{[]byte("init"), nil})

}

func TestAbac_Query(t *testing.T) {
	scc := new(MovieChaincode)
	stub := shim.NewMockStub("movies", scc)

	// Init A=345 B=456
	// checkInit(t, stub, nil)

	tMovieId := "MOVIE1"
	tmovieName := "Mirrot"
	tThaterName := "Multiplex Theater"

	checkInvoke(t, stub, [][]byte{
		[]byte("addMovie"),
		[]byte(tMovieId),
		[]byte(tmovieName),
		[]byte(tThaterName),
	})

}
