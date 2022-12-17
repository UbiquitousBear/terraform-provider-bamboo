package main

import (
	"context"
	"terraform-provider-bamboo/bamboo"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), bamboo.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/UbiquitousBear/bamboo",
	})
}
