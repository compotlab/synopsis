package src

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/compotlab/synopsis/src/utils"
	"github.com/compotlab/synopsis/src/packages"
	"github.com/compotlab/synopsis/src/repository"
	"github.com/compotlab/synopsis/src/repository/downloader"
	"github.com/compotlab/synopsis/src/repository/driver"
	"io"
	"math"
	"os"
	"path"
	"strings"
	"sync"
)

func (repo Repository) Run(pm map[string]map[string]packages.Composer, config *Config) {
	var vcsDriver driver.Driver
	switch repo.Type {
	case "git", "vcs", "composer":
		vcsDriver = &driver.Git{Url: repo.Url, Packages: pm}
		err := vcsDriver.Run()
		if err != nil {
			utils.Logger.Error(err)
		}
		Archivate(pm, config, vcsDriver, repo)
	case "bitbucket", "git-bitbucket":
		utils.Logger.Error("Bitbucket driver hasn't ready yet!")
	default:
		utils.Logger.Error("Has no driver!")
	}
}

func Archivate(pm map[string]map[string]packages.Composer, config *Config, d driver.Driver, repo Repository) {
	if config.File.Archive.Dir != "" {
		ch := make(chan bool, len(pm[d.GetName()]))
		wg := sync.WaitGroup{}
		wg.Add(len(pm[d.GetName()]))
		for _, value := range pm[d.GetName()] {
			go func(r Repository, p packages.Composer) {
				defer func() {
					<-ch
					wg.Done()
				}()
				ch <- true
				// Prepare path and archive path
				gitDownloader := &downloader.Git{
					Name:    p.Name,
					Version: p.Version,
					Url:     r.Url,
					Source:  p.Source,
					DistDir: config.DistDir,
				}
				gitDownloader.Prepare()
				ar := repository.Archivator{
					SourcePath: gitDownloader.SourcePath,
				}
				ar.Prepare()
				// Run git download
				if !gitDownloader.PathExist && !ar.ArchiveExist {
					if err := gitDownloader.Run(); err != nil {
						utils.Logger.Error(err)
					}
				}
				// Run create archive
				if gitDownloader.PathExist && !ar.ArchiveExist {
					if err := ar.Run(); err != nil {
						utils.Logger.Error(err)
					}
				}
				if ar.ArchiveExist {
					af := strings.Replace(ar.TargetPath, gitDownloader.DistDir, "", -1)
					pM := pm[d.GetName()][gitDownloader.Version]
					pM.Dist = map[string]string{
						"type":      ar.Format,
						"url":       config.File.Homepage + "/" + path.Join(config.File.Archive.Dir, af),
						"reference": d.GetReference(),
						"shasum":    hashFile(ar.TargetPath),
					}
					pm[d.GetName()][gitDownloader.Version] = pM
				}
			}(repo, value)
		}
		wg.Wait()
		close(ch)
	}
}

const FILE_CHUNK = 8192

func hashFile(f string) string {
	file, _ := os.Open(f)
	defer file.Close()
	info, _ := file.Stat()
	fileSize := info.Size()
	blocks := uint64(math.Ceil(float64(fileSize) / float64(FILE_CHUNK)))
	hash := sha1.New()
	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(FILE_CHUNK, float64(fileSize-int64(i*FILE_CHUNK))))
		buf := make([]byte, blockSize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
