package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	var sourceAuthor, targetAuthor, sourceRepo, targetRepo string
	var source, target string

	flag.StringVar(&source, "source", "", "source for labels in the format <owner/repository-name>")
	flag.StringVar(&target, "target", "", "target for labels in the format <owner/repository-name>")
	flag.Parse()

	if values := strings.Split(source, "/"); len(values) == 2 {
		sourceAuthor, sourceRepo = values[0], values[1]
	} else {
		log.Println("[ERROR] wrong or missing parameter: source. should be <owner/repo>")
		return
	}
	if values := strings.Split(target, "/"); len(values) == 2 {
		targetAuthor, targetRepo = values[0], values[1]
	} else {
		log.Println("[ERROR] wrong or missing parameter: target. should be <owner/repo>")
		return
	}

	apiClient, err := NewAPIClient(nil)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return
	}

	log.Printf("[INFO] getting a list of source labels: %s%s\n", sourceAuthor, sourceRepo)
	sourceLabels, err := apiClient.LabelList(sourceAuthor, sourceRepo)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return
	}

	log.Printf("[INFO] getting a list of target labels: %s/%s\n", targetAuthor, targetRepo)
	targetLabels, err := apiClient.LabelList(targetAuthor, targetRepo)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return
	}

	log.Printf("[INFO] removing target labels: %s/%s\n", targetAuthor, targetRepo)
	for _, v := range targetLabels {
		log.Println("[INFO] removing target label: ", v.Name)
		if err = apiClient.LabelDelete(targetAuthor, targetRepo, v.Name); err != nil {
			log.Println("[ERROR]: ", err)
			return
		}
	}

	log.Printf("[INFO] creating target labels: %s/%s\n", targetAuthor, targetRepo)
	for _, v := range sourceLabels {
		log.Println("[INFO] creating target label: ", v.Name)
		if err = apiClient.LabelPost(targetAuthor, targetRepo, &v); err != nil {
			log.Println("[ERROR]: ", err)
			return
		}
	}
	log.Printf("[SUCCESS] completed cloning labels from: %s/%s to %s/%s\n", sourceAuthor, sourceRepo, targetAuthor, targetRepo)
}
