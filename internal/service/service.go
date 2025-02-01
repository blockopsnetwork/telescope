// Package service defines a pluggable service for the Flow system.
//
// Services are low-level constructs which run for the lifetime of the Flow
// controller, and are given deeper levels of access to the overall system
// compared to components, such as the individual instances of running
// components.
package service

import (
	"context"
	"fmt"

	"github.com/blockopsnetwork/telescope/internal/component"
	"github.com/blockopsnetwork/telescope/internal/featuregate"
)

// Definition describes an individual Flow service. Services have unique names
// and optional ConfigTypes where they can be configured within the root Flow
// module.
type Definition struct {
	// Name uniquely defines a service.
	Name string

	// ConfigType is an optional config type to configure a
	// service at runtime. The Name of the service is used
	// as the River block name to configure the service.
	// If nil, the service has no runtime configuration.
	//
	// When non-nil, ConfigType must be a struct type with River
	// tags for decoding as a config block.
	ConfigType any

	// DependsOn defines a set of dependencies for a
	// specific service by name. If DependsOn includes an invalid
	// reference to a service (either because of a cyclic dependency,
	// or a named service doesn't exist), it is treated as a fatal
	// error and the root Flow module will exit.
	DependsOn []string

	// Stability is the overall stability level of the service. This is used to
	// make sure the user is not accidentally configuring a service that is not
	// yet stable - users need to explicitly enable less-than-stable services
	// via, for example, a command-line flag. If a service is not stable enough,
	// an attempt to configure it via the controller will fail.
	// This field must be set to a non-zero value.
	Stability featuregate.Stability
}

// Host is a controller for services and Flow components.
type Host interface {
	// GetComponent gets a running component by ID.
	//
	// GetComponent returns [component.ErrComponentNotFound] if a component is
	// not found.
	GetComponent(id component.ID, opts component.InfoOptions) (*component.Info, error)

	// ListComponents lists all running components within a given module.
	//
	// Returns [component.ErrModuleNotFound] if the provided moduleID doesn't
	// exist.
	ListComponents(moduleID string, opts component.InfoOptions) ([]*component.Info, error)

	// GetService gets a running service using its name.
	GetService(name string) (Service, bool)

	// GetServiceConsumers gets the list of services which depend on a service by
	// name.
	GetServiceConsumers(serviceName string) []Consumer

	// NewController returns an unstarted, isolated Controller that a Service
	// can use to instantiate its own components.
	NewController(id string) Controller
}

// Controller is implemented by flow.Flow.
type Controller interface {
	Run(ctx context.Context)
	LoadSource(source []byte, args map[string]any) error
	Ready() bool
}

type Consumer struct {
	Type ConsumerType // Type of consumer.
	ID   string       // Unique identifier for the consumer.

	// Value of the consumer. When Type is ConsumerTypeComponent, this is an
	// instance of [component.Component]. When Type is ConsumerTypeServcice, this
	// is an instance of [Service].
	Value any
}

// ConsumerType represents the type of consumer who is consuming a service.
type ConsumerType int

const (
	// ConsumerTypeInvalid is the default value for ConsumerType.
	ConsumerTypeInvalid ConsumerType = iota

	ConsumerTypeService // ConsumerTypeService represents a service which uses another service.
)

// String returns a string representation of the ConsumerType.
func (ct ConsumerType) String() string {
	switch ct {
	case ConsumerTypeInvalid:
		return "invalid"
	case ConsumerTypeService:
		return "service"
	}

	return fmt.Sprintf("ConsumerType(%d)", ct)
}

// Service is an individual service to run.
type Service interface {
	// Definition returns the Definition of the Service.
	// Definition must always return the same value across all
	// calls.
	Definition() Definition

	// Run starts a Service. Run must block until the provided
	// context is canceled. Returning an error should be treated
	// as a fatal error for the Service.
	Run(ctx context.Context, host Host) error

	// Update updates a Service at runtime. Update is never
	// called if [Definition.ConfigType] is nil. newConfig will
	// be the same type as ConfigType; if ConfigType is a
	// pointer to a type, newConfig will be a pointer to the
	// same type.
	//
	// Update will be called once before Run, and may be called
	// while Run is active.
	Update(newConfig any) error

	// Data returns the Data associated with a Service. Data must always return
	// the same value across multiple calls, as callers are expected to be able
	// to cache the result.
	//
	// The return result of Data must not rely on the runtime config of the
	// service.
	//
	// Data may be invoked before Run.
	Data() any
}
