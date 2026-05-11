// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package eval

import (
	"context"
	"errors"
	"sync"
)

var (
	errInvalidWorkerCount = errors.New("worker count must be >= 1")
	errNilEvaluator       = errors.New("nil evaluator")
)

// CandidateID is a stable identifier for one candidate under evaluation.
type CandidateID string

// ScenarioID is a stable identifier for one scenario under evaluation.
type ScenarioID string

// BatchRequestItem is one typed candidate-scenario evaluation request.
type BatchRequestItem[C any] struct {
	CandidateID CandidateID `json:"candidate_id"`
	ScenarioID  ScenarioID  `json:"scenario_id"`
	Candidate   C           `json:"candidate"`
}

// BatchResultItem stores one typed candidate-scenario evaluation result.
type BatchResultItem[R any] struct {
	CandidateID CandidateID `json:"candidate_id"`
	ScenarioID  ScenarioID  `json:"scenario_id"`
	Result      R           `json:"result"`
	Err         error       `json:"-"`
}

// BatchRequest groups a set of evaluation requests.
type BatchRequest[C any] struct {
	Items []BatchRequestItem[C] `json:"items"`
}

// BatchResult stores results for a batch request.
type BatchResult[R any] struct {
	Items []BatchResultItem[R] `json:"items"`
}

// Evaluator evaluates one candidate-scenario request.
type Evaluator[C, R any] interface {
	Evaluate(ctx context.Context, req BatchRequestItem[C]) (R, error)
}

// Runner dispatches a batch request and returns the batch result.
type Runner[C, R any] interface {
	RunBatch(ctx context.Context, req BatchRequest[C]) (BatchResult[R], error)
}

// WorkerRunner dispatches batched requests onto a bounded worker pool.
type WorkerRunner[C, R any] struct {
	WorkerCount int
	Evaluator   Evaluator[C, R]
}

// RunBatch evaluates the provided requests with bounded concurrency.
func (r WorkerRunner[C, R]) RunBatch(ctx context.Context, req BatchRequest[C]) (BatchResult[R], error) {
	if r.WorkerCount < 1 {
		return BatchResult[R]{}, errInvalidWorkerCount
	}
	if r.Evaluator == nil {
		return BatchResult[R]{}, errNilEvaluator
	}

	results := make([]BatchResultItem[R], len(req.Items))
	done := make([]bool, len(req.Items))
	for i, item := range req.Items {
		results[i].CandidateID = item.CandidateID
		results[i].ScenarioID = item.ScenarioID
	}
	if len(req.Items) == 0 {
		return BatchResult[R]{Items: results}, nil
	}

	workers := r.WorkerCount
	if workers > len(req.Items) {
		workers = len(req.Items)
	}

	jobs := make(chan int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				result, err := r.Evaluator.Evaluate(ctx, req.Items[idx])
				mu.Lock()
				results[idx].Result = result
				results[idx].Err = err
				done[idx] = true
				mu.Unlock()
			}
		}()
	}

dispatch:
	for i := range req.Items {
		select {
		case <-ctx.Done():
			break dispatch
		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()

	if err := ctx.Err(); err != nil {
		for i := range results {
			if !done[i] {
				results[i].Err = err
			}
		}
		return BatchResult[R]{Items: results}, err
	}
	return BatchResult[R]{Items: results}, nil
}
