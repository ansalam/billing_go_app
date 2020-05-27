package main

import (
	"billing_api/pkg/utils"
	"encoding/json"
	"net/http"

	"rsc.io/pdf"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	w.Write([]byte("Hello! It's working!!!\n"))
}

func (app *application) uploadFile(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Upload request hit.")
	if r.Method == "POST" {
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		scanID := r.Form.Get("scan-id")
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			app.infoLog.Println(err)
			app.clientError(w, 400)
		}
		defer file.Close()

		tempFileName, err := utils.SaveUploadedFile(&file, handler)
		if err != nil {
			app.serverError(w, err)
		}
		go app.processRequest(tempFileName, scanID)
		w.Write([]byte("File Upload Successful!!"))
	} else {
		app.clientError(w, 405)
	}
}

func (app *application) processRequest(tempFileName string, scanID string) {
	app.infoLog.Printf("Request Data Processor started for scanID %s\n", scanID)
	PDFFile, err := pdf.Open(tempFileName)
	if err != nil {
		app.infoLog.Println(err)
	}
	_, err = app.DBConn.Insert(app.authenticatorID, scanID, PDFFile.NumPage())
	if err != nil {
		app.infoLog.Println("Insertion failed")
		app.infoLog.Println(err)
		return
	}
	app.infoLog.Println("DB insertion successfull.")
	err = utils.RemoveUploadedFile(tempFileName)
	if err != nil {
		app.infoLog.Println(err)
	}
	app.infoLog.Println("Request Data processor completed successfull.")
}

func (app *application) GetCount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		app.infoLog.Println("GetCount Started.")
		counts, err := app.DBConn.GetCounts(app.authenticatorID)
		if err != nil {
			app.infoLog.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(counts)
		app.infoLog.Println("GetCount successfully completed.")
	} else {
		app.clientError(w, 405)
	}
}
