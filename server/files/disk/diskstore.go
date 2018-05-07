package disk

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/octavore/nagax/util/errors"
)

const (
	KiloByte = 1 << (10 * (iota + 1))
	MegaByte
	uploadMax = 16 * MegaByte
)

type DiskStore struct {
	dataDir string
}

func NewDiskStore(dataDir string) (*DiskStore, error) {
	fi, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dataDir, 0777)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		return &DiskStore{dataDir: dataDir}, nil
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if !fi.IsDir() {
		return nil, errors.New("%s exists but is not a folder", dataDir)
	}
	return &DiskStore{dataDir: dataDir}, nil
}

func (d *DiskStore) cleanPath(filename string) (string, error) {
	absDataDir, err := filepath.Abs(d.dataDir)
	if err != nil {
		return "", errors.Wrap(err)
	}
	fullPath, err := filepath.Abs(filepath.Join(d.dataDir, filename))
	if err != nil {
		return "", errors.Wrap(err)
	}
	if !strings.HasPrefix(fullPath, absDataDir) {
		return "", errors.New("invalid path")
	}
	return fullPath, nil
}

func (d *DiskStore) Get(uuid string) (io.ReadCloser, error) {
	p, err := d.cleanPath(uuid)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return f, nil
}

func (d *DiskStore) Upload(uuid string, r io.Reader) error {
	p, err := d.cleanPath(uuid)
	if err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return errors.Wrap(err)
	}

	rdr := io.LimitReader(r, uploadMax)
	_, err = io.Copy(f, rdr)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (d *DiskStore) Delete(uuid string) error {
	p, err := d.cleanPath(uuid)
	if err != nil {
		return err
	}
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(p)
}
