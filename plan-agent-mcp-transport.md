# MCP Remote Transport Agent Implementation Plan

Adopt the persona from hack/agent-developer.md and follow the Dan Abramov philosophy: DELETE MORE THAN YOU ADD.

## Your Mission
Add support for remote MCP servers using SSE and Streamable HTTP transports.

## Core Changes Required

### 1. Update MCPServer CRD types
- Add new transport types: "sse", "streamable-http" 
- Add optional Headers field for authentication
- Add SessionID field for streamable-http resumption
- Update validation enum

### 2. Update MCPManager transport logic
- Modify ConnectServer method for new transports
- Check mcp-go library for streamable-http support
- Maintain backward compatibility ("http" → SSE)
- Add proper error handling

### 3. Add security and configuration
- HTTPS enforcement for production
- Authentication support (Bearer tokens, API keys)
- Timeout configuration
- Connection health checks

### 4. Update controller logic
- Connection health monitoring
- Network interruption handling
- Status updates with connection info
- Session management

## Implementation Steps

1. **Read mcpmanager.go completely** (minimum 1500 lines)
2. **Read mcpserver_types.go completely**
3. **Check go.mod for mcp-go version** and capabilities
4. **Update CRD schema** with new transport types
5. **Implement transport logic** in mcpmanager.go
6. **Add security validation** and headers support
7. **Create example configurations** in config/samples/
8. **Update documentation** in docs/
9. **Delete redundant code** (minimum 10% reduction)
10. **Run make fmt vet lint test** after each change
11. **Commit every 5-10 minutes** with meaningful messages

## Key Technical Requirements

- Support both SSE (legacy) and streamable-http (new)
- Maintain backward compatibility
- Use secretKeyRef for sensitive headers
- Proper session management for streamable-http
- Graceful reconnection handling

## Transport Mapping
- "http" and "sse" → SSE client (backward compatibility)
- "streamable-http" → New streamable HTTP client
- "stdio" → Existing stdio client (unchanged)

## Dan Abramov Rules to Follow

- READ ENTIRE FILES before making changes
- DELETE MORE CODE than you add
- Use existing patterns from current transports
- Run builds immediately after changes
- Commit frequently (every 5-10 minutes)
- Never create unnecessary files

## Success Criteria

- Can connect to remote MCP servers using SSE transport
- Can connect using streamable-http transport (if library supports)
- Examples added to getting-started.md
- Graceful handling of connection failures
- All existing transports continue working

REMEMBER: You've read the mcpmanager code, you understand the patterns. Extend existing functionality, don't rewrite.