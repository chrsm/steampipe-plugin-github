package gitlab

import (
	"context"
	"os"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/xanzy/go-gitlab"
)

var _groupTable = &plugin.Table{
	Name:        "gitlab_group",
	Description: "Gitlab Groups",

	List: &plugin.ListConfig{
		Hydrate: listGroup,
	},

	Get: &plugin.GetConfig{
		KeyColumns: plugin.SingleColumn("id"),
		Hydrate:    hydGroupDetail,
	},

	Columns: []*plugin.Column{
		{
			Name: "id",
			Type: proto.ColumnType_INT,
		},
		{
			Name: "name",
			Type: proto.ColumnType_STRING,
		},
		{
			Name: "path",
			Type: proto.ColumnType_STRING,
		},
		{
			Name: "description",
			Type: proto.ColumnType_STRING,
		},
		{
			Name: "parent_id",
			Type: proto.ColumnType_INT,
		},
		{
			Name: "created_at",
			Type: proto.ColumnType_DATETIME,
		},
	},
}

func connect(ctx context.Context, d *plugin.QueryData) *gitlab.Client {
	token := os.Getenv("GITLAB_TOKEN")
	cfg := GetConfig(d.Connection)
	if &cfg != nil {
		if cfg.Token != nil {
			token = *cfg.Token
		}
	}

	cli, err := gitlab.NewClient(token)
	if err != nil {
		panic(err)
	}

	return cli
}

func listGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cli := connect(ctx, d)

	opts := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	for {
		var (
			groups []*gitlab.Group
			resp   *gitlab.Response
		)

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error

			groups, resp, err = cli.Groups.ListGroups(opts)

			return err
		})

		if err != nil {
			return nil, err
		}

		for i := range groups {
			d.StreamListItem(ctx, groups[i])
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, nil
}

func hydGroupDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var gid int

	if h.Item != nil {
		group := h.Item.(*gitlab.Group)
		gid = group.ID
	} else {
		gid = int(d.KeyColumnQuals["id"].GetInt64Value())
	}

	cli := connect(ctx, d)

	var (
		group *gitlab.Group
		resp  *gitlab.Response
	)

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error

		group, resp, err = cli.Groups.GetGroup(gid)

		return err
	})

	if err != nil {
		return nil, err
	}

	return group, nil
}
