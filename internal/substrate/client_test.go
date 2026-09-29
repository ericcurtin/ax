// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package substrate

import (
	"context"
	"strings"
	"testing"

	"github.com/agent-substrate/substrate/pkg/proto/ateapipb"
	"google.golang.org/grpc"
)

// tokenLoopControl answers every ListActorTemplates call with the same
// next_page_token.
type tokenLoopControl struct {
	ateapipb.ControlClient
	calls int
}

func (c *tokenLoopControl) ListActorTemplates(ctx context.Context, in *ateapipb.ListActorTemplatesRequest, opts ...grpc.CallOption) (*ateapipb.ListActorTemplatesResponse, error) {
	c.calls++
	return &ateapipb.ListActorTemplatesResponse{NextPageToken: "same"}, nil
}

func TestListActorTemplatesRejectsRepeatedToken(t *testing.T) {
	control := &tokenLoopControl{}
	c := &Client{control: control}

	_, err := c.ListActorTemplates(context.Background(), "default")
	if err == nil || !strings.Contains(err.Error(), "repeated page token") {
		t.Fatalf("expected repeated page token error, got %v", err)
	}
	if control.calls != 2 {
		t.Errorf("expected 2 list calls, got %d", control.calls)
	}
}
