package layout

import (
	"net/http"

	"github.com/allan-lewis/no-geeks-brewing-go/batches"
	"github.com/allan-lewis/no-geeks-brewing-go/oauth"
)

func LayoutHandler(w http.ResponseWriter, r *http.Request) {
	err := LayoutComponent(batches.BatchesComponent(), oauth.AuthComponent()).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}
