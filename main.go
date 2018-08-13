package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handleUpload)

	fs := http.FileServer(http.Dir("testresults"))
	mux.Handle("/", fs)
	// mux.Handle("/", http.StripPrefix("/", fs))

	http.ListenAndServe(":8883", mux)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now(), r.RequestURI)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	testname := r.Header.Get("testname")
	if testname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("testname not specified in headers"))
		return
	}

	nowStr := time.Now().Format("20060102-150405")
	targetFileDir := fmt.Sprintf("./testresults/%s/%s", testname, nowStr)
	_, err := os.OpenFile(targetFileDir, os.O_WRONLY, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		os.MkdirAll(targetFileDir, os.ModePerm)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := decompressArchiveToFile(r.Body, targetFileDir); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decompressArchiveToFile(reqBody io.ReadCloser, targetFileDir string) error {
	gzipReader, err := gzip.NewReader(reqBody)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		name := header.Name

		target := filepath.Join(targetFileDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if e := os.MkdirAll(target, os.ModePerm); e != nil {
				fmt.Println(e)
			}
		case tar.TypeReg:
			file, _ := os.Create(target)
			if _, e := io.Copy(file, tarReader); e != nil {
				fmt.Println(e)
			}
		default:
			fmt.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
			return fmt.Errorf("Unable to figure out type")
		}

	}

	return nil
}
