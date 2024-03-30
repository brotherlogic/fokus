package server

import (
	"sort"

	pb "github.com/brotherlogic/fokus/proto"
	githubridgeclient "github.com/brotherlogic/githubridge/client"

	"context"

	ghbpb "github.com/brotherlogic/githubridge/proto"
)

type Overdue struct {
	client githubridgeclient.GithubridgeClient
}

func (o *Overdue) getName() string {
	return "Overdue"
}

func (o *Overdue) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_CODING_TASKS
}

func (o *Overdue) getFokus(ctx context.Context) (*pb.Focus, error) {
	issues, err := o.client.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(issues.Issues, func(i, j int) bool { return issues.Issues[i].GetId() < issues.Issues[j].GetId() })

	return &pb.Focus{
		Type:   o.getType(),
		Detail: issues.Issues[0].GetTitle(),
	}, nil
}
