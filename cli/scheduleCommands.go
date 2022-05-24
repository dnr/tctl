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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	schedpb "go.temporal.io/api/schedule/v1"
	"go.temporal.io/api/taskqueue/v1"
	workflowpb "go.temporal.io/api/workflow/v1"
	"go.temporal.io/server/common/primitives/timestamp"

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
	frontendClient = cFactory.FrontendClient(c)
	namespace, err = getRequiredGlobalOption(c, FlagNamespace)
	if err != nil {
		return nil, "", "", err
	}
	scheduleID = c.String(FlagScheduleID)
	if scheduleID == "" {
		return nil, "", "", errors.New("empty schedule id")
	}
	return frontendClient, namespace, scheduleID, nil
}

func buildCalendarSpec(s string) (*schedpb.CalendarSpec, error) {
	var cal schedpb.CalendarSpec
	err := json.Unmarshal([]byte(s), &cal)
	if err != nil {
		return nil, err
	}
	return &cal, nil
}

func buildIntervalSpec(s string) (*schedpb.IntervalSpec, error) {
	var interval, phase time.Duration
	var err error
	parts := strings.Split(s, "/")
	if len(parts) > 2 {
		return nil, errors.New("Invalid interval string")
	} else if len(parts) == 2 {
		if phase, err = timestamp.ParseDuration(parts[1]); err != nil {
			return nil, err
		}
	}
	if interval, err = timestamp.ParseDuration(parts[0]); err != nil {
		return nil, err
	}
	return &schedpb.IntervalSpec{Interval: &interval, Phase: &phase}, nil
}

func buildScheduleSpec(c *cli.Context) (*schedpb.ScheduleSpec, error) {
	now := time.Now()

	var out schedpb.ScheduleSpec
	for _, s := range c.StringSlice("calendar") {
		cal, err := buildCalendarSpec(s)
		if err != nil {
			return nil, err
		}
		out.Calendar = append(out.Calendar, cal)
	}
	for _, s := range c.StringSlice("interval") {
		cal, err := buildIntervalSpec(s)
		if err != nil {
			return nil, err
		}
		out.Interval = append(out.Interval, cal)
	}
	if c.IsSet("start-time") {
		t, err := parseTime(c.String("start-time"), time.Time{}, now)
		if err != nil {
			return nil, err
		}
		out.StartTime = timestamp.TimePtr(t)
	}
	if c.IsSet("end-time") {
		t, err := parseTime(c.String("end-time"), time.Time{}, now)
		if err != nil {
			return nil, err
		}
		out.EndTime = timestamp.TimePtr(t)
	}
	if c.IsSet("jitter") {
		d, err := timestamp.ParseDuration(c.String("jitter"))
		if err != nil {
			return nil, err
		}
		out.Jitter = timestamp.DurationPtr(d)
	}
	if c.IsSet("time-zone") {
		// TODO: validate against tzdb
		out.TimezoneName = c.String("time-zone")
	}
	return &out, nil
}

func buildScheduleAction(c *cli.Context) (*schedpb.ScheduleAction, error) {
	// TODO: factor common code out of here and RunWorkflow

	taskQueue := c.String(FlagTaskQueue)
	workflowType := c.String(FlagWorkflowType)
	et := c.Int(FlagWorkflowExecutionTimeout)
	rt := c.Int(FlagWorkflowRunTimeout)
	dt := c.Int(FlagWorkflowTaskTimeout)
	wid := c.String(FlagWorkflowID)
	if len(wid) == 0 {
		wid = uuid.New()
	}
	inputs, err := processJSONInput(c)
	if err != nil {
		return nil, err
	}
	// memo, err := unmarshalMemoFromCLI(c)
	// if err != nil {
	// 	return nil, err
	// }
	// searchAttr, err := unmarshalSearchAttrFromCLI(c)
	// if err != nil {
	// 	return nil, err
	// }
	newWorkflow := &workflowpb.NewWorkflowExecutionInfo{
		WorkflowId:               wid,
		WorkflowType:             &common.WorkflowType{Name: workflowType},
		TaskQueue:                &taskqueue.TaskQueue{Name: taskQueue},
		Input:                    inputs,
		WorkflowExecutionTimeout: timestamp.DurationPtr(time.Second * time.Duration(et)),
		WorkflowRunTimeout:       timestamp.DurationPtr(time.Second * time.Duration(rt)),
		WorkflowTaskTimeout:      timestamp.DurationPtr(time.Second * time.Duration(dt)),
		// TODO: retry policy (not implemented for workflow yet?)
		// TODO: memo
		// Memo: memo,
		// TODO: search attributes
		// SearchAttributes: searchAttr,
		// TODO: header
	}

	return &schedpb.ScheduleAction{
		Action: &schedpb.ScheduleAction_StartWorkflow{
			StartWorkflow: newWorkflow,
		},
	}, nil
}

func buildScheduleState(c *cli.Context) (*schedpb.ScheduleState, error) {
	var out schedpb.ScheduleState
	out.Notes = c.String("initial-notes")
	out.Paused = c.Bool("initial-paused")
	if c.IsSet("remaining-actions") {
		out.LimitedActions = true
		out.RemainingActions = int64(c.Int("remaining-actions"))
	}
	return &out, nil
}

func getOverlapPolicy(c *cli.Context) (enumspb.ScheduleOverlapPolicy, error) {
	i, err := stringToEnum(c.String("overlap-policy"), enumspb.ScheduleOverlapPolicy_value)
	if err != nil {
		return 0, err
	}
	return enumspb.ScheduleOverlapPolicy(i), nil
}

func buildSchedulePolicies(c *cli.Context) (*schedpb.SchedulePolicies, error) {
	var out schedpb.SchedulePolicies
	var err error
	out.OverlapPolicy, err = getOverlapPolicy(c)
	if err != nil {
		return nil, err
	}
	if c.IsSet("catchup-window") {
		d, err := timestamp.ParseDuration(c.String("catchup-window"))
		if err != nil {
			return nil, err
		}
		out.CatchupWindow = timestamp.DurationPtr(d)
	}
	out.PauseOnFailure = c.Bool("pause-on-failure")
	return &out, nil
}

func buildSchedule(c *cli.Context) (*schedpb.Schedule, error) {
	sched := &schedpb.Schedule{}
	var err error
	if sched.Spec, err = buildScheduleSpec(c); err != nil {
		return nil, err
	}
	if sched.Action, err = buildScheduleAction(c); err != nil {
		return nil, err
	}
	if sched.Policies, err = buildSchedulePolicies(c); err != nil {
		return nil, err
	}
	if sched.State, err = buildScheduleState(c); err != nil {
		return nil, err
	}
	return sched, nil
}

func CreateSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	sched, err := buildSchedule(c)
	if err != nil {
		return err
	}

	// TODO: memo and search attributes for schedule itself

	req := &workflowservice.CreateScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Schedule:   sched,
		Identity:   getCliIdentity(),
		RequestId:  uuid.New(),
	}

	resp, err := frontendClient.CreateSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create schedule.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil
}

func UpdateSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	sched, err := buildSchedule(c)
	if err != nil {
		return err
	}

	// TODO: memo and search attributes for schedule itself

	req := &workflowservice.UpdateScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Schedule:   sched,
		Identity:   getCliIdentity(),
		RequestId:  uuid.New(),
	}

	resp, err := frontendClient.UpdateSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create schedule.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil

	return nil
}

func ToggleSchedule(c *cli.Context) error {
	frontendClient, namespace, scheduleID, err := scheduleBaseArgs(c)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	pause, unpause := c.Bool("pause"), c.Bool("unpause")
	if pause && unpause {
		return errors.New("Cannot specify both --pause and --unpause")
	} else if !pause && !unpause {
		return errors.New("Must specify one of --pause and --unpause")
	}
	patch := &schedpb.SchedulePatch{}
	if pause {
		patch.Pause = c.String(FlagReason)
	} else if unpause {
		patch.Unpause = c.String(FlagReason)
	}

	req := &workflowservice.PatchScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Patch:      patch,
		Identity:   getCliIdentity(),
		RequestId:  uuid.New(),
	}
	resp, err := frontendClient.PatchSchedule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to toggle schedule.\n%s", err)
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

	overlap, err := getOverlapPolicy(c)
	if err != nil {
		return err
	}

	req := &workflowservice.PatchScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Patch: &schedpb.SchedulePatch{
			TriggerImmediately: &schedpb.TriggerImmediatelyRequest{
				OverlapPolicy: overlap,
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

	now := time.Now()
	startTime, err := parseTime(c.String("start-time"), time.Time{}, now)
	if err != nil {
		return err
	}
	endTime, err := parseTime(c.String("end-time"), time.Time{}, now)
	if err != nil {
		return err
	}
	overlap, err := getOverlapPolicy(c)
	if err != nil {
		return err
	}

	req := &workflowservice.PatchScheduleRequest{
		Namespace:  namespace,
		ScheduleId: scheduleID,
		Patch: &schedpb.SchedulePatch{
			BackfillRequest: []*schedpb.BackfillRequest{
				&schedpb.BackfillRequest{
					StartTime:     timestamp.TimePtr(startTime),
					EndTime:       timestamp.TimePtr(endTime),
					OverlapPolicy: overlap,
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

func ListSchedules(c *cli.Context) error {
	frontendClient := cFactory.FrontendClient(c)
	namespace, err := getRequiredGlobalOption(c, FlagNamespace)
	if err != nil {
		return err
	}
	ctx, cancel := newContext(c)
	defer cancel()

	req := &workflowservice.ListSchedulesRequest{
		Namespace: namespace,
	}
	resp, err := frontendClient.ListSchedules(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to list schedules.\n%s", err)
	}

	// TODO: make prettier
	prettyPrintJSONObject(resp)

	return nil
}
