package tools

import (
	"fmt"
	"log"

	"github.com/xoba/goutil/tool"
)

func DeleteIssue(args []string) {
	var key string
	flags := tool.FlagsWithDoc("delete", doc)
	flags.StringVar(&key, "key", "", "key of issue to delete")
	parseFlags(flags, args)
	if err := validateAllStringArgs("key", key); err != nil {
		log.Fatal(err)
	}
	resp, err := deleteIssue(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

func deleteIssue(key string) (interface{}, error) {
	var out string
	_, err := apiCall(nil, "DELETE", fmt.Sprintf("issue/%s", key), ExpectedCodeValidator(204), NilParser(&out))
	if err != nil {
		return nil, err
	}
	return out, nil
}
