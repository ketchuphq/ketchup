package tls

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/ketchuphq/ketchup/server/config"
)

const testDomain = "example.com"
const testTime = "2017-01-01T00:51:26Z"
const testPath = testDomain + "-2017-01-01-v000.json"

func testNow() time.Time {
	t, _ := time.Parse(time.RFC3339, testTime)
	return t
}

func setup(t *testing.T) (*Module, string) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(path.Join(dir, tlsDir), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	return &Module{
		Config: &config.Module{
			Config: config.Config{
				DataDir: dir,
			},
		},
	}, dir
}

func TestGetAll(t *testing.T) {
	m, dir := setup(t)

	for _, t := range []string{
		"example2.com-2017-01-01-v000.json",
		"example.com-2017-01-01-v000.json",
		"example.com-2017-01-01-v001.json",
	} {
		_ = ioutil.WriteFile(path.Join(dir, tlsDir, t), []byte{}, os.ModePerm)
	}

	matches, err := m.GetAllRegisteredDomains()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(matches, []string{"example.com", "example2.com"}) {
		t.Errorf("unexpected %s", matches)
	}
}

func TestCurrentTLSPath(t *testing.T) {
	m, dir := setup(t)
	s, err := m.getCurrentRegistrationPath(testDomain)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if s != "" {
		t.Errorf("unexpected string %s", s)
	}

	err = ioutil.WriteFile(path.Join(dir, tlsDir, testPath), []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	s, err = m.getCurrentRegistrationPath(testDomain)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if path.Base(s) != testPath {
		t.Errorf("unexpected string %s", s)
	}

	// should ignore this
	p := testDomain + ".json"
	err = ioutil.WriteFile(path.Join(dir, tlsDir, p), []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	s, err = m.getCurrentRegistrationPath(testDomain)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if path.Base(s) != testPath {
		t.Errorf("unexpected string %s", s)
	}
}

func TestNextTLSPath(t *testing.T) {
	e1 := "example.com-2017-01-01-v000.json"
	e2 := "example.com-2017-01-01-v001.json"
	now = testNow
	m, _ := setup(t)
	s, err := m.getNextRegistrationPath(testDomain)
	if err != nil {
		t.Fatal(err)
	}
	if path.Base(s) != e1 {
		t.Errorf("unexpected next path: %s", path.Base(s))
	}

	err = ioutil.WriteFile(s, []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	s, err = m.getNextRegistrationPath(testDomain)
	if err != nil {
		t.Fatal(err)
	}
	if path.Base(s) != e2 {
		t.Errorf("unexpected next path: %s", path.Base(s))
	}
}
