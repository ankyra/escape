package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type VersionQuery struct {
	LatestVersion   bool
	VersionPrefix   string
	SpecificVersion string
	SpecificTag     string
}

func ParseVersionQuery(v string) (*VersionQuery, error) {
	if v == "" {
		return nil, fmt.Errorf("Empty version query")
	}
	if v == "latest" || v == "@" || v == "v@" {
		return &VersionQuery{
			LatestVersion: true,
		}, nil
	} else if strings.HasPrefix(v, "v") {
		vq := maybeParseVersionQuery(v[1:])
		if vq != nil {
			return vq, nil
		}
	} else if unicode.IsDigit(rune(v[0])) {
		vq := maybeParseVersionQuery(v)
		if vq != nil {
			return vq, nil
		}
	}
	return &VersionQuery{
		SpecificTag: v,
	}, nil
}

func (v *VersionQuery) ToString() string {
	if v.LatestVersion {
		return "latest"
	} else if v.VersionPrefix != "" {
		return "v" + v.VersionPrefix + "@"
	} else if v.SpecificVersion != "" {
		return "v" + v.SpecificVersion
	}
	return v.SpecificTag
}

func (v *VersionQuery) ToVersionSuffix() string {
	sep := "-"
	if v.SpecificTag != "" {
		sep = ":"
	}
	return sep + v.ToString()
}

func maybeParseVersionQuery(versionQuery string) *VersionQuery {
	if versionQuery == "" {
		return nil
	}
	parts := strings.Split(versionQuery, ".")
	if parts[len(parts)-1] == "@" {
		return &VersionQuery{
			VersionPrefix: versionQuery[:len(versionQuery)-1],
		}
	}
	if isSemanticVersion(parts) {
		return &VersionQuery{
			SpecificVersion: versionQuery,
		}
	}
	return nil
}

func isSemanticVersion(parts []string) bool {
	legitVersion := true
	for _, part := range parts {
		i, err := strconv.Atoi(part)
		if err != nil || i < 0 {
			legitVersion = false
			break
		}
	}
	return legitVersion
}
