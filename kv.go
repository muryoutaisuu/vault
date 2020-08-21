package vaulthelper

import (
	"path"

	//"github.com/hashicorp/vault/api"
	pfkv "github.com/postfinance/vault/kv"
)

// Filetype defines the type of the returned value element of vault
type Filetype byte

const (
	CNull  Filetype = 0 // not a valid vault element
	CPath  Filetype = 1 // exists in Vault as a directory or secret
	CKey   Filetype = 2 // Key of a key=value pair in a secret, secret/path/to/secret/CKEY
	CValue Filetype = 3 // Value of a key=value pair
)

func IsPath(pfc *pfkv.Client, vpath string) bool {
	s, err := pfc.List(vpath)
	if err == nil && s != nil {
		return true
	}
	return false
}

func IsKey(pfc *pfkv.Client, vpath string) bool {
	s, err := pfc.Read(vpath)
	if err == nil && s != nil {
		return true
	}
	return false
}

func IsValue(pfc *pfkv.Client, vpath string) bool {
	vpath = path.Dir(vpath) // clip last element
	s, err := pfc.Read(vpath)
	if err == nil && s != nil {
		return true
	}
	return false
}

// GetType returns type of the requested resource
// types will be the defined FileType byte constants on top of this file
func GetType(pfc *pfkv.Client, vpath string) Filetype {
	if IsPath(pfc, vpath) {
		return CPath
	}

	if IsKey(pfc, vpath) {
		return CKey
	}

	if IsValue(pfc, vpath) {
		return CValue
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
	r[CKey] = IsKey(pfc, vpath)

	// if else statement here is needed, case of:
	//   E1             							-> secret
	//   E1/mysecret = 42							-> key=mysecret value=42
	//   E1/            							-> subdir in Vault
	//   E1/subsecret   							-> secret
	//   E1/subsecret/mysecret = 43		-> key=mysecret value=43
	// this would have thrown an error, because for E1/subsecret/mysecret it would
	// have r[CKey] == true AND r[CValue] == true
	// this would cause errors in any further calculations
	if !r[CKey] {
		r[CValue] = IsValue(pfc, vpath)
	} else {
		r[CValue] = false
	}

	return r
}
