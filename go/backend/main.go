package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST")
		//w.Header().Add("Access-Control-Allow-Headers",
		//"Cache-Control, X-Requested-With")
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

	log.Println("Sending IPA ID to worker for APK generation...")
	resp, err = http.PostForm(
		"https://dubhacks-worker-1.ngrok.com/create_generate_apk",
		url.Values{"ipa_id": {sum}, "name": {"Example Name"}})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Output saved to output.apk")
	file, err = os.Create("output.apk")
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":"+port, nil)
}
