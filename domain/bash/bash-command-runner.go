package bash

import (
	"log"
	"os"
	"os/exec"
)

type BashCommandRunnerAgent interface {
	GitClone(repoURL, destinationDir string) (err error)
}

type BashCommandRunnerDomain struct{}

func NewBashDomain() *BashCommandRunnerDomain {
	return &BashCommandRunnerDomain{}
}

func (b *BashCommandRunnerDomain) GitClone(repoURL, destinationDir string) (err error) {

	cmd := exec.Command("git", "clone", repoURL, destinationDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Println("Error executing git clone:", err)
		return
	}

	log.Println("Git clone successful!")
	return
}

func (b *BashCommandRunnerDomain) RunCommand(args []string, filePath string) (err error) {
	// cmd := exec.Command("go", "run", filePath)
	cmd := exec.Command(args[0], args[1], filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("Error running Go script %s: %v\n", filePath, err)
		return
	}

	log.Printf("Go script %s executed successfully.\n", filePath)
	return
}
