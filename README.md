# FileToGo

This is a small but practically useful tool for setting up a web server on a machine in a few seconds, and then you can 
visit its uploading page from other machines to transfer files to it, no needing resorting to ftp, sftp, rsync, QQ, etc.

Although it was built with Go, you can run it without Go environment.

### Installation

You can get the pre-built binary file directly and are ready to run it, on Windows for example:

	curl https://raw.githubusercontent.com/sandbox99/FileToGo/master/bin/fileToGo-win64.exe
	
or you can pull down the source files using your favorite <code>go</code> tool and build <code>main.go</code> by yourself.

### Usage

Start the executable file on the machine you want transfer files to:

	fileToGo-win64.exe (or fileToGo on Unix system)

Once it started up successfully, you can upload files from other machines. Navigate your browser to the url <code>http://<server ip>:9999/</code>

![run demo](./src/main/resource/demo.jpg?raw=true =600x)

Users in favor of shell can upload file through curl as well:

	curl -F file1=@my.txt http://<server ip>:9999/upload

