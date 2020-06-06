package pos

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tiuriandy/ITISAssignment2/config"
	"github.com/tiuriandy/ITISAssignment2/model"
)

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

		_, err := svc.PutObject(params)
		if err != nil {
			fmt.Println(err)
		}

		//create record

		price, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			fmt.Println(err)
		}

		url := "https://%s.amazonaws.com/%s/%s"
		url = fmt.Sprintf(url, "us-east-1", config.AWS_BUCKET_NAME, fileheader.Filename)

		menu := model.Menu{
			Name:  r.FormValue("name"),
			Price: float64(price),
			URL:   url,
		}

		pos.Storage.DB.Create(&menu)

	}
}
