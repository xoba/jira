package tools

import (
	"fmt"
	"log"

	"github.com/xoba/goutil/tool"
)

func CreateIssue(args []string) {
	var project, summary, description, task string
	flags := tool.FlagsWithDoc("create", doc)
	flags.StringVar(&task, "task", "Task", "the type of task")
	flags.StringVar(&project, "project", "", "the project issue is under")
	flags.StringVar(&summary, "summary", "", "summary of issue")
	flags.StringVar(&description, "description", "", "description of issue")
	parseFlags(flags, args)
	if err := validateAllStringArgs("project", project, "summary", summary, "description", description); err != nil {
		log.Fatal(err)
	}
	issue, err := createIssue(project, summary, description, task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(issue)
}

type NewIssue struct {
	Id   string
	Key  string
	Self string
}

func (n NewIssue) String() string {
	return ToString(n)
}

func createIssue(project, summary, description, task string) (*NewIssue, error) {
	content := fmt.Sprintf(`{
    "fields": {
       "project":
       {
          "key": "%s"
       },
       "summary": "%s",
       "description": "%s",
       "issuetype": {
          "name": "%s"
       }
   }
}`, project, summary, description, task)
	var out NewIssue
	parser := JsonParser(&out)
	_, err := apiCall(content, "POST", "issue", ExpectedCodeValidator(201), parser)
	if err != nil {
		return nil, err
	}
	return &out, err
}
