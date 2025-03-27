package layout

import (
	"net/http"

	"github.com/allan-lewis/no-geeks-brewing-go/batches"
	"github.com/allan-lewis/no-geeks-brewing-go/oauth"
)

func LayoutHandler(w http.ResponseWriter, r *http.Request) {
	user := oauth.User(r)
	batchComponent := batches.BatchesComponent(user)
	authComponent := oauth.AuthComponent(user)

	err := LayoutComponent(batchComponent, authComponent).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering base layout", http.StatusInternalServerError)
		return
	}

}
