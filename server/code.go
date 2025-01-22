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

type Code struct {
}

func (h *Code) getName() string {
	return "code"
}

func (h *Code) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_CODING_TASKS
}

func (h *Code) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {

	//Only work on code during commute time
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return nil, status.Errorf(codes.FailedPrecondition, "Not the time for coding")
	}

	if now.Hour() < 5 || now.Hour() > 15 {
		return nil, status.Errorf(codes.FailedPrecondition, "Not the time for coding")
	}

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			foundLabel := false
			for _, label := range issue.GetLabels() {
				if label == "code" {
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
