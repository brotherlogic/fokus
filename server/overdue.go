package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	pb "github.com/brotherlogic/fokus/proto"
	githubridgeclient "github.com/brotherlogic/githubridge/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	ghbpb "github.com/brotherlogic/githubridge/proto"
)

type Overdue struct {
}

func (o *Overdue) getName() string {
	return "Overdue"
}

func (o *Overdue) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_CODING_TASKS
}

func (o *Overdue) getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error) {

	for _, issue := range issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() != "bandcampserver" && issue.GetRepo() != "recordalerting" && issue.GetRepo() != "home" && issue.GetRepo() != "research" {
				if !strings.Contains(issue.GetTitle(), "Incomplete Order") {
					if !strings.HasPrefix(issue.GetTitle(), "CD Rip Need") {
						if time.Unix(issue.GetOpenedDate(), 0).YearDay() < now.YearDay() {
							return &pb.Focus{
								Type:   o.getType(),
								Detail: fmt.Sprintf("%v [%v] -> %v (%v vs %v)", issue.GetTitle(), issue.GetId(), issue.GetState(), time.Unix(issue.GetOpenedDate(), 0).YearDay(), now),
							}, nil
						} else {
							log.Printf("Skipping %v because %v < %v", issue.GetTitle(), time.Unix(issue.GetOpenedDate(), 0), now)
						}
					}
				}
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
