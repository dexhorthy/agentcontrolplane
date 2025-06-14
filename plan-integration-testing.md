# Integration Testing Plan

Adopt the persona from hack/agent-integration-tester.md and thoroughly test all merged features.

## Your Mission
Test the complete integration of all merged features:
- Kind cluster isolation 
- E2E testing framework
- MCP transport enhancements
- End-to-end controller functionality

## Testing Strategy

### 1. Prerequisites Setup
- Verify isolated kind cluster is running
- Check API keys are set: `make check-keys-set`
- Verify controller deployment: `make deploy-local-kind`

### 2. Core Integration Tests
- **Basic Flow**: LLM → Agent → Task workflow
- **MCP Integration**: Test all transport types (stdio, http, sse, streamable-http)
- **E2E Framework**: Run `make test-e2e-framework` 
- **Isolation**: Verify cluster independence

### 3. Advanced Features
- **Human Approval**: Test MCP tools requiring approval
- **Sub-Agent Delegation**: Test agent-to-agent communication
- **Contact Channels**: Test human-as-tool functionality
- **Multiple Models**: Test OpenAI, Anthropic, etc.

### 4. Performance & Reliability
- **Concurrent Tasks**: Multiple tasks running simultaneously
- **Error Handling**: Network failures, timeouts, invalid configs
- **Resource Cleanup**: Proper cleanup after tests
- **Controller Logs**: Verify no error logs during normal operation

## Implementation Steps

1. **Deploy and verify isolated cluster**
2. **Run acp/docs/getting-started.md completely**
3. **Test all MCP transport types with real servers**
4. **Run E2E framework tests**
5. **Test human approval workflows**
6. **Stress test with multiple concurrent tasks**
7. **Document any issues in integration-test-issues.md**
8. **Clean up all test resources**

## Success Criteria

- ✅ All getting-started.md steps complete successfully
- ✅ All E2E framework tests pass
- ✅ All MCP transport types functional
- ✅ No error logs in controller during normal operation
- ✅ Isolated cluster works independently
- ✅ Human approval flows work end-to-end
- ✅ Resource cleanup successful

## Dan Abramov Rules to Follow

- **DELETE MORE THAN YOU ADD**: Remove any redundant test files
- **COMMIT EVERY 5-10 MINUTES**: Document progress frequently
- **READ FULLY**: Understand the complete getting-started guide first
- **BUILD AND TEST**: Run make commands after each verification
- **MAINTAIN 20+ TODO LIST**: Track all test scenarios

## Key Requirements

- Use isolated cluster (not shared/upstream)
- Document issues clearly with reproduction steps
- Test real MCP servers, not just mocks
- Verify resource cleanup between tests
- Check controller logs for any errors
- Test edge cases and error conditions