# Merge Agent Plan

Adopt the persona from hack/agent-merger.md and follow the Dan Abramov philosophy: DELETE MORE THAN YOU ADD.

## Your Mission
Monitor the three development agents and merge their commits incrementally into the manager branch.

## Branches to Monitor
- `acp-kind-isolated-dev` - Kind cluster isolation implementation
- `acp-e2e-framework-dev` - E2E test framework implementation  
- `acp-mcp-transport-dev` - MCP remote transport implementation

## Your Workflow

### 1. Monitor Agent Progress
Check each branch every 2 minutes for new commits:
```bash
git log --oneline -3 acp-kind-isolated-dev
git log --oneline -3 acp-e2e-framework-dev
git log --oneline -3 acp-mcp-transport-dev
```

### 2. Merge Strategy
- Use fast-forward merges when possible
- Merge commits as they become available
- Don't wait for complete features - merge incremental progress
- Handle merge conflicts immediately

### 3. Validation After Each Merge
- Run `make -C acp fmt vet lint test` after every merge
- Deploy controller with `make -C acp deploy-local-kind`
- Test basic functionality with simple resource creation
- Fix any issues before continuing

### 4. File Change Monitoring
Watch for changes in these key areas:
- `Makefile` and `acp/Makefile` (kind-isolated agent)
- `acp/test/e2e/framework/` (e2e-framework agent)
- `acp/api/v1alpha1/mcpserver_types.go` (mcp-transport agent)
- `acp/internal/mcpmanager/` (mcp-transport agent)
- `hack/` scripts (kind-isolated agent)

### 5. Integration Testing
After merging commits from all agents:
- Verify all three features work together
- Test isolated clusters + new e2e framework + remote MCP
- Document any integration issues

## Dan Abramov Rules to Follow

- READ ENTIRE CHANGED FILES before merging (minimum 1500 lines)
- DELETE redundant code during merges
- Build and test after every merge
- Commit merge results frequently
- Don't create unnecessary files during merges

## Merge Commit Messages
Use clear, descriptive messages:
- "Merge: Add kind cluster isolation from acp-kind-isolated-dev"
- "Merge: E2E framework foundation from acp-e2e-framework-dev"
- "Merge: MCP transport updates from acp-mcp-transport-dev"

## Conflict Resolution
- READ both versions completely
- Choose the better implementation (usually the newer one)
- DELETE the redundant code
- Test the resolution immediately

## Success Criteria
- All agent changes merged successfully
- No merge conflicts remain unresolved
- All tests pass after merges
- Controller deploys and functions correctly
- Integration between features verified

REMEMBER: You monitor, merge, and validate. Keep the main branch stable while incorporating all agent progress.