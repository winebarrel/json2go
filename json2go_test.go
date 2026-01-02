package json2go_test

import (
	"os"
	"testing"
)

var (
	testAcc = false
)

func TestMain(m *testing.M) {
	if v := os.Getenv("TEST_ACC"); v == "1" {
		testAcc = true
	}

	m.Run()
}
