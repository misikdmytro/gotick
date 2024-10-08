package gotick

import (
	"context"
	"time"
)

// BackgroundService is an interface that represents a service that can be started and stopped in the background.
type BackgroundService interface {
	// Start starts the service in the background.
	Start(context.Context) error

	// Stop stops the service.
	Stop() error
}

// Job represents a job that can be executed by a scheduler.
type Job interface {
	// ID returns the unique identifier of the job.
	// Should be unique across all jobs and the same ID should be returned for one job instance.
	ID() string

	// Execute executes the job.
	Execute(*JobExecutionContext)
}

// Timeout represents an interface that exposes the timeout (used as extension for a job).
type Timeout interface {
	Timeout() time.Duration
}

// JobSchedule represents a schedule for a job.
type JobSchedule interface {
	// Schedule returns the string representation of the schedule.
	Schedule() string

	// First returns the first time the job should be executed.
	First() time.Time

	// Next returns the next time the job should be executed after the provided time.
	// Nil is returned if the job should not be executed anymore.
	Next(time.Time) *time.Time
}

// MaxDelay represents an interface that exposes the maximum delay (used as extension for a schedule).
type MaxDelay interface {
	MaxDelay() time.Duration
}

// PlannerSubscriber is an interface that represents a subscriber to a planner.
type PlannerSubscriber interface {
	// OnJobExecutionUnplanned is called when a job execution was not planned successfully.
	OnJobExecutionUnplanned(*JobExecutionContext)

	// OnBeforeJobExecution is called before a job execution.
	OnBeforeJobExecution(*JobExecutionContext)

	// OnJobExecuted is called when a job is executed.
	OnJobExecuted(*JobExecutionContext)
}

type Planner interface {
	BackgroundService

	// Subscribe subscribes to the planner updates.
	Subscribe(PlannerSubscriber)

	// Plan plans a job execution.
	Plan(*JobExecutionContext)
}

// SchedulerSubscriber is an interface that represents a subscriber to a scheduler.
type SchedulerSubscriber interface {
	PlannerSubscriber

	// OnStart is called when the scheduler is started.
	OnStart()

	// OnStop is called when the scheduler is stopped.
	OnStop()

	// OnJobExecutionInitiated is called when a job execution is initiated.
	OnJobExecutionInitiated(*JobExecutionContext)

	// OnJobExecutionDelayed is called when a job execution is delayed.
	OnJobExecutionDelayed(*JobExecutionContext)

	// OnJobExecutionSkipped is called when a job execution is skipped.
	OnJobExecutionSkipped(*JobExecutionContext)

	// OnBeforeJobExecutionPlan is called before a job execution is planned.
	OnBeforeJobExecutionPlan(*JobExecutionContext)
}

// Scheduler is an interface that represents a job scheduler.
type Scheduler interface {
	BackgroundService

	// Subscribe subscribes to the scheduler updates.
	Subscribe(SchedulerSubscriber)

	// RegisterJob registers a job in the scheduler.
	// Job should be registered before it can be scheduled.
	RegisterJob(job Job) error

	// ScheduleJob schedules a job with the provided schedule.
	// Returns the schedule ID of the scheduled job.
	ScheduleJob(ctx context.Context, jobID string, schedule JobSchedule) (string, error)

	// UnscheduleJobByJobID unschedules a job by its ID.
	UnscheduleJobByJobID(ctx context.Context, jobID string) error

	// UnscheduleJobByScheduleID unschedules a job by its schedule ID.
	UnscheduleJobByScheduleID(ctx context.Context, scheduleID string) error
}

// SchedulerDriver is an interface that represents a driver (storage) for a scheduler that can schedule jobs, unschedule jobs and plans job executions.
type SchedulerDriver interface {
	// ScheduleJob schedules a job with the provided schedule.
	ScheduleJob(ctx context.Context, job Job, schedule JobSchedule) (string, error)

	// UnscheduleJobByJobID unschedules a job by its ID.
	UnscheduleJobByJobID(ctx context.Context, jobID string) error

	// UnscheduleJobByScheduleID unschedules a job by its schedule ID.
	UnscheduleJobByScheduleID(ctx context.Context, scheduleID string) error

	// NextExecution returns the next job execution that should be executed.
	// Nil is returned if currently there are no more job executions.
	NextExecution(context.Context) *JobPlannedExecution
}
