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
