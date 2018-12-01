/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the land structure, with 4 properties.  Structure tags are used by encoding/json library
type Land struct {
	RLRegistry   string `json:"rlregistry"`
	Extent  int `json:"extent"`
	ParentLandID string `json:"parentlandid"`
	Owner  string `json:"owner"`
	Boundaries [4][2]int `json:"boundaries"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	lands := []Land{
		Land{RLRegistry: "Colombo", Extent: 50, ParentLandID: "nil", Owner: "Tomoko", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Delkanda", Extent: 25, ParentLandID: "nil", Owner: "Brad", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Avissawella", Extent: 75, ParentLandID: "nil", Owner: "Jin Soo", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Homagama", Extent: 40, ParentLandID: "nil", Owner: "Max", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Galle", Extent: 30, ParentLandID: "nil", Owner: "Adriana", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Balapitiya", Extent: 35, ParentLandID: "nil", Owner: "Michel", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Elpitiya", Extent: 45, ParentLandID: "nil", Owner: "Aarav", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Hambantota", Extent: 20, ParentLandID: "nil", Owner: "Pari", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Tangalle", Extent: 25, ParentLandID: "nil", Owner: "Valeria", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Colombo", Extent: 60, ParentLandID: "nil", Owner: "Shotaro", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
	}

	i := 0
	for i < len(lands) {
		fmt.Println("i is ", i)
		landAsBytes, _ := json.Marshal(lands[i])
		APIstub.PutState("LAND"+strconv.Itoa(i), landAsBytes)
		fmt.Println("Added", lands[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryLand" {
		return s.queryLand(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createLand" {
		return s.createLand(APIstub, args)
	} else if function == "queryAllLands" {
		return s.queryAllLands(APIstub)
	} else if function == "changeLandOwner" {
		return s.changeLandOwner(APIstub, args)
	} else if function == "delete" { //obtained from delete a marble
		return s.delete(APIstub, args)
	} else if function == "getHistoryForLand" { //obtained from delete a marble
                return s.getHistoryForLand(APIstub, args)
        } else if function == "forkLand" {
		return s.forkLand(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryLand(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	resultsIterator, err := APIstub.GetHistoryForKey(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}else if response.IsDelete {
			return shim.Error("This land has been deleted")
		}
	}
	landAsBytes, _ := APIstub.GetState(args[0])
        return shim.Success(landAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	lands := []Land{
		Land{RLRegistry: "Colombo", Extent: 50, ParentLandID: "nil", Owner: "Tomoko", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Delkanda", Extent: 25, ParentLandID: "nil", Owner: "Brad", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Avissawella", Extent: 75, ParentLandID: "nil", Owner: "Jin Soo", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Homagama", Extent: 40, ParentLandID: "nil", Owner: "Max", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Galle", Extent: 30, ParentLandID: "nil", Owner: "Adriana", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Balapitiya", Extent: 35, ParentLandID: "nil", Owner: "Michel", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Elpitiya", Extent: 45, ParentLandID: "nil", Owner: "Aarav", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Hambantota", Extent: 20, ParentLandID: "nil", Owner: "Pari", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Tangalle", Extent: 25, ParentLandID: "nil", Owner: "Valeria", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
		Land{RLRegistry: "Colombo", Extent: 60, ParentLandID: "nil", Owner: "Shotaro", Boundaries: [4][2]int{{0,20},{10,20},{10,0},{0,0}}},
	}

	i := 0
	for i < len(lands) {
		fmt.Println("i is ", i)
		landAsBytes, _ := json.Marshal(lands[i])
		APIstub.PutState("LAND"+strconv.Itoa(i), landAsBytes)
		fmt.Println("Added", lands[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createLand(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	size, err := strconv.Atoi(args[2])
        if err != nil {
                return shim.Error("3rd argument must be a numeric string")
        }

	var land = Land{RLRegistry: args[1], Extent: size, ParentLandID: args[3], Owner: args[4]}

	resultsIterator, err := APIstub.GetHistoryForKey(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}else if response.IsDelete {
			return shim.Error("This land has been deleted")
		}
	}

	landAsBytes, _ := json.Marshal(land)
        APIstub.PutState(args[0], landAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) queryAllLands(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "LAND0"
	endKey := "LAND999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	fmt.Printf("- queryAllLands:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeLandOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//from here
	resultsIterator, err := APIstub.GetHistoryForKey(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}else if response.IsDelete {
			return shim.Error("This land has been deleted")
		}else {
			landAsBytes, _ := APIstub.GetState(args[0])
		        land := Land{}

		        json.Unmarshal(landAsBytes, &land)
		        land.Owner = args[1]

		        landAsBytes, _ = json.Marshal(land)
		        APIstub.PutState(args[0], landAsBytes)
		}
	}

	return shim.Success(nil)
	//to here

}

func (s *SmartContract) delete(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var jsonResp string
	var landJSON Land

	if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

	valAsbytes, err := APIstub.GetState(args[0]) //get the marble from chaincode state
        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get state for land" + args[0] + "\"}"
                return shim.Error(jsonResp)
        } else if valAsbytes == nil {
                jsonResp = "{\"Error\":\"Land does not exist: " + args[0] + "\"}"
                return shim.Error(jsonResp)
        }

	err = json.Unmarshal([]byte(valAsbytes), &landJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	err = APIstub.DelState(args[0]) //remove the land from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) forkLand(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}
	//obtain rlregistry, extent, parentlandid of parent land
	var jsonResp string
	var landJSON Land

	valAsbytes, err := APIstub.GetState(args[0]) //get the marble from chaincode state
        if err != nil {
                jsonResp = "{\"Error\":\"Failed to get state for land" + args[0] + "\"}"
                return shim.Error(jsonResp)
        } else if valAsbytes == nil {
                jsonResp = "{\"Error\":\"Land does not exist: " + args[0] + "\"}"
                return shim.Error(jsonResp)
        }

	err = json.Unmarshal([]byte(valAsbytes), &landJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	size1, err := strconv.Atoi(args[3])
        if err != nil {
                return shim.Error("4th argument must be a numeric string")
        }

	size2, err := strconv.Atoi(args[6])
        if err != nil {
                return shim.Error("7th argument must be a numeric string")
        }

	xdivcentre, err := strconv.Atoi(args[8])
	if err != nil {
                return shim.Error("9th argument must be a numeric string")
        }

        ydivcentre, err := strconv.Atoi(args[9])
        if err != nil {
                return shim.Error("10th argument must be a numeric string")
        }

	if (size1+size2>landJSON.Extent){
		return shim.Error("Sum of extents of partitioned lands is greater than the extent of original land")
	}else{
		if (args[7]=="v"){
			//create 2 child lands
			var childland1 = Land{RLRegistry: landJSON.RLRegistry, Extent: size1, ParentLandID: args[0], Owner: args[2], Boundaries: [4][2]int{{landJSON.Boundaries[0][0],landJSON.Boundaries[0][1]},{xdivcentre,landJSON.Boundaries[0][1]},{xdivcentre,landJSON.Boundaries[3][1]},{landJSON.Boundaries[3][0],landJSON.Boundaries[3][1]}}}
			landAsBytes1, _ := json.Marshal(childland1)
		        APIstub.PutState(args[1], landAsBytes1)

			var childland2 = Land{RLRegistry: landJSON.RLRegistry, Extent: size2, ParentLandID: args[0], Owner: args[5], Boundaries: [4][2]int{{xdivcentre,landJSON.Boundaries[0][1]},{landJSON.Boundaries[1][0],landJSON.Boundaries[1][1]},{landJSON.Boundaries[2][0],landJSON.Boundaries[2][1]},{xdivcentre,landJSON.Boundaries[3][1]}}}
		        landAsBytes2, _ := json.Marshal(childland2)
		        APIstub.PutState(args[4], landAsBytes2)
		}else if (args[7]=="h"){
			//create 2 child lands
			var childland1 = Land{RLRegistry: landJSON.RLRegistry, Extent: size1, ParentLandID: args[0], Owner: args[2], Boundaries: [4][2]int{{landJSON.Boundaries[0][0],landJSON.Boundaries[0][1]},{landJSON.Boundaries[1][0],landJSON.Boundaries[1][1]},{landJSON.Boundaries[1][0],ydivcentre},{landJSON.Boundaries[3][0],ydivcentre}}}
			landAsBytes1, _ := json.Marshal(childland1)
		        APIstub.PutState(args[1], landAsBytes1)

			var childland2 = Land{RLRegistry: landJSON.RLRegistry, Extent: size2, ParentLandID: args[0], Owner: args[5], Boundaries: [4][2]int{{landJSON.Boundaries[0][0],ydivcentre},{landJSON.Boundaries[1][0],ydivcentre},{landJSON.Boundaries[2][0],landJSON.Boundaries[2][1]},{landJSON.Boundaries[3][0],landJSON.Boundaries[3][1]}}}
		        landAsBytes2, _ := json.Marshal(childland2)
		        APIstub.PutState(args[4], landAsBytes2)
		}
	}

	//delete parent land
	err = APIstub.DelState(args[0]) //remove the land from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)

}

func (s *SmartContract) getHistoryForLand(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	landName := args[0]

	fmt.Printf("- start getHistoryForLand: %s\n", landName)

	resultsIterator, err := APIstub.GetHistoryForKey(landName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForLand returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
