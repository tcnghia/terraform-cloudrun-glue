/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package interceptors

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/clog/gcp"
	"google.golang.org/grpc/metadata"
)

func TestTrace(t *testing.T) {
	slog.SetDefault(slog.New(gcp.NewHandler(slog.LevelDebug)))
	// This ensures the metadata server is not called at all during tests.
	md := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected metadata request: %v", r)
	}))
	defer md.Close()
	t.Setenv("GCE_METADATA_HOST", md.URL)
	slog.SetDefault(slog.New(gcp.NewHandler(slog.LevelDebug)))
	for _, c := range []struct {
		name      string
		env       string
		wantTrace bool
	}{
		{"no env set", "", false},
		{"env set", "my-project", true},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("GOOGLE_CLOUD_PROJECT", c.env)

			md := metadata.New(map[string]string{"x-cloud-trace-context": "trace/id/yay"})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			_, _ = ClogTraceContextUnaryInterceptor(ctx, nil, nil, func(ctx context.Context, req any) (any, error) {
				clog.InfoContext(ctx, "hello world")
				if found := ctx.Value(traceContextKey("trace")) != nil; found != c.wantTrace {
					t.Fatalf("got trace context %t, want %t", found, c.wantTrace)
					if c.wantTrace {
						if trace := ctx.Value("trace"); !strings.Contains(trace.(string), "/"+c.env+"/") {
							t.Errorf("got trace context %q, want %q", trace, c.env)
						}
					}
				}
				return nil, nil
			})
		})
	}
}