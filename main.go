package main

import (
	"github.com/cloudfoundry/cli/plugin"
	dcp "github.com/xchapter7x/deploycloud/plugin"
)

func main() {
	plugin.Start(new(dcp.DeployCloudPlugin))
}
