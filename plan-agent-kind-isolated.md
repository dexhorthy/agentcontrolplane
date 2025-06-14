# Kind-Isolated Agent Implementation Plan

Adopt the persona from hack/agent-developer.md and follow the Dan Abramov philosophy: DELETE MORE THAN YOU ADD.

## Your Mission
Implement isolated kind clusters for each developer agent to prevent conflicts and enable true parallel development.

## Core Changes Required

### 1. Update Makefile setup target
- Modify the `setup` target to create unique clusters per branch
- Add dynamic port allocation to prevent conflicts
- Export kubeconfig to worktree-local location

### 2. Create find_free_port.sh script
- Implement port scanning function in `hack/find_free_port.sh`
- Must find available ports in range 10000-11000
- Used by setup target for dynamic port allocation

### 3. Update acp/Makefile
- Modify to use local kubeconfig when available
- Add cluster name detection from worktree
- Update deploy-local-kind target with cluster validation

### 4. Update cleanup process
- Modify `cleanup_coding_workers.sh` to include cluster cleanup
- Add cluster deletion function
- Make cleanup idempotent

## Implementation Steps

1. **Read existing Makefiles completely** (minimum 1500 lines total)
2. **Create hack/find_free_port.sh** with port scanning logic
3. **Update root Makefile setup target** with cluster creation
4. **Update acp/Makefile** for local kubeconfig support
5. **Update cleanup script** with cluster deletion
6. **Test with single worktree** to verify isolation
7. **Delete redundant code** (minimum 10% reduction)
8. **Run make fmt vet lint test** after each change
9. **Commit every 5-10 minutes** with meaningful messages

## Key Technical Requirements

- Each cluster named: `acp-$(git branch --show-current)`
- Kubeconfig stored in `.kube/config` within worktree
- Dynamic port allocation to prevent conflicts
- Resource limits added to manager pod
- Proper error handling for missing clusters

## Dan Abramov Rules to Follow

- READ ENTIRE FILES before making changes
- DELETE MORE CODE than you add
- Use existing patterns, don't invent new ones
- Run builds immediately after changes
- Commit frequently (every 5-10 minutes)
- Never create unnecessary files

## Success Criteria

- Each agent can deploy without affecting others
- Clusters are automatically created and cleaned up
- Resource usage is reasonable
- Existing workflows continue to work
- Clear error messages when resources are insufficient

REMEMBER: You've read the full files, you understand everything. Trust your knowledge and execute with confidence.