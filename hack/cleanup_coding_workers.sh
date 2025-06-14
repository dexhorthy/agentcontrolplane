#!/bin/bash
# cleanup_coding_workers.sh - Cleans up worktree environments and tmux sessions
# Usage: ./cleanup_coding_workers.sh [suffix] [--tmux-only|--worktrees-only]

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to log messages
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" >&2
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARN:${NC} $1"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO:${NC} $1"
}

# Parse arguments  
CLEANUP_MODE="${1:-all}"

# Configuration
REPO_NAME="agentcontrolplane"
WORKTREES_BASE="$HOME/.humanlayer/worktrees"

# Configuration
TMUX_SESSION="acp-agents"
declare -a BRANCH_NAMES=(
    "acp-kind-isolated-claude"
    "acp-e2e-framework-claude" 
    "acp-mcp-transport-claude"
    "acp-kind-isolated-cb"
    "acp-e2e-framework-cb"
    "acp-mcp-transport-cb"
    "acp-merge-claude"
    "acp-merge-cb"
)

# Function to kill tmux session
cleanup_tmux() {
    if tmux has-session -t "$TMUX_SESSION" 2>/dev/null; then
        log "Killing tmux session: $TMUX_SESSION"
        tmux kill-session -t "$TMUX_SESSION"
    else
        info "Tmux session not found: $TMUX_SESSION"
    fi
}

# Function to remove worktree
remove_worktree() {
    local branch_name=$1
    local worktree_dir="${WORKTREES_BASE}/${REPO_NAME}_${branch_name}"
    
    if [ -d "$worktree_dir" ]; then
        log "Removing worktree: $worktree_dir"
        # Fix permissions before removal to handle any permission issues
        chmod -R 755 "$worktree_dir" 2>/dev/null || warn "Failed to fix permissions on $worktree_dir"
        git worktree remove --force "$worktree_dir" 2>/dev/null || {
            warn "Failed to remove worktree with git, removing directory manually"
            rm -rf "$worktree_dir"
        }
    else
        info "Worktree not found: $worktree_dir"
    fi
}

# Function to delete branch
delete_branch() {
    local branch_name=$1
    
    if git show-ref --verify --quiet "refs/heads/${branch_name}"; then
        log "Deleting branch: $branch_name"
        git branch -D "$branch_name" 2>/dev/null || warn "Failed to delete branch: $branch_name"
    else
        info "Branch not found: $branch_name"
    fi
}

# Function to cleanup all worktrees
cleanup_worktrees() {
    for branch_name in "${BRANCH_NAMES[@]}"; do
        remove_worktree "$branch_name"
        delete_branch "$branch_name"
    done
    
    # Prune worktree list
    log "Pruning git worktree list..."
    git worktree prune
}

# Function to show usage
usage() {
    echo "Usage: $0 [--tmux-only|--worktrees-only]"
    echo
    echo "Options:"
    echo "  --tmux-only     - Only clean up tmux session"
    echo "  --worktrees-only - Only clean up worktrees and branches"
    echo
    echo "Examples:"
    echo "  $0                    # Clean up everything (acp-agents session and all branches)"
    echo "  $0 --tmux-only       # Only clean up tmux session"
    echo "  $0 --worktrees-only  # Only clean up worktrees and branches"
}

# Main execution
main() {
    log "Starting cleanup_coding_workers.sh"
    
    if [ "$CLEANUP_MODE" == "--help" ] || [ "$CLEANUP_MODE" == "-h" ]; then
        usage
        exit 0
    fi
    
    # Status report before cleanup
    info "=== Current Status ==="
    echo "Tmux sessions:"
    tmux list-sessions 2>/dev/null | grep "acp-agents" || echo "  None found"
    echo
    echo "Git worktrees:"
    git worktree list | grep -E "acp-.*-(claude|cb)" || echo "  None found"
    echo
    
    # Perform cleanup based on mode
    case "$CLEANUP_MODE" in
        --tmux-only)
            log "Cleaning up tmux sessions only..."
            cleanup_tmux
            ;;
        --worktrees-only)
            log "Cleaning up worktrees only..."
            cleanup_worktrees
            ;;
        all|"")
            log "Cleaning up everything..."
            cleanup_tmux
            cleanup_worktrees
            ;;
        *)
            error "Unknown cleanup mode: $CLEANUP_MODE"
            usage
            exit 1
            ;;
    esac
    
    # Status report after cleanup
    info "=== Status After Cleanup ==="
    echo "Tmux sessions:"
    tmux list-sessions 2>/dev/null | grep "acp-agents" || echo "  None found"
    echo
    echo "Git worktrees:"
    git worktree list | grep -E "acp-.*-(claude|cb)" || echo "  None found"
    echo
    
    log "✅ Cleanup completed successfully!"
}

# Run main
main