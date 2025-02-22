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
	"fmt"

	"github.com/urfave/cli/v2"
)

func doSort(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() != 1 {
		return fmt.Errorf("expecting one argument, got %d", args.Len())
	}

	// Just need to parse the profile and write it back. The parsing function is
	// sorting the blocks.
	profiles, err := ParseProfiles(args.Get(0))
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

var sortCommand = cli.Command{

	Name:      "sort",
	Usage:     "sort coverage blocks by file then start position",
	ArgsUsage: "profile",
	Action:    doSort,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "output",
			Value: "-",
			Usage: "output file",
		},
	},
}
