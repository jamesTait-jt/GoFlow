package goflow

import (
	"github.com/jamesTait-jt/goflow/pkg/broker"
	"github.com/jamesTait-jt/goflow/pkg/store"
	"github.com/jamesTait-jt/goflow/pkg/task"
)

// Option configures various aspects of GoFlow. It can be used to set
// any customizable part of GoFlow.
type Option interface {
	apply(*options)
}

type options struct {
	taskHandlerStore KVStore[string, task.Handler]
	resultsStore     KVStore[string, task.Result]
	taskBroker       Broker[task.Task]
	resultBroker     Broker[task.Result]
}

// defaultOptions returns a set of default options for configuring GoFlow.
// It provides an in-memory task handler store, an in-memory results store,
// and a channel-based task broker with a buffer size of 10.
//
// These defaults are suitable for most simple use cases, but they can be
// customized further by passing specific options to the GoFlow constructor.
func defaultOptions() options {
	defaultTaskBrokerChannelSize := 10

	return options{
		taskHandlerStore: store.NewInMemoryKVStore[string, task.Handler](),
		resultsStore:     store.NewInMemoryKVStore[string, task.Result](),
		taskBroker:       broker.NewChannelBroker[task.Task](defaultTaskBrokerChannelSize),
		resultBroker:     broker.NewChannelBroker[task.Result](0),
	}
}

type taskHandlerStoreOption struct {
	TaskHandlerStore KVStore[string, task.Handler]
}

func (t taskHandlerStoreOption) apply(opts *options) {
	opts.taskHandlerStore = t.TaskHandlerStore
}

// WithTaskHandlerStore allows the user to provide a custom task handler store,
// which will be used to register and retrieve task handlers for GoFlow.
func WithTaskHandlerStore(taskHandlerStore KVStore[string, task.Handler]) Option {
	return taskHandlerStoreOption{TaskHandlerStore: taskHandlerStore}
}

type resultsStoreOption struct {
	ResultsStore KVStore[string, task.Result]
}

func (r resultsStoreOption) apply(opts *options) {
	opts.resultsStore = r.ResultsStore
}

// WithResultsStore allows the user to provide a custom results store, which will be used
// to persist and retrieve task results in GoFlow.
func WithResultsStore(resultsStore KVStore[string, task.Result]) Option {
	return resultsStoreOption{ResultsStore: resultsStore}
}

type taskBrokerOption struct {
	TaskBroker Broker[task.Task]
}

func (t taskBrokerOption) apply(opts *options) {
	opts.taskBroker = t.TaskBroker
}

// WithTaskBroker allows the user to set a custom task broker, which is responsible
// for managing task submission and distribution to workers in GoFlow.
func WithTaskBroker(taskBroker Broker[task.Task]) Option {
	return taskBrokerOption{TaskBroker: taskBroker}
}

type resultBrokerOption struct {
	ResultBroker Broker[task.Result]
}

func (t resultBrokerOption) apply(opts *options) {
	opts.resultBroker = t.ResultBroker
}

// WithResultBroker allows the user to set a custom result broker, which is provided
// to give workers a way to send their results back
func WithResultBroker(taskBroker Broker[task.Result]) Option {
	return resultBrokerOption{ResultBroker: taskBroker}
}
