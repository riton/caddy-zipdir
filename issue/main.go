package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zip"
)

var (
	// generated using "dd if=/dev/zero of=./file1 bs=1G count=2" or equivalent
	fileNames = []string{
		"./file1",
		"./file2",
	}
)

func main() {
	outFileW, err := os.Create("./archive.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer outFileW.Close()

	zipWriter := zip.NewWriter(outFileW)

	for _, filename := range fileNames {

		fileToZip, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer fileToZip.Close()

		// Get the file information
		info, err := fileToZip.Stat()
		if err != nil {
			log.Fatal(err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Fatal(err)
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = filepath.Base(filename)

		// Change to store to avoid compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Store

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("archive.zip file created")
}
