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

type RecordAdd struct {
}

func (r *RecordAdd) getName() string {
	return "RecordAdd"
}

func (r *RecordAdd) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_RECORD_ADDER
}

func (r *RecordAdd) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {
	if now.Weekday() != time.Saturday && now.Weekday() != time.Sunday {
		if now.Hour() < 17 {
			return nil, status.Errorf(codes.InvalidArgument, "Unable to find a suitable issue")
		}
	}

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() == "recordadder" {
				return &pb.Focus{
					Type:   r.getType(),
					Detail: fmt.Sprintf("%v [%v]", issue.GetTitle(), issue.GetId()),
				}, nil
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
