package server

import (
	"fmt"
	"sort"
	"time"

	pb "github.com/brotherlogic/fokus/proto"
	githubridgeclient "github.com/brotherlogic/githubridge/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	ghbpb "github.com/brotherlogic/githubridge/proto"
)

type Home struct {
	client githubridgeclient.GithubridgeClient
}

func (h *Home) getName() string {
	return "Home"
}

func (h *Home) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_HOME_TASKS
}

func (h *Home) getFokus(ctx context.Context) (*pb.Focus, error) {
	if time.Now().Weekday() != time.Saturday && time.Now().Weekday() != time.Sunday {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	if time.Now().Hour() < 13 && time.Now().Hour() >= 16 {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	issues, err := h.client.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(issues.Issues, func(i, j int) bool { return issues.Issues[i].GetOpenedDate() < issues.Issues[j].GetOpenedDate() })

	for _, issue := range issues.Issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() == "home" {
				if time.Unix(issue.GetOpenedDate(), 0).YearDay() < time.Now().In(location).YearDay() {
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
