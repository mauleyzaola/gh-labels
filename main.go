package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	var sourceAuthor, targetAuthor string
	var sourceRepo, targetRepo string

	flag.StringVar(&sourceAuthor, "source-author", "", "name of source author")
	flag.StringVar(&targetAuthor, "target-author", "", "name of target author")
	flag.StringVar(&sourceRepo, "source-repo", "", "name of source repository")
	flag.StringVar(&targetRepo, "target-repo", "", "name of target repository")
	flag.Parse()

	if sourceRepo == "" {
		log.Println("[ERROR] missing parameter: source-repo")
		return
	}

	if targetRepo == "" {
		log.Println("[ERROR] missing parameter: target-repo")
		return
	}

	if targetAuthor == "" {
		log.Println("[ERROR] missing parameter: source-author")
		return
	}

	if targetAuthor == "" {
		log.Println("[ERROR] missing parameter: target-author")
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
		//v.Id=0
		if err = apiClient.LabelPost(targetAuthor, targetRepo, &v); err != nil {
			log.Println("[ERROR]: ", err)
			return
		}
	}
	log.Printf("[SUCCESS] completed cloning labels from: %s/%s to %s/%s\n", sourceAuthor, sourceRepo, targetAuthor, targetRepo)
}
