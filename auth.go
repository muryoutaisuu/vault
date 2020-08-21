package vaulthelper

import (
	"errors"

	"github.com/hashicorp/vault/api"
)

// ApproleLogin performs an Approlelogin on a Client with an approleid and
// returns the accessToken
func ApproleLogin(c *api.Client, approleId string) (accessToken string, err error) {
	// Prepare data for Login request
  data := map[string]interface{}{
    "role_id": approleId,
  }
	// Perform Login
  resp, err := c.Logical().Write("auth/approle/login", data)
  if err != nil {
    return "", err
  }
  if resp.Auth == nil {
    return "", errors.New("no auth info returned")
  }
  return resp.Auth.ClientToken, nil
}
