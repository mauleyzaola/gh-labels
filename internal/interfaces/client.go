package interfaces

import "github.com/mauleyzaola/gh-labels/internal/types"

//go:generate moq -out ../../mocks/client.go -pkg mocks . Client
type Client interface {
	Create(dst types.RepoInfo, label types.Label) (*types.Label, error)
	Delete(dst types.RepoInfo, labelName string) error
	List(info types.RepoInfo) ([]types.Label, error)
}
