CREATE TABLE projects (
  id UUID PRIMARY KEY,
  slug TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE environments (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  name TEXT NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE shadow_targets (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  source_environment_id UUID NOT NULL,
  target_environment_id UUID NOT NULL,
  name TEXT NOT NULL,
  status VARCHAR(32) NOT NULL,
  replay_strategy TEXT NOT NULL,
  diff_strategy TEXT NOT NULL,
  config JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE traffic_sessions (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  source_environment_id UUID NOT NULL,
  shadow_target_id UUID,
  external_session_key TEXT,
  status VARCHAR(32) NOT NULL,
  started_at TIMESTAMPTZ,
  ended_at TIMESTAMPTZ,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE traffic_requests (
  id UUID PRIMARY KEY,
  session_id UUID NOT NULL,
  sequence_no INTEGER NOT NULL,
  captured_at TIMESTAMPTZ NOT NULL,
  method VARCHAR(16) NOT NULL,
  scheme VARCHAR(16),
  host TEXT,
  path TEXT NOT NULL,
  query_string TEXT,
  headers JSONB,
  content_length_bytes INTEGER,
  body_hash TEXT,
  request_fingerprint TEXT,
  client_ip_hash TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE transform_rule_sets (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  name TEXT NOT NULL,
  version INTEGER NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE transform_rules (
  id UUID PRIMARY KEY,
  rule_set_id UUID NOT NULL,
  order_no INTEGER NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT true,
  match_criteria JSONB,
  action_type TEXT NOT NULL,
  action_config JSONB,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE replay_jobs (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  shadow_target_id UUID NOT NULL,
  transform_rule_set_id UUID,
  requested_by_api_key_id UUID,
  status VARCHAR(32) NOT NULL,
  requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  params JSONB,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE replay_tasks (
  id UUID PRIMARY KEY,
  replay_job_id UUID NOT NULL,
  task_type VARCHAR(32) NOT NULL,
  traffic_session_id UUID,
  batch_key TEXT,
  status VARCHAR(32) NOT NULL,
  attempt INTEGER NOT NULL DEFAULT 0,
  scheduled_at TIMESTAMPTZ,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  error_message TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE replay_results (
  id UUID PRIMARY KEY,
  replay_task_id UUID NOT NULL,
  traffic_request_id UUID,
  status VARCHAR(32) NOT NULL,
  target_status_code INTEGER,
  latency_ms INTEGER,
  response_size_bytes INTEGER,
  response_hash TEXT,
  error_class TEXT,
  error_message TEXT,
  finished_at TIMESTAMPTZ,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE diff_results (
  id UUID PRIMARY KEY,
  replay_result_id UUID NOT NULL,
  status VARCHAR(32) NOT NULL,
  diff_strategy TEXT NOT NULL,
  summary JSONB,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE diff_metrics (
  id UUID PRIMARY KEY,
  diff_result_id UUID NOT NULL,
  metric_key TEXT NOT NULL,
  metric_value_numeric DOUBLE PRECISION,
  metric_value_text TEXT,
  units TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE api_keys (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  name TEXT NOT NULL,
  prefix TEXT,
  key_hash TEXT NOT NULL,
  status VARCHAR(32) NOT NULL,
  scopes JSONB,
  last_used_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE audit_logs (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  actor_type VARCHAR(32) NOT NULL,
  actor_api_key_id UUID,
  action TEXT NOT NULL,
  entity_type TEXT,
  entity_id UUID,
  request_id UUID,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE projects
  ADD CONSTRAINT projects_slug_non_empty_check CHECK (length(slug) > 0);

CREATE UNIQUE INDEX idx_projects_slug_active_unique ON projects (slug) WHERE deleted_at IS NULL;

ALTER TABLE environments
  ADD CONSTRAINT environments_name_non_empty_check CHECK (length(name) > 0);

CREATE UNIQUE INDEX idx_environments_project_name_active_unique ON environments (project_id, name) WHERE deleted_at IS NULL;

ALTER TABLE environments
  ADD CONSTRAINT environments_status_check CHECK (status IN ('active', 'disabled'));

ALTER TABLE shadow_targets
  ADD CONSTRAINT shadow_targets_name_non_empty_check CHECK (length(name) > 0);

CREATE UNIQUE INDEX idx_shadow_targets_project_name_active_unique ON shadow_targets (project_id, name) WHERE deleted_at IS NULL;

ALTER TABLE shadow_targets
  ADD CONSTRAINT shadow_targets_status_check CHECK (status IN ('active', 'paused', 'disabled'));

ALTER TABLE shadow_targets
  ADD CONSTRAINT shadow_targets_source_target_diff_check CHECK (source_environment_id <> target_environment_id);

ALTER TABLE traffic_sessions
  ADD CONSTRAINT traffic_sessions_status_check CHECK (status IN ('open', 'closed', 'failed'));

ALTER TABLE traffic_requests
  ADD CONSTRAINT traffic_requests_session_sequence_unique UNIQUE (session_id, sequence_no);

ALTER TABLE traffic_requests
  ADD CONSTRAINT traffic_requests_sequence_no_check CHECK (sequence_no >= 0);

ALTER TABLE traffic_requests
  ADD CONSTRAINT traffic_requests_content_length_bytes_check CHECK (content_length_bytes IS NULL OR content_length_bytes >= 0);

ALTER TABLE transform_rule_sets
  ADD CONSTRAINT transform_rule_sets_version_check CHECK (version > 0);

CREATE UNIQUE INDEX idx_transform_rule_sets_project_name_version_active_unique ON transform_rule_sets (project_id, name, version) WHERE deleted_at IS NULL;

ALTER TABLE transform_rule_sets
  ADD CONSTRAINT transform_rule_sets_status_check CHECK (status IN ('draft', 'active', 'archived'));

ALTER TABLE transform_rules
  ADD CONSTRAINT transform_rules_rule_set_order_unique UNIQUE (rule_set_id, order_no);

ALTER TABLE transform_rules
  ADD CONSTRAINT transform_rules_order_no_check CHECK (order_no >= 0);

ALTER TABLE replay_jobs
  ADD CONSTRAINT replay_jobs_status_check CHECK (status IN ('queued', 'running', 'succeeded', 'failed', 'canceled'));

ALTER TABLE replay_tasks
  ADD CONSTRAINT replay_tasks_type_check CHECK (task_type IN ('session', 'batch'));

ALTER TABLE replay_tasks
  ADD CONSTRAINT replay_tasks_status_check CHECK (status IN ('queued', 'running', 'succeeded', 'failed', 'canceled'));

ALTER TABLE replay_tasks
  ADD CONSTRAINT replay_tasks_attempt_check CHECK (attempt >= 0);

ALTER TABLE replay_results
  ADD CONSTRAINT replay_results_status_check CHECK (status IN ('succeeded', 'failed', 'skipped'));

ALTER TABLE diff_results
  ADD CONSTRAINT diff_results_status_check CHECK (status IN ('matched', 'mismatched', 'error', 'skipped'));

ALTER TABLE diff_metrics
  ADD CONSTRAINT diff_metrics_value_present_check CHECK (
    metric_value_numeric IS NOT NULL OR metric_value_text IS NOT NULL
  );

ALTER TABLE api_keys
  ADD CONSTRAINT api_keys_name_non_empty_check CHECK (length(name) > 0);

CREATE UNIQUE INDEX idx_api_keys_project_name_active_unique ON api_keys (project_id, name) WHERE deleted_at IS NULL;

ALTER TABLE api_keys
  ADD CONSTRAINT api_keys_key_hash_unique UNIQUE (key_hash);

ALTER TABLE api_keys
  ADD CONSTRAINT api_keys_status_check CHECK (status IN ('active', 'revoked'));

ALTER TABLE audit_logs
  ADD CONSTRAINT audit_logs_actor_type_check CHECK (actor_type IN ('api_key', 'system'));

ALTER TABLE environments
  ADD CONSTRAINT environments_project_id_fkey
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE shadow_targets
  ADD CONSTRAINT shadow_targets_project_id_fkey
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE shadow_targets
  ADD CONSTRAINT shadow_targets_source_environment_id_fkey
  FOREIGN KEY (source_environment_id) REFERENCES environments (id) ON DELETE RESTRICT;

ALTER TABLE shadow_targets
  ADD CONSTRAINT shadow_targets_target_environment_id_fkey
  FOREIGN KEY (target_environment_id) REFERENCES environments (id) ON DELETE RESTRICT;

ALTER TABLE traffic_sessions
  ADD CONSTRAINT traffic_sessions_project_id_fkey
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE traffic_sessions
  ADD CONSTRAINT traffic_sessions_source_environment_id_fkey
  FOREIGN KEY (source_environment_id) REFERENCES environments (id) ON DELETE RESTRICT;

ALTER TABLE traffic_sessions
  ADD CONSTRAINT traffic_sessions_shadow_target_id_fkey
  FOREIGN KEY (shadow_target_id) REFERENCES shadow_targets (id) ON DELETE RESTRICT;

ALTER TABLE traffic_requests
  ADD CONSTRAINT traffic_requests_session_id_fkey
  FOREIGN KEY (session_id) REFERENCES traffic_sessions (id) ON DELETE RESTRICT;

ALTER TABLE transform_rule_sets
  ADD CONSTRAINT transform_rule_sets_project_id_fkey
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE transform_rules
  ADD CONSTRAINT transform_rules_rule_set_id_fkey
  FOREIGN KEY (rule_set_id) REFERENCES transform_rule_sets (id) ON DELETE RESTRICT;

ALTER TABLE replay_jobs
  ADD CONSTRAINT replay_jobs_project_id_fkey
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE replay_jobs
  ADD CONSTRAINT replay_jobs_shadow_target_id_fkey
  FOREIGN KEY (shadow_target_id) REFERENCES shadow_targets (id) ON DELETE RESTRICT;

ALTER TABLE replay_jobs
  ADD CONSTRAINT replay_jobs_transform_rule_set_id_fkey
  FOREIGN KEY (transform_rule_set_id) REFERENCES transform_rule_sets (id) ON DELETE RESTRICT;

ALTER TABLE replay_jobs
  ADD CONSTRAINT replay_jobs_requested_by_api_key_id_fkey
  FOREIGN KEY (requested_by_api_key_id) REFERENCES api_keys (id) ON DELETE RESTRICT;

ALTER TABLE replay_tasks
  ADD CONSTRAINT replay_tasks_replay_job_id_fkey
  FOREIGN KEY (replay_job_id) REFERENCES replay_jobs (id) ON DELETE RESTRICT;

ALTER TABLE replay_tasks
  ADD CONSTRAINT replay_tasks_traffic_session_id_fkey
  FOREIGN KEY (traffic_session_id) REFERENCES traffic_sessions (id) ON DELETE RESTRICT;

ALTER TABLE replay_results
  ADD CONSTRAINT replay_results_replay_task_id_fkey
  FOREIGN KEY (replay_task_id) REFERENCES replay_tasks (id) ON DELETE RESTRICT;

ALTER TABLE replay_results
  ADD CONSTRAINT replay_results_traffic_request_id_fkey
  FOREIGN KEY (traffic_request_id) REFERENCES traffic_requests (id) ON DELETE RESTRICT;

ALTER TABLE diff_results
  ADD CONSTRAINT diff_results_replay_result_id_fkey
  FOREIGN KEY (replay_result_id) REFERENCES replay_results (id) ON DELETE RESTRICT;

ALTER TABLE diff_metrics
  ADD CONSTRAINT diff_metrics_diff_result_id_fkey
  FOREIGN KEY (diff_result_id) REFERENCES diff_results (id) ON DELETE RESTRICT;

ALTER TABLE api_keys
  ADD CONSTRAINT api_keys_project_id_fkey
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE audit_logs
  ADD CONSTRAINT audit_logs_project_id_fkey
  FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE RESTRICT;

ALTER TABLE audit_logs
  ADD CONSTRAINT audit_logs_actor_api_key_id_fkey
  FOREIGN KEY (actor_api_key_id) REFERENCES api_keys (id) ON DELETE RESTRICT;

CREATE INDEX idx_environments_project_id ON environments (project_id);

CREATE INDEX idx_shadow_targets_project_id ON shadow_targets (project_id);

CREATE INDEX idx_shadow_targets_source_environment_id ON shadow_targets (source_environment_id);

CREATE INDEX idx_shadow_targets_target_environment_id ON shadow_targets (target_environment_id);

CREATE INDEX idx_traffic_sessions_project_id ON traffic_sessions (project_id);

CREATE INDEX idx_traffic_sessions_source_environment_id ON traffic_sessions (source_environment_id);

CREATE INDEX idx_traffic_sessions_shadow_target_id ON traffic_sessions (shadow_target_id);

CREATE INDEX idx_traffic_sessions_started_at ON traffic_sessions (started_at);

CREATE INDEX idx_traffic_requests_session_id ON traffic_requests (session_id);

CREATE INDEX idx_traffic_requests_captured_at ON traffic_requests (captured_at);

CREATE INDEX idx_traffic_requests_request_fingerprint ON traffic_requests (request_fingerprint);

CREATE INDEX idx_transform_rule_sets_project_id ON transform_rule_sets (project_id);

CREATE INDEX idx_transform_rules_rule_set_id ON transform_rules (rule_set_id);

CREATE INDEX idx_replay_jobs_project_id ON replay_jobs (project_id);

CREATE INDEX idx_replay_jobs_shadow_target_id ON replay_jobs (shadow_target_id);

CREATE INDEX idx_replay_jobs_status ON replay_jobs (status);

CREATE INDEX idx_replay_jobs_requested_at ON replay_jobs (requested_at);

CREATE INDEX idx_replay_tasks_replay_job_id ON replay_tasks (replay_job_id);

CREATE INDEX idx_replay_tasks_status ON replay_tasks (status);

CREATE INDEX idx_replay_tasks_traffic_session_id ON replay_tasks (traffic_session_id);

CREATE INDEX idx_replay_results_replay_task_id ON replay_results (replay_task_id);

CREATE INDEX idx_replay_results_traffic_request_id ON replay_results (traffic_request_id);

CREATE INDEX idx_diff_results_replay_result_id ON diff_results (replay_result_id);

CREATE INDEX idx_diff_metrics_diff_result_id ON diff_metrics (diff_result_id);

CREATE INDEX idx_diff_metrics_metric_key ON diff_metrics (metric_key);

CREATE INDEX idx_api_keys_project_id ON api_keys (project_id);

CREATE INDEX idx_api_keys_status ON api_keys (status);

CREATE INDEX idx_audit_logs_project_id_created_at ON audit_logs (project_id, created_at);

CREATE INDEX idx_audit_logs_actor_api_key_id ON audit_logs (actor_api_key_id);

CREATE INDEX idx_audit_logs_entity_type_entity_id ON audit_logs (entity_type, entity_id);
