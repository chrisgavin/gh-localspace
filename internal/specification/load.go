package specification

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func Load(path string) (*Specification, error) {
	specification := &Specification{
		Base: ".devcontainer/devcontainer.json",
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(specification); err != nil {
		return nil, err
	}

	if strings.HasPrefix(specification.Root, "~/") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Wrap(err, "error getting user home directory")
		}
		specification.Root = filepath.Join(dirname, specification.Root[2:])
	}

	return specification, nil
}
