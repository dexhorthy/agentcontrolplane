# KUBECONFIG Isolation Test

Test that worker maintains its isolated cluster even when new KIND clusters are created in manager worktree.

## Tasks:
1. Verify worker uses its own ./.kube/config
2. Test isolation from manager cluster changes
3. Confirm proper environment setup