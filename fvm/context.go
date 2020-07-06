package fvm

import (
	"github.com/dapperlabs/flow-go/model/flow"
)

// A Context defines a set execution parameters used by the virtual machine.
type Context struct {
	ASTCache                         ASTCache
	Blocks                           Blocks
	Metrics                          *MetricsCollector
	GasLimit                         uint64
	BlockHeader                      *flow.Header
	SignatureVerificationEnabled     bool
	FeePaymentsEnabled               bool
	RestrictedAccountCreationEnabled bool
	RestrictedDeploymentEnabled      bool
}

// NewContext initializes a new execution context with the provided options.
func NewContext(opts ...Option) Context {
	return newContext(defaultContext(), opts...)
}

// NewContextFromParent spawns a child execution context with the provided options.
func NewContextFromParent(parent Context, opts ...Option) Context {
	return newContext(parent, opts...)
}

func newContext(ctx Context, opts ...Option) Context {
	for _, applyOption := range opts {
		ctx = applyOption(ctx)
	}

	return ctx
}

const AccountKeyWeightThreshold = 1000

const defaultGasLimit = 100000

func defaultContext() Context {
	return Context{
		ASTCache:                         nil,
		Blocks:                           nil,
		Metrics:                          nil,
		GasLimit:                         defaultGasLimit,
		BlockHeader:                      nil,
		SignatureVerificationEnabled:     true,
		FeePaymentsEnabled:               true,
		RestrictedAccountCreationEnabled: true,
		RestrictedDeploymentEnabled:      true,
	}
}

// An Option sets a configuration parameter for a virtual machine context.
type Option func(ctx Context) Context

// WithASTCache sets the AST cache for a virtual machine context.
func WithASTCache(cache ASTCache) Option {
	return func(ctx Context) Context {
		ctx.ASTCache = cache
		return ctx
	}
}

// WithGasLimit sets the gas limit for a virtual machine context.
func WithGasLimit(limit uint64) Option {
	return func(ctx Context) Context {
		ctx.GasLimit = limit
		return ctx
	}
}

// WithBlockHeader sets the block header for a virtual machine context.
//
// The VM uses the header to provide current block information to the Cadence runtime,
// as well as to seed the pseudorandom number generator.
func WithBlockHeader(header *flow.Header) Option {
	return func(ctx Context) Context {
		ctx.BlockHeader = header
		return ctx
	}
}

// WithBlocks sets the block storage provider for a virtual machine context.
//
// The VM uses the block storage provider to provide historical block information to
// the Cadence runtime.
func WithBlocks(blocks Blocks) Option {
	return func(ctx Context) Context {
		ctx.Blocks = blocks
		return ctx
	}
}

// WithMetricsCollector sets the metrics collector for a virtual machine context.
//
// A metrics collector is used to gather metrics reported by the Cadence runtime.
func WithMetricsCollector(mc *MetricsCollector) Option {
	return func(ctx Context) Context {
		ctx.Metrics = mc
		return ctx
	}
}

// WithSignatureVerification enables or disables signature verification and sequence
// number checks for a virtual machine context.
func WithSignatureVerification(enabled bool) Option {
	return func(ctx Context) Context {
		ctx.SignatureVerificationEnabled = enabled
		return ctx
	}
}

// WithFeePayments enables or disables fee payments for a virtual machine context.
func WithFeePayments(enabled bool) Option {
	return func(ctx Context) Context {
		ctx.FeePaymentsEnabled = enabled
		return ctx
	}
}

// WithRestrictedDeployment enables or disables restricted contract deployment for a
// virtual machine context.
func WithRestrictedDeployment(enabled bool) Option {
	return func(ctx Context) Context {
		ctx.RestrictedDeploymentEnabled = enabled
		return ctx
	}
}

// WithRestrictedAccountCreation enables or disables restricted account creation for a
// virtual machine context
func WithRestrictedAccountCreation(enabled bool) Option {
	return func(ctx Context) Context {
		ctx.RestrictedAccountCreationEnabled = enabled
		return ctx
	}
}