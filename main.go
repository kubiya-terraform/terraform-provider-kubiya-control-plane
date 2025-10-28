package main

import (
	"context"
	"log"

	"kubiya-control-plane/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const (
	version = "dev"
	address = "hashicorp.com/kubiya/kubiya-control-plane"
)

func main() {
	ctx := context.Background()
	kubiyaProvider := provider.New(version)

	opts := providerserver.ServeOpts{
		Address: address,
	}

	if err := providerserver.Serve(ctx, kubiyaProvider, opts); err != nil {
		log.Fatal(err.Error())
	}
}
