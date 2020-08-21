package vaulthelper

import (
	"github.com/hashicorp/vault/api"
	pfvault "github.com/postfinance/vault/kv"
)

// GetClient returns a postfinance vault client.
func GetClient(conf *api.Config, approleid string) (*pfvault.Client, error) {
  // Create new vault client with vault configuration
  vc, err := api.NewClient(conf) // VaultClient
  if err != nil {
    return nil, err
  }
  // Login with approleToken, get accessToken
  accessToken, err := ApproleLogin(vc, approleid)
  if err != nil {
    return nil, err
  }
  // Set accessToken
  vc.SetToken(accessToken)
  // Create new pfvault client with vault client
  pfc, err := pfvault.New(vc, "secret/") //PostFinanceClient
  if err != nil {
    return nil, err
  }
	// Return PostFinanceClient
  return pfc, err
}
