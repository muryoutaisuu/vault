package vaulthelper

import (
	"github.com/hashicorp/vault/api"
	pfkv "github.com/postfinance/vaultkv"
)

// GetClient returns a postfinance vault client.
func GetClient(conf *api.Config, approleid string) (*pfkv.Client, error) {
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
	// Create new pfkv client with vault client
	pfc, err := pfkv.New(vc, "secret/") //PostFinanceClient
	if err != nil {
		return nil, err
	}
	// Return PostFinanceClient
	return pfc, err
}
