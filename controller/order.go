package pos

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tiuriandy/ITISAssignment2/model"
)

func (pos PosEngine) OrderCreate(w http.ResponseWriter, r *http.Request) {
	menus := []model.Menu{}
	pos.Storage.DB.Find(&menus)

	if r.Method == "GET" {
		viewPage := "view/newOrder.html"
		assetsUrl := "http://" + r.Host + "/assets/"

		t, _ := template.ParseFiles(viewPage)

		data := map[string]interface{}{
			"assets": assetsUrl,
			"Menus":  menus,
		}

		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "new_order", data)
	} else {
		r.ParseForm()

		db := pos.Storage.DB.Begin()
		defer db.RollbackUnlessCommitted()

		record := model.Order{}

		currentTime := time.Now()
		record.Date = currentTime

		hours, minutes, seconds := time.Now().Clock()
		record.OrderNumber = "#" + currentTime.Format("0601")
		record.OrderNumber = record.OrderNumber + fmt.Sprintf("%d%02d%02d", hours, minutes, seconds)

		total := 0.0
		temp := 0
		for i := 0; i < len(menus); i++ {
			temp, _ = strconv.Atoi(r.FormValue(strconv.Itoa(menus[i].ID) + "-amount"))
			temp = temp * int(menus[i].Price)
			total += float64(temp)
			temp = 0
		}
		record.Total = total
		db.Create(&record)

		amount := 0
		recordLines := model.OrderLine{}
		for i := 0; i < len(menus); i++ {
			amount, _ = strconv.Atoi(r.FormValue(strconv.Itoa(menus[i].ID) + "-amount"))
			if amount == 0 {
				continue
			} else {
				recordLines = model.OrderLine{
					OrderID:  record.ID,
					MenuID:   menus[i].ID,
					Amount:   amount,
					Subtotal: float64(amount) * menus[i].Price,
				}
				db.Create(&recordLines)
			}
		}

		db.Commit()

		http.Redirect(w, r, "http://"+r.Host+"/home", http.StatusSeeOther)

	}
}

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
