// +build integration

package pkg

import (
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/packages"
)

// note: this test requires and uses system git and makes network requests.
func TestGitOperations(t *testing.T) {
	m := &Module{}
	svc := service.New(m)
	defer svc.StartForTest()()
	sanMarzanoURL := "https://github.com/ketchuphq/san-marzano.git"
	gitDir := path.Join(m.Config.Config.DataDir, "themes", "san-marzano")

	// TestClone
	err := m.Clone(&packages.Package{
		Name:   proto.String("san-marzano"),
		VcsUrl: proto.String(sanMarzanoURL),
	}, "themes")
	assert.NoError(t, err)

	cmd := exec.Command("git", "rev-parse", "HEAD^")
	cmd.Dir = gitDir
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, string(out))
	prevRef := strings.TrimSpace(string(out))

	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = gitDir
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err, string(out))
	curRef := strings.TrimSpace(string(out))

	// TestFetchDir errors if commit is not newer
	err = FetchDir(gitDir, prevRef)
	assert.EqualError(t, err, prevRef+" is not a descendant of the current head")

	err = FetchDir(gitDir, curRef)
	assert.EqualError(t, err, curRef+" is not a descendant of the current head")

	// TestFetchDir success
	cmd = exec.Command("git", "checkout", prevRef)
	cmd.Dir = gitDir
	err = cmd.Run()
	assert.NoError(t, err)

	err = FetchDir(gitDir, curRef)
	assert.NoError(t, err)

	// TestGetLatestRef
	ref, err := GetLatestRef(sanMarzanoURL)
	assert.NoError(t, err)
	assert.Equal(t, curRef, ref)
}
