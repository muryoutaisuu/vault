package main

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
	vh "github.com/muryoutaisuu/vaulthelper"
	pfvault "github.com/postfinance/vault/kv"
)

func main() {
	c, err := getClient()
	if err != nil {
		fmt.Printf("Encountered error: %v\n")
	} else {
		paths := [...]string{"/", "hello", "hello/foo", "subdir", "subdir/mury", "subdir/mury/foo2", "asdf", "subdir/asdf", "subdir/mury/asdf"}
		for _, v := range paths {
			fmt.Println("")
			fmt.Printf("'%s' IsPath: %v\n", v, vh.IsPath(c, v))
			fmt.Printf("'%s' IsSecret: %v\n", v, vh.IsSecret(c, v))
			fmt.Printf("'%s' IsKey: %v\n", v, vh.IsKey(c, v))
		}
	}
}

func getClient() (*pfvault.Client, error) {
	conf := api.DefaultConfig()
	conf.Address = "http://127.0.0.1:8200"
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	o, err := ioutil.ReadFile(filepath.Join(u.HomeDir, ".vault-roleid"))
	if err != nil {
		return nil, err
	}
	approleId := strings.TrimSuffix(string(o), "\n")
	return vh.GetClient(conf, approleId)
}
