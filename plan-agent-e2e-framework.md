# E2E Test Framework Agent Implementation Plan

Adopt the persona from hack/agent-developer.md and follow the Dan Abramov philosophy: DELETE MORE THAN YOU ADD.

## Your Mission
Create a new e2e testing framework using envtest and real controllers instead of shell commands.

## Core Changes Required

### 1. Create framework package structure
- New directory: `acp/test/e2e/framework/`
- Framework setup with envtest and real controllers
- Test helpers and utilities

### 2. Implement TestFramework struct
- Setup envtest environment
- Start all controllers (LLM, Agent, Task, MCPServer, ToolCall)
- Provide k8s client for tests

### 3. Create basic test implementation
- Test LLM → Agent → Task flow
- Use Go assertions instead of shell commands
- Proper namespace isolation per test

### 4. Suite configuration
- BeforeSuite/AfterSuite setup
- Proper timeout and interval configuration
- Controller manager lifecycle

## Implementation Steps

1. **Read existing controller tests completely** (minimum 1500 lines)
2. **Read getting-started.md** to understand manual flow
3. **Create framework directory structure**
4. **Implement TestFramework struct** with envtest
5. **Add all controller setup** to framework
6. **Create basic integration test** (LLM → Agent → Task)
7. **Add suite setup and teardown**
8. **Delete redundant shell-based tests** (minimum 10% reduction)
9. **Run make fmt vet lint test** after each change
10. **Commit every 5-10 minutes** with meaningful messages

## Key Technical Requirements

- Use envtest for isolated API server
- Run real controllers in background
- Go client for resource creation
- Proper Eventually() assertions
- Unique namespace per test
- Real LLM/HumanLayer API calls with env vars

## Framework Structure
```
acp/test/e2e/framework/
├── framework.go            (TestFramework implementation)
├── suite_test.go          (BeforeSuite/AfterSuite)
└── basic_test.go          (First integration test)
```

## Dan Abramov Rules to Follow

- READ ENTIRE FILES before making changes
- DELETE MORE CODE than you add
- Use existing controller test patterns
- Run builds immediately after changes
- Commit frequently (every 5-10 minutes)
- Never create unnecessary files

## Success Criteria

- Framework starts successfully with all controllers
- Basic test (LLM → Agent → Task) passes
- Tests run faster than shell-based equivalent
- Easy to add new test cases
- Clear error messages on failures

REMEMBER: You've read the controller tests, you understand the patterns. Build on existing code, don't reinvent.