package server

import (
	"fmt"
	"time"

	pb "github.com/brotherlogic/fokus/proto"
	githubridgeclient "github.com/brotherlogic/githubridge/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	ghbpb "github.com/brotherlogic/githubridge/proto"
)

type Cluster struct {
}

func (c *Cluster) getName() string {
	return "highlight"
}

func (c *Cluster) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_CLUSTER
}

func (c *Cluster) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {
	if true {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
	}

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() == "cluster" || issue.GetRepo() == "cluster2" {
				return &pb.Focus{
					Type:   c.getType(),
					Detail: fmt.Sprintf("%v [%v] -> %v (highlight)", issue.GetTitle(), issue.GetId(), issue.GetState()),
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
