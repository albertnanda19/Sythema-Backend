BEGIN;

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS diff_metrics;
DROP TABLE IF EXISTS diff_results;
DROP TABLE IF EXISTS replay_results;
DROP TABLE IF EXISTS replay_tasks;
DROP TABLE IF EXISTS replay_jobs;
DROP TABLE IF EXISTS transform_rules;
DROP TABLE IF EXISTS transform_rule_sets;
DROP TABLE IF EXISTS traffic_requests;
DROP TABLE IF EXISTS traffic_sessions;
DROP TABLE IF EXISTS shadow_targets;
DROP TABLE IF EXISTS environments;
DROP TABLE IF EXISTS projects;

COMMIT;
