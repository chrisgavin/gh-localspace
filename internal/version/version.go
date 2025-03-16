package version

import (
	"runtime/debug"
)

type VersionInfo struct {
	Version string
	Commit  string
}

func Version() string {
	version := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	return version
}

func Commit() string {
	commit := ""
	dirty := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}
			if setting.Key == "vcs.modified" && setting.Value == "true" {
				dirty = "-dirty"
			}
		}
	}
	return commit + dirty
}
