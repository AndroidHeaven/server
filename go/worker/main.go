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

	artifact, err := os.Open(fmt.Sprintf("%s/artifacts/artifact.tar.gz",
		tmpWorkDir))
	if err != nil {
		log.Println(err)
		return
	}

	defer artifact.Close()
	if _, err := io.Copy(w, artifact); err != nil {
		log.Println(err)
		return
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Println(err)
		return
	}
}

func createCompileAPK(w http.ResponseWriter, r *http.Request) {
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

	ipaID := r.Form.Get("ipa_id")
	name := r.Form.Get("name")

	cmd := exec.Command(fmt.Sprintf("%s/generate_apk.sh", cwd), ipaID, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = tmpWorkDir
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}

	artifact, err := os.Open(fmt.Sprintf("%s/artifacts/output.apk", cwd))
	if err != nil {
		log.Println(err)
		return
	}

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