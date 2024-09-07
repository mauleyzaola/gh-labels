package cli

import (
	"log"

	"github.com/mauleyzaola/gh-labels/internal/interfaces"

	"github.com/mauleyzaola/gh-labels/internal/types"
)

func CopyLabels(client interfaces.Client, src, dst types.RepoInfo) error {
	log.Printf("[INFO] getting a list of source labels: %s%s\n", src.Username, src.Repository)
	sourceLabels, err := client.List(src)
	if err != nil {
		return err
	}

	log.Printf("[INFO] getting a list of target labels: %s/%s\n", dst.Username, dst.Repository)
	targetLabels, err := client.List(dst)
	if err != nil {
		return err
	}

	log.Printf("[INFO] removing target labels: %s/%s\n", dst.Username, dst.Repository)
	for _, v := range targetLabels {
		log.Println("[INFO] removing target label: ", v.Name)
		if err = client.Delete(dst, v.Name); err != nil {
			return err
		}
	}

	log.Printf("[INFO] creating target labels: %s/%s\n", dst.Username, dst.Repository)
	for _, v := range sourceLabels {
		log.Println("[INFO] creating target label: ", v.Name)
		if err = client.Create(dst, v); err != nil {
			return err
		}
	}
	log.Printf("[SUCCESS] completed cloning labels from: %s/%s to %s/%s\n",
		src.Username, src.Repository, dst.Username, dst.Repository)
	return nil
}
