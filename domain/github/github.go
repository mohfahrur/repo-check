package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type GithubAgent interface {
	GetRepoContent(username, repo string) (contents Contents, Cerr error)
	DownloadAndRunGoFile(name, fileURL string)
}

type GithubDomain struct {
}

func NewGithubDomain() *GithubDomain {
	return &GithubDomain{}
}

type Contents []struct {
	Name string `json:"name"`
}

func (d *GithubDomain) GetRepoContent(username, repo string) (contents Contents, Cerr error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", username, repo)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making API request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("%s GitHub API returned an error: %s\n", username, resp.Status)
		return
	}

	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error creating download directory: %v\n", err)
		return
	}
	json.Unmarshal(respB, &contents)

	return
}

func (d *GithubDomain) DownloadFile(username, fileName, fileURL string) (downloadPath string, err error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Printf("Error downloading file %s: %v\n", fileName, err)
		return
	}
	defer resp.Body.Close()

	downloadDir := "downloaded_files/" + username
	err = os.MkdirAll(downloadDir, 0755)
	if err != nil {
		log.Printf("Error creating download directory: %v\n", err)
		return
	}
	downloadPath = filepath.Join(downloadDir, fileName)
	file, err := os.Create(downloadPath)
	if err != nil {
		log.Printf("Error creating file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Error copying file %s: %v\n", fileName, err)
		return
	}
	return
}
