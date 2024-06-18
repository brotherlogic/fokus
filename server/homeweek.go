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

type HomeWeek struct {
}

func (h *HomeWeek) getName() string {
	return "Homeweek"
}

func (h *HomeWeek) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_HOME_TASKS
}

func (h *HomeWeek) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time) (*pb.Focus, error) {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	if now.Hour() < 20 || now.Hour() >= 21 {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for home tasks")
	}

	issues, err := client.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(issues.Issues, func(i, j int) bool { return issues.Issues[i].GetOpenedDate() < issues.Issues[j].GetOpenedDate() })

	for _, issue := range issues.Issues {
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
