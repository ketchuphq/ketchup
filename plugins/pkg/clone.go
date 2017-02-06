package pkg

import (
	"io"
	"os"
	"path"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/octavore/ketchup/proto/ketchup/packages"
)

// Clone the given repo to {data_dir}/{themes}/{name}
func (m *Module) Clone(p *packages.Package) error {
	packagePath := m.Config.DataPath(path.Join(themeDir, p.GetName()), "")
	return m.clone(packagePath, p.GetVcsUrl())
}

func (m *Module) clone(dest, url string) error {
	r, err := git.NewFilesystemRepository(path.Join(dest, ".git"))
	if err != nil {
		return err
	}
	err = r.Clone(&git.CloneOptions{URL: url})
	if err != nil {
		return err
	}

	ref, err := r.Reference("refs/remotes/origin/master", true)
	if err != nil {
		return err
	}

	commit, err := r.Commit(ref.Hash())
	if err != nil {
		return err
	}

	iter, err := commit.Files()
	if err != nil {
		return err
	}
	defer iter.Close()

	return iter.ForEach(func(f *object.File) error {
		pth := path.Join(dest, f.Name)
		err := os.MkdirAll(path.Dir(pth), 0700)
		if err != nil {
			return err
		}
		g, err := os.OpenFile(pth, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode)
		if err != nil {
			return err
		}
		defer g.Close()
		rdr, err := f.Reader()
		if err != nil {
			return err
		}
		defer rdr.Close()
		_, err = io.Copy(g, rdr)
		return err
	})
}
