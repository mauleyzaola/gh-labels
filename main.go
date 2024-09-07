package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mauleyzaola/gh-labels/internal/api"
	"github.com/mauleyzaola/gh-labels/internal/cli"
	"github.com/mauleyzaola/gh-labels/internal/types"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	if err := run(); err != nil {
		log.Println("[error]", err)
	}
}

func run() error {
	var sourceAuthor, targetAuthor, sourceRepo, targetRepo string
	var source, target string

	flag.StringVar(&source, "source", "", "source for labels in the format <owner/repository-name>")
	flag.StringVar(&target, "target", "", "target for labels in the format <owner/repository-name>")
	flag.Parse()

	if values := strings.Split(source, "/"); len(values) == 2 {
		sourceAuthor, sourceRepo = values[0], values[1]
	} else {
		return fmt.Errorf("[ERROR] wrong or missing parameter: source. should be <owner/repo>")
	}
	if values := strings.Split(target, "/"); len(values) == 2 {
		targetAuthor, targetRepo = values[0], values[1]
	} else {
		return fmt.Errorf("[ERROR] wrong or missing parameter: target. should be <owner/repo>")
	}

	client := api.New(os.Getenv("TOKEN"))
	src := types.RepoInfo{
		Repository: sourceRepo,
		Username:   sourceAuthor,
	}
	dst := types.RepoInfo{
		Repository: targetRepo,
		Username:   targetAuthor,
	}
	return cli.CopyLabels(client, src, dst)
}
