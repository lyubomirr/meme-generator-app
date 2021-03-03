package web

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
)

func (s *apiServer) test (w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", header.Filename)
	fmt.Printf("File Size: %+v\n", header.Size)
	fmt.Printf("MIME Header: %+v\n", header.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fileBytes)
}
