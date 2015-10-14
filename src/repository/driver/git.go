package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/compotlab/synopsis/src/packages"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Git struct {
	Name              string
	RepoDir           string
	Version           string
	VersionNormalized string
	Reference         string
	Source            map[string]string
	Dist              map[string]string
	Url               string
	Packages          map[string]map[string]packages.Composer
}

func (git *Git) Run() error {
	err := git.prepareRepoDir()
	if err != nil {
		return err
	} else {
		err = git.prepareMainBranch()
		if err != nil {
			return err
		}
		_, ok := git.Packages[git.Name]
		if !ok {
			git.Packages[git.Name] = make(map[string]packages.Composer)
			err = git.prepareTags()
			if err != nil {
				return err
			}
			err = git.prepareBranches()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (git *Git) GetName() string {
	return git.Name
}

func (git *Git) GetSource() map[string]string {
	return git.Source
}

func (git *Git) GetReference() string {
	return git.Reference
}

func (git *Git) prepareRepoDir() error {
	re := regexp.MustCompile("[^a-z0-9.]")
	dir := re.ReplaceAllString(git.Url, "-")
	git.RepoDir = os.Getenv("HOME") + packages.COMPOSER_CACHE_VCS + dir
	_, err := os.Stat(git.RepoDir)
	if err != nil {
		cmd := exec.Command("git", "clone", "--mirror", git.Url, git.RepoDir)
		cmd.Dir = os.Getenv("HOME") + packages.COMPOSER_CACHE_VCS
		_, err = cmd.CombinedOutput()
		if err != nil {
			return errors.New(fmt.Sprintf("git clone --mirror %s %s. %s", git.Url, git.RepoDir, err))
		}
	} else {
		cmd := exec.Command("git", "remote", "set-url", "origin", git.Url)
		cmd.Dir = git.RepoDir
		_, err = cmd.CombinedOutput()
		if err != nil {
			return errors.New(fmt.Sprintf("git remote set-url origin %s. %s", git.Url, err))
		}
		cmd = exec.Command("git", "remote", "update", "--prune", "origin")
		cmd.Dir = git.RepoDir
		_, err = cmd.CombinedOutput()
		if err != nil {
			return errors.New(fmt.Sprintf("git remote update --prune origin. %s. %s", git.Url, err))
		}
	}
	return nil
}

func (git *Git) prepareMainBranch() error {
	cmd := exec.Command("git", "branch", "--no-color")
	cmd.Dir = git.RepoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("git branch --no-color %s", err))
	}
	re := regexp.MustCompile("\\* +(\\S+)")
	response := re.FindStringSubmatch(string(out))
	composer, err := git.getComposerInformation(response[1])
	if err != nil {
		return err
	}
	git.Name = composer.Name
	return nil
}

func (git *Git) prepareTags() error {
	cmd := exec.Command("git", "show-ref", "--tags")
	cmd.Dir = git.RepoDir
	out, _ := cmd.CombinedOutput()
	for _, tag := range strings.SplitAfter(string(out), "\n") {
		re := regexp.MustCompile("^([a-f0-9]{40}) refs/tags/(\\S+)")
		response := re.FindStringSubmatch(tag)
		if response != nil {
			git.Version = packages.PrepareTagVersion(response[2])
			nVersion := packages.VersionNormalizedTag(response[2])
			git.VersionNormalized = packages.PrepareTagVersionNormalized(nVersion)
			git.Reference = response[1]
			p, err := git.getComposerInformation(response[1])
			if err != nil {
				return err
			}
			git.Packages[git.Name][git.Version] = p
		}
	}
	return nil
}

func (git *Git) prepareBranches() error {
	cmd := exec.Command("git", "branch", "--no-color", "--no-abbrev", "-v")
	cmd.Dir = git.RepoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("git branches --no-color --no-abbrev -v. %s", err))
	}
	for _, branch := range strings.SplitAfter(string(out), "\n") {
		re := regexp.MustCompile("(?:\\* )? *(\\S+) *([a-f0-9]+)(?: .*)?")
		response := re.FindStringSubmatch(branch)
		if response != nil {
			git.Version = packages.PrepareBranchVersion(response[1])
			git.VersionNormalized = packages.VersionNormalizedBranch(response[1])
			git.Reference = response[2]
			p, err := git.getComposerInformation(response[2])
			if err != nil {
				return err
			}
			git.Packages[git.Name][git.Version] = p
		}
	}
	return nil
}

func (git *Git) getComposerInformation(ref string) (packages.Composer, error) {
	co := new(packages.Composer)
	cmd := exec.Command("git", "show", ref+":composer.json")
	cmd.Dir = git.RepoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return *co, errors.New(fmt.Sprintf("git show %s. %s", ref+":composer.json", err))
	}
	json.Unmarshal(out, co)
	if co.Time == "" {
		cmd = exec.Command("git", "log", "-1", "--format=%at", ref)
		cmd.Dir = git.RepoDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			return *co, errors.New(fmt.Sprintf("git log -1 --format. %s", err))
		}
		t, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
		co.Time = time.Unix(t, 0).String()
	}
	git.Source = map[string]string{
		"type":      "git",
		"url":       git.Url,
		"reference": git.Reference,
	}
	co.Source = git.Source
	co.Version = git.Version
	co.VersionNormalized = git.VersionNormalized
	return *co, nil
}
