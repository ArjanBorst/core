package helpers

import (
	"io/ioutil"
	"os"
)

func SaveToFile(path, filename string, jsonData []byte) {

	jsonFile, err := os.Create(path + filename)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}

func LoadFromFile(path, filename string) []byte {
	file, _ := os.Open(path + filename)
	data, _ := ioutil.ReadAll(file)

	return data
}
