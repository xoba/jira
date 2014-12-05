package tools

import (
	"fmt"
	"log"

	"github.com/xoba/goutil/tool"
)

func ListIssues(args []string) {
	var reporter, assignee, key, value string
	flags := tool.FlagsWithDoc("list", doc)
	flags.StringVar(&reporter, "reporter", "", "the reporter to list issues for")
	flags.StringVar(&assignee, "assignee", "", "the assignee to list issues for")
	parseFlags(flags, args)
	switch {
	case len(reporter) > 0:
		key = "reporter"
		value = reporter
	case len(assignee) > 0:
		key = "assignee"
		value = assignee
	}
	if err := validateAllStringArgs("assignee/reporter", value); err != nil {
		log.Fatal(err)
	}
	jql := fmt.Sprintf("%s=%s", key, value)
	fmt.Printf("results for %q:\n", jql)
	resp, err := searchIssues(jql, 0)
	check(err)
	for _, r := range resp.Issues {
		fmt.Printf("%s (in %s): %q\n", r.Key, r.Fields.Project.Key, r.Fields.Summary)
	}
}

func searchIssues(jql string, start int) (*IssueSearchResponse, error) {
	var out IssueSearchResponse
	parser := JsonParser(&out)
	_, err := apiCall(nil, "GET", fmt.Sprintf("search?jql=%s&startAt=%d", jql, start), ExpectedCodeValidator(200), parser)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

type ResponseMeta struct {
	MaxResults int
	StartAt    int
	Total      int
}

type IssueSearchResponse struct {
	ResponseMeta
	Issues []Issue
}

func (i IssueSearchResponse) String() string {
	return ToString(i)
}

type Issue struct {
	Id     string
	Key    string
	Self   string
	Fields Fields
}

type Fields struct {
	Created   string
	IssueType IssueType
	Reporter  Reporter
	Project   Project
	Status    Status
	Summary   string
	Updated   string
}

type IssueType struct {
	Description string
	IconUrl     string
	Id          string
	Name        string
	Self        string
	Subtask     bool
}

type Reporter struct {
	Active       bool
	AvatarUrls   map[string]string
	DisplayName  string
	EmailAddress string
	Name         string
	Self         string
}

type Status struct {
	Description string
	IconUrl     string
	Id          string
	Name        string
	Self        string
}

type Project struct {
	AvatarUrls map[string]string
	Id         string
	Key        string
	Name       string
	Self       string
}
