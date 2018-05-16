package tls

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestGetAll(t *testing.T) {
	m, stop := setup(t)
	defer stop()
	for _, t := range []string{
		"example2.com-2017-01-01-v000.json",
		"example.com-2017-01-01-v000.json",
		"example.com-2017-01-01-v001.json",
	} {
		p := path.Join(m.Config.DataPath(tlsDir, ""), t)
		_ = ioutil.WriteFile(p, []byte{}, os.ModePerm)
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
	m, stop := setup(t)
	defer stop()
	s, err := m.getCurrentRegistrationPath(testDomain)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if s != "" {
		t.Errorf("unexpected string %s", s)
	}

	p1 := path.Join(m.Config.DataPath(tlsDir, ""), testPath)
	err = ioutil.WriteFile(p1, []byte{}, os.ModePerm)
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
	p2 := path.Join(m.Config.DataPath(tlsDir, ""), testDomain+".json")
	err = ioutil.WriteFile(p2, []byte{}, os.ModePerm)
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
	m, stop := setup(t)
	defer stop()
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
