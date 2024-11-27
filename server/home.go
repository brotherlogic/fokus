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

type Home struct {
}

func (h *Home) getName() string {
	return "Home"
}

func (h *Home) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_HOME_TASKS
}

func (h *Home) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {
	if now.Weekday() != time.Saturday && now.Weekday() != time.Sunday {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	if now.Hour() < 21 && now.Hour() >= 22 {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() == "home" {
				if time.Unix(issue.GetOpenedDate(), 0).YearDay() < now.YearDay() {
					return &pb.Focus{
						Type:   h.getType(),
						Detail: fmt.Sprintf("%v [%v] -> %v (weekend)", issue.GetTitle(), issue.GetId(), issue.GetState()),
					}, nil
				}
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
