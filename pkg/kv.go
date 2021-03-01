package vaulthelper

import (
	"fmt"
	"path"

	//"github.com/hashicorp/vault/api"
	pfkv "github.com/postfinance/vault/kv"
)

// Filetype defines the type of the returned value element of vault
type Filetype byte

const (
	CNull   Filetype = 0 // not a valid vault element
	CPath   Filetype = 1 // exists in Vault as a directory
	CSecret Filetype = 2 // exists in Vault as a secret
	CKey    Filetype = 3 // Key of a key=value pair in a secret, secret/path/to/secret/CKEY
)

func IsPath(pfc *pfkv.Client, vpath string) bool {
	s, err := pfc.List(vpath)
	if err == nil && s != nil {
		return true
	}
	return false
}

func IsSecret(pfc *pfkv.Client, vpath string) bool {
	s, err := pfc.Read(vpath)
	if err == nil && s != nil {
		return true
	}
	return false
}

func IsKey(pfc *pfkv.Client, vpath string) bool {
	k := path.Base(vpath)
	vpath = path.Dir(vpath) // clip last element
	s, err := pfc.Read(vpath)
	if err == nil && s != nil {
		if _, ok := s[k]; ok {
			return true
		}
	}
	return false
}

func GetValueFromKey(pfc *pfkv.Client, vpath string) (string, error) {
	if IsKey(pfc, vpath) {
		s, _ := pfc.Read(path.Dir(vpath))
		return s[path.Base(vpath)].(string), nil
	}
	return "", fmt.Errorf("Key '%s' or Secret '%s' not found\n", path.Base(vpath), path.Dir(vpath))
}

// GetType returns type of the requested resource
// types will be the defined FileType byte constants on top of this file
func GetType(pfc *pfkv.Client, vpath string) Filetype {
	if IsPath(pfc, vpath) {
		return CPath
	}

	if IsSecret(pfc, vpath) {
		return CSecret
	}

	if IsKey(pfc, vpath) {
		return CKey
	}

	return CNull
}

// GetTypes returns similar to GetType the types of the requested resources
// imagine following situation:
//   secret/foo
//   secret/foo/
//   secret/foo/bar
// here foo is a secret as well as a subdirectory. It should be possible, to
// get both those types
func GetTypes(pfc *pfkv.Client, vpath string) map[Filetype]bool {
	r := make(map[Filetype]bool)

	r[CPath] = IsPath(pfc, vpath)
	r[CSecret] = IsSecret(pfc, vpath)

	// if else statement here is needed, case of:
	//   E1                           -> secret
	//   E1/mysecret = 42             -> key=mysecret value=42
	//   E1/                          -> subdir in Vault
	//   E1/subsecret                 -> secret
	//   E1/subsecret/mysecret = 43   -> key=mysecret value=43
	// this would have thrown an error, because for E1/subsecret/mysecret it would
	// have r[CKey] == true AND r[CValue] == true
	// this would cause errors in any further calculations
	if !r[CSecret] {
		r[CKey] = IsKey(pfc, vpath)
	} else {
		r[CKey] = false
	}

	return r
}
