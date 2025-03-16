package secrets

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetSecrets(repository string) (map[string]string, error) {
	githubClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating GitHub client")
	}

	secrets := make(map[string]string)

	// We get the list of secret names from the API rather than trying to filter the list direct from Codespaces.
	page := 1
	for {
		secretsResponse := &secretsResponse{}
		parameters := url.Values{}
		parameters.Set("per_page", "100")
		parameters.Set("page", fmt.Sprintf("%d", page))
		err = githubClient.Get(fmt.Sprintf("repos/%s/codespaces/secrets?%s", repository, parameters.Encode()), secretsResponse)
		if err != nil {
			return nil, errors.Wrap(err, "error fetching secrets")
		}
		if len(secretsResponse.Secrets) == 0 {
			break
		}
		for _, secret := range secretsResponse.Secrets {
			secrets[secret.Name] = ""
		}
		page++
	}

	createCodespaceRequest := &createCodespaceRequest{
		MultiRepoPermissionsOptOut: true,
	}
	createCodespaceRequestBody, err := json.Marshal(createCodespaceRequest)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling create codespace request")
	}
	createCodespaceResponse := &createCodespaceResponse{}

	err = githubClient.Post(fmt.Sprintf("repos/%s/codespaces", repository), bytes.NewReader(createCodespaceRequestBody), createCodespaceResponse)
	if err != nil {
		return nil, errors.Wrap(err, "error creating codespace")
	}
	defer func() {
		err := githubClient.Delete(fmt.Sprintf("user/codespaces/%s", createCodespaceResponse.Name), nil)
		if err != nil {
			logrus.Warnf("error deleting codespace %s: %v\n", createCodespaceResponse.Name, err)
		}
	}()

	stdout, stderr, err := gh.Exec("codespace", "ssh", "--codespace", createCodespaceResponse.Name, "--", "cat", "/workspaces/.codespaces/shared/.env-secrets")
	if err != nil {
		return nil, errors.Wrapf(err, "error executing codespace ssh: %s", stderr.String())
	}

	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		nameValue := strings.SplitN(line, "=", 2)
		if _, ok := secrets[nameValue[0]]; ok {
			value, err := base64.StdEncoding.DecodeString(nameValue[1])
			if err != nil {
				return nil, errors.Wrapf(err, "error decoding secret %s", nameValue[0])
			}
			secrets[nameValue[0]] = string(value)
		}
	}

	return secrets, nil
}

type secret struct {
	Name string `json:"name"`
}

type secretsResponse struct {
	Secrets []secret `json:"secrets"`
}

type createCodespaceRequest struct {
	MultiRepoPermissionsOptOut bool `json:"multi_repo_permissions_opt_out"`
}

type createCodespaceResponse struct {
	Name string `json:"name"`
}
