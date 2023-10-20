package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohfahrur/repo-check/domain/bash"
	"github.com/mohfahrur/repo-check/domain/gcp"
	"github.com/mohfahrur/repo-check/domain/github"
	validator "github.com/mohfahrur/repo-check/usecase/validate-repository"
	// "github.com/robfig/cron/v3"
)

func main() {
	log.SetFlags(log.Llongfile)
	bashDomain := bash.NewBashDomain()

	gcpDomain := gcp.NewSheetDomain("./credential.json")
	githubDomain := github.NewGithubDomain()
	validator := validator.NewValidator(*bashDomain,
		*gcpDomain,
		*githubDomain)

	err := validator.ValidateFromSpreadsheet("1d31pmojFTvqrGghkqn3ubW2Mkqk8lyVIkN3QOxC7IQA",
		"Form Responses 1",
		"c:/Users/fa162/go/src/github.com/mohfahrur/repo-check",
		"learn-gin-framework")
	if err != nil {
		log.Printf("Validation error: %v\n", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/validate", func(c *gin.Context) {
		err := validator.ValidateFromSpreadsheet("1d31pmojFTvqrGghkqn3ubW2Mkqk8lyVIkN3QOxC7IQA",
			"Form Responses 1",
			"c:/Users/fa162/go/src/github.com/mohfahrur/repo-check",
			"learn-gin-framework")
		if err != nil {
			log.Printf("Validation error: %v\n", err)
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.Run(":5000")
}
