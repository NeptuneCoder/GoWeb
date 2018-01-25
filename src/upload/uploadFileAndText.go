package upload

import (
	"net/http"
	"fmt"
)

func UploadFileAndText(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "POST" {
		fmt.Println("contentType",r.Header.Get("Content-type"))
		r.ParseMultipartForm(32 << 20);
		//data,err := ioutil.ReadAll(r.Body)
		//if err != nil{
		//	return
		//}
		fileName := r.MultipartForm.File["file[]"]
		fmt.Println("file",fileName)
		//if (len(request.Value) != 0) {
		//
		//
		//}
		//if (len(request.File) != 0) {
		//files := request.File
		//for k,v :=  range files  {
		//	fmt.Println("test upload file",k,v)
		//}

		//}
	}
}