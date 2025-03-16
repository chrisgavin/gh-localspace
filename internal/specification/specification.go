package specification

import jsonpatch "github.com/evanphx/json-patch"

type Specification struct {
	Repository           string          `json:"repository"`
	Base                 string          `json:"base"`
	Root                 string          `json:"root"`
	Patches              jsonpatch.Patch `json:"patches"`
	ImpersonateCodespace bool            `json:"impersonate_codespace"`
}
