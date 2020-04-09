package utils

import (
	"fmt"
	"testing"
)

const filename = "env6-default-nginx-base-67bd6995f6-2020.04.03"

func TestIsNeedBackup(t *testing.T) {

	fmt.Println(IsNeedBackup(filename, 6))
}
