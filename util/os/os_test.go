package os

import (
	"log"
	"testing"
)

func TestFileExist(t *testing.T) {
	log.Println(FileExist("/home/frank/Documents/project"))
	log.Println(DirExist("/home/frank/Documents/project"))
	log.Println(FileExist("/home/frank/Documents/project/volumes/README.md"))
	log.Println(DirExist("/home/frank/Documents/project/volumes/README.md"))
}
