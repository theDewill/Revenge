package Arc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Arc_Handler struct{}

func (AH *Arc_Handler) FetchServices() interface{} {

	//Later implement to chunk the details in the services.json and use it
	fmt.Println("Fetcing Services....")
	Sfile, err := os.Open("services.json")

	if err != nil {
		fmt.Print("error opening the file")
	}
	defer Sfile.Close()

	fileBytes, err := ioutil.ReadAll(Sfile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %s", err)
	}

	var result map[string]interface{}
	json.Unmarshal(fileBytes, &result)

	fmt.Print(result)

	return "nothing"

}
