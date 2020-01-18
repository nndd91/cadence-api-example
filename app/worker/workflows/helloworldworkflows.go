package workflows

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const TaskListName = "helloWorldGroup"
const SignalName = "helloWorldSignal"

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	activity.Register(helloworldActivity)
	workflow.Register(waitForAgeResponse)
}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
}

func helloworldActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started")
	return "Hello " + name + "! How old are you!", nil
}

func waitForAgeResponse(ctx workflow.Context) (string, error) {
	signalName := SignalName

	selector := workflow.NewSelector(ctx)

	// First Selector to receive Signal
	var ageResult int

	for {
		signalChan := workflow.GetSignalChannel(ctx, signalName)
		selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
			c.Receive(ctx, &ageResult)
			workflow.GetLogger(ctx).Info("Received age results from signal!", zap.String("signal", signalName), zap.Int("value", ageResult))
		})
		workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
		// Wait for signal
		selector.Select(ctx)

		// Exception Handler for if signal received is invalid!
		if ageResult > 0 {
			if ageResult < 150 {
				return fmt.Sprintf("%v years old", ageResult), nil
			} else {
				fmt.Println("You can't be that old!", zap.Int("Signal", ageResult))

				return "", errors.New("invalid age")
			}
		}
	}
}

func Workflow(ctx workflow.Context, name string) error {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("helloworld workflow started")
	var activityResult string
	err := workflow.ExecuteActivity(ctx, helloworldActivity, name).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	var ageResult string
	err = workflow.ExecuteChildWorkflow(ctx, waitForAgeResponse).Get(ctx, &ageResult)
	if err != nil {
		logger.Error("Childworkflow failed.", zap.Error(err))
		return err
	}
	logger.Info("Workflow completed.", zap.String("NameResult", activityResult), zap.String("AgeResult", ageResult))

	return nil
}
