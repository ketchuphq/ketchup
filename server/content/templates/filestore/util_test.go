package filestore

import (
	"io/ioutil"
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v3/osfs"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func TestGetCurrentRef(t *testing.T) {
	// setup a repo
	gitDir, err := ioutil.TempDir("", "test-repo-")
	assert.Nil(t, err)

	dir := osfs.New(gitDir)
	s, err := filesystem.NewStorage(dir)
	assert.Nil(t, err)

	repo, err := git.Init(s, dir)
	assert.Nil(t, err)

	work, err := repo.Worktree()
	assert.Nil(t, err)

	// make a commit
	hsh, err := work.Commit("hello", &git.CommitOptions{Author: &object.Signature{Name: "test"}})
	assert.Nil(t, err)

	ref, err := getCurrentRef(gitDir)
	assert.Nil(t, err)
	assert.Equal(t, hsh.String(), ref)
}
