package gcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetAgent interface {
	ReadSpreadsheet(spreadsheetID, sheetName string) ([]string, error)
}

type SheetDomain struct {
	SheetService *sheets.Service
}

func NewSheetDomain(creds string) *SheetDomain {
	ctx := context.Background()

	service, err := sheets.NewService(ctx,
		option.WithScopes(sheets.SpreadsheetsScope),
		option.WithCredentialsFile(creds))

	if err != nil {
		log.Fatalf("Unable to create Google Sheets service: %v", err)
	}
	return &SheetDomain{
		SheetService: service,
	}
}

func (s *SheetDomain) ReadSpreadsheet(spreadsheetID, sheetName string) (data []string, err error) {

	readRange := fmt.Sprintf("%s!A2:D90", sheetName)
	resp, err := s.SheetService.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var repoAccountLinks, accountNames []string
	for _, row := range resp.Values {
		if len(row) > 0 {
			repoLink := row[3].(string)
			repolinkArray := strings.Split(repoLink, "github.com/")
			if len(repolinkArray) < 2 {
				log.Println("failed:", repolinkArray)
				continue
			}
			githubLinks := strings.Split(repolinkArray[1], "/")
			githubAccountName := strings.ReplaceAll(githubLinks[0], "?tab=repositories", "")
			repoAccountLink := "https://github.com/" + githubAccountName
			repoAccountLinks = append(repoAccountLinks, repoAccountLink)
			accountNames = append(accountNames, githubAccountName)
		}
	}
	return accountNames, nil
}
