package retrieve_data

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var ticker *time.Ticker
var tick int

func saveJSON() {
	for range ticker.C {
		data := mongo.FetchDocs(context.Background(), "teams", bson.M{})
		path := fmt.Sprintf("./json_data/data-tick-%d.json", tick)
		jsonData, err := convertBSONArrayToJSONArray(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		storeJSONToFile(jsonData, path)
		tick++
	}
}

func StartSaving() {
	ticker = utils.GetTicker()
	utils.CreateDirectoryIfNotExists("json_data")
	go saveJSON()
}

func convertBSONArrayToJSONArray(bsonArray []bson.M) ([]byte, error) {
	var jsonArray []map[string]interface{}

	for _, bsonDoc := range bsonArray {
		delete(bsonDoc, "publicKey")
		delete(bsonDoc, "password")

		jsonArray = append(jsonArray, bsonDoc)
	}

	jsonData, err := json.Marshal(jsonArray)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON array: %w", err)
	}

	return jsonData, nil
}

func storeJSONToFile(jsonData []byte, filePath string) error {
	jsonMap := map[string]interface{}{"data": json.RawMessage(jsonData)}

	finalJSONData, err := json.MarshalIndent(jsonMap, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal final JSON: %w", err)
	}

	err = os.WriteFile(filePath, finalJSONData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	fmt.Printf("JSON data stored in file: %s\n", filePath)
	return nil
}
