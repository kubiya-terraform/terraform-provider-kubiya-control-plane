
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create environment for worker queue testing
resource "controlplane_environment" "test_env" {
  name        = "test-env-for-worker-queues"
  description = "Environment for worker queue testing"
}

# Test 1: Minimal worker queue (required fields only)
resource "controlplane_worker_queue" "minimal" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-worker-queue-minimal"
}

# Test 2: Full worker queue with all optional fields
resource "controlplane_worker_queue" "full" {
  environment_id    = controlplane_environment.test_env.id
  name              = "test-worker-queue-full"
  display_name      = "Test Worker Queue Full"
  description       = "Comprehensive test worker queue with all fields configured"
  status            = "active"
  max_workers       = 10
  heartbeat_interval = 60

  tags = ["production", "critical", "monitored"]

  settings = {
    priority          = "high"
    auto_scale        = "true"
    min_workers       = "2"
    scale_up_threshold = "80"
    scale_down_threshold = "20"
  }
}

# Test 3: Worker queue with high worker limit
resource "controlplane_worker_queue" "high_capacity" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-high-capacity"
  display_name   = "High Capacity Queue"
  description    = "Worker queue with high worker limit"
  max_workers    = 100
}

# Test 4: Worker queue with unlimited workers
resource "controlplane_worker_queue" "unlimited" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-unlimited"
  description    = "Worker queue with unlimited workers (null max_workers)"
  # max_workers is null by default = unlimited
}

# Test 5: Worker queue with short heartbeat interval
resource "controlplane_worker_queue" "short_heartbeat" {
  environment_id     = controlplane_environment.test_env.id
  name               = "test-wq-short-heartbeat"
  description        = "Worker queue with short heartbeat interval"
  heartbeat_interval = 10  # Minimum: 10 seconds
}

# Test 6: Worker queue with long heartbeat interval
resource "controlplane_worker_queue" "long_heartbeat" {
  environment_id     = controlplane_environment.test_env.id
  name               = "test-wq-long-heartbeat"
  description        = "Worker queue with long heartbeat interval"
  heartbeat_interval = 300  # Maximum: 300 seconds (5 minutes)
}

# Test 7: Worker queue with inactive status
resource "controlplane_worker_queue" "inactive" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-inactive"
  description    = "Worker queue with inactive status"
  status         = "inactive"
}

# Test 8: Worker queue with paused status
resource "controlplane_worker_queue" "paused" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-paused"
  description    = "Worker queue with paused status"
  status         = "paused"
}

# Test 9: Worker queue with tags only
resource "controlplane_worker_queue" "with_tags" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-with-tags"
  description    = "Worker queue with tags"

  tags = ["test", "automated", "low-priority"]
}

# Test 10: Worker queue with settings only
resource "controlplane_worker_queue" "with_settings" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-with-settings"
  description    = "Worker queue with custom settings"

  settings = {
    region            = "us-west-2"
    availability_zone = "us-west-2a"
    instance_type     = "t3.medium"
    enable_spot       = "true"
  }
}

# Test 11: Worker queue with complex settings
resource "controlplane_worker_queue" "complex_settings" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-complex-settings"
  description    = "Worker queue with complex settings"

  settings = {
    # Autoscaling configuration
    auto_scale_enabled    = "true"
    min_workers           = "1"
    max_workers_limit     = "50"
    scale_up_threshold    = "85"
    scale_down_threshold  = "15"
    scale_up_cooldown     = "300"
    scale_down_cooldown   = "600"

    # Resource configuration
    cpu_limit             = "2000m"
    memory_limit          = "4Gi"
    storage_limit         = "20Gi"

    # Network configuration
    network_policy        = "default"
    egress_enabled        = "true"

    # Monitoring
    metrics_enabled       = "true"
    log_level             = "info"
  }
}

# Test 12: Worker queue with empty optional fields
resource "controlplane_worker_queue" "empty_optionals" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-empty-optionals"
  display_name   = ""
  description    = ""
  tags           = []
  settings       = {}
}

# Test 13: Worker queue for update testing
resource "controlplane_worker_queue" "for_update" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-wq-for-update"
  description    = "Initial description"
  status         = "active"
  max_workers    = 5

  settings = {
    version = "1"
  }
}

# Create second environment for testing multiple environments
resource "controlplane_environment" "test_env2" {
  name        = "test-env-2-for-worker-queues"
  description = "Second environment for worker queue testing"
}

# Test 14: Worker queue in different environment
resource "controlplane_worker_queue" "different_env" {
  environment_id = controlplane_environment.test_env2.id
  name           = "test-wq-different-env"
  description    = "Worker queue in a different environment"
}

# Data source tests
data "controlplane_worker_queue" "minimal_lookup" {
  id = controlplane_worker_queue.minimal.id
}

data "controlplane_worker_queue" "full_lookup" {
  id = controlplane_worker_queue.full.id
}

data "controlplane_worker_queue" "inactive_lookup" {
  id = controlplane_worker_queue.inactive.id
}

# Test worker_queues data source (list all queues in an environment)
data "controlplane_worker_queues" "test_env_queues" {
  environment_id = controlplane_environment.test_env.id
}

# Outputs for test validation
output "minimal_worker_queue_id" {
  value = controlplane_worker_queue.minimal.id
}

output "minimal_worker_queue_name" {
  value = controlplane_worker_queue.minimal.name
}

output "minimal_worker_queue_environment_id" {
  value = controlplane_worker_queue.minimal.environment_id
}

output "minimal_worker_queue_created_at" {
  value = controlplane_worker_queue.minimal.created_at
}

output "full_worker_queue_id" {
  value = controlplane_worker_queue.full.id
}

output "full_worker_queue_name" {
  value = controlplane_worker_queue.full.name
}

output "full_worker_queue_display_name" {
  value = controlplane_worker_queue.full.display_name
}

output "full_worker_queue_description" {
  value = controlplane_worker_queue.full.description
}

output "full_worker_queue_status" {
  value = controlplane_worker_queue.full.status
}

output "full_worker_queue_max_workers" {
  value = controlplane_worker_queue.full.max_workers
}

output "full_worker_queue_heartbeat_interval" {
  value = controlplane_worker_queue.full.heartbeat_interval
}

output "full_worker_queue_tags" {
  value = controlplane_worker_queue.full.tags
}

output "full_worker_queue_settings" {
  value = controlplane_worker_queue.full.settings
}

output "full_worker_queue_active_workers" {
  value = controlplane_worker_queue.full.active_workers
}

output "full_worker_queue_task_queue_name" {
  value = controlplane_worker_queue.full.task_queue_name
}

output "full_worker_queue_created_at" {
  value = controlplane_worker_queue.full.created_at
}

output "full_worker_queue_updated_at" {
  value = controlplane_worker_queue.full.updated_at
}

output "high_capacity_worker_queue_id" {
  value = controlplane_worker_queue.high_capacity.id
}

output "high_capacity_worker_queue_max_workers" {
  value = controlplane_worker_queue.high_capacity.max_workers
}

output "unlimited_worker_queue_id" {
  value = controlplane_worker_queue.unlimited.id
}

output "unlimited_worker_queue_max_workers" {
  value = controlplane_worker_queue.unlimited.max_workers
}

output "short_heartbeat_worker_queue_id" {
  value = controlplane_worker_queue.short_heartbeat.id
}

output "short_heartbeat_worker_queue_heartbeat_interval" {
  value = controlplane_worker_queue.short_heartbeat.heartbeat_interval
}

output "long_heartbeat_worker_queue_id" {
  value = controlplane_worker_queue.long_heartbeat.id
}

output "long_heartbeat_worker_queue_heartbeat_interval" {
  value = controlplane_worker_queue.long_heartbeat.heartbeat_interval
}

output "inactive_worker_queue_id" {
  value = controlplane_worker_queue.inactive.id
}

output "inactive_worker_queue_status" {
  value = controlplane_worker_queue.inactive.status
}

output "paused_worker_queue_id" {
  value = controlplane_worker_queue.paused.id
}

output "paused_worker_queue_status" {
  value = controlplane_worker_queue.paused.status
}

output "with_tags_worker_queue_id" {
  value = controlplane_worker_queue.with_tags.id
}

output "with_tags_worker_queue_tags" {
  value = controlplane_worker_queue.with_tags.tags
}

output "with_settings_worker_queue_id" {
  value = controlplane_worker_queue.with_settings.id
}

output "with_settings_worker_queue_settings" {
  value = controlplane_worker_queue.with_settings.settings
}

output "complex_settings_worker_queue_id" {
  value = controlplane_worker_queue.complex_settings.id
}

output "complex_settings_worker_queue_settings" {
  value = controlplane_worker_queue.complex_settings.settings
}

output "for_update_worker_queue_id" {
  value = controlplane_worker_queue.for_update.id
}

output "different_env_worker_queue_id" {
  value = controlplane_worker_queue.different_env.id
}

output "different_env_worker_queue_environment_id" {
  value = controlplane_worker_queue.different_env.environment_id
}

# Data source outputs for validation
output "data_minimal_name" {
  value = data.controlplane_worker_queue.minimal_lookup.name
}

output "data_minimal_environment_id" {
  value = data.controlplane_worker_queue.minimal_lookup.environment_id
}

output "data_full_description" {
  value = data.controlplane_worker_queue.full_lookup.description
}

output "data_full_max_workers" {
  value = data.controlplane_worker_queue.full_lookup.max_workers
}

output "data_full_tags" {
  value = data.controlplane_worker_queue.full_lookup.tags
}

output "data_inactive_status" {
  value = data.controlplane_worker_queue.inactive_lookup.status
}

# Worker queues data source outputs
output "test_env_queues_count" {
  value = length(data.controlplane_worker_queues.test_env_queues.queues)
}

output "test_env_queues_list" {
  value = [for q in data.controlplane_worker_queues.test_env_queues.queues : q.name]
}
