package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/evilmonkeyinc/markdownconverter"
	"github.com/evilmonkeyinc/markdownconverter/http"
	"github.com/evilmonkeyinc/markdownconverter/slack"
	flag "github.com/spf13/pflag"
)

const (
	cmdHelp    string = "help"
	cmdVersion string = "version"
)

var (
	// Arch build identifier
	Arch string = ""
	// Command is the expected name of the build
	Command string = "markdownconverter"
	// OS build identifier
	OS string = ""
	// Version build identifier
	Version string = "dev"

	errFormatUndefined   error = fmt.Errorf("format undefined")
	errFormatUnexpected  error = fmt.Errorf("unexpected format")
	errInputUndefined    error = fmt.Errorf("input undefined")
	errInputFailedRead   error = fmt.Errorf("failed to read input")
	errOutputFailedOpen  error = fmt.Errorf("failed to open output")
	errOutputFailedWrite error = fmt.Errorf("failed to write output")
	errParseFailed       error = fmt.Errorf("failed to parse")
)

func loadConverters() (map[string]markdownconverter.Converter, []string) {
	converters := make(map[string]markdownconverter.Converter)
	available := make([]string, 0)

	// TODO: use plugins to find other converters in same directory as tool?

	slackConverter := slack.New()
	available = append(available, slackConverter.Format())
	converters[slackConverter.Format()] = slackConverter

	httpConverter := http.New()
	available = append(available, httpConverter.Format())
	converters[httpConverter.Format()] = httpConverter

	return converters, available
}

func printHelp(writer *os.File, flagset *flag.FlagSet) {
	flagset.SetOutput(writer)
	fmt.Fprintf(writer, "%s is a tool for converting markdown to other formats\n\n", Command)
	fmt.Fprintf(writer, "Usage:\n\n")
	fmt.Fprintf(writer, "  %s [format] [input] [output]\n", Command)
	fmt.Fprintf(writer, "\nExample:\n\n")
	fmt.Fprintf(writer, `  %s slack "[evilmonkeyinc](https://github.com/evilmonkeyinc)"`+"\n", Command)
	fmt.Fprintf(writer, `  > <https://github.com/evilmonkeyinc|evilmonkeyinc>`+"\n")
	fmt.Fprintf(writer, "\nOptions:\n\n")
	flagset.PrintDefaults()
}

func outputError(err error) {
	fmt.Fprintf(os.Stderr, "failed: %s\n", err.Error())
	os.Exit(1)
}

func main() {
	var format, input, output string

	flagset := flag.NewFlagSet("", flag.ContinueOnError)
	flagset.Usage = func() {}

	flagset.StringVarP(&format, "format", "f", "", "The output format")
	flagset.StringVarP(&input, "input", "i", "", "The input source file")
	flagset.StringVarP(&output, "output", "o", "", "The output destination file. optional")
	if err := flagset.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			printHelp(os.Stderr, flagset)
			outputError(err)
			return
		}
		printHelp(os.Stdout, flagset)
		return
	}

	switch flagset.Arg(0) {
	case cmdHelp:
		printHelp(os.Stdout, flagset)
		return
	case cmdVersion:
		fmt.Fprintf(os.Stdout, "version %s %s/%s\n", Version, OS, Arch)
		return
	default:
		break
	}

	arg := 0
	if format == "" {
		format = flagset.Arg(arg)
		arg++
	}
	if input == "" {
		input = flagset.Arg(arg)
		arg++
	}
	if output == "" {
		output = flagset.Arg(arg)
		arg++
	}

	if format == "" {
		outputError(errFormatUndefined)
	}

	converters, available := loadConverters()

	if converter, ok := converters[format]; ok {
		inputBytes, err := handleInput(input)
		if err != nil {
			outputError(err)
		}

		outputBytes, err := converter.Parse(inputBytes)
		if err != nil {
			outputError(fmt.Errorf("%w %s", errParseFailed, err))
		}
		if err = handleOutput(output, outputBytes); err != nil {
			outputError(err)
		}
	} else {
		outputError(fmt.Errorf("%w '%s', expected: (%v)", errFormatUnexpected, format, strings.Join(available, ", ")))
	}
	os.Exit(0)
}

func handleInput(filename string) ([]byte, error) {
	if filename == "" {
		return nil, errInputUndefined
	}

	filepath := filepath.Ext(filename)
	if filepath == "" {
		return []byte(filename), nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

func handleOutput(filename string, content []byte) error {
	if filename == "" {
		if _, err := os.Stdout.Write(content); err != nil {
			return err
		}
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}
