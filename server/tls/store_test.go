package tls

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

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
