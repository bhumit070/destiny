package config

import (
	"strings"
	"testing"
)

func ValuesHaveTrailingSlash(t *testing.T) {
	for extension, folder := range FolderGroups {
		if !strings.HasSuffix(folder, "/") {
			t.Errorf("Extension %s does not have a trailing slash in folder path", extension)
		}
	}
}
