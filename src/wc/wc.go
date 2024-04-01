package main

import (
	"args"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func exitWithErrorAndCode(err error, code int) {
	log.Println(err)
	os.Exit(code)
}

func getDefaultArguments() []args.Arg {
	return []args.Arg{
		{
			Name:         "w",
			Description:  "for counting words",
			DefaultValue: true,
		},
		{
			Name:         "l",
			Description:  "for counting lines",
			DefaultValue: false,
		},
		{
			Name:         "m",
			Description:  "for counting characters",
			DefaultValue: false,
		},
	}
}

func checkArguments(args map[string]any) error {
	lPresence, wPresence, mPresence := false, false, false
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "l":
			lPresence = true
		case "w":
			wPresence = true
		case "m":
			mPresence = true
		}
	})

	argsCount := func(bools ...bool) int {
		result := 0
		for _, v := range bools {
			if v {
				result++
			}
		}
		return result
	}(lPresence, wPresence, mPresence)

	if argsCount > 1 {
		return fmt.Errorf("only one parameter can be specified at once")
	}

	if argsCount == 1 {
		args["l"] = lPresence
		args["w"] = wPresence
		args["m"] = mPresence
	}

	return nil
}

type Processor func(reader *bufio.Reader) (int, error)

func scanFile(reader *bufio.Reader, splitter bufio.SplitFunc) (int, error) {
	result := 0

	scanner := bufio.NewScanner(reader)
	scanner.Split(splitter)
	for scanner.Scan() {
		result++
	}

	return result, scanner.Err()
}

func processSymbols(reader *bufio.Reader) (int, error) {
	return scanFile(reader, bufio.ScanRunes)
}

func processLines(reader *bufio.Reader) (int, error) {
	return scanFile(reader, bufio.ScanLines)
}

func processWords(reader *bufio.Reader) (int, error) {
	return scanFile(reader, bufio.ScanWords)
}

func getWorkFunc(args map[string]any) Processor {
	if args["w"].(bool) {
		return processWords
	}
	if args["l"].(bool) {
		return processLines
	}
	return processSymbols
}

func processFile(path string, processor Processor) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	fs, err := file.Stat()
	if err != nil {
		return 0, err
	}

	if fs.IsDir() {
		return 0, fmt.Errorf("%s is directory", path)
	}

	defer file.Close()

	return processor(bufio.NewReader(file))
}

func main() {
	log.SetFlags(log.Lshortfile)

	params, rest, err := args.ParseArgs(getDefaultArguments()...)
	if err != nil {
		exitWithErrorAndCode(err, 1)
	}

	if len(rest) == 0 {
		exitWithErrorAndCode(fmt.Errorf("nothing to count"), 2)
	}

	err = checkArguments(params)
	if err != nil {
		exitWithErrorAndCode(err, 3)
	}

	type ProcessorResult struct {
		Value int
		Name  string
		Err   error
	}

	processor := getWorkFunc(params)
	resultChannel := make(chan ProcessorResult)

	for _, path := range rest {
		go func() {
			count, err := processFile(path, processor)
			resultChannel <- ProcessorResult{Value: count, Name: path, Err: err}
		}()
	}

	for i := 0; i < len(rest); i++ {
		if v := <-resultChannel; v.Err != nil {
			log.Println(v.Err)
		} else {
			fmt.Printf("%d\t%s\n", v.Value, v.Name)
		}
	}
}
