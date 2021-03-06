package tools

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const doc = `
you could put exports like these into your ~/.bashrc file for convenience:

  export JIRA_USERNAME=joe.smith
  export JIRA_PASSWORD=abc123
  export JIRA_URL=http://www.example.com/jira
`

func JsonParser(i interface{}) ResponseParser {
	return func(resp *http.Response) (interface{}, error) {
		d := json.NewDecoder(resp.Body)
		if err := d.Decode(i); err != nil {
			return nil, err
		}
		return i, nil
	}
}

func NilParser(i interface{}) ResponseParser {
	return func(resp *http.Response) (interface{}, error) {
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			return nil, err
		}
		switch x := i.(type) {
		case *string:
			*x = "ok"
		}
		return i, nil
	}
}

func ExpectedCodeValidator(code int) StatusCodeValidator {
	return func(c int) bool {
		return c == code
	}
}

type StatusCodeValidator func(int) bool
type ResponseParser func(*http.Response) (interface{}, error)

func apiCall(content interface{}, method, path string, val StatusCodeValidator, parser ResponseParser) (interface{}, error) {
	var r io.Reader
	if content != nil {
		switch c := content.(type) {
		case []byte: // it's already json
			r = bytes.NewReader(c)
		case string: // it's already json
			r = strings.NewReader(c)
		default:
			buf, err := json.Marshal(c)
			if err != nil {
				return nil, err
			}
			r = bytes.NewReader(buf)
		}
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/rest/api/2/%s", _url, path), r)
	if err != nil {
		return nil, err
	}
	auth(req)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !val(resp.StatusCode) {
		return nil, fmt.Errorf("bad status: %q\n", resp.Status)
	}
	return parser(resp)
}

func client() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			auth(req)
			return nil
		},
	}
}

func validateAllStringArgs(args ...string) error {
	return errorJoiner(func() (errs []error) {
		if err := validateStringArgs("username", _username, "password", _password, "url", _url); err != nil {
			errs = append(errs, err)
		}
		if err := validateStringArgs(args...); err != nil {
			errs = append(errs, err)
		}
		return
	})
}

func errorJoiner(f func() []error) error {
	if errs := f(); len(errs) == 0 {
		return nil
	} else {
		var out []string
		for _, e := range errs {
			out = append(out, e.Error())
		}
		return fmt.Errorf(strings.Join(out, ", "))
	}
}

func validateStringArgs(args ...string) error {
	return errorJoiner(func() (errs []error) {
		for i := 0; i < len(args); i += 2 {
			name := args[i]
			value := args[i+1]
			if len(value) == 0 {
				errs = append(errs, fmt.Errorf("-%s needs setting", name))
			}
		}
		return
	})
}

var _username, _password, _url string

func auth(req *http.Request) {
	req.SetBasicAuth(_username, _password)
}

func init() {
	if p := os.Getenv("JIRA_PASSWORD"); len(p) > 0 {
		_password = p
	}
	if p := os.Getenv("JIRA_USERNAME"); len(p) > 0 {
		_username = p
	}
	if p := os.Getenv("JIRA_URL"); len(p) > 0 {
		_url = p
	}
}

func parseFlags(flags *flag.FlagSet, args []string) {
	var username, password, url string
	flags.StringVar(&username, "username", _username, "authentication username")
	flags.StringVar(&password, "password", _password, "authentication password")
	flags.StringVar(&url, "url", _url, "base jira url")
	flags.Parse(args)
	if len(username) > 0 {
		_username = username
	}
	if len(password) > 0 {
		_password = password
	}
	if len(url) > 0 {
		_url = url
	}
	for strings.HasSuffix(_url, "/") {
		_url = _url[:len(_url)-1]
	}
}

func ToString(i interface{}) string {
	buf, err := json.Marshal(i)
	check(err)
	return string(buf)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
