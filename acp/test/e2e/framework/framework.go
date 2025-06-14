/*
Copyright 2025 the Agent Control Plane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"context"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	acp "github.com/humanlayer/agentcontrolplane/acp/api/v1alpha1"
	"github.com/humanlayer/agentcontrolplane/acp/internal/controller/agent"
	"github.com/humanlayer/agentcontrolplane/acp/internal/controller/llm"
	"github.com/humanlayer/agentcontrolplane/acp/internal/controller/mcpserver"
	"github.com/humanlayer/agentcontrolplane/acp/internal/controller/task"
	"github.com/humanlayer/agentcontrolplane/acp/internal/controller/toolcall"
)

// TestFramework provides a unified test environment with all controllers running
type TestFramework struct {
	TestEnv *envtest.Environment
	Config  *rest.Config
	Client  client.Client
	Manager manager.Manager
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewTestFramework creates a new test framework instance
func NewTestFramework() *TestFramework {
	return &TestFramework{}
}

// Start initializes envtest environment and starts all controllers
func (f *TestFramework) Start() error {
	// Create context
	f.ctx, f.cancel = context.WithCancel(context.Background())

	// Setup envtest environment
	f.TestEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	// Find binary directory dynamically
	if binDir := getFirstFoundEnvTestBinaryDir(); binDir != "" {
		f.TestEnv.BinaryAssetsDirectory = binDir
	}

	// Start test environment
	var err error
	f.Config, err = f.TestEnv.Start()
	if err != nil {
		return err
	}

	// Add schemes
	err = acp.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	// Create client
	f.Client, err = client.New(f.Config, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return err
	}

	// Create manager
	f.Manager, err = ctrl.NewManager(f.Config, ctrl.Options{
		Scheme:                 scheme.Scheme,
		HealthProbeBindAddress: "0", // Disable health probes
	})
	if err != nil {
		return err
	}

	// Setup all controllers
	if err = f.setupControllers(); err != nil {
		return err
	}

	// Start manager in background
	go func() {
		if err := f.Manager.Start(f.ctx); err != nil {
			// Log error but don't fail the test immediately
			// Tests should handle controller startup verification
			_ = err // Explicitly ignore error for tests
		}
	}()

	// Wait a moment for controllers to start
	time.Sleep(100 * time.Millisecond)

	return nil
}

// Stop cleans up the test environment
func (f *TestFramework) Stop() error {
	if f.cancel != nil {
		f.cancel()
	}
	if f.TestEnv != nil {
		return f.TestEnv.Stop()
	}
	return nil
}

// GetClient returns the Kubernetes client
func (f *TestFramework) GetClient() client.Client {
	return f.Client
}

// GetContext returns the context
func (f *TestFramework) GetContext() context.Context {
	return f.ctx
}

// setupControllers initializes all required controllers
func (f *TestFramework) setupControllers() error {
	// Setup LLM controller
	if err := (&llm.LLMReconciler{
		Client: f.Manager.GetClient(),
		Scheme: f.Manager.GetScheme(),
	}).SetupWithManager(f.Manager); err != nil {
		return err
	}

	// Setup Agent controller
	if err := (&agent.AgentReconciler{
		Client: f.Manager.GetClient(),
		Scheme: f.Manager.GetScheme(),
	}).SetupWithManager(f.Manager); err != nil {
		return err
	}

	// Setup Task controller
	if err := (&task.TaskReconciler{
		Client: f.Manager.GetClient(),
		Scheme: f.Manager.GetScheme(),
	}).SetupWithManager(f.Manager); err != nil {
		return err
	}

	// Setup MCPServer controller
	if err := (&mcpserver.MCPServerReconciler{
		Client: f.Manager.GetClient(),
		Scheme: f.Manager.GetScheme(),
	}).SetupWithManager(f.Manager); err != nil {
		return err
	}

	// Setup ToolCall controller
	if err := (&toolcall.ToolCallReconciler{
		Client: f.Manager.GetClient(),
		Scheme: f.Manager.GetScheme(),
	}).SetupWithManager(f.Manager); err != nil {
		return err
	}

	return nil
}

// Helper function from existing patterns
func getFirstFoundEnvTestBinaryDir() string {
	// This mirrors the pattern from the existing suite_test.go files
	basePath := filepath.Join("..", "..", "..", "bin", "k8s")
	entries, err := filepath.Glob(filepath.Join(basePath, "*"))
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		return entry
	}
	return ""
}
