package config_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"go.avenga.cloud/couper/gateway/config"
)

func TestWriteGateway(t *testing.T) {
	conf := config.Gateway{
		Applications: []*config.Application{
			{
				BasePath: "/hans/v1",
				Name:     "wurst",
				Path: []*config.Path{
					{
						Backend: &config.Backend{Kind: "proxy"},
						Pattern: "/proxy/"},
				},
			},
		},
	}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&conf, f.Body())
	fmt.Printf("%s", f.Bytes())
}