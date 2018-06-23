package c14

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dutchcoders/goftp"
	"github.com/nyhu/mamouth/entity"
)

const (
	BASE_URL = "https://api.online.net/api/v1"
	token    = "2e9bfae5eac96f8ab6a32b07b979de26258d7a3d"
	api      = "https://api.online.net"
)

func GetSafe(name string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", api+"/api/v1/storage/c14/safe", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	ho, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var body []map[string]interface{}

	err = json.Unmarshal(ho, &body)
	if err != nil {
		var part map[string]interface{}

		err = json.Unmarshal(ho, &part)
		if err != nil {
			panic("Big Fucking Error 1")
		}
		body = append(body, part)
	}

	for i := range body {
		if body[i]["name"] == name {
			return body[i]["$ref"].(string), nil
		}
	}

	panic("Big Fucking Error 2")

	return "", nil
}

func GetArchive(name string, safe string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", api+safe+"/archive", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	ho, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var body []map[string]interface{}

	err = json.Unmarshal(ho, &body)
	if err != nil {
		var part map[string]interface{}

		err = json.Unmarshal(ho, &part)
		if err != nil {
			panic("Big Fucking Error 3")
		}
		body = append(body, part)
	}

	for i := range body {
		if body[i]["name"] == name {
			if body[i]["status"] != "active" {
				return "not ready", nil
			}
			return body[i]["$ref"].(string), nil
		}
	}

	panic("Big Fucking Error 4")

	return "", nil
}

func CreateSafe(name string) (string, error) {
	client := &http.Client{}
	data := entity.CreateSafe{
		Name: name,
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", api+"/api/v1/storage/c14/safe",
		strings.NewReader(string(encoded)))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 201 {
		ho, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		var body map[string]interface{}

		err = json.Unmarshal(ho, &body)
		if err != nil {
			return "", err
		}

		if len(body) != 2 {
			panic("Big Fucking Error 5")
		}

		if body["error"] == "This name is already used" {
			path, err := GetSafe(name)
			if err != nil {
				panic("Big Fucking Error 6")
			}

			return path, nil
		} else {
			panic("Big Fucking Error 7")
		}
	}

	location, err := res.Location()
	if err != nil {
		return "", err
	}

	return location.Path, nil
}

func CreateArchive(name string, safe string) (string, error) {
	client := &http.Client{}
	data := entity.CreateArchive{
		Name:        name,
		Description: "description",
		Protocols:   []string{"ftp"},
		Platforms:   []string{"1"},
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", api+safe+"/archive",
		strings.NewReader(string(encoded)))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 201 {
		ho, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		var body map[string]interface{}

		err = json.Unmarshal(ho, &body)
		if err != nil {
			panic("Big Fucking Error 8")
		}

		if len(body) != 2 {
			panic("Big Fucking Error 9")
		}

		if body["error"] == "This name is already used" {
			path := "not ready"

			for path == "not ready" {
				path, err = GetArchive(name, safe)
				if err != nil {
					panic("Big Fucking Error 10")
				}
				if path == "not ready" {
					time.Sleep(time.Second)
				}
			}
			time.Sleep(1 * time.Minute)
			return path, nil
		} else {
			panic("Big Archive error: " + body["error"].(string))
		}
	}

	path := "not ready"

	for path == "not ready" {
		path, err = GetArchive(name, safe)
		if err != nil {
			panic("Big Fucking Error 12")
		}
		if path == "not ready" {
			time.Sleep(time.Second)
		}
	}
	time.Sleep(1 * time.Minute)
	return path, nil
}

func GetBucket(archive string) (entity.Bucket, error) {
	var decoded entity.Bucket

	client := &http.Client{}
	req, err := http.NewRequest("GET", api+archive+"/bucket", nil)
	if err != nil {
		return decoded, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return decoded, err
	}

	ho, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return decoded, err
	}

	r := bytes.NewReader(ho)
	decoder := json.NewDecoder(r)

	err = decoder.Decode(&decoded)
	if err != nil {
		return decoded, err
	}

	return decoded, nil
}

func ConnectToBucket(credential entity.Credentials) (*goftp.FTP, error) {
	var err error
	var ftp *goftp.FTP

	split := strings.Split(credential.Uri, "@")

	if ftp, err = goftp.Connect(split[1]); err != nil {
		return ftp, err
	}

	if err = ftp.Login(credential.Login, credential.Password); err != nil {
		return ftp, err
	}

	return ftp, nil
}

func SendToBucket(ftp *goftp.FTP, path string) error {
	var err error
	var file *os.File

	if file, err = os.Open(path); err != nil {
		return err
	}

	split := strings.Split(path, "/")
	index := len(split) - 1
	if index < 0 {
		panic("Big Fucking Error 13")
	}
	fileName := "/" + split[index]

	if err := ftp.Stor(fileName, file); err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	//var files []string
	//if files, err = ftp.List(""); err != nil {
	//	panic(err)
	//}
	//fmt.Println("Directory listing:", files)

	return nil
}

func Freeze(archive string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", api+archive+"/archive", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return err
}

func GrosPorc(pathRelayChan chan string, ftp *goftp.FTP, now int64, signal chan os.Signal, archive string) {
	for {
		select {
		case path := <-pathRelayChan:
			fmt.Println("new message : " + path)

			err := SendToBucket(ftp, path)
			if err != nil {
				fmt.Println(err)
				return
			}

			if now < time.Now().Truncate(5*time.Minute).Unix() {
				go Freeze(archive)
				return
			}

		case <-signal:
			return
		}
	}
}

func Relay(topicName string, pathRelayChan chan string, signal chan os.Signal) {
	safe, err := CreateSafe(topicName)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		now := time.Now().Truncate(5 * time.Minute).Unix()
		archiveName := fmt.Sprintf("%v", now)

		archive, err := CreateArchive(archiveName, safe)
		if err != nil {
			fmt.Println(err)
			return
		}

		bucket, err := GetBucket(archive)
		if err != nil {
			fmt.Println(err)
			return
		}
		credentials := entity.NewCredentials(bucket.Credentials[0].(map[string]interface{}))

		ftp, err := ConnectToBucket(credentials)
		if err != nil {
			fmt.Println(err)
			return
		}

		GrosPorc(pathRelayChan, ftp, now, signal, archive)
		ftp.Close()
	}
}

func GetAllArchives(safeId string) ([]entity.Archive, error) {
	client := &http.Client{}
	var decoded []entity.Archive

	req, err := http.NewRequest("GET", BASE_URL+"/storage/c14/safe/"+safeId+"/archive", nil)
	if err != nil {
		return decoded, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return decoded, err
	}

	temp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return decoded, err
	}
	r := bytes.NewReader(temp)
	decoder := json.NewDecoder(r)

	err = decoder.Decode(&decoded)
	if err != nil {
		return decoded, err
	}
	return decoded, nil
}
