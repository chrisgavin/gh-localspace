package specification

import (
	"encoding/json"
	"os"
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

	return specification, nil
}
