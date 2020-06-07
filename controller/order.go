package pos

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tiuriandy/ITISAssignment2/model"
)

func (pos PosEngine) OrderGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	record := model.Order{}
	pos.Storage.DB.Preload("Lines").Preload("Lines.Menu").First(&record, id)

	viewPage := "view/getOrder.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets": assetsUrl,
		"Order":  record,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "getOrder", data)
}

func (pos PosEngine) Orders(w http.ResponseWriter, r *http.Request) {
	records := []model.Order{}
	pos.Storage.DB.Find(&records)

	viewPage := "view/orders.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets": assetsUrl,
		"Orders": records,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "orders", data)
}
