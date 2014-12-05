package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xoba/goutil/tool"
)

func InstallGitHooks(args []string) {
	var repo, pattern string
	flags := tool.FlagsWithDoc("install", doc)
	flags.StringVar(&repo, "repo", ".", "path of git repository to install hooks to")
	flags.StringVar(&pattern, "pattern", "", "url pattern with %s for commit hash")
	parseFlags(flags, args)
	if err := validateStringArgs("pattern", pattern, "repo", repo); err != nil {
		log.Fatal(err)
	}
	for _, h := range strings.Split("commit-msg,post-commit", ",") {
		v := filepath.Join(repo, fmt.Sprintf(".git/hooks/%s", h))
		out, err := os.Create(v)
		check(err)
		fmt.Fprintln(out, "#!/bin/bash")
		fmt.Fprintf(out, "jira hook.%s %s $@\n", h, pattern)
		check(out.Close())
		os.Chmod(v, os.ModePerm)
	}
}

func CommitHook(args []string) {
}

func PostCommitHook(args []string) {
	var lines []string
	f := strings.NewReader(commitMessage())
	var key string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
		if len(key) == 0 {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				key = parts[0]
			}
		}
	}
	check(s.Err())
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("see "+args[0], commitHash()))
	if len(key) == 0 {
		log.Println("oops, couldn't find a jira issue!")
		return
	}
	log.Printf("going to comment on jira issue %q\n", key)
	resp, err := commentOnIssue(key, strings.Join(lines, "\n"))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp)
}

func commitHash() string {
	return commitFormat("%H")
}

func commitMessage() string {
	return fmt.Sprintf("%s\n%s\n", commitFormat("%s"), commitFormat("%b"))
}

func commitFormat(f string) string {
	buf := new(bytes.Buffer)
	cmd := exec.Command("git", "log", "-n", "1", "HEAD", "--format=format:"+f)
	cmd.Stdout = buf
	check(cmd.Run())
	return buf.String()
}
