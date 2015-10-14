package packages

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Packages struct {
	Package map[string]map[string]Composer `json:"packages"`
}

type Composer struct {
	Name              string                 `json:"name"`
	Type              string                 `json:"type"`
	Description       string                 `json:"description"`
	Keywords          []string               `json:"keywords"`
	Homepage          string                 `json:"homepage"`
	Authors           []map[string]string    `json:"authors"`
	Require           map[string]string      `json:"require"`
	Version           string                 `json:"version"`
	VersionNormalized string                 `json:"version_normalized"`
	Source            map[string]string      `json:"source"`
	Dist              map[string]string      `json:"dist"`
	Time              string                 `json:"time"`
	Extra             map[string]interface{} `json:"extra"`
	InstallSource     string                 `json:"installation-source"`
	Stability         bool                   `json:"-"`
}

const COMPOSER_CACHE_VCS = "/.composer/cache/vcs/"

func init() {
	ccv := path.Join(os.Getenv("HOME"), COMPOSER_CACHE_VCS)
	if err := os.MkdirAll(ccv, 0777); err != nil {
		log.Fatal(err)
	}
}

func (packages *Packages) ToJson(output string) error {
	j, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(path.Join(output, "packages.json"), j, 0755); err != nil {
		return err
	}
	return nil
}
