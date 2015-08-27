package service

import (
	"gost/config"
	"testing"
)

func TestServiceBase(t *testing.T) {
	config.InitTestsDatabase()

	sess, col := Connect("testCollection")
	if sess == nil || col == nil {
		t.Fatal("Cannot connect to the mongodb service")
	}
}
