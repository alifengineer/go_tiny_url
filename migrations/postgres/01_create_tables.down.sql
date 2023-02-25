DROP INDEX IF EXISTS "idx_passcode_user_id";
DROP TABLE IF EXISTS "passcode";
DROP INDEX IF EXISTS "idx_session_user_id";
DROP TABLE IF EXISTS "session";
DROP TABLE IF EXISTS "user_info";
DROP TABLE IF EXISTS "user_relation";
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "role_permission";
DROP TABLE IF EXISTS "permission_scope";
ALTER TABLE "permission" DROP CONSTRAINT IF EXISTS "fk_permission_parent_id";
DROP TABLE IF EXISTS "permission";
DROP TABLE IF EXISTS "scope";
DROP TABLE IF EXISTS "role";
DROP TABLE IF EXISTS "client";
DROP TYPE IF EXISTS "login_strategies";
DROP TABLE IF EXISTS "user_info_field";
DROP TABLE IF EXISTS "relation";
DROP TYPE IF EXISTS "relation_types";
DROP TABLE IF EXISTS "client_type";
DROP TYPE IF EXISTS "confirm_strategies";
DROP TABLE IF EXISTS "client_platform";
DROP TABLE IF EXISTS "project";