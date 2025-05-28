package main

import (
	"context"
	"fmt"
	"log"
	"time"

	github_client "github.com/brotherlogic/githubridge/client"
	ghbpb "github.com/brotherlogic/githubridge/proto"
)

type triage struct {
	c github_client.GithubridgeClient
}

func buildTriage(ctx context.Context) (*triage, error) {
	c, err := github_client.GetClientExternal("madeup")

	if err != nil {
		return nil, err
	}

	// Create a new triage instance
	t := &triage{
		c: c,
	}

	return t, nil
}

func needsTriage(issue *ghbpb.GithubIssue) bool {
	return true
}

func (t *triage) getIssues(ctx context.Context) ([]*ghbpb.GithubIssue, error) {
	var issues []*ghbpb.GithubIssue

	// Fetch issues from the GitHub client
	resp, err := t.c.GetIssues(ctx, &ghbpb.GetIssuesRequest{})
	if err != nil {
		return nil, err
	}

	for _, issue := range resp.GetIssues() {
		if needsTriage(issue) {
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	t, err := buildTriage(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	issues, err := t.getIssues(ctx)
	if err != nil {
		log.Fatalf("Failed to get issues: %v", err)
	}

	fmt.Printf("Found %d issues that need triage:\n\n", len(issues))
	cancel()

	for i, issue := range issues {
		fmt.Printf("Issue %v: %v\n", i, issue.GetTitle())

		return
	}
}
