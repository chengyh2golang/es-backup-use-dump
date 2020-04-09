package esops

import (
	"fmt"
	"testing"
)


func TestFetchIndices(t *testing.T) {
	strings, err := FetchIndices("192.168.250.9", "9200")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(strings)
}




