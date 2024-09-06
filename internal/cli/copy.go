package cli

import (
	"github.com/mauleyzaola/gh-labels/internal/api"
	"log"
)

func CopyLabels(apiClient *api.Client, sourceAuthor, sourceRepo, targetAuthor, targetRepo string) error {
	log.Printf("[INFO] getting a list of source labels: %s%s\n", sourceAuthor, sourceRepo)
	sourceLabels, err := apiClient.LabelList(sourceAuthor, sourceRepo)
	if err != nil {
		return err
	}

	log.Printf("[INFO] getting a list of target labels: %s/%s\n", targetAuthor, targetRepo)
	targetLabels, err := apiClient.LabelList(targetAuthor, targetRepo)
	if err != nil {
		return err
	}

	log.Printf("[INFO] removing target labels: %s/%s\n", targetAuthor, targetRepo)
	for _, v := range targetLabels {
		log.Println("[INFO] removing target label: ", v.Name)
		if err = apiClient.LabelDelete(targetAuthor, targetRepo, v.Name); err != nil {
			return err
		}
	}

	log.Printf("[INFO] creating target labels: %s/%s\n", targetAuthor, targetRepo)
	for _, v := range sourceLabels {
		log.Println("[INFO] creating target label: ", v.Name)
		if err = apiClient.LabelPost(targetAuthor, targetRepo, &v); err != nil {
			return err
		}
	}
	log.Printf("[SUCCESS] completed cloning labels from: %s/%s to %s/%s\n", sourceAuthor, sourceRepo, targetAuthor, targetRepo)
	return nil
}
