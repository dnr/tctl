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
	"errors"
	"fmt"
	"time"

	enumspb "go.temporal.io/api/enums/v1"
	schedpb "go.temporal.io/api/schedule/v1"

	"github.com/pborman/uuid"
	"github.com/urfave/cli/v2"
	"go.temporal.io/api/workflowservice/v1"
)

func scheduleBaseArgs(c *cli.Context) (
	frontendClient workflowservice.WorkflowServiceClient,
	namespace string,
	scheduleID string,
	err error,
) {
	frontendClient := cFactory.FrontendClient(c)
	namespace, err := getRequiredGlobalOption(c, FlagNamespace)
	if err != nil {
		return nil, "", "", err
	}
	scheduleID := c.String(FlagScheduleID)
	if scheduleID == "" {
		return nil, "", "", errors.New("empty schedule id")
	}
	return frontendClient, namespace, scheduleID, nil
}

func CreateSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	return nil
}

func UpdateSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	return nil
}

func PauseSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	req := &workflowservice.PatchScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Patch: &schedpb.SchedulePatch{
			// TODO: get from flag
			Pause: true,
		},
		Identity:  getCliIdentity(),
		RequestId: uuid.New(),
	}
	resp, err := frontendClient.PatchSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to pause schedule.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil
}

func TriggerSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	req := &workflowservice.PatchScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Patch: &schedpb.SchedulePatch{
			TriggerImmediately: &schedpb.TriggerImmediatelyRequest{
				// TODO: get from flag
				OverlapPolicy: enumspb.SCHEDULE_OVERLAP_POLICY_UNSPECIFIED,
			},
		},
		Identity:  getCliIdentity(),
		RequestId: uuid.New(),
	}
	resp, err := frontendClient.PatchSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to trigger schedule.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil
}

func BackfillSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	req := &workflowservice.PatchScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Patch: &schedpb.SchedulePatch{
			BackfillRequest: []*schedpb.BackfillRequest{
				&schedpb.BackfillRequest{
					// TODO: Get these from flags
					StartTime:     &time.Time{},
					EndTime:       &time.Time{},
					OverlapPolicy: enumspb.SCHEDULE_OVERLAP_POLICY_UNSPECIFIED,
				},
			},
		},
		Identity:  getCliIdentity(),
		RequestId: uuid.New(),
	}
	resp, err := frontendClient.PatchSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to backfill schedule.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil
}

func DescribeSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
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

func DeleteSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	req := &workflowservice.DeleteScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Identity:   getCliIdentity(),
	}
	resp, err := frontendClient.DeleteSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete schedule.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil
}
