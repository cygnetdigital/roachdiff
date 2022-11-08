package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// OpenFile at a particular gitref
func OpenFile(gitDir, subPath, gitRef string) ([]byte, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:./%s", gitRef, subPath))
	cmd.Dir = gitDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if bytes.Contains(out, []byte("exists on disk, but not in")) {
			return []byte{}, nil
		}
		if bytes.Contains(out, []byte(fmt.Sprintf("does not exist in '%s'", gitRef))) {
			return []byte{}, nil
		}
		return nil, fmt.Errorf("failed to get file: %w: %s", err, string(out))
	}

	return out, nil
}
