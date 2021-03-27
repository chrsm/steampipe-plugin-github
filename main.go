package main

import (
	"github.com/chrsm/steampipe-plugin-gitlab/gitlab"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: gitlab.Plugin})
}
