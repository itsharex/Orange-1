DROP TABLE IF EXISTS "payments";
CREATE TABLE "payments" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "project_id" integer NOT NULL,
  "stage" text(50) NOT NULL,
  "amount" real NOT NULL,
  "percentage" real,
  "plan_date" date NOT NULL,
  "status" text(20) NOT NULL,
  "actual_date" date,
  "method" text(30),
  "remark" text(255),
  "user_id" integer NOT NULL,
  "create_time" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" datetime DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_payments_plan_date" ON "payments" ("plan_date");
CREATE INDEX "idx_payments_project" ON "payments" ("project_id");
CREATE INDEX "idx_payments_status" ON "payments" ("status");
