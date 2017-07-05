package pkg

import (
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func GetLatestRef(vcsURL string) (string, error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: vcsURL})
	if err != nil {
		return "", err
	}
	refs, err := r.References()
	if err != nil {
		return "", err
	}

	masterRef := ""
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if ref.Name().Short() == "origin/master" {
				masterRef = ref.Hash().String()
			}
		}
		return nil
	})
	return masterRef, err
}
