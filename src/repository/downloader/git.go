package downloader

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type Git struct {
	Name       string
	Version    string
	Url        string
	Source     map[string]string
	DistDir    string
	SourcePath string
	PathExist  bool
}

func (git *Git) Prepare() {
	re := regexp.MustCompile("[^a-z0-9-_]")
	dir := re.ReplaceAllString(git.Name, "-")
	h := sha1.New()
	h.Write([]byte(git.Source["reference"]))
	ref := hex.EncodeToString(h.Sum(nil))
	ref = ref[:6]
	temp := dir + "-" + git.Version + "-" + ref
	git.SourcePath = path.Join(git.DistDir, strings.Replace(temp, "/", "-", -1))
	git.PathExist = false
	if _, err := os.Stat(git.SourcePath); err == nil {
		git.PathExist = true
	}
}

func (git *Git) Run() error {
	cmd := exec.Command("git", "clone", "--no-checkout", git.Url, git.SourcePath)
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New(fmt.Sprintf("git clone --no-checkout %s %s. %s", git.Url, git.SourcePath, err))
	}

	cmd = exec.Command("git", "remote", "add", "composer", git.Url)
	cmd.Dir = git.SourcePath
	if _, err := cmd.CombinedOutput(); err != nil {
		return errors.New(fmt.Sprintf("git remote add composer %s. %s", git.Url, err))
	}

	cmd = exec.Command("git", "fecth", "composer")
	cmd.Dir = git.SourcePath
	cmd.CombinedOutput()

	cmd = exec.Command("git", "checkout", "-b", git.Version, "composer/"+git.Version)
	cmd.Dir = git.SourcePath
	cmd.CombinedOutput()

	cmd = exec.Command("git", "reset", "--hard", git.Source["reference"])
	cmd.Dir = git.SourcePath
	cmd.CombinedOutput()

	if _, err := os.Stat(git.SourcePath); err != nil {
		return err
	}
	git.PathExist = true
	return nil
}
