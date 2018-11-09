package main

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	configuration["lall"] = "fasel"
	log.Infof("started\n")
	os.Exit(m.Run())

}
