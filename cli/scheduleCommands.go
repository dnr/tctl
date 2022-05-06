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
	"fmt"

	// schedpb "go.temporal.io/api/schedule/v1"

	"github.com/urfave/cli/v2"
	"go.temporal.io/api/workflowservice/v1"
)

// CreateSchedule
func CreateSchedule(c *cli.Context) error {
	return nil
}

// DescribeSchedule
func DescribeSchedule(c *cli.Context) error {
	frontendClient := cFactory.FrontendClient(c)
	namespace, err := getRequiredGlobalOption(c, FlagNamespace)
	if err != nil {
		return err
	}

	scheduleID := c.String(FlagScheduleID)

	ctx, cancel := newContext(c)
	defer cancel()

	req := &workflowservice.DescribeScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
	}
	resp, err := frontendClient.DescribeSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to describe schedule.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil
}
