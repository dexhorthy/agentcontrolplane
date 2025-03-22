# Testing Practices

This document outlines our testing practices and patterns, with a focus on readability, maintainability, and reliability.

## Core Testing Principles

1. **Test Helper Structs**: Encapsulate test setup and teardown in well-named structs
2. **Clear Test Steps**: Use descriptive "By" blocks to document test progression
3. **Consistent Cleanup**: Always clean up resources before and after tests
4. **Descriptive Names**: Use clear, intention-revealing names for tests and variables

## Test Structure

### Helper Structs Pattern

We use helper structs to manage test resources. Each struct handles its own lifecycle:

```go
// TestTool manages the setup and teardown of test tools
type TestTool struct {
    Name      string
    Namespace string
    Tool      *v1alpha1.Tool
    Ctx       context.Context
}

func (t *TestTool) Setup() {
    // Create and configure the resource
}

func (t *TestTool) TearDown() {
    // Clean up the resource
}
```

Benefits:
- Encapsulates setup/teardown logic
- Makes tests more readable
- Reduces code duplication
- Makes resource management explicit

### Test Steps Documentation

We use Ginkgo's "By" blocks to document test progression:

```go
It("should successfully execute a function", func() {
    By("Creating the resource")
    testResource.Setup()

    By("Reconciling the resource")
    result, err := reconciler.Reconcile(ctx, req)
    
    By("Verifying the resource status")
    Expect(err).NotTo(HaveOccurred())
})
```

Benefits:
- Makes test flow clear and readable
- Provides better error context
- Documents test intentions
- Helps with debugging failures

### BeforeEach/AfterEach Pattern

```go
BeforeEach(func() {
    // Initialize context
    ctx, cancel = context.WithCancel(context.TODO())
    
    // Clean up any existing resources
    By("Cleaning up existing resources")
    testResource.TearDown()
    
    // Set up fresh test resources
    By("Setting up test resources")
    testResource.Setup()
})

AfterEach(func() {
    // Clean up
    cancel()
    testResource.TearDown()
})
```

Benefits:
- Ensures clean test environment
- Prevents resource conflicts
- Makes tests independent
- Handles cleanup reliably

## Testing Patterns

### Controller Tests

For controller tests, we follow these patterns:

1. **Resource Setup**:
```go
testTool := &TestTool{
    Name:      "test-tool",
    Namespace: "default",
    Ctx:       ctx,
}
```

2. **Reconciler Setup**:
```go
reconciler := &ToolReconciler{
    Client:   k8sClient,
    Scheme:   scheme,
    Recorder: eventRecorder,
}
```

3. **Status Verification**:
```go
By("Verifying the resource status")
updatedResource := &v1alpha1.Tool{}
err := k8sClient.Get(ctx, key, updatedResource)
Expect(err).NotTo(HaveOccurred())
Expect(updatedResource.Status.Ready).To(BeTrue())
```

4. **Event Verification**:
```go
By("Verifying events were emitted")
testutils.ExpectEvent(eventRecorder).ToEmitEventContaining("Expected Event")
```

### Error Case Testing

For error cases, we verify both the error and the resource status:

```go
It("should handle invalid input", func() {
    By("Creating resource with invalid config")
    testResource.SetupWithInvalidConfig()

    By("Verifying error is returned")
    _, err := reconciler.Reconcile(ctx, req)
    Expect(err).To(HaveOccurred())

    By("Verifying error status is set")
    updatedResource := &v1alpha1.Resource{}
    Expect(k8sClient.Get(ctx, key, updatedResource)).To(Succeed())
    Expect(updatedResource.Status.Status).To(Equal("Error"))
})
```

## Running Tests

### Unit Tests
```bash
make test
```

### End-to-End Tests
```bash
make test-e2e
```

### Coverage
```bash
make test-coverage
```

## Best Practices

1. **Independent Tests**: Each test should be independent and not rely on state from other tests

2. **Clear Assertions**: Use descriptive matchers:
```go
Expect(resource.Status.Ready).To(BeTrue())
Expect(resource.Status.Error).To(BeEmpty())
```

3. **Resource Cleanup**: Always clean up resources, even if tests fail:
```go
AfterEach(func() {
    By("Cleaning up resources")
    testResource.TearDown()
})
```

4. **Context Management**: Always handle context properly:
```go
var ctx context.Context
var cancel context.CancelFunc

BeforeEach(func() {
    ctx, cancel = context.WithCancel(context.TODO())
})

AfterEach(func() {
    cancel()
})
```

5. **Descriptive Test Names**: Use clear, behavior-focused test names:
```go
It("should successfully reconcile when all dependencies are ready", func() {
    // test code
})
```

6. **Event Verification**: Use the ExpectEvent helper for consistent event checking:
```go
// Import the helper
testutils "github.com/humanlayer/smallchain/kubechain/test/utils"

// Use it in tests
testutils.ExpectEvent(eventRecorder).ToEmitEventContaining("ValidationSucceeded")
```

## Why This Approach Works

1. **Maintainability**: Helper structs and clear test steps make tests easy to maintain
2. **Reliability**: Consistent cleanup prevents test interference
3. **Readability**: "By" blocks and clear structure make tests self-documenting
4. **Debuggability**: Clear test steps make it easy to identify failure points
5. **Reusability**: Helper structs and utilities can be reused across test suites

Remember: Tests are documentation. Write them with the same care as production code.