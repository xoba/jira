package main

import (
	"xoba/jira/tools"

	"github.com/xoba/goutil"
	"github.com/xoba/goutil/tool"
)

func main() {
	tool.Run()
}

func init() {
	goutil.PlatformInit()
	add := func(name, desc string, f func([]string)) {
		tool.Register(tool.Named(name+","+desc, tool.RunFunc(f)))
	}
	add("comment", "comment on a jira issue", tools.Comment)
	add("create", "create a jira issue", tools.CreateIssue)
	add("delete", "delete a jira issue", tools.DeleteIssue)
	add("list", "list jira issues", tools.ListIssues)
	add("hook.install", "install git hooks", tools.InstallGitHooks)
	add("hook.commit-msg", "git commit hook", tools.CommitHook)
	add("hook.post-commit", "git postcommit hook", tools.PostCommitHook)
}
