// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package eval

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

type seededFakeEvaluator struct {
	seed int
}

func (e seededFakeEvaluator) Evaluate(_ context.Context, req BatchRequestItem[int]) (int, error) {
	return e.seed + req.Candidate + len(req.ScenarioID), nil
}

type concurrencyEvaluator struct {
	sleep time.Duration

	mu          sync.Mutex
	inFlight    int
	maxInFlight int
}

func (e *concurrencyEvaluator) Evaluate(_ context.Context, req BatchRequestItem[int]) (int, error) {
	e.mu.Lock()
	e.inFlight++
	if e.inFlight > e.maxInFlight {
		e.maxInFlight = e.inFlight
	}
	e.mu.Unlock()

	time.Sleep(e.sleep)

	e.mu.Lock()
	e.inFlight--
	e.mu.Unlock()
	return req.Candidate, nil
}

func (e *concurrencyEvaluator) maxSeen() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.maxInFlight
}

type cancelEvaluator struct {
	sleep time.Duration
}

func (e cancelEvaluator) Evaluate(ctx context.Context, _ BatchRequestItem[int]) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-time.After(e.sleep):
	}
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return 1, nil
	}
}

func TestRunBatchPreservesRequestsAndIsDeterministic(t *testing.T) {
	runner := WorkerRunner[int, int]{
		WorkerCount: 3,
		Evaluator:   seededFakeEvaluator{seed: 7},
	}
	req := BatchRequest[int]{
		Items: []BatchRequestItem[int]{
			{CandidateID: "cand-3", ScenarioID: "train-0", Candidate: 3},
			{CandidateID: "cand-5", ScenarioID: "train-1", Candidate: 5},
			{CandidateID: "cand-8", ScenarioID: "val-0", Candidate: 8},
			{CandidateID: "cand-13", ScenarioID: "test-0", Candidate: 13},
		},
	}

	got1, err := runner.RunBatch(context.Background(), req)
	if err != nil {
		t.Fatalf("RunBatch #1: %v", err)
	}
	got2, err := runner.RunBatch(context.Background(), req)
	if err != nil {
		t.Fatalf("RunBatch #2: %v", err)
	}

	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("determinism mismatch:\n got1: %#v\n got2: %#v", got1, got2)
	}
	if len(got1.Items) != len(req.Items) {
		t.Fatalf("result count mismatch: got %d want %d", len(got1.Items), len(req.Items))
	}
	for i, item := range req.Items {
		if got1.Items[i].CandidateID != item.CandidateID {
			t.Fatalf("candidate id mismatch at %d: got %q want %q", i, got1.Items[i].CandidateID, item.CandidateID)
		}
		if got1.Items[i].ScenarioID != item.ScenarioID {
			t.Fatalf("scenario id mismatch at %d: got %q want %q", i, got1.Items[i].ScenarioID, item.ScenarioID)
		}
		if got1.Items[i].Err != nil {
			t.Fatalf("unexpected error at %d: %v", i, got1.Items[i].Err)
		}
	}
}

func TestRunBatchRespectsWorkerLimit(t *testing.T) {
	evaluator := &concurrencyEvaluator{sleep: 20 * time.Millisecond}
	runner := WorkerRunner[int, int]{
		WorkerCount: 2,
		Evaluator:   evaluator,
	}
	req := BatchRequest[int]{
		Items: []BatchRequestItem[int]{
			{CandidateID: "c0", ScenarioID: "s", Candidate: 0},
			{CandidateID: "c1", ScenarioID: "s", Candidate: 1},
			{CandidateID: "c2", ScenarioID: "s", Candidate: 2},
			{CandidateID: "c3", ScenarioID: "s", Candidate: 3},
			{CandidateID: "c4", ScenarioID: "s", Candidate: 4},
			{CandidateID: "c5", ScenarioID: "s", Candidate: 5},
		},
	}

	if _, err := runner.RunBatch(context.Background(), req); err != nil {
		t.Fatalf("RunBatch: %v", err)
	}
	if got := evaluator.maxSeen(); got > runner.WorkerCount {
		t.Fatalf("max in-flight workers exceeded limit: got %d limit %d", got, runner.WorkerCount)
	}
}

func TestRunBatchContextCancellationMarksUnstartedWork(t *testing.T) {
	runner := WorkerRunner[int, int]{
		WorkerCount: 2,
		Evaluator:   cancelEvaluator{sleep: 100 * time.Millisecond},
	}
	req := BatchRequest[int]{
		Items: []BatchRequestItem[int]{
			{CandidateID: "c0", ScenarioID: "s0", Candidate: 0},
			{CandidateID: "c1", ScenarioID: "s1", Candidate: 1},
			{CandidateID: "c2", ScenarioID: "s2", Candidate: 2},
			{CandidateID: "c3", ScenarioID: "s3", Candidate: 3},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	got, err := runner.RunBatch(ctx, req)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("RunBatch error = %v, want %v", err, context.DeadlineExceeded)
	}
	if len(got.Items) != len(req.Items) {
		t.Fatalf("result count mismatch: got %d want %d", len(got.Items), len(req.Items))
	}
	for i, item := range got.Items {
		if !errors.Is(item.Err, context.DeadlineExceeded) {
			t.Fatalf("result %d error = %v, want %v", i, item.Err, context.DeadlineExceeded)
		}
	}
}

func TestRunBatchRequiresWorkerCount(t *testing.T) {
	runner := WorkerRunner[int, int]{
		WorkerCount: 0,
		Evaluator:   seededFakeEvaluator{seed: 1},
	}
	_, err := runner.RunBatch(context.Background(), BatchRequest[int]{})
	if !errors.Is(err, errInvalidWorkerCount) {
		t.Fatalf("RunBatch error = %v, want %v", err, errInvalidWorkerCount)
	}
}

func TestRunBatchRequiresEvaluator(t *testing.T) {
	runner := WorkerRunner[int, int]{
		WorkerCount: 1,
	}
	_, err := runner.RunBatch(context.Background(), BatchRequest[int]{})
	if !errors.Is(err, errNilEvaluator) {
		t.Fatalf("RunBatch error = %v, want %v", err, errNilEvaluator)
	}
}
