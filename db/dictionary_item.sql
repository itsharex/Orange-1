DROP TABLE IF EXISTS "dictionary_item";
CREATE TABLE "dictionary_item" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "dictionary_id" integer NOT NULL,
  "label" text(50) NOT NULL,
  "value" text(50) NOT NULL,
  "sort" integer NOT NULL DEFAULT 0,
  "status" integer(1) NOT NULL DEFAULT 1,
  "remark" text(255),
  "create_time" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" datetime DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_dictionary_id" ON "dictionary_item" ("dictionary_id");
CREATE INDEX "idx_dictionary_id_label" ON "dictionary_item" ("dictionary_id", "label");
