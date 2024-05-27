package server

import (
	"fmt"
	"sort"
	"strings"

	pb "github.com/brotherlogic/fokus/proto"
	githubridgeclient "github.com/brotherlogic/githubridge/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	sort.SliceStable(issues.Issues, func(i, j int) bool { return issues.Issues[i].GetOpenedDate() < issues.Issues[j].GetOpenedDate() })

	for _, issue := range issues.Issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() != "bandcampserver" && issue.GetRepo() != "recordalerting" {
				if !strings.Contains(issue.GetTitle(), "Incomplete Order") {
					return &pb.Focus{
						Type:   o.getType(),
						Detail: fmt.Sprintf("%v [%v] -> %v (%v)", issue.GetTitle(), issue.GetId(), issue.GetState(), issue.GetRepo()),
					}, nil
				}
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
