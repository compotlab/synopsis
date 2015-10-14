package packages

import (
	"regexp"
	"strings"
)

var modifierRegexp = `[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?`

func VersionNormalizedTag(version string) string {
	var v string
	var index int
	version = strings.TrimSpace(version)
	if "dev-" == strings.ToLower(version[:4]) {
		return "dev-" + version[4:]
	}
	// match classical number version
	re := regexp.MustCompile(`^v?(\d{1,3})(\.\d+)?(\.\d+)?(\.\d+)?` + modifierRegexp)
	response := re.FindStringSubmatch(version)
	if response != nil {
		index = 5
		for i := 1; i < 5; i++ {
			if response[i] != "" {
				v += response[i]
			} else {
				v += ".0"
			}
		}
	}
	if index > 0 {
		if len(response) >= index && response[index] != "" {
			if response[index] != "stable" {
				v += "-" + expandStability(response[index])
				if len(response) >= index+1 && response[index+1] != "" {
					v += response[index+1]
				}
				if len(response) >= index+2 && response[index+2] != "" {
					v += "-dev"
				}
			}
		}
		return v
	}
	// match dev branches
	re = regexp.MustCompile("(.*?)[.-]?dev$")
	response = re.FindStringSubmatch(version)
	if response != nil {
		return VersionNormalizedBranch(response[1])
	}
	// return default version
	return version
}

func VersionNormalizedBranch(version string) string {
	version = strings.TrimSpace(version)
	// match master-like branch
	re := regexp.MustCompile(`^(?:dev-)?(?:master|trunk|default)$`)
	response := re.FindStringSubmatch(version)
	if response != nil {
		return "9999999-dev"
	}
	// match dev branches
	re = regexp.MustCompile(`^v?(\d+)(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?$`)
	response = re.FindStringSubmatch(version)
	if response != nil {
		var v string
		for i := 1; i < 5; i++ {
			if response[i] != "" {
				response[i] = strings.Replace(response[i], "*", "x", -1)
				response[i] = strings.Replace(response[i], "X", "x", -1)
				v += response[i]
			} else {
				v += ".x"
			}
		}
		return strings.Replace(v, "x", "9999999", -1) + "-dev"
	}
	// return default version
	return "dev-" + version
}

func PrepareTagVersion(version string) string {
	reg, _ := regexp.Compile("[.-]?dev$")
	v := reg.ReplaceAllString(version, "")
	return v
}

func PrepareTagVersionNormalized(nVersion string) string {
	reg, _ := regexp.Compile("(^dev-|[.-]?dev$)")
	v := reg.ReplaceAllString(nVersion, "")
	return v
}

func PrepareBranchVersion(version string) string {
	nVersion := VersionNormalizedBranch(version)
	var v string
	if "dev-" == nVersion[:4] || "9999999-dev" == nVersion {
		v = "dev-" + version
	} else {
		reg, _ := regexp.Compile(`(\.9{7})+`)
		v = reg.ReplaceAllString(nVersion, ".x")
	}
	return v
}

func expandStability(stability string) string {
	stability = strings.ToLower(stability)
	switch stability {
	case "a":
		return "alpha"
	case "b":
		return "beta"
	case "p":
	case "pl":
		return "patch"
	case "rc":
		return "RC"
	}
	return stability
}
