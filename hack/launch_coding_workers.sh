#!/bin/bash
# launch_coding_workers.sh - Sets up parallel work environments for executing code
# Usage: ./launch_coding_workers.sh [suffix]

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

# Get suffix argument
SUFFIX="${1:-$(date +%s)}"

log "Using suffix: $SUFFIX"

# Configuration
REPO_NAME="agentcontrolplane"
WORKTREES_BASE="$HOME/.humanlayer/worktrees"
TMUX_SESSION="acp-agents"

# Define plan files and their configurations
declare -a PLAN_FILES=(
    "plan-agent-kind-isolated.md"
    "plan-agent-e2e-framework.md"
    "plan-agent-mcp-transport.md"
)

declare -a CLAUDE_BRANCH_NAMES=(
    "acp-kind-isolated-claude"
    "acp-e2e-framework-claude"
    "acp-mcp-transport-claude"
)

declare -a CB_BRANCH_NAMES=(
    "acp-kind-isolated-cb"
    "acp-e2e-framework-cb"
    "acp-mcp-transport-cb"
)

# Merge agent configuration
MERGE_PLAN="plan-merge-agent.md"
CLAUDE_MERGE_BRANCH="acp-merge-claude"
CB_MERGE_BRANCH="acp-merge-cb"

# Function to create worktree
create_worktree() {
    local branch_name=$1
    local plan_file=$2
    local worktree_dir="${WORKTREES_BASE}/${REPO_NAME}_${branch_name}"
    
    log "Creating worktree for $branch_name..."
    
    # Use create_worktree.sh if available
    if [ -f "hack/create_worktree.sh" ]; then
        ./hack/create_worktree.sh "$branch_name"
    else
        # Fallback to manual creation
        if [ ! -d "$WORKTREES_BASE" ]; then
            mkdir -p "$WORKTREES_BASE"
        fi
        
        if [ -d "$worktree_dir" ]; then
            warn "Worktree already exists: $worktree_dir"
            return 0
        fi
        
        git worktree add -b "$branch_name" "$worktree_dir" HEAD
        
        # Copy .claude directory
        if [ -d ".claude" ]; then
            cp -r .claude "$worktree_dir/"
        fi
    fi
    
    # Copy plan file
    cp "$plan_file" "$worktree_dir/"
    
    # Create prompt.md file
    cat > "$worktree_dir/prompt.md" << EOF
Adopt the persona from hack/agent-developer.md

Your task is to implement the features described in $plan_file

Key requirements:
- Read and understand the plan in $plan_file
- Follow the Dan Abramov methodology
- Commit your changes every 5-10 minutes
- Run tests frequently
- Delete more code than you add
- Keep a 20+ item TODO list

Start by reading the plan file and understanding the task ahead.
EOF
    
    log "Worktree created: $worktree_dir"
}

# Function to launch tmux window for agent
launch_agent_window() {
    local window_num=$1
    local branch_name=$2
    local plan_file=$3
    local launch_claude=${4:-true}
    local agent_type=${5:-"claude"}
    local base_name="$(basename "$plan_file" .md | sed 's/plan-agent-//' | sed 's/plan-merge-agent/merge/')"
    local window_name="${base_name}-${agent_type}"
    local worktree_dir="${WORKTREES_BASE}/${REPO_NAME}_${branch_name}"
    
    log "Launching window $window_num: $window_name (claude: $launch_claude)"
    
    # Create window
    if [ "$window_num" -eq 1 ]; then
        tmux new-session -d -s "$TMUX_SESSION" -n "$window_name" -c "$worktree_dir"
    else
        tmux new-window -t "$TMUX_SESSION:$window_num" -n "$window_name" -c "$worktree_dir"
    fi
    
    if [ "$launch_claude" = "true" ]; then
        # Launch Claude Code directly in the window
        tmux send-keys -t "$TMUX_SESSION:$window_num" "claude \"\$(cat prompt.md)\"" C-m
        # Send newline to accept trust directory prompt
        sleep 1
        tmux send-keys -t "$TMUX_SESSION:$window_num" C-m
    else
        # Just show ready message for manual agent launch
        tmux send-keys -t "$TMUX_SESSION:$window_num" "echo 'Ready for manual agent launch'" C-m
        tmux send-keys -t "$TMUX_SESSION:$window_num" "echo 'Branch: $branch_name'" C-m
        tmux send-keys -t "$TMUX_SESSION:$window_num" "echo 'Plan: $plan_file'" C-m
        tmux send-keys -t "$TMUX_SESSION:$window_num" "echo 'To launch agent: claude \"\$(cat prompt.md)\"'" C-m
    fi
}

# Main execution
main() {
    log "Starting launch_coding_workers.sh with suffix: $SUFFIX"
    
    # Check prerequisites
    if ! command -v tmux &> /dev/null; then
        error "tmux is not installed"
        exit 1
    fi
    
    if ! command -v claude &> /dev/null; then
        error "claude CLI is not installed"
        exit 1
    fi
    
    # Kill existing session if it exists
    if tmux has-session -t "$TMUX_SESSION" 2>/dev/null; then
        warn "Killing existing tmux session: $TMUX_SESSION"
        tmux kill-session -t "$TMUX_SESSION"
    fi
    
    # Create worktrees for all Claude agents
    log "Creating Claude agent worktrees..."
    for i in "${!PLAN_FILES[@]}"; do
        create_worktree "${CLAUDE_BRANCH_NAMES[$i]}" "${PLAN_FILES[$i]}"
    done
    
    # Create worktrees for all CB agents
    log "Creating CB agent worktrees..."
    for i in "${!PLAN_FILES[@]}"; do
        create_worktree "${CB_BRANCH_NAMES[$i]}" "${PLAN_FILES[$i]}"
    done
    
    # Create merge agent worktrees
    log "Creating merge agent worktrees..."
    create_worktree "$CLAUDE_MERGE_BRANCH" "$MERGE_PLAN"
    create_worktree "$CB_MERGE_BRANCH" "$MERGE_PLAN"
    
    # Create Claude merge agent prompt
    local claude_merge_worktree="${WORKTREES_BASE}/${REPO_NAME}_${CLAUDE_MERGE_BRANCH}"
    cat > "$claude_merge_worktree/prompt.md" << EOF
Adopt the persona from hack/agent-merger.md

Your task is to merge the work from the following Claude agent branches into the current branch:
${CLAUDE_BRANCH_NAMES[@]}

Key requirements:
- Read the plan in $MERGE_PLAN
- Monitor agent branches for commits every 2 minutes
- Merge changes in dependency order
- Resolve conflicts appropriately
- Maintain clean build state
- Commit merged changes

Start by reading the merge plan and checking the status of all agent branches.
EOF

    # Create CB merge agent prompt
    local cb_merge_worktree="${WORKTREES_BASE}/${REPO_NAME}_${CB_MERGE_BRANCH}"
    cat > "$cb_merge_worktree/prompt.md" << EOF
Adopt the persona from hack/agent-merger.md

Your task is to merge the work from the following CB agent branches into the current branch:
${CB_BRANCH_NAMES[@]}

Key requirements:
- Read the plan in $MERGE_PLAN
- Monitor agent branches for commits every 2 minutes
- Merge changes in dependency order
- Resolve conflicts appropriately
- Maintain clean build state
- Commit merged changes

Start by reading the merge plan and checking the status of all agent branches.
EOF
    
    # Launch agent windows
    log "Launching tmux session: $TMUX_SESSION"
    local window_num=1
    
    # Launch Claude agents
    for i in "${!PLAN_FILES[@]}"; do
        launch_agent_window "$window_num" "${CLAUDE_BRANCH_NAMES[$i]}" "${PLAN_FILES[$i]}" "true" "claude"
        ((window_num++))
    done
    
    # Launch CB agents
    for i in "${!PLAN_FILES[@]}"; do
        launch_agent_window "$window_num" "${CB_BRANCH_NAMES[$i]}" "${PLAN_FILES[$i]}" "false" "cb"
        ((window_num++))
    done
    
    # Launch merge agents
    launch_agent_window "$window_num" "$CLAUDE_MERGE_BRANCH" "$MERGE_PLAN" "true" "claude-merge"
    ((window_num++))
    launch_agent_window "$window_num" "$CB_MERGE_BRANCH" "$MERGE_PLAN" "true" "cb-merge"
    
    # Summary
    log "✅ All coding workers launched successfully!"
    echo
    echo "Session: $TMUX_SESSION"
    echo "Claude agents (windows 1-3):"
    for i in "${!PLAN_FILES[@]}"; do
        local task_name=$(basename "${PLAN_FILES[$i]}" .md | sed 's/plan-agent-//')
        echo "  - Window $((i+1)): ${task_name}-claude (${CLAUDE_BRANCH_NAMES[$i]})"
    done
    echo "CB agents (windows 4-6):"
    for i in "${!PLAN_FILES[@]}"; do
        local task_name=$(basename "${PLAN_FILES[$i]}" .md | sed 's/plan-agent-//')
        echo "  - Window $((i+4)): ${task_name}-cb (${CB_BRANCH_NAMES[$i]})"
    done
    echo "Merge agents:"
    echo "  - Window 7: merge-claude ($CLAUDE_MERGE_BRANCH)"
    echo "  - Window 8: merge-cb ($CB_MERGE_BRANCH)"
    echo
    echo "To attach to the session:"
    echo "  tmux attach -t $TMUX_SESSION"
    echo
    echo "To switch between windows:"
    echo "  Ctrl-b [window-number]"
    echo
    echo "To clean up later:"
    echo "  ./cleanup_coding_workers.sh"
}

# Run main
main