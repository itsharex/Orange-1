DROP TABLE IF EXISTS "projects";
CREATE TABLE "projects" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text(100) NOT NULL,
  "company" text(100) NOT NULL,
  "total_amount" real NOT NULL,
  "received_amount" real DEFAULT 0,
  "status" text(20) NOT NULL,
  "type" text(50) NOT NULL,
  "contract_number" text(50),
  "contract_date" date,
  "payment_method" text(30),
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "description" text,
  "user_id" integer NOT NULL,
  "create_time" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" datetime DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_projects_user_id" ON "projects" ("user_id");
