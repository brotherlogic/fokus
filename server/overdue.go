package server

import (
	"fmt"
	"log"
	"sort"
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
	client githubridgeclient.GithubridgeClient
}

func (o *Overdue) getName() string {
	return "Overdue"
}

func (o *Overdue) getType() pb.Focus_FocusType {
	return pb.Focus_FOCUS_ON_CODING_TASKS
}

func (o *Overdue) getFokus(ctx context.Context) (*pb.Focus, error) {
	issues, err := o.client.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	sort.SliceStable(issues.Issues, func(i, j int) bool { return issues.Issues[i].GetOpenedDate() < issues.Issues[j].GetOpenedDate() })

	for _, issue := range issues.Issues {
		if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
			if issue.GetRepo() != "bandcampserver" && issue.GetRepo() != "recordalerting" && issue.GetRepo() != "home" {
				if !strings.Contains(issue.GetTitle(), "Incomplete Order") {
					// We can't rely on America/Los_Angeles being present it seems; ignore Daylight savbings
					location := time.FixedZone("UTC-8", -8*60*60)
					if err != nil {
						return nil, err
					}
					if time.Unix(issue.GetOpenedDate(), 0).YearDay() < time.Now().In(location).YearDay() {
						return &pb.Focus{
							Type:   o.getType(),
							Detail: fmt.Sprintf("%v [%v] -> %v (%v vs %v)", issue.GetTitle(), issue.GetId(), issue.GetState(), time.Unix(issue.GetOpenedDate(), 0).YearDay(), time.Now().In(location)),
						}, nil
					} else {
						log.Printf("Skipping %v because %v < %v", issue.GetTitle(), time.Unix(issue.GetOpenedDate(), 0), time.Now())
					}
				}
			}
		}
	}
	return nil, status.Errorf(codes.InvalidArgument, "Unable to locate an open issue")
}
