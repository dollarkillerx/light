package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestIsPublic(t *testing.T) {
	r := "Name"
	fmt.Println(IsPublic(r))
	r = "name"
	fmt.Println(IsPublic(r))

	rc := map[string]string{}
	fmt.Println(rc["jex"])
}

func TestDisID(t *testing.T) {
	fmt.Println(DistributedID())
}

func TestGZIP(t *testing.T) {
	rc := []byte("hello world")
	zip, err := Zip(rc)
	if err != nil {
		log.Fatalln(err)
		return
	}

	unzip, err := Unzip(zip)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(unzip))
}
