# OpenAPI Specification Changes Summary

**Comparison between:**
- Newer Version: `/home/dima/Projects/work/kubiya/agent-control-plane/openapi.last.json` (2025-11-16)
- Older Version: `/home/dima/Projects/work/kubiya/agent-control-plane/openapi.new.json` (2025-11-11)

**Summary Statistics:**
- Total endpoints in OLD: 127
- Total endpoints in NEW: 141
- Total schemas in OLD: 115
- Total schemas in NEW: 116
- New endpoints added: 14
- Removed endpoints: 0
- Modified endpoints: 1
- New schemas: 1
- Modified schemas: 5

---

## 1. NEW ENDPOINTS ADDED (14 endpoints)

### Context Graph API (12 endpoints)

A complete new Context Graph API has been added with comprehensive graph database capabilities.

#### 1.1 GET `/api/v1/context-graph/api/v1/graph/integrations`
- **Summary:** Get Integrations
- **Description:** Get all available integrations for the organization
- **Query Parameters:**
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for reading integrations

#### 1.2 GET `/api/v1/context-graph/api/v1/graph/labels`
- **Summary:** Get Labels
- **Description:** Get all node labels in the context graph
- **Query Parameters:**
  - `integration` (optional)
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for reading graph node labels

#### 1.3 GET `/api/v1/context-graph/api/v1/graph/nodes`
- **Summary:** Get All Nodes
- **Description:** Get all nodes in the organization. Optionally filter by integration label
- **Query Parameters:**
  - `integration` (optional)
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for reading all graph nodes

#### 1.4 GET `/api/v1/context-graph/api/v1/graph/nodes/{node_id}`
- **Summary:** Get Node
- **Description:** Get a specific node by its ID
- **Path Parameters:**
  - `node_id` (string, required)
- **Query Parameters:**
  - `integration` (optional)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for reading a specific graph node

#### 1.5 GET `/api/v1/context-graph/api/v1/graph/nodes/{node_id}/relationships`
- **Summary:** Get Relationships
- **Description:** Get relationships for a specific node
- **Path Parameters:**
  - `node_id` (string, required)
- **Query Parameters:**
  - `direction` (string, optional)
  - `relationship_type` (optional)
  - `integration` (optional)
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for reading node relationships

#### 1.6 POST `/api/v1/context-graph/api/v1/graph/nodes/search`
- **Summary:** Search Nodes
- **Description:** Search for nodes in the context graph
- **Request Body:** NodeSearchRequest
  - `label` (optional): Node label to filter by
  - `property_name` (optional): Property name to filter by
  - `property_value` (optional): Property value to match
- **Query Parameters:**
  - `integration` (optional)
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for searching nodes

#### 1.7 POST `/api/v1/context-graph/api/v1/graph/nodes/search/text`
- **Summary:** Search By Text
- **Description:** Search nodes by text pattern in a property
- **Request Body:** TextSearchRequest
  - `property_name`: Property name to search in
  - `search_text`: Text to search for
  - `label` (optional): Node label to filter by
- **Query Parameters:**
  - `integration` (optional)
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for text search

#### 1.8 POST `/api/v1/context-graph/api/v1/graph/query`
- **Summary:** Execute Query
- **Description:** Execute a custom Cypher query (read-only). Query is automatically scoped to organization's data
- **Request Body:** CustomQueryRequest
  - `query`: Cypher query to execute (read-only)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for executing custom queries

#### 1.9 GET `/api/v1/context-graph/api/v1/graph/relationship-types`
- **Summary:** Get Relationship Types
- **Description:** Get all relationship types in the context graph
- **Query Parameters:**
  - `integration` (optional)
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for reading relationship types

#### 1.10 GET `/api/v1/context-graph/api/v1/graph/stats`
- **Summary:** Get Stats
- **Description:** Get statistics about the context graph
- **Query Parameters:**
  - `integration` (optional)
  - `skip` (integer, optional, default: 0)
  - `limit` (integer, optional, default: 100)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for graph statistics

#### 1.11 POST `/api/v1/context-graph/api/v1/graph/subgraph`
- **Summary:** Get Subgraph
- **Description:** Get a subgraph starting from a node
- **Request Body:** SubgraphRequest
  - `node_id`: Starting node ID
  - `depth`: Traversal depth (1-5)
  - `relationship_types` (optional): List of relationship types to follow
- **Query Parameters:**
  - `integration` (optional)
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New data source for retrieving subgraphs

#### 1.12 GET `/api/v1/context-graph/health`
- **Summary:** Health Check
- **Description:** Health check endpoint for Context Graph API
- **Responses:**
  - 200: Successful Response
- **Terraform Impact:** Can be used for health checking

### Execution Management (1 endpoint)

#### 1.13 POST `/api/v1/executions/{execution_id}/cancel_workflow`
- **Summary:** Cancel Workflow
- **Description:** Cancel a specific workflow tool call within an execution. This cancels only the workflow, allowing the agent to continue running (different from /cancel which stops entire execution)
- **Path Parameters:**
  - `execution_id` (string, required): The agent execution ID
- **Request Body:** CancelWorkflowRequest (REQUIRED)
  - `workflow_message_id` (string, required): The unique workflow message ID to cancel
- **Responses:**
  - 200: Successful Response
  - 422: Validation Error
- **Terraform Impact:** New action/function for canceling workflows within executions

### WebSocket Stats (1 endpoint)

#### 1.14 GET `/api/v1/ws/stats`
- **Summary:** Websocket Stats
- **Description:** Get WebSocket connection statistics. Returns statistics about active connections and message throughput
- **Responses:**
  - 200: Successful Response
- **Terraform Impact:** New data source for WebSocket statistics

---

## 2. REMOVED ENDPOINTS

**None** - No endpoints were removed between versions.

---

## 3. MODIFIED ENDPOINTS (1 endpoint)

### 3.1 GET `/api/v1/analytics/aem/summary`
- **Change Type:** Description update
- **What Changed:** The description text was updated to reflect provider-agnostic model tier classification
- **OLD Description:**
  ```
  Get Agentic Engineering Minutes (AEM) summary.

  Returns:
  - Total AEM consumed
  - Total AEM cost
  - Breakdown by model tier (Opus, Sonnet, Haiku)
  - Average runtime, model weight, tool complexity
  ```
- **NEW Description:**
  ```
  Get Agentic Engineering Minutes (AEM) summary.

  Returns:
  - Total AEM consumed
  - Total AEM cost
  - Breakdown by model tier (Premium, Mid, Basic) - provider-agnostic classification
  - Average runtime, model weight, tool complexity
  ```
- **Terraform Impact:** Documentation should be updated to reflect the new model tier naming (Premium, Mid, Basic instead of Opus, Sonnet, Haiku). The actual response structure may have changed to use these new tier names.

---

## 4. NEW SCHEMAS (1 schema)

### 4.1 CancelWorkflowRequest
```json
{
  "properties": {
    "workflow_message_id": {
      "type": "string",
      "title": "Workflow Message Id"
    }
  },
  "type": "object",
  "required": [
    "workflow_message_id"
  ],
  "title": "CancelWorkflowRequest",
  "description": "Request to cancel a specific workflow"
}
```
- **Usage:** Request body for POST `/api/v1/executions/{execution_id}/cancel_workflow`
- **Terraform Impact:** New schema for canceling workflows

---

## 5. MODIFIED SCHEMAS (5 schemas)

### 5.1 TeamCreate
- **Change Type:** New field added
- **New Field:** `runtime`
  - **Type:** string (nullable)
  - **Description:** Runtime type for team leader: 'default' (Agno) or 'claude_code' (Claude Code SDK)
  - **Default:** "default"
  - **Required:** No
- **Terraform Impact:**
  - Add optional `runtime` attribute to team resource
  - Valid values: "default", "claude_code"
  - Default value: "default"

### 5.2 TeamResponse
- **Change Type:** New field added
- **New Field:** `runtime`
  - **Type:** string
  - **Description:** Runtime type for team leader: 'default' (Agno) or 'claude_code' (Claude Code SDK)
  - **Default:** "default"
  - **Required:** No (has default)
- **Terraform Impact:**
  - Add `runtime` attribute to team data source and resource outputs
  - This field will be present in all team responses

### 5.3 TeamWithAgentsResponse
- **Change Type:** New field added
- **New Field:** `runtime`
  - **Type:** string
  - **Description:** Runtime type for team leader: 'default' (Agno) or 'claude_code' (Claude Code SDK)
  - **Default:** "default"
  - **Required:** No (has default)
- **Terraform Impact:**
  - Add `runtime` attribute when reading teams with agents
  - Update schema to include this field

### 5.4 WorkerQueueCreate
- **Change Type:** Field default value changed and description updated
- **Modified Field:** `heartbeat_interval`
  - **OLD Default:** 30 seconds
  - **NEW Default:** 60 seconds
  - **OLD Description:** "Seconds between heartbeats"
  - **NEW Description:** "Seconds between heartbeats (lightweight)"
  - **Type:** integer (unchanged)
  - **Constraints:** min=10, max=300 (unchanged)
- **Terraform Impact:**
  - Update default value for `heartbeat_interval` from 30 to 60
  - Update documentation to reflect "lightweight" heartbeats
  - Existing configurations with explicit values are not affected
  - Configurations relying on default will now use 60 instead of 30

### 5.5 WorkerConfigResponse
- **Change Type:** Field default value changed
- **Modified Field:** `heartbeat_interval`
  - **OLD Default:** 30 seconds
  - **NEW Default:** 60 seconds
  - **Type:** integer (unchanged)
- **Terraform Impact:**
  - Update schema default for worker config responses
  - Existing worker queues will maintain their configured values
  - New worker queues will default to 60 seconds

---

## 6. OTHER SIGNIFICANT CHANGES

### 6.1 API Version
- **No Change:** Version remains at "0.3.0"

### 6.2 Overall Statistics
- **Endpoints:** Increased from 127 to 141 (+14 endpoints, +11%)
- **Schemas:** Increased from 115 to 116 (+1 schema)

### 6.3 Feature Additions
The major addition is a comprehensive **Context Graph API** that provides:
- Graph database query capabilities
- Node and relationship management
- Integration with organization data
- Cypher query support
- Text and property-based search
- Subgraph traversal
- Statistics and analytics

---

## 7. TERRAFORM PROVIDER IMPACT SUMMARY

### 7.1 New Data Sources Required
1. `kubiya_context_graph_integrations` - List integrations
2. `kubiya_context_graph_labels` - Get node labels
3. `kubiya_context_graph_nodes` - List all nodes
4. `kubiya_context_graph_node` - Get specific node by ID
5. `kubiya_context_graph_node_relationships` - Get node relationships
6. `kubiya_context_graph_relationship_types` - List relationship types
7. `kubiya_context_graph_stats` - Get graph statistics
8. `kubiya_websocket_stats` - Get WebSocket statistics

### 7.2 New Data Sources with POST Methods (Read Operations)
These endpoints use POST but are read-only operations (searching/querying):
1. `kubiya_context_graph_nodes_search` - Search nodes by properties
2. `kubiya_context_graph_nodes_text_search` - Text-based node search
3. `kubiya_context_graph_query` - Execute custom Cypher queries
4. `kubiya_context_graph_subgraph` - Get subgraph from a node

### 7.3 Resource Updates Required
1. **Team Resource (`kubiya_team`)**
   - Add optional `runtime` attribute (string, default: "default")
   - Valid values: "default", "claude_code"
   - Add to schema, create, read, and update operations

2. **Worker Queue Resource (`kubiya_worker_queue`)**
   - Update `heartbeat_interval` default from 30 to 60
   - Update documentation to reflect "lightweight" heartbeats
   - No breaking changes (existing explicit values preserved)

### 7.4 New Resource Actions/Functions
1. **Execution Cancel Workflow**
   - New action to cancel specific workflow within execution
   - Requires: execution_id, workflow_message_id
   - Different from existing cancel execution functionality

### 7.5 Documentation Updates Required
1. Update AEM analytics documentation to reflect new model tier names:
   - OLD: Opus, Sonnet, Haiku
   - NEW: Premium, Mid, Basic (provider-agnostic)

### 7.6 Breaking Changes
**None** - All changes are additive or have backward-compatible defaults.

### 7.7 Priority Implementation Order
1. **High Priority:** Team runtime field (affects core resource)
2. **High Priority:** Context Graph data sources (major new feature)
3. **Medium Priority:** Worker queue heartbeat default change
4. **Medium Priority:** Cancel workflow action
5. **Low Priority:** WebSocket stats data source
6. **Low Priority:** Documentation updates for AEM model tiers
