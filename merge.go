// Copyright (c) 2017 Intel Corporation
// Copyright (c) 2022 Oscar Picas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type mergeContext struct {
	// Where to write the combined result file.
	out io.Writer
	// Ensure we keep the same Mode throughout the input files.
	mode string
}

func parseMode(line string) (string, error) {
	const p = "mode: "

	if !strings.HasPrefix(line, p) || line == p {
		return "", fmt.Errorf("merge: bad mode line: %v", line)
	}
	mode := line[len(p):]

	if mode != "set" && mode != "count" && mode != "atomic" {
		return "", fmt.Errorf("merge: bad mode line: %v", line)
	}

	return mode, nil
}

func compatibleModes(a, b string) bool {
	if a == "set" && b != "set" {
		return false
	}
	if a != "set" && b == "set" {
		return false
	}
	return true
}

func (ctx *mergeContext) addFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("problem closing file: %s", f.Name())
		}
	}(f)

	// Parse the coverage mode.
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return err
	}
	line := scanner.Text()
	mode, err := parseMode(line)
	if err != nil {
		return err
	}
	if ctx.mode == "" {
		ctx.mode = mode
		_, err := ctx.out.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("merge: problem writing buffer: %w", err)
		}
		_, err = ctx.out.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("merge: problem writing buffer: %w", err)
		}
	} else if !compatibleModes(ctx.mode, mode) {
		return fmt.Errorf("merge: inconsistency in the profiles modes, got both %v and %v", ctx.mode, mode)
	}

	// Copy the list of block information after skipping the line we just parsed.
	_, err = f.Seek(int64(len([]byte(line))+1), 0)
	if err != nil {
		return fmt.Errorf("merge: internal problem: %w", err)
	}
	_, err = io.Copy(ctx.out, f)
	if err != nil {
		return fmt.Errorf("merge: internal copy problem: %w", err)
	}

	return nil
}

func merge(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() < 2 {
		return fmt.Errorf("expecting at least two arguments, got %d", args.Len())
	}
	aggregatedProfile, err := ioutil.TempFile("", "cover-tool-")
	if err != nil {
		return err
	}
	defer func() {
		_ = aggregatedProfile.Close()
		err := os.Remove(aggregatedProfile.Name())
		if err != nil {
			log.Println("problem cleaning temporary files")
		}
	}()

	parse := &mergeContext{
		out: aggregatedProfile,
	}

	for _, arg := range args.Slice() {
		if err := parse.addFile(arg); err != nil {
			return err
		}
	}

	// Loading the aggregated profile will merge counters for the same blocks.
	profiles, err := ParseProfiles(aggregatedProfile.Name())
	if err != nil {
		return err
	}

	// Write out the new profile
	output := ctx.String("output")
	if err := WriteProfilesToFile(output, profiles); err != nil {
		return err
	}

	return nil
}

var mergeCommand = cli.Command{
	Name:      "merge",
	Usage:     "merge two or more profiles",
	ArgsUsage: "profile [profile]*",
	Action:    merge,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "output",
			Value: "-",
			Usage: "output file",
		},
	},
}
