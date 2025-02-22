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
	"os"

	"github.com/urfave/cli/v2"
)

const name = "covertool"
const usage = "Perform various operations on Go coverage files"

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = "0.4.1"
	app.Commands = []*cli.Command{
		&mergeCommand,
		&reportCommand,
		&sortCommand,
	}

	if err := app.Run(os.Args); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "%v", err)
		if err != nil {
			os.Exit(2)
		}
	}
}
