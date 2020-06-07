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
	r.HandleFunc("/home", pos.Home)

	r.HandleFunc("/menu/create", pos.MenuCreate)
	r.HandleFunc("/menu", pos.Menu)
	r.HandleFunc("/menu/{id:[0-9]+}", pos.MenuGet)
	r.HandleFunc("/menu/{id:[0-9]+}/delete", pos.MenuDelete)

	r.HandleFunc("/orders", pos.Orders)
	r.HandleFunc("/order/{id:[0-9]+}", pos.OrderGet)
	r.HandleFunc("/orders/create", pos.OrderCreate)
}

func (pos PosEngine) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		viewPage := "view/login.html"
		assetsUrl := "http://" + r.Host + "/assets/"

		t, _ := template.ParseFiles(viewPage)

		data := map[string]interface{}{
			"assets": assetsUrl,
		}

		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "login", data)
	} else {
		r.ParseForm()

		if r.FormValue("username") == "admin" && r.FormValue("password") == "admin" {
			http.Redirect(w, r, "http://"+r.Host+"/home", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "http://"+r.Host, http.StatusSeeOther)
		}
	}

}

func (pos PosEngine) Home(w http.ResponseWriter, r *http.Request) {
	viewPage := "view/home.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets": assetsUrl,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "home", data)
}
