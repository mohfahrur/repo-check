package validator

import (
	"fmt"
	"log"
	"strings"

	"github.com/mohfahrur/repo-check/domain/bash"
	"github.com/mohfahrur/repo-check/domain/gcp"
	"github.com/mohfahrur/repo-check/domain/github"
)

type ValidatorAgent interface {
	ValidateFromSpreadsheet(spreadsheetID, sheetName, destination, repoName string) error
}
type Validator struct {
	BashCommandRunnerDomain bash.BashCommandRunnerDomain
	SheetDomain             gcp.SheetDomain
	GithubDomain            github.GithubDomain
}

func NewValidator(commandRunner bash.BashCommandRunnerDomain,
	sheetService gcp.SheetDomain,
	githubDomain github.GithubDomain) *Validator {
	return &Validator{
		BashCommandRunnerDomain: commandRunner,
		SheetDomain:             sheetService,
		GithubDomain:            githubDomain,
	}
}

func (v *Validator) ValidateFromSpreadsheet(spreadsheetID, sheetName, destination, repoName string) (
	err error) {

	repoData, err := v.SheetDomain.ReadSpreadsheet(spreadsheetID, sheetName)
	if err != nil {
		log.Println(err)
		return
	}
	for _, username := range repoData {
		// repoLink := repoLink + "/" + repoName
		// destination = fmt.Sprintf("%s/fp%d", destination, k)
		contents, err := v.GithubDomain.GetRepoContent(username, repoName)
		if err != nil {
			log.Println(username, err)
			continue
		}
		for _, content := range contents {
			name := content.Name
			if strings.HasSuffix(name, ".go") {
				log.Println(username, "found")
			}
		}
	}

	fmt.Printf("Validation completed")
	return nil
}
