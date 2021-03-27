package gitlab

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	//"github.com/xanzy/go-gitlab"
)

// Plugin returns this plugin
func Plugin(context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-gitlab",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"gitlab_group": _groupTable,
			//"github_repository":       tableGitHubRepository(),
			//"github_repository_issue": tableGitHubRepositoryIssue(),
			//"github_team":             tableGitHubTeam(),
			//"github_user":             tableGitHubUser(),
		},
	}
	return p
}
