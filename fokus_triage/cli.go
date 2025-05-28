package main

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	github_client "github.com/brotherlogic/githubridge/client"
	ghbpb "github.com/brotherlogic/githubridge/proto"
)

type triage struct {
	c github_client.GithubridgeClient
}

var (
	issueTypes     = []string{"type-process", "type-code", "type-research"}
	priorities     = []string{"priority-p0", "priority-p1", "priority-p2"}
	processTShirt  = []string{"Small", "Medium", "Large", "X-Large"}
	researchTShirt = []string{"Large", "X-Large"}
	codeBreakdown  = []string{"Sized", "Need-Breakdown"}
)

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
	hasType := false
	for _, label := range issue.GetLabels() {
		if slices.Contains(issueTypes, label) {
			fmt.Printf("Issue %s already has type %s\n", issue.GetTitle(), label)
			hasType = true
			break
		}
	}

	return !(hasType)
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
		} else {
			fmt.Printf("%v does not need triage\n", issue.GetTitle())
		}
	}

	return issues, nil
}

func (t *triage) setLabel(issue *ghbpb.GithubIssue, label string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("Setting %v to %v\n", issue.GetTitle(), label)

	// Set the label on the issue
	_, err := t.c.AddLabel(ctx, &ghbpb.AddLabelRequest{
		User:  issue.GetUser(),
		Repo:  issue.GetRepo(),
		Id:    int32(issue.GetId()),
		Label: label,
	})
	if err != nil {
		log.Printf("Failed to set label %s on issue %s: %v", label, issue.GetTitle(), err)
	}
}

func (t *triage) runLabel(issue *ghbpb.GithubIssue, options []string, signifier string, rlabel string) {

	// Check if we have a required label
	if rlabel != "" {
		found := false
		for _, label := range issue.GetLabels() {
			if label == rlabel {
				found = true
				break
			}
		}

		if !found {
			return
		}
	}

	for _, label := range issue.GetLabels() {
		for _, op := range options {
			if label == op {
				return
			}
		}
	}

	fmt.Printf("Choose %v:\n\n", signifier)
	for i, typ := range options {
		fmt.Printf("%d: %s\n", i, typ)
	}
	var choice int
	fmt.Scanf("%d", &choice)
	if choice < 0 || choice >= len(options) {
		fmt.Println("Invalid choice, exiting.")
		return
	}

	t.setLabel(issue, options[choice])
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
		fmt.Printf("Issue %v: %v (%v)\n", i, issue.GetTitle(), issue.GetLabels())

		t.runLabel(issue, issueTypes, "type", "")
		t.runLabel(issue, priorities, "priority", "")
		t.runLabel(issue, processTShirt, "size", "type-process")
		t.runLabel(issue, researchTShirt, "size", "type-research")
		t.runLabel(issue, codeBreakdown, "sized", "type-code")

		return
	}
}
