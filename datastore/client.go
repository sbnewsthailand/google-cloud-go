// Copyright 2017 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datastore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore/internal"
	cloudinternal "cloud.google.com/go/internal"
	"cloud.google.com/go/internal/trace"
	"cloud.google.com/go/internal/version"
	gax "github.com/googleapis/gax-go/v2"
	pb "google.golang.org/genproto/googleapis/datastore/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// datastoreClient is a wrapper for the pb.DatastoreClient that includes gRPC
// metadata to be sent in each request for server-side traffic management.
type datastoreClient struct {
	// Embed so we still implement the DatastoreClient interface,
	// if the interface adds more methods.
	pb.DatastoreClient

	c  pb.DatastoreClient
	md metadata.MD
}

func newDatastoreClient(conn grpc.ClientConnInterface, projectID string) pb.DatastoreClient {
	return &datastoreClient{
		c: pb.NewDatastoreClient(conn),
		md: metadata.Pairs(
			resourcePrefixHeader, "projects/"+projectID,
			"x-goog-api-client", fmt.Sprintf("gl-go/%s gccl/%s grpc/", version.Go(), internal.Version)),
	}
}

func (dc *datastoreClient) Lookup(ctx context.Context, in *pb.LookupRequest, opts ...grpc.CallOption) (res *pb.LookupResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/datastore.datastoreClient.Lookup")
	defer func() { trace.EndSpan(ctx, err) }()

	err = dc.invoke(ctx, func(ctx context.Context) error {
		res, err = dc.c.Lookup(ctx, in, opts...)
		return err
	})
	return res, err
}

func (dc *datastoreClient) RunQuery(ctx context.Context, in *pb.RunQueryRequest, opts ...grpc.CallOption) (res *pb.RunQueryResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/datastore.datastoreClient.RunQuery")
	defer func() { trace.EndSpan(ctx, err) }()

	err = dc.invoke(ctx, func(ctx context.Context) error {
		res, err = dc.c.RunQuery(ctx, in, opts...)
		return err
	})
	return res, err
}

func (dc *datastoreClient) RunAggregationQuery(ctx context.Context, in *pb.RunAggregationQueryRequest, opts ...grpc.CallOption) (res *pb.RunAggregationQueryResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/datastore.datastoreClient.RunAggregationQuery")
	defer func() { trace.EndSpan(ctx, err) }()

	err = dc.invoke(ctx, func(ctx context.Context) error {
		res, err = dc.c.RunAggregationQuery(ctx, in, opts...)
		return err
	})
	return res, nil
}

func (dc *datastoreClient) BeginTransaction(ctx context.Context, in *pb.BeginTransactionRequest, opts ...grpc.CallOption) (res *pb.BeginTransactionResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/datastore.datastoreClient.BeginTransaction")
	defer func() { trace.EndSpan(ctx, err) }()

	err = dc.invoke(ctx, func(ctx context.Context) error {
		res, err = dc.c.BeginTransaction(ctx, in, opts...)
		return err
	})
	return res, err
}

func (dc *datastoreClient) Commit(ctx context.Context, in *pb.CommitRequest, opts ...grpc.CallOption) (res *pb.CommitResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/datastore.datastoreClient.Commit")
	defer func() { trace.EndSpan(ctx, err) }()

	err = dc.invoke(ctx, func(ctx context.Context) error {
		res, err = dc.c.Commit(ctx, in, opts...)
		return err
	})
	return res, err
}

func (dc *datastoreClient) Rollback(ctx context.Context, in *pb.RollbackRequest, opts ...grpc.CallOption) (res *pb.RollbackResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/datastore.datastoreClient.Rollback")
	defer func() { trace.EndSpan(ctx, err) }()

	err = dc.invoke(ctx, func(ctx context.Context) error {
		res, err = dc.c.Rollback(ctx, in, opts...)
		return err
	})
	return res, err
}

func (dc *datastoreClient) AllocateIds(ctx context.Context, in *pb.AllocateIdsRequest, opts ...grpc.CallOption) (res *pb.AllocateIdsResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/datastore.datastoreClient.AllocateIds")
	defer func() { trace.EndSpan(ctx, err) }()

	err = dc.invoke(ctx, func(ctx context.Context) error {
		res, err = dc.c.AllocateIds(ctx, in, opts...)
		return err
	})
	return res, err
}

func (dc *datastoreClient) invoke(ctx context.Context, f func(ctx context.Context) error) error {
	ctx = metadata.NewOutgoingContext(ctx, dc.md)
	return cloudinternal.Retry(ctx, gax.Backoff{Initial: 100 * time.Millisecond}, func() (stop bool, err error) {
		err = f(ctx)
		return !shouldRetry(err), err
	})
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	// Only retry on UNAVAILABLE as per https://aip.dev/194. Other errors from
	// https://cloud.google.com/datastore/docs/concepts/errors may be retried
	// by the user if desired, but are not retried by the clientg.
	return s.Code() == codes.Unavailable
}
