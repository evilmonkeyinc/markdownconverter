package test

import (
	"bufio"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	helpText string = "markdownconverter is a tool for converting markdown to other formats\n\nUsage:\n\n  markdownconverter [format] [input] [output]\n\nExample:\n\n  markdownconverter slack \"[evilmonkeyinc](https://github.com/evilmonkeyinc)\"\n  > <https://github.com/evilmonkeyinc|evilmonkeyinc>\n\nOptions:\n\n  -f, --format string   The output format\n  -i, --input string    The input source file\n  -o, --output string   The output destination file. optional\n"
)

func runCommand(arg ...string) (string, error) {

	runArgs := []string{
		"run",
		"../cmd/main.go",
	}
	runArgs = append(runArgs, arg...)

	cmd := exec.Command("go", runArgs...)

	stdErr, _ := cmd.StderrPipe()
	stdOut, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(stdErr)
	errMsg := ""
	for scanner.Scan() {
		errMsg += scanner.Text() + "\n"
	}
	if errMsg != "" {
		return errMsg, nil
	}

	scanner = bufio.NewScanner(stdOut)
	outMsg := ""
	for scanner.Scan() {
		outMsg += scanner.Text() + "\n"
	}

	return outMsg, nil
}

func Test_IntegrationTests(t *testing.T) {

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "exec_version",
			args:     []string{"version"},
			expected: "version dev /\n",
		},
		{
			name:     "exec_help",
			args:     []string{"help"},
			expected: helpText,
		},
		{
			name:     "exec_help_flag",
			args:     []string{"--help"},
			expected: helpText,
		},
		{
			name:     "exec_h_flag",
			args:     []string{"-h"},
			expected: helpText,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := runCommand(test.args...)
			if err != nil {
				assert.Fail(t, "execution failed", "test '%s' failed, %s", test.name, err.Error())
				t.FailNow()
			}

			assert.Equal(t, test.expected, actual)
		})
	}

}
