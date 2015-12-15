package plugin

import "github.com/cloudfoundry/cli/plugin"

type (
	CmdRunner interface {
		Run() error
	}
	logger interface {
		Fatalln(...interface{})
		Println(...interface{})
	}
	//DeployCloudPlugin - deploy cloud plugin object
	DeployCloudPlugin struct {
		conn       plugin.CliConnection
		list       *bool
		nomanifest *bool
		run        *string
		show       *string
		org        *string
		repo       *string
		branch     *string
		url        *string
		token      *string
		cfuser     *string
		cfpass     *string
		configFile *string
		Errors     []error
	}
)
