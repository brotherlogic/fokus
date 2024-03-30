package server

import githubridgeclient "github.com/brotherlogic/githubridge/client"
import pb "github.com/brotherlogic/fokus/proto"
import ghbpb "github.com/brotherlogic/githubridge/proto"
import "context"

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
		return nil,err
	}
}