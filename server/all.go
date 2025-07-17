package server

import (
	"fmt"
	"log"
	"time"

	pb "github.com/brotherlogic/fokus/proto"
	githubridgeclient "github.com/brotherlogic/githubridge/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	ghbpb "github.com/brotherlogic/githubridge/proto"
)

type All struct {
}

func (h *All) getName() string {
	return "all"
}

func (h *All) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_ALL
}

func (h *All) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {

	log.Printf("Getting all from %v issues", len(issues))

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			return &pb.Focus{
				Type:   h.getType(),
				Detail: fmt.Sprintf("%v [%v] -> %v (fokus)", issue.GetTitle(), issue.GetId(), issue.GetState()),
			}, nil
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
