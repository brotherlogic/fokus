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

type RecordAdd struct {
}

func (r *RecordAdd) getName() string {
	return "RecordAdd"
}

func (r *RecordAdd) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_RECORD_ADDER
}

func (r *RecordAdd) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time) (*pb.Focus, error) {
	issues, err := client.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(issues.Issues, func(i, j int) bool { return issues.Issues[i].GetOpenedDate() < issues.Issues[j].GetOpenedDate() })

	for _, issue := range issues.Issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() == "recordadder" {
				if time.Unix(issue.GetOpenedDate(), 0).YearDay() < now.YearDay() {
					return &pb.Focus{
						Type:   r.getType(),
						Detail: fmt.Sprintf("%v [%v]", issue.GetTitle(), issue.GetId()),
					}, nil
				}
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
