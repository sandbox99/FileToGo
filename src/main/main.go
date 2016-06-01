package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var ipexp *regexp.Regexp

func init() {
	ipexp = regexp.MustCompile("([0-9]+\\.){3}\\d+")
}

func main() {

	http.HandleFunc("/upload", doUpload)
	http.HandleFunc("/blank", showBlankPage)
	http.HandleFunc("/", showUploadPage)

	port := 9999
	addrs, err := getLocalAddresses()
	if err != nil {
		log.Printf("error getting local IP addresses, just ignore: %v\n", err)
		log.Printf("will listen at port %d", port)
	} else {
		log.Printf("will listen at the following addresses:\n")
		for _, addr := range addrs {
			log.Printf("  http://%s:%d/", addr, port)
		}
	}

	log.Fatal(http.ListenAndServe(":9999", nil))
	log.Println("ok")
}

func getLocalAddresses() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, addr := range addrs {
		t := ipexp.FindString(addr.String())
		if t != "" {
			res = append(res, t)
		}
	}
	return res, nil
}

func readPageContent(res http.ResponseWriter) (int64, error) {
	file, err := os.Open("resource/index.html")
	if err != nil {
		panic("read file failed.")
		return 0, err
	}
	defer file.Close()
	return io.Copy(res, file)
}

func showUploadPage(res http.ResponseWriter, request *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(res, UPLOAD_PAGE_CONT)

	//we can read page content from resource/index.html when develop
	//readPageContent(res)
}

func showBlankPage(res http.ResponseWriter, request *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(res, wrapSimpleHtml(""))
}

func doUpload(res http.ResponseWriter, request *http.Request) {
	log.Println("========call doUpload========")

	res.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	if request.Method == "GET" {
		io.WriteString(res, wrapSimpleHtml("url accepts method POST only"))
		return
	}

	file, header, err := request.FormFile("file1")
	if err != nil {
		fmt.Println("Error while call FormFile", err)
		io.WriteString(res, wrapSimpleHtml("发生了错误"))
		return
	} else {
		destfilename := header.Filename
		log.Println("File name is:", destfilename)
		destfile, err := os.Create(destfilename)
		if err != nil {
			log.Fatalln("create file failed:", err, destfilename)
			io.WriteString(res, wrapSimpleHtml("发生了错误:"+err.Error()))
			return
		}
		defer destfile.Close()
		cnt, err := io.Copy(destfile, file)
		if err != nil {
			log.Fatalln("Error while call Copy:", err)
			io.WriteString(res, wrapSimpleHtml("发生了错误"))
			return
		}
		log.Println("Written", cnt, "bytes")
		io.WriteString(res, wrapSimpleHtml(fmt.Sprintf("成功(%d bytes)", cnt)))
	}
}

func wrapSimpleHtml(str string) string {
	var template string = `<!DOCTYPE html>
<html>
	<head>
		<title></title>
		<style type="text/css">body{font-family:Microsoft Yahei;}</style>
	</head>
	<body>MSG</body>
</html>`
	return strings.Replace(template, "MSG", str, 1)
}

//value from file resource/index.html
const UPLOAD_PAGE_CONT = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>File Transfer Helper</title>
		<style type="text/css">
			body,input{font-family:Microsoft Yahei;}
			.phead{overflow:hidden;}.phead h2{float:left;}.phead h5{float: right;padding: 8px 0 0;}
			.line{padding:10px;border-bottom:1px solid lightgreen;}
			.first.line{border-top:1px solid lightgreen;}
			form{float:left;}.container{width:800px;margin:20px auto;}
			iframe{height:23px;margin-left: 20px;width:430px;}
			input#btnAdd{margin: 10px 0 0 13px;}
		</style>
		<script type="text/javascript">
			var idx=2;
			function doInit(){
				document.getElementById("btnAdd").onclick=function(){
					appendAnother();
				};
			}
			function appendAnother(){
				idx+=1;
				var frag=document.createDocumentFragment();
				var line=document.createElement("div");
				line.className="line";frag.appendChild(line);
				var form=document.createElement("form");
				line.appendChild(form);
				form.action="upload";form.method="POST";form.target="fresult"+idx;form.enctype="multipart/form-data";
				var input1=document.createElement("input"),input2=document.createElement("input");
				input1.type="file";input1.name="file1";input2.type="submit";				
				form.appendChild(input1);form.appendChild(input2);
				var mframe=document.createElement("iframe");
				line.appendChild(mframe);mframe.id="fresult"+idx;
				mframe.name="fresult"+idx;mframe.src="blank";mframe.frameBorder=0;mframe.scrolling="no";mframe.marginHeight="0";
				mframe.marginWidth="0";
				var btn=document.getElementById("btnAdd");
				btn.parentNode.insertBefore(frag,btn);
			}
		</script>
	</head>
	<body onload="doInit()">
		<div class="container">
			<div class="phead">
				<h2>FileToGo</h2>
				<h5>Simple and fast file transfer helper, <a href="http://github.com/sandbox99/FileToGo">fork on Github</a></h5>
			</div>
			<div class="line first">
				<form action="upload" method="POST" target="fresult" enctype="multipart/form-data">
					<input type="file" name="file1"><input type="submit">
				</form>
				<iframe name="fresult" id="fresult" src="blank" frameborder="0" scrolling="no" marginheight="0" marginwidth="0"></iframe>
			</div>
			<div class="line">
				<form action="upload" method="POST" target="fresult2" enctype="multipart/form-data">
					<input type="file" name="file1"><input type="submit">
				</form>
				<iframe name="fresult2" id="fresult2" src="blank" frameborder="0" scrolling="no" marginheight="0" marginwidth="0"></iframe>
			</div>
			<input type="button" id="btnAdd" value="添加">
		</div>
	</body>
</html>`
