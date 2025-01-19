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

type Tasks struct {
}

func (h *Tasks) getName() string {
	return "tasks"
}

func (h *Tasks) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_TASKS
}

func (h *Tasks) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			foundLabel := false
			for _, label := range issue.GetLabels() {
				if label == "task" {
					foundLabel = true
				}
			}
			if foundLabel {
				return &pb.Focus{
					Type:   h.getType(),
					Detail: fmt.Sprintf("%v [%v]", issue.GetTitle(), issue.GetId()),
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
