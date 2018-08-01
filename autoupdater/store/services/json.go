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
	cleanFile := filepath.Clean(file)
	createDirIfNotExists(cleanFile)

	jsonFile, err := os.OpenFile(cleanFile, os.O_RDONLY|os.O_CREATE, 0744)
	if err != nil{
		logrus.Fatalf("unable to recreate state %s",err)
	}

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil{
		logrus.Fatalf("unable to read %s. %s", file, err)
	}

	return &jsonStore{
		file: file,
		services: decodeServices(jsonBytes),
	}
}

func (s *jsonStore) Get(service string) *Service{
	return s.services.get(service)
}

func (s *jsonStore) Store(service *Service) {
	logrus.Infof("writing status file %s", s.file)
	s.services.set(service)
	encoded, err := json.MarshalIndent(s.services,"","	 ")
	if err != nil {
		logrus.Fatalf("cannot encoded status to json %s", err)
	}

	err = ioutil.WriteFile(s.file, encoded, 0744)
	if err != nil {
		logrus.Fatalf("cannot write status file %s", err)
	}
}

func createDirIfNotExists(path string){
	dir:= filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0744)
	}
}

func decodeServices(jsonBytes []byte) Services {
	services := Services{}

	if len(jsonBytes) == 0 {
		return services
	}

	err := json.Unmarshal(jsonBytes, &services)
	if err != nil{
		logrus.Fatalf("unable to decode json file %s",err)
	}

	return services
}