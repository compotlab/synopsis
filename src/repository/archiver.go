package repository

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Archivator struct {
	Format       string
	TargetPath   string
	SourcePath   string
	ArchiveExist bool
}

func (archivator *Archivator) Prepare() {
	if archivator.Format == "" {
		archivator.Format = "zip"
	}
	archivator.TargetPath = archivator.SourcePath + "." + archivator.Format
	archivator.ArchiveExist = false
	if _, err := os.Stat(archivator.TargetPath); err == nil {
		archivator.ArchiveExist = true
	}
}

func (archivator *Archivator) Run() error {
	zipFile, err := os.Create(archivator.TargetPath)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		zipWriter.Close()
		os.RemoveAll(archivator.SourcePath)
	}()
	err = filepath.Walk(archivator.SourcePath, func(p string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			buffer, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			file := strings.Replace(p, archivator.SourcePath, "", -1)
			file = strings.Trim(file, "/")
			f, err := zipWriter.Create(file)
			if err != nil {
				return err
			}
			if _, err = f.Write([]byte(buffer)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	archivator.ArchiveExist = true
	return nil
}
