package batches

type Batch struct {
	ID       string    `json:"batchId"`
	Name     string    `json:"batchName"`
	Style    string    `json:"batchStyle"`
	Number   int       `json:"batchNumber"`
	Date     int64	   `json:"batchDate"`
	Status   string    `json:"batchStatus"`
	URL      string    `json:"batchUrl"`
}
