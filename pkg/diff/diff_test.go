package diff

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"testing"
)

//go:embed testdata/tests.sql
var testsSQL []byte

func TestDiff(t *testing.T) {
	tests := parseTestData(t, testsSQL)

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			diff, err := NewDiffer(tc.Orig, tc.New).Run()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			str := strings.TrimSpace(diff.String())

			if str != tc.Output {
				fmt.Printf("Wanted: %q\n", tc.Output)
				fmt.Printf("Got:    %q\n", str)
				t.Fatalf("expected diff to match")
			}
		})
	}
}

func parseTestData(t *testing.T, contents []byte) []*testcase {
	t.Helper()

	scanner := bufio.NewScanner(bytes.NewBuffer(contents))

	var tests []*testcase
	var tc *testcase

	for scanner.Scan() {
		line := scanner.Text()

		// new test
		if strings.Contains(line, "---- Test:") {
			tc = &testcase{
				Name:  strings.TrimSpace(strings.TrimPrefix(line, "---- Test: ")),
				phase: 1,
			}
			continue
		}

		if tc == nil {
			continue
		}

		if strings.Contains(line, "----") {
			if tc.phase == 3 {
				tests = append(tests, tc)
			}
			tc.phase++
			continue
		}

		if tc.phase == 1 {
			tc.Orig = join(tc.Orig, line)
		}
		if tc.phase == 2 {
			tc.New = join(tc.New, line)
		}
		if tc.phase == 3 {
			tc.Output = join(tc.Output, line)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	return tests
}

func join(a, b string) string {
	if a == "" {
		return b
	}
	return fmt.Sprintf("%s\n%s", a, b)
}

type testcase struct {
	phase  int
	Name   string
	Orig   string
	New    string
	Output string
}
