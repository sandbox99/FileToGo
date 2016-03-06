package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main(){
	http.HandleFunc("/upload",doUpload)
	http.HandleFunc("/",showUploadPage)
	
	fmt.Println("Will listen at port 9999")
	http.ListenAndServe(":9999",nil)
	
}

func showUploadPage(res http.ResponseWriter, request *http.Request){
	
}

func doUpload(res http.ResponseWriter, request *http.Request){
	fmt.Println("call doUpload...")
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	
	file,header,err:=request.FormFile("file1")
	if err!=nil{
		fmt.Println(err)
		
	}else{
		fmt.Println(header)
		destfilename:=header.Filename
		destfile,err:=os.Create(destfilename)
		if err!=nil{
			fmt.Println("create file failed:",err,destfilename)
			io.WriteString(res,"bad occurred")
			return
		}
		defer destfile.Close()
		cnt,err:=io.Copy(destfile,file)
		if err!=nil{
			fmt.Println("Write file failed:",err)
			io.WriteString(res,"bad occurred")
			return
		}
		fmt.Println("Written ",cnt)
		
	}
	
	io.WriteString(res,
		`<DOCTYPE html>
<html>
	<head>
		<title>Hello</title>
	</head>
	<body>
		Hello World!
		<form action="upload" method="POST" target="fresult" enctype="multipart/form-data">
			<input type="file" name="file1">
			<input type="submit">
		</form>
		<iframe id="fresult"></iframe>
	</body>
</html>`,
	)
}





