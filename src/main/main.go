package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"net/http"
)

func main(){
	
	http.HandleFunc("/upload",doUpload)
	http.HandleFunc("/",showUploadPage)
	http.HandleFunc("/blank",showBlankPage)
	
	fmt.Println("Will listen at port 9999")
	http.ListenAndServe(":9999",nil)
	
}

func readPageContent(res http.ResponseWriter) (int64,error){
	file,err:=os.Open("resource/index.html")
	if err!=nil{
		panic("read file failed.")
		return 0,err
	}
	defer file.Close()
	return io.Copy(res,file)
	
}

func showUploadPage(res http.ResponseWriter, request *http.Request){
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	
	readPageContent(res)
}

func showBlankPage(res http.ResponseWriter, request *http.Request){
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	
	io.WriteString(res,wrapSimpleHtml(""))
}

func doUpload(res http.ResponseWriter, request *http.Request){
	fmt.Println("========call doUpload========")
	
	res.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)
	
	if request.Method=="GET"{
		io.WriteString(res,wrapSimpleHtml("url accepts method POST only"))
		return
	}
	
	file,header,err:=request.FormFile("file1")
	if err!=nil{
		fmt.Println("Error while call FormFile",err)
		io.WriteString(res,wrapSimpleHtml("发生了错误"))
		return
	}else{
		
		destfilename:=header.Filename
		fmt.Println("File name is:",destfilename)
		destfile,err:=os.Create(destfilename)
		if err!=nil{
			fmt.Println("create file failed:", err, destfilename)
			io.WriteString(res,wrapSimpleHtml("发生了错误:"+err.Error()))
			return
		}
		defer destfile.Close()
		cnt,err:=io.Copy(destfile,file)
		if err!=nil{
			fmt.Println("Error while call Copy:",err)
			io.WriteString(res,wrapSimpleHtml("发生了错误"))
			return
		}
		fmt.Println("Written",cnt,"bytes")
		io.WriteString(res,wrapSimpleHtml(fmt.Sprintf("成功(%d bytes)",cnt)))
	}
}

func wrapSimpleHtml(str string)string{
	var template string=`<!DOCTYPE html>
<html>
	<head>
		<title></title>
		<style type="text/css">body{font-family:Microsoft Yahei;}</style>
	</head>
	<body>MSG</body>
</html>`
	return strings.Replace(template,"MSG",str,1)
}
