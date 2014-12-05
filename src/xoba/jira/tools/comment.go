package tools

import (
	"fmt"
	"log"

	"github.com/xoba/goutil/tool"
)

func Comment(args []string) {
	var key, comment string
	flags := tool.FlagsWithDoc("comment", doc)
	flags.StringVar(&key, "key", "", "key of issue to comment on")
	flags.StringVar(&comment, "comment", "", "the comment")
	parseFlags(flags, args)
	if err := validateAllStringArgs("key", key, "comment", comment); err != nil {
		log.Fatal(err)
	}
	resp, err := commentOnIssue(key, comment)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

type CommentResponse struct {
	Self string
}

func (c CommentResponse) String() string {
	return ToString(c)
}

func commentOnIssue(key, comment string) (*CommentResponse, error) {
	content := map[string]interface{}{
		"body": comment,
	}
	var out CommentResponse
	parser := JsonParser(&out)
	_, err := apiCall(content, "POST", fmt.Sprintf("issue/%s/comment", key), ExpectedCodeValidator(201), parser)
	if err != nil {
		return nil, err
	}
	return &out, err
}
