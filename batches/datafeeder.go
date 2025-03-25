package batches

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func KickOff(channel chan Batch, authToken string) {
	time.Sleep(5 * time.Second)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	makeRequests(channel, authToken)

	for range ticker.C {
		makeRequests(channel, authToken)
	}
}

func makeRequests(channel chan Batch, authToken string) {
	var keepGoing bool = true
	var lastID string

	for keepGoing {
		lastID = makeRequest(lastID, channel, authToken)

		keepGoing = lastID != ""
	}
}

func makeRequest(lastID string, channel chan Batch, authToken string) string {
	var url = "https://api.brewfather.app/v2/batches?complete=true&limit=10&start_after=" + lastID

	log.Printf("Gettting batches %v", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return ""
	}

	req.Header.Set("Authorization", "Basic "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status: %d", resp.StatusCode)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return ""
	}

	if len(body) == 0 {
		log.Println("No data in the response body")
		return ""
	}

	var batches []brewfatherBatch

	err = json.Unmarshal(body, &batches)
	if err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		return ""
	}

	log.Printf("Read %v batch(es)", len(batches))

	if len(batches) == 0 {
		return ""
	}

	for _, b := range batches {
		var batch = Batch{
			ID:     b.ID,
			Name:   b.Name,
			Style:  b.Recipe.Style.Name,
			Number: b.BatchNumber,
			Date:   b.BrewDate,
			Status: b.Status,
			URL:    b.Share,
		}

		channel <- batch
	}

	return batches[len(batches)-1].ID
}

// BrewfatherBatch represents a batch with its associated fields
type brewfatherBatch struct {
	ID          string            `json:"_id"`
	Name        string            `json:"name"`
	BrewDate    int64             `json:"brewDate"` // Use pointer to handle null values
	Status      string            `json:"status"`
	Share       string            `json:"_share"`
	BatchNumber int               `json:"batchNo"`
	Recipe      *brewfatherRecipe `json:"recipe"`
}

// BrewfatherRecipe represents the recipe details of the batch
type brewfatherRecipe struct {
	Name  string           `json:"name"`
	Style *brewfatherStyle `json:"style"`
}

// BrewfatherStyle represents the style of the beer
type brewfatherStyle struct {
	Name string `json:"name"`
}
