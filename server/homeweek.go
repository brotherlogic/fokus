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

type HomeWeek struct {
}

func (h *HomeWeek) getName() string {
	return "Homeweek"
}

func (h *HomeWeek) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_HOME_TASKS
}

func (h *HomeWeek) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	if now.Hour() != 14 && now.Hour() != 15 && now.Hour() != 20 {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() == "home" {
				if time.Unix(issue.GetOpenedDate(), 0).YearDay() < now.YearDay() {
					return &pb.Focus{
						Type:   h.getType(),
						Detail: fmt.Sprintf("%v [%v] -> %v (home)", issue.GetTitle(), issue.GetId(), issue.GetState()),
					}, nil
				}
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
