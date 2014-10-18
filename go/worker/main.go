package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"code.google.com/p/go-uuid/uuid"
)

func createCompileIPA(w http.ResponseWriter, r *http.Request) {
	log.Println("Compiling IPA started...")
	tmpDir := fmt.Sprintf("/tmp/%s", uuid.NewRandom().String())
	tmpWorkDir := fmt.Sprintf("%s/work", tmpDir)
	tmpFilename := fmt.Sprintf("%s/%s", tmpDir, uuid.NewRandom().String())

	if err := os.Mkdir(tmpDir, 0777); err != nil {
		log.Println(err)
		return
	}
	if err := os.Mkdir(tmpWorkDir, 0777); err != nil {
		log.Println(err)
		return
	}

	tmpFile, err := os.Create(tmpFilename)
	if err != nil {
		log.Println(err)
		return
	}
	defer tmpFile.Close()

	defer r.Body.Close()
	if _, err := io.Copy(tmpFile, r.Body); err != nil {
		log.Println(err)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	cmd := exec.Command(fmt.Sprintf("%s/compile_ipa.sh", cwd), tmpFilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = tmpWorkDir
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}
	log.Println("Compiling IPA finished and archive saved...")

	artifact, err := os.Open(fmt.Sprintf("%s/artifacts/artifact.tar.gz",
		tmpWorkDir))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Sending over the wire in response...")
	defer artifact.Close()
	if _, err := io.Copy(w, artifact); err != nil {
		log.Println(err)
		return
	}
}

func createCompileAPK(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting generation of APK...")
	tmpWorkDir := fmt.Sprintf("/tmp/%s", uuid.NewRandom().String())
	if err := os.Mkdir(tmpWorkDir, 0777); err != nil {
		log.Println(err)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	r.ParseForm()
	ipaID := r.Form.Get("ipa_id")
	name := r.Form.Get("name")

	log.Println("Running generate_apk.sh script")
	cmd := exec.Command(fmt.Sprintf("%s/generate_apk.sh", cwd),
		fmt.Sprintf("%s/android_base", cwd), ipaID, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = tmpWorkDir
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}

	log.Println("Script run completed, archiving artifact")
	artifact, err := os.Open(fmt.Sprintf("%s/artifacts/output.apk", cwd))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Writing artifact to response")
	if _, err := io.Copy(w, artifact); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/create_compile_ipa", createCompileIPA)
	http.HandleFunc("/create_generate_apk", createCompileAPK)
	http.ListenAndServe(":"+port, nil)
}
