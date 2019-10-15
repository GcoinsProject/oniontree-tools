package oniontree

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	cairnName string = ".oniontree"
	entryExt  string = ".yaml"
)

var (
	ErrCairnNotFound = errors.New("not an oniontree repository")
	ErrIdExists      = errors.New("id exists")
	ErrIdNotExists   = errors.New("id not exists")
)

type OnionTree struct {
	dir string
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

func (o OnionTree) Add(id string, data []byte) error {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return err
	}
	pth := path.Join(unsorted, id+entryExt)
	if isFile(pth) {
		return ErrIdExists
	}
	file, err := os.Create(pth)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return file.Close()
}

func (o OnionTree) Edit(id string, data []byte) error {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return err
	}
	pth := path.Join(unsorted, id+entryExt)
	if !isFile(pth) {
		return ErrIdNotExists
	}
	file, err := os.Create(pth)
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
		pth := path.Join(unsorted, id+entryExt)
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

func (o OnionTree) Get(id string) ([]byte, error) {
	unsorted, err := o.getUnsortedDir()
	if err != nil {
		return nil, err
	}
	pth := path.Join(unsorted, id+entryExt)
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
	return &OnionTree{dir: dir}
}
