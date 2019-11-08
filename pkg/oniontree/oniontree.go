package oniontree

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/onionltd/oniontree-tools/pkg/types/service"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	cairnName string = ".oniontree"
)

var (
	ErrCairnNotFound = errors.New("not an oniontree repository")
	ErrIdExists      = errors.New("id exists")
	ErrIdNotExists   = errors.New("id not exists")
)

type OnionTree struct {
	dir    string
	format string
}

func (o OnionTree) Init() error {
	for _, subdir := range []string{".", "unsorted", "tagged"} {
		if err := os.Mkdir(path.Join(o.dir, subdir), 0755); err != nil {
			if os.IsExist(err) {
				continue
			}
			return err
		}
	}
	pth := path.Join(o.dir, cairnName)
	cairnFile, err := os.Create(pth)
	if err != nil {
		return err
	}
	return cairnFile.Close()
}

func (o OnionTree) Add(id string, s service.Service) error {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return err
	}
	pth := path.Join(unsorted, o.idToFilename(id))
	if isFile(pth) {
		return ErrIdExists
	}
	file, err := os.Create(pth)
	if err != nil {
		return err
	}
	data, err := o.marshalData(s)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return file.Close()
}

func (o OnionTree) Edit(id string, s service.Service) error {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return err
	}
	pth := path.Join(unsorted, o.idToFilename(id))
	if !isFile(pth) {
		return ErrIdNotExists
	}
	file, err := os.Create(pth)
	if err != nil {
		return err
	}
	data, err := o.marshalData(s)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return file.Close()
}

func (o OnionTree) Tag(id string, tags []string) error {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return err
	}
	tagged, err := o.getTaggedDir()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		pth := path.Join(unsorted, o.idToFilename(id))
		if !isFile(pth) {
			return ErrIdNotExists
		}
		pthTag := path.Join(tagged, tag)
		// Create tag directory, ignore error if it already exists.
		if err := os.Mkdir(pthTag, 0755); err != nil {
			if !os.IsExist(err) {
				return err
			}
		}
		pthRel, err := filepath.Rel(pthTag, pth)
		if err != nil {
			return err
		}
		if err := os.Symlink(pthRel, path.Join(pthTag, path.Base(pth))); err != nil {
			if !os.IsExist(err) {
				return err
			}
		}
	}
	return nil
}

func (o OnionTree) Get(id string) (service.Service, error) {
	data, err := o.GetRaw(id)
	if err != nil {
		return service.Service{}, err
	}
	s := service.Service{}
	if err := o.unmarshalData(data, &s); err != nil {
		return service.Service{}, err
	}
	return s, nil
}

func (o OnionTree) GetRaw(id string) ([]byte, error) {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return nil, err
	}
	pth := path.Join(unsorted, o.idToFilename(id))
	if !isFile(pth) {
		return nil, ErrIdNotExists
	}
	file, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := []byte{}
	buff := make([]byte, 15535)
	for {
		num, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		data = append(data, buff[:num]...)
	}
	return data, nil
}

func (o OnionTree) List() ([]string, error) {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(unsorted)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	files, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	for idx, _ := range files {
		files[idx] = o.filenameToId(files[idx])
	}
	return files, nil
}

func (o OnionTree) ListTags() ([]string, error) {
	tagged, err := o.getTaggedDir()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(tagged)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	files, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (o OnionTree) getUnsortedDir() (string, error) {
	root, err := findRootDir(o.dir)
	if err != nil {
		return "", err
	}
	return path.Join(root, "unsorted"), nil
}

func (o OnionTree) getTaggedDir() (string, error) {
	root, err := findRootDir(o.dir)
	if err != nil {
		return "", err
	}
	return path.Join(root, "tagged"), nil
}

func (o OnionTree) marshalData(data interface{}) (b []byte, err error) {
	switch o.format {
	case "yaml":
		b, err = yaml.Marshal(data)
	default:
		panic("unsupported format")
	}
	return
}

func (o OnionTree) unmarshalData(b []byte, data interface{}) (err error) {
	switch o.format {
	case "yaml":
		err = yaml.Unmarshal(b, data)
	default:
		panic("unsupported format")
	}
	return
}

func (o OnionTree) idToFilename(id string) string {
	return fmt.Sprintf("%s.%s", id, o.format)
}

func (o OnionTree) filenameToId(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func findRootDir(dir string) (string, error) {
	for {
		pth := path.Join(dir, cairnName)
		match, err := filepath.Glob(pth)
		if err != nil {
			return "", err
		}
		if len(match) > 0 {
			return path.Dir(pth), nil
		}
		if dir == "/" {
			break
		}
		dir = path.Join(dir, "..")
	}
	return "", ErrCairnNotFound
}

func New(dir string) *OnionTree {
	return &OnionTree{dir: dir, format: "yaml"}
}
