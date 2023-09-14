// Copyright 2016-2023, Pulumi Corporation.
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

package tests

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/backend/httpstate"
	"github.com/pulumi/pulumi/pkg/v3/backend/httpstate/client"
	ptesting "github.com/pulumi/pulumi/sdk/v3/go/common/testing"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/stretchr/testify/assert"
)

func TestPolicyNewNonInteractive(t *testing.T) {
	t.Parallel()
	e := ptesting.NewEnvironment(t)
	defer deleteIfNotFailed(e)
	e.RunCommand("pulumi", "policy", "new", "aws-typescript", "--force", "--generate-only")
}

func TestPremiumPolicyAuth(t *testing.T) {
	t.Parallel()
	if os.Getenv("PULUMI_ACCESS_TOKEN") == "" {
		t.Skipf("Skipping: PULUMI_ACCESS_TOKEN is not set")
	}

	// TODO remove this
	accessToken := os.Getenv("PULUMI_ACCESS_TOKEN")
	client := client.NewClient(httpstate.DefaultURL(), accessToken, false, cmdutil.Diag())
	username, organizations, err := client.GetPulumiAccountDetails(context.TODO())
	assert.NoError(t, err)
	assert.Empty(t, username)
	assert.Empty(t, organizations)
	e := ptesting.NewEnvironment(t)
	defer deleteIfNotFailed(e)

	e.RunCommand("pulumi", "login")
	defer e.RunCommand("pulumi", "logout")
	// Remove `PULUMI_ACCESS_TOKEN` so that `pulumi policy new` automatically sets it for installation.
	for i, elem := range e.Env {
		if strings.HasPrefix(elem, "PULUMI_ACCESS_TOKEN=") {
			_ = i
			// e.Env = append(e.Env[:i], e.Env[i+1:]...)
			break
		}
	}

	e.RunCommand("pulumi", "policy", "new", "kubernetes-premium-policies-typescript", "--force")
}
