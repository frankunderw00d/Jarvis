package archive

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	tarFileName = "frank.tar"
	zipFileName = "frank.zip"
)

// example code
func TestTar(t *testing.T) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nFrank\nUnderwood\nMr.Frank"},
		{"todo.txt", "Get animal handling licence."},
	}

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}

	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}

	r := bytes.NewReader(buf.Bytes())
	tr := tar.NewReader(r)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Contexts of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			log.Fatalln(err)
		}
		fmt.Println()
	}
}

// tar compress code
func TestTarCompress(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	filePath := filepath.Join(wd, tarFileName)

	err = TarCompress(filePath, filepath.Join("/home/frank/Desktop/aaa", "gopher.txt"),
		filepath.Join("/home/frank/Desktop/aaa", "readme.txt"),
		filepath.Join("/home/frank/Desktop/aaa", "todo.txt"))
	if err != nil {
		log.Fatalln(err)
	}
}

// tar decompress code
func TestTarDecompress(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	filePath := filepath.Join(wd, tarFileName)

	err = TarDecompress(filePath, filepath.Join(wd, "decompress"))

	if err != nil {
		log.Fatalln(err)
	}
}

func TestZip(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	filePath := filepath.Join(wd, zipFileName)

	err = ZipCompress(filePath, filepath.Join("/home/frank/Desktop/aaa", "gopher.txt"),
		filepath.Join("/home/frank/Desktop/aaa", "readme.txt"),
		filepath.Join("/home/frank/Desktop/aaa", "todo.txt"))

	if err != nil {
		log.Fatalln(err)
	}
}

func TestZipDepress(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	filePath := filepath.Join(wd, zipFileName)

	err = ZipDecompress(filePath, filepath.Join(wd, "decompress"))
	if err != nil {
		log.Fatalln(err)
	}
}
