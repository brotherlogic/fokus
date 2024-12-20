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

type Highlight struct {
}

func (h *Highlight) getName() string {
	return "highlight"
}

func (h *Highlight) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_HIGHLIGHT
}

func (h *Highlight) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {

	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for highlight tasks")
	}

	log.Printf("Evaluating time: %v", now.Hour())
	if (now.Hour() < 6 || now.Hour() >= 7) &&
		(now.Hour() < 15 || now.Hour() >= 17) {
		return nil, status.Errorf(codes.FailedPrecondition, "Not ready for highlight tasks")
	}

	for _, issue := range issues {
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
