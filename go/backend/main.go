package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func initDB() error {
	os.Mkdir("./db", 0755)
	os.Mkdir("./db/apps", 0755)
	return nil
}

func saveArtifact(f *os.File) error {
	return nil
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST")
		return
	}

	log.Println("Upload received! Processing started...")
	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Sending upload to worker...")
	resp, err := http.Post(
		"https://dubhacks-worker-1.ngrok.com/create_compile_ipa",
		"application/octet-stream", file)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Received IPA build artifact from worker...")

	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		log.Println(err)
		return
	}

	sum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))

	log.Println("Saving build artifact...")
	if err := os.Mkdir(fmt.Sprintf("./db/apps/%s", sum), 0755); err != nil {
		if err == os.ErrExist {
			log.Println("App already exists!")
			return
		}
		log.Println(err)
		return
	}
	artifact, err := os.Create(fmt.Sprintf("./db/apps/%s/ios.tar.gz", sum))
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := io.Copy(artifact, bytes.NewBuffer(buf.Bytes())); err != nil {
		log.Println(err)
		return
	}

	log.Println("Sending IPA ID to worker for APK generation...")
	resp, err = http.PostForm(
		"https://dubhacks-worker-1.ngrok.com/create_generate_apk",
		url.Values{"ipa_id": {sum}, "name": {"Andrios"}})
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Saving output to ./db/apps/%s/android.apk\n", sum)
	f, err := os.Create(fmt.Sprintf("./db/apps/%s/android.apk", sum))
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Println(err)
		return
	}
}

const frameHtml = `<!DOCTYPE html>
<head>
<title>Andrios</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
body, html {
	margin: 0;
	padding: 0;
	height: 100%;
}
</style>
</head>
<body>
<iframe src="{{.}}" style="height: 100%; width: 100%;">
</iframe>
</body>
`

func frame(w http.ResponseWriter, r *http.Request) {
	//u := r.URL.Query().Get("url")
	tmpl, err := template.New("frame").Parse(frameHtml)
	if err != nil {
		log.Println(err)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "frame", "http://173.250.177.106:8080/guacamole/client.xhtml?id=c%2Fdubhacks"); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := initDB(); err != nil {
		panic(err)
	}

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/frame", frame)
	http.ListenAndServe(":"+port, nil)
}
