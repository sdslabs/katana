package leaderboardservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/gofiber/fiber/v2"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
)

func Leaderboard(c *fiber.Ctx, str string) error {
	var sortedData []types.CTFTeam
	if str == "" {
		data := mongo.FetchDocs(context.Background(), "teams", bson.M{})
		sort.Slice(data, func(i, j int) bool {
			return data[i]["score"].(int32) > data[j]["score"].(int32)
		})
		for _, team := range data {
			var temp types.CTFTeam
			bsonBytes, _ := bson.Marshal(team)
			bson.Unmarshal(bsonBytes, &temp)
			sortedData = append(sortedData, temp)
		}
	} else {
		tick, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return c.SendString(fmt.Sprintf("%s is not an integer", str))
		}
		filename := fmt.Sprintf("./json_data/data-tick-%d.json", tick)
		if _, err := os.Stat(filename); err != nil {
			return c.SendString("Tick hasn't occured yet. Please ask for a tick that has occured")
		} else {
			content, err := os.ReadFile(filename)
			if err != nil {
				return c.SendString("Error while reading file")
			}
			type temp struct {
				Data []types.CTFTeam `json:"data"`
			}
			var tmp temp
			if err := json.Unmarshal(content, &tmp); err != nil {
				return c.SendString("Error in unmarshaling json")
			}
			sortedData = append(sortedData, tmp.Data...)
			sort.Slice(sortedData, func(i, j int) bool {
				return sortedData[i].Score > sortedData[j].Score
			})
		}
	}
	var challenges []string
	var teams []string
	for _, challenge := range sortedData[0].Challenges {
		challenges = append(challenges, challenge.ChallengeName)
	}
	for _, team := range sortedData {
		teams = append(teams, team.Name)
	}
	columns := len(challenges) + 2
	rows := len(teams) + 1
	tempDisplayData := make([][]string, rows)
	for i := range tempDisplayData {
		tempDisplayData[i] = make([]string, columns)
	}
	for i := 0; i < columns-2; i++ {
		tempDisplayData[0][i+1] = challenges[i]
	}
	for i := 0; i < rows-1; i++ {
		tempDisplayData[i+1][0] = teams[i]
	}
	tempDisplayData[0][columns-1] = "Score"

	for i := 0; i < rows-1; i++ {
		for j := 0; j < columns-2; j++ {
			tempDisplayData[i+1][j+1] = fmt.Sprintf("Attacks: %d, Defences: %d, Uptime %f", sortedData[i].Challenges[j].Attacks, sortedData[i].Challenges[j].Defences, sortedData[i].Challenges[j].Uptime)
		}
		tempDisplayData[i+1][columns-1] = strconv.Itoa(sortedData[i].Score)
	}
	return c.SendString(createMatrixString(tempDisplayData))
}

func createMatrixString(matrix [][]string) string {
	var result string
	buf := new(bytes.Buffer)
	defer buf.Reset()
	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)

	for _, row := range matrix {
		for _, value := range row {
			fmt.Fprintf(w, "%s\t", value)
		}
		fmt.Fprintln(w, "")
	}

	w.Flush()

	result = buf.String()

	return result
}
