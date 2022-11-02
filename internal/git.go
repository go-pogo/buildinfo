package internal

import (
	"bytes"
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-pogo/errors"
)

const DefaultTag = "0.0.0"

func LatestTag(ctx context.Context) (string, string, error) {
	tags, err := execGit(ctx, "tag", "--sort=-v:refname")
	if err != nil {
		return "", "", err
	}

	tags = bytes.TrimLeftFunc(tags, unicode.IsSpace)
	if i := bytes.IndexRune(tags, '\n'); i > 0 {
		tags = tags[:i]
	}

	return cutTag(tags)
}

func CurrentTag(ctx context.Context) (string, string, error) {
	tag, err := execGit(ctx, "describe", "--tags", "--abbrev=0")
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) &&
			(bytes.Contains(exitErr.Stderr, []byte("No names found")) ||
				bytes.Contains(exitErr.Stderr, []byte("cannot describe"))) {
			// repository does not contain any tags
			return DefaultTag, "", nil
		}
		return "", "", err
	}

	return cutTag(tag)
}

func TagDetails(ctx context.Context, tag string) (string, time.Time, error) {
	details, err := execGit(ctx, "log", "-1", "--pretty=%h,%ct", tag)
	if err != nil {
		return "", time.Time{}, err
	}

	details = bytes.TrimSpace(details)
	s := strings.SplitN(string(details), ",", 2)

	ts, err := strconv.ParseInt(s[1], 10, 64)
	return s[0], time.Unix(ts, 0), err
}

func execGit(ctx context.Context, args ...string) ([]byte, error) {
	// log.Println("git", strings.Join(args, " "))
	out, err := exec.CommandContext(ctx, "git", args...).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			err = &ExitError{cause: exitErr}
		}
		return nil, errors.WithStack(err)
	}
	return out, nil
}

func cutTag(b []byte) (string, string, error) {
	b = bytes.TrimSpace(b)
	tag, remains, _ := strings.Cut(string(b), "-")
	return tag, remains, nil

}

type ExitError struct {
	cause *exec.ExitError
}

func (e *ExitError) Unwrap() error { return e.cause }
func (e *ExitError) Error() string {
	if len(e.cause.Stderr) != 0 {
		return string(e.cause.Stderr)
	}
	return e.cause.Error()
}
