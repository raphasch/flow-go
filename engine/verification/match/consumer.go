package match

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/jobqueue"
	"github.com/onflow/flow-go/storage"
)

type ChunkJob struct {
	Chunk *flow.Chunk
}

func (j *ChunkJob) ID() module.JobID {
	return module.JobID(fmt.Sprintf("%v", j.Chunk.ID()))
}

type ChunksJob struct {
	chunks storage.ChunksQueue
}

func (j ChunksJob) AtIndex(index int64) (storage.Job, error) {
	chunk, err := j.chunks.AtIndex(index)
	if err != nil {
		return nil, fmt.Errorf("could not read chunk: %w", err)
	}
	return &ChunkJob{Chunk: chunk}, nil
}

type Worker struct {
	engine *Engine
}

func (w *Worker) Run(job storage.Job) {
	chunkjob, _ := job.(*ChunkJob)
	w.engine.ProcessMyChunk(chunkjob.Chunk)
}

type ChunkConsumer struct {
	consumer module.JobConsumer
}

func NewChunkConsumer(log zerolog.Logger, processedIndex storage.ConsumerProgress, chunksQueue storage.ChunksQueue, engine *Engine, maxProcessing int64, maxFinished int64) (*ChunkConsumer, jobqueue.Worker) {
	worker := &Worker{engine: engine}
	jobs := &ChunksJob{chunks: chunksQueue}
	consumer := jobqueue.NewConsumer(
		log, jobs, processedIndex, worker, maxProcessing, maxFinished,
	)
	return &ChunkConsumer{
		consumer: consumer,
	}, worker
}

func (c *ChunkConsumer) Ready() <-chan struct{} {
	err := c.consumer.Start()
	if err != nil {
		panic(fmt.Errorf("could not start the chunk consumer for match engine: %w", err))
	}

	ready := make(chan struct{})
	close(ready)
	return ready
}

func (c *ChunkConsumer) Done() <-chan struct{} {
	c.consumer.Stop()

	ready := make(chan struct{})
	close(ready)
	return ready
}