package pkg

import (
	"io"
	"os"
	"path"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
	"github.com/ketchuphq/ketchup/util/errors"
)

// Clone the given package to {data_dir}/{dataSubdir}/{package name}
func (m *Module) Clone(p *packages.Package, dataSubdir string) error {
	packagePath := m.Config.DataPath(path.Join(dataSubdir, p.GetName()), "")
	return CloneToDir(packagePath, p.GetVcsUrl())
}

func getRepoMasterIterator(repo *git.Repository) (*object.FileIter, error) {
	ref, err := repo.Reference("refs/remotes/origin/master", true)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	iter, err := commit.Files()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	// defer iter.Close()
	return iter, nil
}

type setDatar interface {
	SetData(v *string)
}

// Clone the given url to the given dest folder
func CloneToDir(dest, url string) error {
	r, err := git.PlainClone(dest, false, &git.CloneOptions{URL: url})
	if err != nil {
		return errors.Wrap(err)
	}

	iter, err := getRepoMasterIterator(r)
	if err != nil {
		return errors.Wrap(err)
	}
	defer iter.Close()

	return iter.ForEach(func(f *object.File) error {
		pth := path.Join(dest, f.Name)
		err := os.MkdirAll(path.Dir(pth), 0700)
		if err != nil {
			return errors.Wrap(err)
		}
		mode, err := f.Mode.ToOSFileMode()
		if err != nil {
			return errors.Wrap(err)
		}
		g, err := os.OpenFile(pth, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mode)
		if err != nil {
			return errors.Wrap(err)
		}
		defer g.Close()
		rdr, err := f.Reader()
		if err != nil {
			return errors.Wrap(err)
		}
		defer rdr.Close()
		_, err = io.Copy(g, rdr)
		return errors.Wrap(err)
	})
}

// CheckForUpdates in the given repo at {data_dir}/{dataSubdir}/{packageName}
func (m *Module) CheckForUpdates(packageName, dataSubdir string) (headRef, latestRef string, err error) {
	packagePath := m.Config.DataPath(path.Join(dataSubdir, packageName), "")
	repo, err := git.PlainOpen(packagePath)
	if err != nil {
		return "", "", errors.Wrap(err)
	}
	head, err := repo.Head()
	if err != nil {
		return "", "", errors.Wrap(err)
	}
	latest, err := repo.Reference("refs/remotes/origin/master", true)
	if err != nil {
		return "", "", errors.Wrap(err)
	}
	headRef = head.Hash().String()
	latestRef = latest.Hash().String()
	return
}
