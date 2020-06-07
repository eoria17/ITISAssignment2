package pos

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/tiuriandy/ITISAssignment2/config"
	"github.com/tiuriandy/ITISAssignment2/model"
)

func (pos PosEngine) MenuDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	record := model.Menu{}
	pos.Storage.DB.First(&record, id)

	creds := credentials.NewStaticCredentials(config.AWS_KEY, config.AWS_SECRET, config.AWS_TOKEN)
	creds.Get()

	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)

	svc := s3.New(session.New(), cfg)

	x := strings.Split(record.ImageURL, "/")
	deletePath := "/media/" + x[len(x)-1]
	delParams := &s3.DeleteObjectInput{
		Bucket: aws.String(config.AWS_BUCKET_NAME),
		Key:    aws.String(deletePath),
	}

	_, err := svc.DeleteObject(delParams)
	if err != nil {
		fmt.Println(err)
	}

	pos.Storage.DB.Delete(&record)

	http.Redirect(w, r, "http://"+r.Host+"/menu", http.StatusSeeOther)
}

func (pos PosEngine) MenuGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	record := model.Menu{}
	pos.Storage.DB.First(&record, id)

	if r.Method == "GET" {

		viewPage := "view/getMenu.html"
		assetsUrl := "http://" + r.Host + "/assets/"

		t, _ := template.ParseFiles(viewPage)

		data := map[string]interface{}{
			"assets": assetsUrl,
			"Menu":   record,
		}

		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "getMenu", data)
	} else {
		//upload new pic
		//AWS S3 connection
		creds := credentials.NewStaticCredentials(config.AWS_KEY, config.AWS_SECRET, config.AWS_TOKEN)
		creds.Get()

		cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)

		svc := s3.New(session.New(), cfg)

		file, fileheader, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		size := fileheader.Size

		r.ParseMultipartForm(size)

		buffer := make([]byte, size)
		file.Read(buffer)

		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)

		path := "/media/" + fileheader.Filename
		params := &s3.PutObjectInput{
			Bucket:        aws.String(config.AWS_BUCKET_NAME),
			Key:           aws.String(path),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
		}

		_, err = svc.PutObject(params)
		if err != nil {
			fmt.Println(err)
		}

		//create record

		price, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			fmt.Println(err)
		}

		url := "https://%s.s3.amazonaws.com%s"
		url = fmt.Sprintf(url, config.AWS_BUCKET_NAME, path)

		x := strings.Split(record.ImageURL, "/")
		deletePath := "/media/" + x[len(x)-1]
		delParams := &s3.DeleteObjectInput{
			Bucket: aws.String(config.AWS_BUCKET_NAME),
			Key:    aws.String(deletePath),
		}

		_, err = svc.DeleteObject(delParams)
		if err != nil {
			fmt.Println(err)
		}

		record.Price = float64(price)
		record.ImageURL = url

		pos.Storage.DB.Save(&record)

		http.Redirect(w, r, "http://"+r.Host+"/menu", http.StatusSeeOther)
	}

}

func (pos PosEngine) Menu(w http.ResponseWriter, r *http.Request) {
	menus := []model.Menu{}
	pos.Storage.DB.Find(&menus)

	viewPage := "view/menus.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets": assetsUrl,
		"Menus":  menus,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "menus", data)
}

func (pos PosEngine) MenuCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		viewPage := "view/newMenus.html"
		assetsUrl := "http://" + r.Host + "/assets/"

		t, _ := template.ParseFiles(viewPage)

		data := map[string]interface{}{
			"assets": assetsUrl,
		}

		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "new_menus", data)
	} else {

		//AWS S3 connection
		creds := credentials.NewStaticCredentials(config.AWS_KEY, config.AWS_SECRET, config.AWS_TOKEN)
		creds.Get()

		cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)

		svc := s3.New(session.New(), cfg)

		file, fileheader, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		size := fileheader.Size

		r.ParseMultipartForm(size)

		buffer := make([]byte, size)
		file.Read(buffer)

		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)

		path := "/media/" + fileheader.Filename
		params := &s3.PutObjectInput{
			Bucket:        aws.String(config.AWS_BUCKET_NAME),
			Key:           aws.String(path),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
		}

		_, err = svc.PutObject(params)
		if err != nil {
			fmt.Println(err)
		}

		//create record

		price, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			fmt.Println(err)
		}

		url := "https://%s.s3.amazonaws.com%s"
		url = fmt.Sprintf(url, config.AWS_BUCKET_NAME, path)

		menu := model.Menu{
			Name:     r.FormValue("name"),
			Price:    float64(price),
			ImageURL: url,
		}

		pos.Storage.DB.Create(&menu)
		http.Redirect(w, r, "http://"+r.Host+"/home", http.StatusSeeOther)

	}
}
