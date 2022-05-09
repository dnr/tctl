// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cli

import (
	"github.com/temporalio/tctl-kit/pkg/flags"
	"github.com/urfave/cli/v2"
)

func newScheduleCommands() []*cli.Command {
	sid := &cli.StringFlag{
		Name:     FlagScheduleID,
		Aliases:  FlagScheduleIDAlias,
		Usage:    "Schedule id",
		Required: true,
	}

	return []*cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create a new schedule",
			Flags: []cli.Flag{
				sid,
			},
			Action: CreateSchedule,
		},
		{
			Name:    "update",
			Aliases: []string{"up"},
			Usage:   "Updates a schedule with a new definition",
			Flags: []cli.Flag{
				sid,
			},
			Action: UpdateSchedule,
		},
		{
			Name:  "pause",
			Usage: "Pauses or unpauses a schedule",
			Flags: []cli.Flag{
				sid,
			},
			Action: PauseSchedule,
		},
		{
			Name:  "trigger",
			Usage: "Triggers an immediate action",
			Flags: []cli.Flag{
				sid,
			},
			Action: TriggerSchedule,
		},
		{
			Name:  "backfill",
			Usage: "Backfills a past time range of actions",
			Flags: []cli.Flag{
				sid,
			},
			Action: BackfillSchedule,
		},
		{
			Name:    "describe",
			Aliases: []string{"d"},
			Usage:   "Get schedule configuration and current state",
			Flags: append([]cli.Flag{
				sid,
			}, flags.FlagsForRendering...),
			Action: DescribeSchedule,
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Usage:   "Deletes a schedule",
			Flags: []cli.Flag{
				sid,
			},
			Action: DeleteSchedule,
		},
	}
}
