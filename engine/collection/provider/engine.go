// Package provider implements an engine for providing access to resources held
// by the collection node, including collections, collection guarantees, and
// transactions.
package provider

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/dapperlabs/flow-go/engine"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/model/flow/filter"
	"github.com/dapperlabs/flow-go/model/messages"
	"github.com/dapperlabs/flow-go/module"
	"github.com/dapperlabs/flow-go/module/mempool"
	"github.com/dapperlabs/flow-go/module/trace"
	"github.com/dapperlabs/flow-go/network"
	"github.com/dapperlabs/flow-go/protocol"
	"github.com/dapperlabs/flow-go/storage"
	"github.com/dapperlabs/flow-go/utils/logging"
)

// Engine is the collection provider engine, which provides access to resources
// held by the collection node.
type Engine struct {
	unit         *engine.Unit
	log          zerolog.Logger
	tracer       trace.Tracer
	con          network.Conduit
	me           module.Local
	state        protocol.State
	pool         mempool.Transactions
	collections  storage.Collections
	transactions storage.Transactions
}

func New(log zerolog.Logger, net module.Network, state protocol.State, tracer trace.Tracer, me module.Local, pool mempool.Transactions, collections storage.Collections, transactions storage.Transactions) (*Engine, error) {
	e := &Engine{
		unit:         engine.NewUnit(),
		log:          log.With().Str("engine", "provider").Logger(),
		tracer:       tracer,
		me:           me,
		state:        state,
		pool:         pool,
		collections:  collections,
		transactions: transactions,
	}

	con, err := net.Register(engine.CollectionProvider, e)
	if err != nil {
		return nil, fmt.Errorf("could not register engine: %w", err)
	}

	e.con = con

	return e, nil
}

// Ready returns a ready channel that is closed once the engine has fully
// started.
func (e *Engine) Ready() <-chan struct{} {
	return e.unit.Ready()
}

// Done returns a done channel that is closed once the engine has fully stopped.
func (e *Engine) Done() <-chan struct{} {
	return e.unit.Done()
}

// SubmitLocal submits an event originating on the local node.
func (e *Engine) SubmitLocal(event interface{}) {
	e.Submit(e.me.NodeID(), event)
}

// Submit submits the given event from the node with the given origin ID
// for processing in a non-blocking manner. It returns instantly and logs
// a potential processing error internally when done.
func (e *Engine) Submit(originID flow.Identifier, event interface{}) {
	e.unit.Launch(func() {
		err := e.Process(originID, event)
		if err != nil {
			e.log.Error().Err(err).Msg("could not process submitted event")
		}
	})
}

// ProcessLocal processes an event originating on the local node.
func (e *Engine) ProcessLocal(event interface{}) error {
	return e.Process(e.me.NodeID(), event)
}

// Process processes the given event from the node with the given origin ID in
// a blocking manner. It returns the potential processing error when done.
func (e *Engine) Process(originID flow.Identifier, event interface{}) error {
	return e.unit.Do(func() error {
		return e.process(originID, event)
	})
}

// process processes events for the provider engine on the collection node.
func (e *Engine) process(originID flow.Identifier, event interface{}) error {
	switch ev := event.(type) {
	case *messages.CollectionRequest:
		return e.onCollectionRequest(originID, ev)
	case *messages.SubmitCollectionGuarantee:
		return e.onSubmitCollectionGuarantee(originID, ev)
	case *messages.TransactionRequest:
		return e.onTransactionRequest(originID, ev)
	case *messages.TransactionResponse:
		return e.onTransactionResponse(originID, ev)
	default:
		return fmt.Errorf("invalid event type (%T)", event)
	}
}

func (e *Engine) onCollectionRequest(originID flow.Identifier, req *messages.CollectionRequest) error {
	coll, err := e.collections.ByID(req.ID)
	if err != nil {
		return fmt.Errorf("could not retrieve requested collection: %w", err)
	}

	res := &messages.CollectionResponse{
		Collection: *coll,
	}
	err = e.con.Submit(res, originID)
	if err != nil {
		return fmt.Errorf("could not respond to collection requester: %w", err)
	}

	return nil
}

func (e *Engine) onSubmitCollectionGuarantee(originID flow.Identifier, req *messages.SubmitCollectionGuarantee) error {
	if originID != e.me.NodeID() {
		return fmt.Errorf("invalid remote request to submit collection guarantee [%x]", req.Guarantee.ID())
	}

	return e.SubmitCollectionGuarantee(&req.Guarantee)
}

// onTransactionRequest handles requests for individual transactions.
func (e *Engine) onTransactionRequest(originID flow.Identifier, req *messages.TransactionRequest) error {

	// check the mempool first
	if e.pool.Has(req.ID) {
		tx, err := e.pool.ByID(req.ID)
		if err != nil {
			return fmt.Errorf("could not get transaction from pool: %w", err)
		}

		res := &messages.TransactionResponse{Transaction: *tx}

		err = e.con.Submit(res, originID)
		if err != nil {
			return fmt.Errorf("could not submit transaction resopnse: %w", err)
		}
	}

	// if it isn't in the mempool, check persistent storage
	tx, err := e.transactions.ByID(req.ID)
	if err != nil {
		return fmt.Errorf("could not get transaction from db: %w", err)
	}

	res := &messages.TransactionResponse{Transaction: *tx}

	err = e.con.Submit(res, originID)
	if err != nil {
		return fmt.Errorf("could not submit transaction response: %w", err)
	}

	return nil
}

// onTransactionResponse handles responses for requests we have made for a
// transaction by adding the transaction
func (e *Engine) onTransactionResponse(originID flow.Identifier, res *messages.TransactionResponse) error {

	err := e.pool.Add(&res.Transaction)
	if err != nil {
		e.log.Debug().
			Err(err).
			Hex("tx_id", logging.ID(res.Transaction.ID())).
			Msg("could not add transaction to mempool")
	}

	return nil
}

// SubmitCollectionGuarantee submits the collection guarantee to all
// consensus nodes.
func (e *Engine) SubmitCollectionGuarantee(guarantee *flow.CollectionGuarantee) error {
	defer e.tracer.FinishSpan(guarantee.ID())
	consensusNodes, err := e.state.Final().Identities(filter.HasRole(flow.RoleConsensus))
	if err != nil {
		return fmt.Errorf("could not get consensus consensusNodes: %w", err)
	}

	err = e.con.Submit(guarantee, consensusNodes.NodeIDs()...)
	if err != nil {
		return fmt.Errorf("could not submit collection guarantee: %w", err)
	}

	return nil
}
