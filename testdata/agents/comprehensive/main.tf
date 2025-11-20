
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Comprehensive test combining multiple agent scenarios
# Used by TestAgentComprehensive in test/resources/agent_test.go

# Test 1: Minimal agent (required fields only)
resource "controlplane_agent" "minimal" {
  name = "test-agent-minimal"
}

# Test 2: Full agent with all optional fields
resource "controlplane_agent" "full" {
  name        = "test-agent-full"
  description = "Comprehensive test agent with all fields configured"
  status      = "idle"
  model_id    = "gpt-4"
  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
    top_p       = 0.9
  })
  runtime      = "default"
  capabilities = ["code_execution", "file_operations", "web_search", "data_analysis"]
  configuration = jsonencode({
    max_retries    = 3
    timeout        = 300
    retry_delay    = 5
    enable_logging = true
    settings = {
      verbose = true
      debug   = false
    }
  })
}

# Test 3: Agent with claude_code runtime
resource "controlplane_agent" "claude_code" {
  name        = "test-agent-claude-code"
  description = "Agent using Claude Code SDK runtime"
  runtime     = "claude_code"
  model_id    = "claude-3-5-sonnet-20241022"
  llm_config = jsonencode({
    temperature = 0.5
  })
  capabilities = ["advanced_reasoning", "code_generation"]
}

# Test 4: Agent with team assignment
resource "controlplane_team" "test_team" {
  name        = "test-team-for-agent"
  description = "Team for agent assignment testing"
}

resource "controlplane_agent" "with_team" {
  name         = "test-agent-with-team"
  description  = "Agent assigned to a team"
  team_id      = controlplane_team.test_team.id
  capabilities = ["team_collaboration"]
}

# Test 5: Agent for update testing
resource "controlplane_agent" "for_update" {
  name        = "test-agent-for-update"
  description = "Initial description"
  runtime     = "default"
  capabilities = ["basic"]
  configuration = jsonencode({
    version = 1
  })
}

# Data source tests
data "controlplane_agent" "minimal_lookup" {
  id = controlplane_agent.minimal.id
}

data "controlplane_agent" "full_lookup" {
  id = controlplane_agent.full.id
}

data "controlplane_agent" "claude_code_lookup" {
  id = controlplane_agent.claude_code.id
}

data "controlplane_agent" "with_team_lookup" {
  id = controlplane_agent.with_team.id
}

# Outputs for test validation
output "minimal_agent_id" {
  value = controlplane_agent.minimal.id
}

output "minimal_agent_name" {
  value = controlplane_agent.minimal.name
}

output "full_agent_id" {
  value = controlplane_agent.full.id
}

output "full_agent_name" {
  value = controlplane_agent.full.name
}

output "full_agent_description" {
  value = controlplane_agent.full.description
}

output "full_agent_runtime" {
  value = controlplane_agent.full.runtime
}

output "full_agent_model_id" {
  value = controlplane_agent.full.model_id
}

output "claude_code_agent_id" {
  value = controlplane_agent.claude_code.id
}

output "claude_code_agent_runtime" {
  value = controlplane_agent.claude_code.runtime
}

output "with_team_agent_id" {
  value = controlplane_agent.with_team.id
}

output "with_team_agent_team_id" {
  value = controlplane_agent.with_team.team_id
}

output "for_update_agent_id" {
  value = controlplane_agent.for_update.id
}

# Data source outputs for validation
output "data_minimal_name" {
  value = data.controlplane_agent.minimal_lookup.name
}

output "data_full_description" {
  value = data.controlplane_agent.full_lookup.description
}

output "data_full_capabilities" {
  value = data.controlplane_agent.full_lookup.capabilities
}

output "data_claude_code_runtime" {
  value = data.controlplane_agent.claude_code_lookup.runtime
}
