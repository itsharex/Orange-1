-- 通知表
DROP TABLE IF EXISTS "notifications";
CREATE TABLE IF NOT EXISTS "notifications" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "title" VARCHAR(100) NOT NULL,
    "content" TEXT NOT NULL,
    "type" INTEGER DEFAULT 1, -- 1:system, 2:activity, 3:private
    "sender_id" INTEGER NOT NULL,
    "is_global" INTEGER DEFAULT 0, -- 0:No, 1:Yes
    "create_time" DATETIME DEFAULT CURRENT_TIMESTAMP,
    "update_time" DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_notifications_sender_id" ON "notifications"("sender_id");
CREATE INDEX IF NOT EXISTS "idx_notifications_is_global" ON "notifications"("is_global");
