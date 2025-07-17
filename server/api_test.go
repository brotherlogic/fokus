package server

import (
	"context"
	"testing"

	ghbclient "github.com/brotherlogic/githubridge/client"

	pb "github.com/brotherlogic/fokus/proto"
	pbgh "github.com/brotherlogic/githubridge/proto"
)

func getTestServer(ctx context.Context, t *testing.T, issues ...*pbgh.GithubIssue) *Server {
	ghbc := ghbclient.GetTestClient()

	for _, issue := range issues {
		cissue, err := ghbc.CreateIssue(ctx, &pbgh.CreateIssueRequest{
			Title: issue.GetTitle(),
			Repo:  issue.GetRepo(),
			User:  "test-user",
		})

		if err != nil {
			t.Fatalf("Unable to seed issue: %v", err)
		}

		for _, label := range issue.GetLabels() {
			_, err = ghbc.AddLabel(ctx, &pbgh.AddLabelRequest{
				Repo:  issue.GetRepo(),
				Id:    int32(cissue.GetIssueId()),
				Label: label,
			})
			if err != nil {
				t.Fatalf("Unable to add label: %v", err)
			}
		}
	}

	return &Server{
		modules: []Fokusable{
			&All{},
		},
		client: ghbc,
	}
}

func TestGetLabel_Success(t *testing.T) {
	ctx := context.Background()

	s := getTestServer(ctx, t,
		&pbgh.GithubIssue{Title: "Should Receive", Repo: "home", Labels: []string{"type-process"}},
		&pbgh.GithubIssue{Title: "Shuold not receive", Repo: "home", Labels: []string{"type-nothing"}},
		&pbgh.GithubIssue{Title: "Should not receive", Repo: "gome", Labels: []string{""}},
	)

	res, err := s.GetFokus(ctx, &pb.GetFokusRequest{Label: "type-process"})
	if err != nil {
		t.Fatalf("Failed to get fokus: %v", err)
	}

	if res.GetFocus().GetDetail() != "Should Receive [1] -> ISSUE_STATE_OPEN (fokus)" {
		t.Errorf("Expected 'Should Receive', got '%s'", res.GetFocus().GetDetail())
	}
}

func TestGetLabel_Failure(t *testing.T) {
	ctx := context.Background()

	s := getTestServer(ctx, t,
		&pbgh.GithubIssue{Title: "Should Receive", Repo: "home", Labels: []string{"type-code"}},
		&pbgh.GithubIssue{Title: "Shuold not receive", Repo: "home", Labels: []string{"type-nothing"}},
		&pbgh.GithubIssue{Title: "Should not receive", Repo: "gome", Labels: []string{""}},
	)

	r, err := s.GetFokus(ctx, &pb.GetFokusRequest{Label: "type-process"})
	if err == nil {
		t.Fatalf("Should have failed: %v -> %v", err, r)
	}

}
