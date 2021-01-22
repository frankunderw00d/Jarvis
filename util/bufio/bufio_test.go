package bufio

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)

func TestBufio(t *testing.T) {
	l, err := SplitFile("gopher.txt", '\n')
	if err != nil {
		log.Fatalln(err)
	}

	for _, str := range l {
		log.Printf("[%s]", strings.Trim(str, "\n"))
	}
}

func TestBufioScanner(t *testing.T) {
	file, err := os.OpenFile("gopher.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	data := Scan(file)
	for _, str := range data {
		log.Printf("[%s]", str)
	}
}

func TestBufioWriter(t *testing.T) {
	n, err := Write("gopher.txt", "12345", "123456", "1234567", "12345678", "123456789")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%d", n)
}

func TestBufioReaderWriter(t *testing.T) {
	readFile, err := os.OpenFile("gopher.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer readFile.Close()
	reader := bufio.NewReader(readFile)

	writeFile, err := os.OpenFile("copy.txt", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer writeFile.Close()
	writer := bufio.NewWriter(writeFile)

	rw := bufio.NewReadWriter(reader, writer)

	var e error
	for {
		str, err := rw.ReadString('\n')
		if len(str) > 0 {
			n, err := rw.WriteString(str)
			if err != nil {
				e = err
				break
			}
			log.Printf("[%d]", n)
		}

		if err != nil {
			log.Println(err)
			break
		}
	}

	if e != nil {
		log.Println(e)
	}

	if err := rw.Flush(); err != nil {
		log.Fatalln(err)
	}
}
