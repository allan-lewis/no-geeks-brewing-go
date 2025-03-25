package index

import (
	"net/http"
	"github.com/allan-lewis/no-geeks-brewing-go/batches"
)

// Serve the index page with embedded JS for live reload
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Create a slice to hold the values from the map
	var values []batches.Batch

	// Iterate over the map and append the values to the slice
	// for _, value := range batchesMap {
	// 	values = append(values, value)
	// }

	// sort.Slice(values, func(i, j int) bool {
	// 	return values[i].Number < values[j].Number
	// })

	err := Index(batches.Batches(values)).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
