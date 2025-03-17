package paths

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

func SpecificationsPath() string {
	return filepath.Join(xdg.ConfigHome, "gh-localspace", "specifications")
}

func SpecificationPath(name string) string {
	return filepath.Join(SpecificationsPath(), name, "specification.json")
}

func DevcontainersPath() string {
	return filepath.Join(xdg.ConfigHome, "Code", "User", "globalStorage", "ms-vscode-remote.remote-containers", "configs")
}

func DevcontainerPath(name string) string {
	return filepath.Join(DevcontainersPath(), name)
}
