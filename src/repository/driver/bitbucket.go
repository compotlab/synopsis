package driver

import (
	"encoding/json"
	"github.com/compotlab/synopsis/src/packages"
	"io/ioutil"
	"net/http"
	"regexp"
)

type BitBucket struct {
	Owner      string
	Repository string
	Version    string
	Reference  string
	Source     map[string]string
	Dist       map[string]string
	Url        string
	Packages   map[string]map[string]packages.Composer
}

func (bucket *BitBucket) Run() error {
	bucket.prepareRepository()
	if err := bucket.prepareMainBranch(); err != nil {
		return err
	}
	return nil
}

func (bucket *BitBucket) GetName() string {
	return bucket.Owner
}

func (bucket *BitBucket) GetSource() map[string]string {
	return bucket.GetSource()
}

func (bucket *BitBucket) GetReference() string {
	return bucket.Reference
}

func (bucket *BitBucket) prepareRepository() {
	re := regexp.MustCompile("^https?://bitbucket\\.org/([^/]+)/(.+?)\\.git$")
	response := re.FindStringSubmatch(bucket.Url)
	bucket.Owner, bucket.Repository = response[1], response[2]
}

func (bucket *BitBucket) prepareMainBranch() error {
	url := "https://api.bitbucket.org/1.0/repositories/" + bucket.Owner + "/" + bucket.Repository
	var r interface{}
	_, err := httpGet(url, r)
	if err != nil {
		return err
	}
	return nil
}

func (bucket *BitBucket) getComposerInformation(commit string) packages.Composer {
	return packages.Composer{}
}

func (bucket *BitBucket) prepareTags() {

}

func (bucket *BitBucket) prepareBranches() {

}

func (bucket *BitBucket) prepareSource(commit string) {

}

func httpGet(url string, result interface{}) (interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &result)
	return result, nil
}
