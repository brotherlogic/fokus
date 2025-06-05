package server

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ghbclient "github.com/brotherlogic/githubridge/client"
	githubridgeclient "github.com/brotherlogic/githubridge/client"
	ghbpb "github.com/brotherlogic/githubridge/proto"

	pb "github.com/brotherlogic/fokus/proto"
)

type Fokusable interface {
	getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time, issues []*ghbpb.GithubIssue) (*pb.Focus, error)
	getName() string
	getType() pb.Focus_FocusType
}

type Server struct {
	modules []Fokusable
	client  githubridgeclient.GithubridgeClient
}

func NewServer() *Server {
	client, err := ghbclient.GetClientInternal()
	if err != nil {
		log.Fatalf("Unable to reach GHB")
	}
	return &Server{
		modules: []Fokusable{
			&Cluster{},
			&RecordAdd{},
			&Code{},
			&Tasks{},
			&Highlight{},
			&HomeWeek{},
			&Home{},
			&Overdue{},
			&Fokus{}},
		client: client,
	}
}

func (s *Server) trimToActionable(ctx context.Context, issues []*ghbpb.GithubIssue) ([]*ghbpb.GithubIssue, error) {
	var validIssues []*ghbpb.GithubIssue

	for _, issue := range issues {
		// See if we've manually tagged the issue as blocked
		labels, err := s.client.GetLabels(ctx, &ghbpb.GetLabelsRequest{
			User: "brotherlogic",
			Repo: issue.GetRepo(),
			Id:   int32(issue.GetId()),
		})
		if err != nil {
			return nil, err
		}
		blocked := false
		for _, label := range labels.GetLabels() {
			if strings.ToLower(label) == "blocked" {
				blocked = true
			}
		}

		// Skip this issue if it's blocked
		if blocked {
			continue
		}

		comments, err := s.client.GetComments(ctx, &ghbpb.GetCommentsRequest{
			Repo: issue.GetRepo(),
			User: issue.GetUser(),
			Id:   int32(issue.GetId()),
		})
		if err != nil {
			return nil, err
		}

		sort.SliceStable(comments.Comments, func(i, j int) bool {
			return comments.Comments[i].GetTimestamp() > comments.Comments[j].GetTimestamp()
		})

		if len(comments.Comments) > 0 && (strings.HasPrefix(comments.Comments[0].GetText(), "Block on") || strings.HasPrefix(comments.Comments[0].GetText(), "Blocked on")) {
			elems := strings.Split(strings.Fields(comments.Comments[0].GetText())[2], "/")
			number, err := strconv.ParseInt(elems[2], 10, 32)
			if err != nil {
				return nil, err
			}

			open := false
			for _, issue := range issues {
				if issue.GetRepo() == elems[1] &&
					issue.GetUser() == elems[0] &&
					issue.GetId() == number && issue.GetState() == ghbpb.IssueState_ISSUE_STATE_OPEN {
					open = true
				}
			}

			if !open {
				validIssues = append(validIssues, issue)

				// Also post that the issue is unblocked with fire and forget
				s.client.CommentOnIssue(ctx, &ghbpb.CommentOnIssueRequest{
					User:    issue.GetUser(),
					Repo:    issue.GetRepo(),
					Id:      int32(issue.GetId()),
					Comment: "Unblocked",
				})
			}
		} else {
			validIssues = append(validIssues, issue)
		}
	}

	return validIssues, nil
}

func (s *Server) GetFokus(ctx context.Context, req *pb.GetFokusRequest) (*pb.GetFokusResponse, error) {
	log.Printf("Getting Fokus")

	// We can't rely on America/Los_Angeles being present it seems; ignore Daylight savings
	location := time.FixedZone("UTC-8", -8*60*60)
	t := time.Now().In(location)

	rissues, err := s.client.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	issues, err := s.trimToActionable(ctx, rissues.Issues)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(issues, func(i, j int) bool { return issues[i].GetOpenedDate() < issues[j].GetOpenedDate() })

	for _, m := range s.modules {
		focus, err := m.getFokus(ctx, s.client, t, issues)
		log.Printf("%v -> %v", m.getName(), err)
		if err == nil && focus != nil {
			return &pb.GetFokusResponse{
				Focus:     focus,
				GivenTime: t.Unix(),
			}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "Could not find focus task")
}
