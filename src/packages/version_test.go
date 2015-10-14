package packages

import (
	"testing"
)

type TestCase struct {
	Version string
	Expect  string
}

var testCaseVersionNormalizedTag []TestCase
var testCasePrepareTagVersion []TestCase
var testCasePrepareTagVersionNormalized []TestCase

var testCaseVersionNormalizedBranch []TestCase
var testCasePrepareBranchVersion []TestCase

func init() {
	testCaseVersionNormalizedTag = []TestCase{
		TestCase{"1.0.0", "1.0.0.0"},
		TestCase{"v0.2.0", "0.2.0.0"},
		TestCase{"v1.0.0-RC1", "1.0.0.0-RC1"},
		TestCase{"v1.2.0-alpha2", "1.2.0.0-alpha2"},
	}
	testCasePrepareTagVersion = []TestCase{
		TestCase{"1.0.0", "1.0.0"},
		TestCase{"v0.2.0", "v0.2.0"},
		TestCase{"v1.0.0-RC1", "v1.0.0-RC1"},
		TestCase{"v1.2.0-alpha2", "v1.2.0-alpha2"},
	}
	testCasePrepareTagVersionNormalized = []TestCase{
		TestCase{"1.0.0.0", "1.0.0.0"},
		TestCase{"0.2.0.0", "0.2.0.0"},
		TestCase{"v1.0.0-RC1", "v1.0.0-RC1"},
		TestCase{"v1.2.0-alpha2", "v1.2.0-alpha2"},
	}

	testCaseVersionNormalizedBranch = []TestCase{
		TestCase{"master", "9999999-dev"},
		TestCase{"dev", "dev-dev"},
		TestCase{"develop-2.3.0", "dev-develop-2.3.0"},
		TestCase{"pr/171", "dev-pr/171"},
		TestCase{"revert-354-master", "dev-revert-354-master"},
		TestCase{"v2.0", "2.0.9999999.9999999-dev"},
		TestCase{"issue/126", "dev-issue/126"},
		TestCase{"dev-1.7", "dev-dev-1.7"},
	}
	testCasePrepareBranchVersion = []TestCase{
		TestCase{"master", "dev-master"},
		TestCase{"dev", "dev-dev"},
		TestCase{"develop-2.3.0", "dev-develop-2.3.0"},
		TestCase{"pr/171", "dev-pr/171"},
		TestCase{"revert-354-master", "dev-revert-354-master"},
		TestCase{"v2.0", "2.0.x-dev"},
		TestCase{"issue/126", "dev-issue/126"},
		TestCase{"dev-1.7", "dev-dev-1.7"},
	}
}

func TestNormalize(t *testing.T) {
	for _, tCase := range testCaseVersionNormalizedTag {
		res := VersionNormalizedTag(tCase.Version)
		if res != tCase.Expect {
			t.Fatalf("%s != %s", res, tCase.Expect)
		}
	}
}

func TestPrepareTagVersion(t *testing.T) {
	for _, tCase := range testCasePrepareTagVersion {
		res := PrepareTagVersion(tCase.Version)
		if res != tCase.Expect {
			t.Fatalf("%s != %s", res, tCase.Expect)
		}
	}
}

func TestPrepareTagVersionNormalized(t *testing.T) {
	for _, tCase := range testCasePrepareTagVersionNormalized {
		res := PrepareTagVersionNormalized(tCase.Version)
		if res != tCase.Expect {
			t.Fatalf("%s != %s", res, tCase.Expect)
		}
	}
}

func TestNormalizeBranch(t *testing.T) {
	for _, tCase := range testCaseVersionNormalizedBranch {
		res := VersionNormalizedBranch(tCase.Version)
		if res != tCase.Expect {
			t.Fatalf("%s != %s", res, tCase.Expect)
		}
	}
}

func TestPrepareBranchVersion(t *testing.T) {
	for _, tCase := range testCasePrepareBranchVersion {
		res := PrepareBranchVersion(tCase.Version)
		if res != tCase.Expect {
			t.Fatalf("%s != %s", res, tCase.Expect)
		}
	}
}
