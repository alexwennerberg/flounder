package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

var c Config // global var to hold static configuration

const ( // todo make configurable
	userFilesPath = "./files"
)

type File struct {
	Creator     string
	Name        string
	UpdatedTime string
}

func getUsers() ([]string, error) {
	return []string{"me", "other guy"}, nil
}

func getIndexFiles() ([]*File, error) { // cache this function
	result := []*File{}
	err := filepath.Walk(userFilesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failure accessing a path %q: %v\n", path, err)
			return err // think about
		}
		// make this do what it should
		result = append(result, &File{
			Name:        info.Name(),
			Creator:     "alex",
			UpdatedTime: "123123",
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	// sort
	// truncate
	return result, nil
} // todo clean up paths

func getUserFiles(user string) ([]*File, error) {
	result := []*File{}
	files, err := ioutil.ReadDir(path.Join(userFilesPath, user))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		result = append(result, &File{
			Name:        file.Name(),
			Creator:     user,
			UpdatedTime: "123123",
		})
	}
	return result, nil
}

func main() {
	configPath := flag.String("c", "flounder.toml", "path to config file")
	var err error
	c, err = getConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	runHTTPServer()
	// runGeminiServer()
	// go log.Fatal(gmi.ListenAndServe(":8080", nil))
}
