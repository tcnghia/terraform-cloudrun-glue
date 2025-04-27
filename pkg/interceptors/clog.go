/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package interceptors

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	gcemetadata "cloud.google.com/go/compute/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const CloudTraceHeader = "X-Cloud-Trace-Context"

var (
	projectID  string
	lookupOnce sync.Once
)

type traceContextKey string

func extractProjectID(ctx context.Context) string {
	// Get the project ID from the environment if specified
	fromEnv := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if fromEnv != "" {
		projectID = fromEnv
		return projectID
	}
	lookupOnce.Do(func() {
		p, err := gcemetadata.ProjectIDWithContext(ctx)
		if err == nil {
			projectID = p
		}
	})
	return projectID
}

func gcpTraceURI(projectID, traceHeader string) string {
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		return fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
	}
	return ""
}

// ClogTraceContextUnaryInterceptor is a gRPC unary interceptor that adds the
// Cloud Trace ID to the context. This is used to correlate the structured
// logs with the request log.
func ClogTraceContextUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	projectID := extractProjectID(ctx)
	if projectID != "" {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			traceHeader := ""
			if x := md.Get(CloudTraceHeader); len(x) > 0 {
				traceHeader = x[0]
			}
			if trace := gcpTraceURI(projectID, traceHeader); trace != "" {
				ctx = context.WithValue(ctx, traceContextKey("trace"), trace)
			}
		}
	}
	return handler(ctx, req)
}