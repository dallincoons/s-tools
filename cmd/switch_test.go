package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestEnvFileExists(t *testing.T) {
	dir, _ := os.Getwd()

	_, err := RunSwitch(dir, "doesntmatter")

	if err == nil {
		log.Fatalf("error should have been thrown for missing environment file")
	}
}

func TestReplaceDatabaseNameInEnvFile(t *testing.T) {
	ioutil.WriteFile("../test_fixtures/.env", []byte("DB_DATABASE=fakedb1"), 0644)

	_, err := RunSwitch("../test_fixtures", "newfakedb1")

	if (err != nil) {
		log.Fatal(err)
	}

	contents, _ := ioutil.ReadFile("../test_fixtures/.env")

	if string(contents) != "DB_DATABASE=newfakedb1" {
		log.Fatalf("contents does not equal DB_DATABASE=newfakedb1, received %s", string(contents))
	}
}
