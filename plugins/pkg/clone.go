package pkg

import (
	"io"
	"os"
	"path"

	"github.com/octavore/nagax/util/errors"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
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

// Fetch does a git fetch to the given ref
func FetchDir(dir, ref string) error {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return errors.Wrap(err)
	}

	// ensure worktree is clean
	wt, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err)
	}
	// todo: need support for gitignore
	status, err := wt.Status()
	if err != nil {
		return errors.Wrap(err)
	}
	if !status.IsClean() {
		return errors.New("%s is not clean; refusing to update. %s", dir)
	}

	// fetch
	err = r.Fetch(&git.FetchOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err)
	}

	// get current head
	head, err := r.Head()
	if err != nil {
		return errors.Wrap(err)
	}
	headHash := head.Hash().String

	// get commit for the desired ref
	remoteHash := plumbing.NewHash(ref)
	// check that head is an ancestor of remoteRef
	remoteCommit, err := r.CommitObject(remoteHash)
	if err != nil {
		return errors.Wrap(err)
	}
	iter := remoteCommit.Parents()
	found := false
	for {
		commit, err := iter.Next()
		if err == io.EOF {
			break
		}
		if commit.Hash.String() == headHash() {
			found = true
			break
		}
	}
	if !found {
		return errors.New("%s is not a descendant of the current head", ref)
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Hash: remoteCommit.Hash,
	})
	return errors.Wrap(err)
}
