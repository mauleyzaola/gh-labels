package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/mauleyzaola/gh-labels/internal/types"

	"github.com/mauleyzaola/gh-labels/internal/api"
	"github.com/mauleyzaola/gh-labels/internal/cli"
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

	client, err := api.New(os.Getenv("TOKEN"))
	if err != nil {
		log.Println("[ERROR]: ", err)
		return
	}
	src := types.RepoInfo{
		Repository: sourceRepo,
		Username:   sourceAuthor,
	}
	dst := types.RepoInfo{
		Repository: targetRepo,
		Username:   targetAuthor,
	}
	if err = cli.CopyLabels(client, src, dst); err != nil {
		log.Println("[ERROR]: ", err)
	}
}
