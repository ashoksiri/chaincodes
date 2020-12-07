package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Movie struct {
	MovieId     string    `json:"movieId"`
	Name        string    `json:"name"`
	TheaterName string    `json:"thearerName"`
	ShowCount   int64     `json:"showCount"`
	CreatedTime time.Time `json:"createdTime"`
	Created     string    `json:"created"`
	UpdatedTime time.Time `json:"updatedTime"`
	Updated     string    `json:"updated"`
	DocType     int64     `json:"docType"`
}

// Fill Default Values for Movie
func (movie *Movie) FillDefaults() {

	movie.ShowCount = 0
	movie.CreatedTime = time.Now()
	movie.Created = time.Now().Format("2006-01-02")
	movie.UpdatedTime = time.Now()
	movie.Updated = time.Now().Format("2006-01-02")
	movie.DocType = 1
}

// ====================================================
// Structure of Show
// ====================================================
type Show struct {
	ShowId      string    `json:"showId"`
	MovieId     string    `json:"movieId"`
	TicketSold  int64     `json:"tickerSold"`
	IsFull      bool      `json:"isFull"`
	SeatCount   int64     `json:"seatCount"`
	FilledSeats int64     `json:"filledSeats"`
	FreeSeats   int64     `json:"freeSeats"`
	CreatedTime time.Time `json:"createdTime"`
	Created     string    `json:"created"`
	UpdatedTime time.Time `json:"updatedTime"`
	Updated     string    `json:"updated"`
	DocType     int64     `json:"docType"`
}

// Fill Default Values for Show
func (show *Show) FillDefaults() {
	show.TicketSold = 0
	show.IsFull = false
	show.SeatCount = 100
	show.FreeSeats = 100
	show.FilledSeats = 0
	show.CreatedTime = time.Now()
	show.Created = time.Now().Format("2006-01-02")
	show.UpdatedTime = time.Now()
	show.Updated = time.Now().Format("2006-01-02")
	show.DocType = 2
}

// ====================================================
// Structure of Seat
// ====================================================
type Seat struct {
	SeatId        string    `json:"seatId"`
	ShowId        string    `json:"showId"`
	PopCorn       bool      `json:"popcorn"`
	WaterBottle   bool      `json:"waterBottle"`
	Soda          bool      `json:"soda"`
	ExchangeCount int64     `json:"exchangeCount"`
	ExchangeId    int64     `json:"exchangeId"`
	CreatedTime   time.Time `json:"createdTime"`
	Created       string    `json:"created"`
	UpdatedTime   time.Time `json:"updatedTime"`
	Updated       string    `json:"updated"`
	DocType       int64     `json:"docType"`
}

// Fill Default Values for Seat
func (seat *Seat) FillDefaults() {

	seat.ExchangeCount = 0
	seat.CreatedTime = time.Now()
	seat.Created = time.Now().Format("2006-01-02")
	seat.UpdatedTime = time.Now()
	seat.Updated = time.Now().Format("2006-01-02")
	seat.DocType = 3
	seat.ExchangeId = rand.Int63n(1000)
	seat.WaterBottle = true
	seat.PopCorn = true
}

// MovieChaincode example simple Chaincode implementation
type MovieChaincode struct {
}

func (t *MovieChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Movie ChaineCode Iniated Successfully.....")
	return shim.Success(nil)
}

// ====================================================
// Invokes the Transation Using Args
// ====================================================
func (t *MovieChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "addShow" {
		return t.AddShow(stub, args)
	} else if function == "addMovie" {
		return t.AddMovie(stub, args)
	} else if function == "purchaseTickets" {
		return t.purchaseTickets(stub, args)
	} else if function == "exchange" {
		return t.exchange(stub, args)
	} else if function == "getMovieShows" {
		return t.getMovieShows(stub, args)
	} else if function == "getMovies" {
		return t.getMovies(stub, args)
	} else if function == "getSeat" {
		return t.getSeat(stub, args)
	} else if function == "getTheaters" {
		return t.getTheaters(stub, args)
	} else {
		return shim.Error("called invalid Function...")
	}

}

// ====================================================
// This Method will Add Movie Shows to BlockChain NetWork
// ====================================================
func (t *MovieChaincode) AddShow(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	if len(args) < 2 {
		return shim.Error("please provide Show Id and Thearer Name..")
	}
	ShowId := args[0]
	MovieId := args[1]
	movie, err := stub.GetState(MovieId)

	sshow, err := stub.GetState(ShowId)

	if sshow != nil {
		return shim.Error("Show Must Not be Duplicate...")
	}

	if err != nil {
		return shim.Error("Movie Not Found....")
	}

	today := time.Now().Format("2006-01-02")

	queryString := `{
		"selector": {
			"movieId": {
				"$eq" : "%s"
			},
			"created": {
				"$eq" : "%s"
			},
			"docType": {
				"$eq" : 2
			}
		}
	}`

	queryString = fmt.Sprintf(queryString, MovieId, today)

	resultIterator, err := stub.GetQueryResult(queryString)

	defer resultIterator.Close()

	if err != nil {
		return shim.Error(err.Error())
	}
	count := 0
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()

		fmt.Println(queryResponse)

		show := Show{}
		json.Unmarshal(queryResponse.Value, &show)

		if show.MovieId == MovieId && show.ShowId == ShowId && show.Created == today {
			return shim.Error("Show already Placed....")
		}

		if err != nil {
			return shim.Error(err.Error())
		}
		if queryResponse != nil {
			count++
		}

	}

	if count >= 5 {
		return shim.Error(fmt.Sprintf("Reached Show Limits Per Movie...."))
	}

	show := Show{ShowId: ShowId, MovieId: MovieId}
	show.FillDefaults()
	show_bytes, err := json.Marshal(show)

	if err != nil {
		return shim.Error("Error Occured During Marshal")
	}

	err = stub.PutState(args[0], show_bytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	M := Movie{}
	json.Unmarshal(movie, &M)

	M.ShowCount = M.ShowCount + 1
	movie, err = json.Marshal(M)

	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(M.MovieId, movie)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ====================================================
// This Method will Add Movie to BlockChain NetWork
// ====================================================
func (t *MovieChaincode) AddMovie(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	if len(args) < 3 {
		return shim.Error("please provide Movie Id, Name and Thearer Name..")
	}

	MovieId := args[0]
	Name := args[1]
	TheaterName := args[2]

	if err != nil {
		return shim.Error(err.Error())
	}

	tmovie, err := stub.GetState(MovieId)

	if tmovie != nil {
		M := Movie{}
		json.Unmarshal(tmovie, &M)
		return shim.Error(fmt.Sprintf("Movie %s Already exist", M.Name))
	}

	movie := Movie{MovieId: MovieId, Name: Name, TheaterName: TheaterName}
	movie.FillDefaults()

	movie_bytes, err := json.Marshal(movie)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(MovieId, movie_bytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *MovieChaincode) purchaseTickets(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) < 2 {
		return shim.Error("Insufficient args to execute transaction")
	}

	ShowId := args[0]
	TicketId := args[1]

	show, err := stub.GetState(ShowId)

	if err != nil {
		return shim.Error("Show Not Found")
	}

	tshow := Show{}
	json.Unmarshal(show, &tshow)

	if tshow.TicketSold >= 100 {
		return shim.Error("Show House Full...")
	}

	SeatId := ShowId + "-" + TicketId

	TSeat, err := stub.GetState(SeatId)

	if TSeat != nil {
		return shim.Error("This Ticket is Already Booked")
	}

	seat := Seat{SeatId: SeatId, ShowId: ShowId}
	seat.FillDefaults()

	seatBytes, err := json.Marshal(seat)

	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(SeatId, seatBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	tshow.FilledSeats = tshow.FilledSeats + 1
	tshow.FreeSeats = tshow.FreeSeats - 1
	tshow.TicketSold = tshow.TicketSold + 1

	if tshow.FilledSeats == 100 {
		tshow.IsFull = true
	}

	show, err = json.Marshal(tshow)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(tshow.ShowId, show)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ====================================================
// This Method will Exchange the waterBottle with Soda
// ====================================================
func (t *MovieChaincode) exchange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	if len(args) < 1 {
		return shim.Error("Insufficient args...")
	}

	SeatId := args[0]

	seat, err := stub.GetState(SeatId)

	tseat := Seat{}
	json.Unmarshal(seat, &tseat)

	if err != nil {
		return shim.Error(err.Error())
	}
	if tseat.ExchangeCount >= 100 {
		return shim.Error("EXchange limit exeeded...")
	}

	tseat.WaterBottle = false
	tseat.Soda = true
	tseat.ExchangeCount = tseat.ExchangeCount + 1

	seat, err = json.Marshal(tseat)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(SeatId, seat)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ====================================================
// This Method is a utility Function
// ====================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// ====================================================
// This Method will List all movies in Block Chain Network
// ====================================================
func (t *MovieChaincode) getMovies(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := `
	{
		"selector": {
			"docType": {
				"$eq": 1
			}
		}
	}
	`
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	buffer, err := constructQueryResponseFromIterator(resultsIterator)

	return shim.Success(buffer.Bytes())
}

// ====================================================
// This Method will List All Shows of A movie
// ====================================================

func (t *MovieChaincode) getMovieShows(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Invalid Movie Id ....")
	}

	queryString := `
	{
		"selector": {
			"movieId": {
				"$eq": "%s"
			},
			"docType": {
				"$eq": 2
			}
		}
	}
	`
	queryString = fmt.Sprintf(queryString, args[0])

	fmt.Println(queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	buffer, err := constructQueryResponseFromIterator(resultsIterator)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(buffer.Bytes())
}

// ====================================================
// This Method will Retrive the Seat Details
// ====================================================
func (t *MovieChaincode) getSeat(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("InSufficient Args for Transactions ....")
	}

	TicketId := args[0]
	ShowId := args[1]

	SeatId := ShowId + "-" + TicketId

	seat, err := stub.GetState(SeatId)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(seat)
}

// ====================================================
// This Method Retrive all Theaters from Block Chain Network
// ====================================================
func (t *MovieChaincode) getTheaters(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryString := `
	{
		"selector": {
			"docType": {
				"$eq": 1
			}
		}
	}
	`
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	if err != nil {
		return shim.Error(err.Error())
	}

	theators := []string{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println(queryResponse)

		movie := Movie{}
		json.Unmarshal(queryResponse.Value, &movie)

		theators = append(theators, movie.TheaterName)

	}

	theators_bytes, err := json.Marshal(theators)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(theators_bytes)
}

// ====================================================
// This Main Method to intialte
// ====================================================
func main() {

	err := shim.Start(new(MovieChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}
