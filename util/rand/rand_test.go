package rand

import (
	"log"
	"testing"
)

func TestRand(t *testing.T) {
	log.Printf("int   : %d", Int(10))
	log.Printf("int8  : %d", Int8(10))
	log.Printf("int16 : %d", Int16(10))
	log.Printf("int32 : %d", Int32(10))
	log.Printf("int64 : %d", Int64(10))
	log.Printf("float32 : %f", Float32(10))
	log.Printf("float64: %f", Float64(10))

	log.Printf("random string 10 : %s", RandomString(10, SeedUCL))
	log.Printf("random string 20 : %s", RandomString(20, SeedLCL, SeedUCL))
	log.Printf("random string 30 : %s", RandomString(30))
}
