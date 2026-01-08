DROP TABLE IF EXISTS "dictionaries";
CREATE TABLE "dictionaries" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "code" text(50) NOT NULL,
  "name" text(50) NOT NULL,
  "status" integer(1) NOT NULL DEFAULT 1,
  "remark" text(255),
  "create_time" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" datetime DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX "idx_dictionaries_code" ON "dictionaries" ("code");
