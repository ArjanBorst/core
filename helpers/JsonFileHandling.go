package helpers

import (
	"encoding/json"
	"errors"
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

func CreateFileIfNotExist(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {

		file, err := os.Create(filename)
		if err != nil {
			return errors.New("ERROR CREATING NEW FILE")
		}
		defer file.Close()

		return nil
	} else if err != nil {
		return errors.New("ERROR ACCESSING FILE")
	}

	return nil
}

func LoadDataFromFile(path, name string, target interface{}) error {
	filename := path + name

	if err := CreateFileIfNotExist(filename); err != nil {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return errors.New("ERROR LOADING FILE")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("ERROR READING DATA FROM FILE")
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return errors.New("ERROR UNMARSHALLING JSON")
	}

	return nil
}

func SaveDataToFile(path, name string, target interface{}) error {
	filename := path + name

	if err := CreateFileIfNotExist(filename); err != nil {
		return nil
	}

	jsonData, err := json.Marshal(target)
	if err != nil {
		return errors.New("ERROR MARSHALLING JSON")
	}

	file, err := os.Create(path + name)
	if err != nil {
		return errors.New("ERROR OPENING FILE")
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return errors.New("ERROR SAVING DATA TO FILE")
	}

	return nil
}
