package utils

import (
	"fmt"
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
		panic(err)
		return
	}
	fmt.Println(zip)
	unzip, err := Unzip(zip)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(unzip))
}
