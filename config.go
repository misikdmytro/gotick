package gotick

import "time"

type SchedulerConfiguration struct {
	// How often the scheduler should check for jobs to run.
	PollInterval time.Duration

	// Time for which a job can be planned in advance.
	MaxPlanAhead time.Duration

	// Planner factory function.
	PlannerFactory func() Planner

	// Driver factory function.
	DriverFactory func() SchedulerDriver
}

type SchedulerOption func(*SchedulerConfiguration)

func DefaultConfig(options ...SchedulerOption) SchedulerConfiguration {
	config := SchedulerConfiguration{
		PollInterval: 1 * time.Second,
		MaxPlanAhead: 1 * time.Minute,
	}

	config.PlannerFactory = func() Planner {
		return newPlanner(1)
	}

	config.DriverFactory = func() SchedulerDriver {
		return nil
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

func WithPollInterval(interval time.Duration) SchedulerOption {
	return func(config *SchedulerConfiguration) {
		config.PollInterval = interval
	}
}

func WithMaxPlanAhead(maxPlanAhead time.Duration) SchedulerOption {
	return func(config *SchedulerConfiguration) {
		config.MaxPlanAhead = maxPlanAhead
	}
}

func WithDefaultPlannerFactory(threads int) SchedulerOption {
	return func(config *SchedulerConfiguration) {
		config.PlannerFactory = func() Planner {
			return newPlanner(threads)
		}
	}
}

func WithPlannerFactory(factory func() Planner) SchedulerOption {
	return func(config *SchedulerConfiguration) {
		config.PlannerFactory = factory
	}
}

func WithDriverFactory(factory func() SchedulerDriver) SchedulerOption {
	return func(config *SchedulerConfiguration) {
		config.DriverFactory = factory
	}
}
