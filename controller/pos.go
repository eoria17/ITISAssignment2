package pos

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tiuriandy/ITISAssignment2/storage"
)

type PosEngine struct {
	Storage *storage.Storage
}

func (pos PosEngine) Route(r *mux.Router) {
	r.HandleFunc("/", pos.Login)
}

func (pos PosEngine) Login(w http.ResponseWriter, r *http.Request) {
	viewPage := "view/index.html"

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{}

	t.ExecuteTemplate(w, "login", data)
	w.WriteHeader(http.StatusOK)
}
