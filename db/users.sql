DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "username" text(50) NOT NULL,
  "password" text(100) NOT NULL,
  "name" text(50) NOT NULL,
  "email" text(100),
  "phone" text(20),
  "avatar" text(255),
  "role" text(20) NOT NULL DEFAULT 'user',
  "department" text(50),
  "position" text(50),
  "status" integer(1) NOT NULL DEFAULT 1,
  "last_login_time" datetime,
  "create_time" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_time" datetime DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
