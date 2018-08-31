package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/leboncoin/terraform-provider-opennebula/opennebula"
	//"./opennebula" // to test local
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opennebula.Provider,
	})
}
