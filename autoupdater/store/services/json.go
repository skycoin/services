package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type jsonStore struct {
	file string
	services Services
}

func newJsonStore(file string) *jsonStore {
	var services Services

	cleanFile := filepath.Clean(file)
	createDirIfNotExists(cleanFile)

	jsonFile, err := os.OpenFile(cleanFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil{
		logrus.Fatalf("unable to recreate state %s",err)
	}

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil{
		logrus.Fatalf("unable to read %s. %s", file, err)
	}

	err = json.Unmarshal(jsonBytes, &services)
	if err != nil{
		logrus.Fatalf("unable to decode json file %s",err)
	}

	return &jsonStore{
		file: file,
		services: services,
	}
}

func (s *jsonStore) Get(service string) *Service{
	return &Service{}
}

func (s *jsonStore) Store(service *Service) {

}

func createDirIfNotExists(path string){
	dir:= filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModeDir)
	}
}