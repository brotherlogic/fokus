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

type Highlight struct {
	client githubridgeclient.GithubridgeClient
}

func (h *Highlight) getName() string {
	return "highlight"
}

func (h *Highlight) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_HIGHLIGHT
}

func (h *Highlight) getFokus(ctx context.Context) (*pb.Focus, error) {
	// We can't rely on America/Los_Angeles being present it seems; ignore Daylight savbings
	location := time.FixedZone("UTC-8", -8*60*60)

	if time.Now().In(location).Weekday() == time.Saturday || time.Now().In(location).Weekday() == time.Sunday {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for highlight tasks")
	}

	if (time.Now().In(location).Hour() < 6 || time.Now().In(location).Hour() >= 7) &&
		(time.Now().In(location).Hour() < 15 || time.Now().In(location).Hour() >= 17) {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for highlight tasks")
	}

	issues, err := h.client.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(issues.Issues, func(i, j int) bool { return issues.Issues[i].GetOpenedDate() < issues.Issues[j].GetOpenedDate() })

	for _, issue := range issues.Issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() == "gramophile" || issue.GetRepo() == "blog" {
				return &pb.Focus{
					Type:   h.getType(),
					Detail: fmt.Sprintf("%v [%v] -> %v (highlight)", issue.GetTitle(), issue.GetId(), issue.GetState()),
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
