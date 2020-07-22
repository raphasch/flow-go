package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dapperlabs/flow-go/module/trace"
)

type VerificationCollector struct {
	tracer          *trace.OpenTracer
	storagePerChunk prometheus.Gauge // storage per chunk

	// Finder Engine
	rcvReceiptsTotal         prometheus.Counter // total execution receipts arrived at finder engine
	sntExecutionResultsTotal prometheus.Counter // total execution results processed by finder engine

	// Match Engine
	rcvExecutionResultsTotal prometheus.Counter // total execution results received by match engine
	sntVerifiableChunksTotal prometheus.Counter // total chunks matched by match engine and sent to verifier engine
	rcvChunkDataPackTotal    prometheus.Counter // total chunk data packs received by match engine
	reqChunkDataPackTotal    prometheus.Counter // total number of chunk data packs requested by match engine

	// Verifier Engine
	rcvVerifiableChunksTotal prometheus.Counter // total verifiable chunks received by verifier engine
	sntResultApprovalsTotal  prometheus.Counter // total result approvals sent by verifier engine

}

func NewVerificationCollector(tracer *trace.OpenTracer, registerer prometheus.Registerer) *VerificationCollector {

	vc := &VerificationCollector{
		tracer: tracer,

		// Finder Engine
		rcvReceiptsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "execution_receipt_received_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemFinderEngine,
			Help:      "total number of execution receipts received by finder engine",
		}),

		sntExecutionResultsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "execution_result_sent_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemFinderEngine,
			Help:      "total number of execution results sent by finder engine to match engine",
		}),

		// Match Engine
		rcvExecutionResultsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "execution_result_received_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemMatchEngine,
			Help:      "total number of execution results received by match engine from finder engine",
		}),

		sntVerifiableChunksTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "verifiable_chunk_sent_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemMatchEngine,
			Help:      "total number of verifiable chunks sent by match engine to verifier engine",
		}),

		rcvChunkDataPackTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "chunk_data_pack_received_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemMatchEngine,
			Help:      "total number of chunk data packs received by match engine",
		}),

		reqChunkDataPackTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "chunk_data_pack_requested_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemMatchEngine,
			Help:      "total number of chunk data packs requested by match engine",
		}),

		// Verifier Engine
		rcvVerifiableChunksTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "verifiable_chunk_received_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemVerifierEngine,
			Help:      "total number verifiable chunks received by verifier engine from match engine",
		}),

		sntResultApprovalsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name:      "result_approvals_total",
			Namespace: namespaceVerification,
			Subsystem: subsystemVerifierEngine,
			Help:      "total number of emitted result approvals by verifier engine",
		}),

		// Storage
		storagePerChunk: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:      "storage_latest_chunk_size_bytes",
			Namespace: namespaceVerification,
			Help:      "latest ingested chunk resources storage (bytes)",
		}),
	}

	registerAllFields(vc, registerer)

	return vc
}

// OnExecutionReceiptReceived is called whenever a new execution receipt arrives
// at Finder engine. It increments total number of received receipts.
func (vc *VerificationCollector) OnExecutionReceiptReceived() {
	vc.rcvReceiptsTotal.Inc()
}

// OnExecutionResultSent is called whenever a new execution result is sent by
// Finder engine to the match engine. It increments total number of sent execution results.
func (vc *VerificationCollector) OnExecutionResultSent() {
	vc.sntExecutionResultsTotal.Inc()
}

// OnExecutionResultReceived is called whenever a new execution result is successfully received
// by Match engine from Finder engine.
// It increments the total number of received execution results.
func (vc *VerificationCollector) OnExecutionResultReceived() {
	vc.rcvExecutionResultsTotal.Inc()
}

// OnVerifiableChunkSent is called on a successful submission of matched chunk
// by Match engine to Verifier engine.
// It increments the total number of chunks matched by match engine.
func (vc *VerificationCollector) OnVerifiableChunkSent() {
	vc.sntVerifiableChunksTotal.Inc()
}

// OnChunkDataPackReceived is called on a receiving a chunk data pack by Match engine
// It increments the total number of chunk data packs received.
func (vc *VerificationCollector) OnChunkDataPackReceived() {
	vc.rcvChunkDataPackTotal.Inc()
}

// OnChunkDataPackRequested is called on requesting a chunk data pack by Match engine
// It increments the total number of chunk data packs requested.
func (vc *VerificationCollector) OnChunkDataPackRequested() {
	vc.reqChunkDataPackTotal.Inc()
}

// OnVerifiableChunkReceived is called whenever a verifiable chunk is received by Verifier engine
// from Match engine.It increments the total number of sent verifiable chunks.
func (vc *VerificationCollector) OnVerifiableChunkReceived() {
	vc.rcvVerifiableChunksTotal.Inc()
}

// OnResultApproval is called whenever a result approval for is emitted to consensus nodes.
// It increases the total number of result approvals.
func (vc *VerificationCollector) OnResultApproval() {
	// increases the counter of disseminated result approvals
	// fo by one. Each result approval corresponds to a single chunk of the block
	// the approvals disseminated by verifier engine
	vc.sntResultApprovalsTotal.Inc()

}

// LogVerifiableChunkSize is called whenever a verifiable chunk is shaped for a specific
// chunk. It adds the size of the verifiable chunk to the histogram. A verifiable chunk is assumed
// to capture all the resources needed to verify a chunk.
// The purpose of this function is to track the overall chunk resources size on disk.
// Todo wire this up to do monitoring (3183)
func (vc *VerificationCollector) LogVerifiableChunkSize(size float64) {
	vc.storagePerChunk.Set(size)
}
