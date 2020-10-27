package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
	vh "github.com/muryoutaisuu/vaulthelper"
	pfvault "github.com/postfinance/vault/kv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s VAULT_ADDR VAULT_TOKEN PATH...\n", os.Args[0])
		fmt.Printf("A good PATH List would be: \"/\" \"hello\" \"hello/foo\" \"subdir\" \"subdir/mury\" \"subdir/mury/foo2\" \"asdf\" \"subdir/asdf\" \"subdir/mury/asdf\"")
		fmt.Printf("To get VAULT_TOKEN, do: vault write auth/approle/login role_id=$ROLEID\n")
		os.Exit(1)
	}
	fmt.Println(os.Args[0])
	c, err := getClient()
	c.SetToken(os.Args[2])
	if err != nil {
		fmt.Printf("Encountered error: %v\n", err)
	} else {
		fmt.Printf("pfvc: \"%v\"\n", c)
		//fmt.Printf("pfvc.client: \"%v\"\n", c.Client())
		//paths := [...]string{"/", "hello", "hello/foo", "subdir", "subdir/mury", "subdir/mury/foo2", "asdf", "subdir/asdf", "subdir/mury/asdf"}
		paths := os.Args[3:]
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
	conf.Address = os.Args[1]
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	o, err := ioutil.ReadFile(filepath.Join(u.HomeDir, ".vault-roleid"))
	if err != nil {
		return nil, err
	}
	approleId := strings.TrimSuffix(string(o), "\n")
	fmt.Printf("approleId=\"%v\"\n", approleId)
	return vh.GetClient(conf, approleId)
}
