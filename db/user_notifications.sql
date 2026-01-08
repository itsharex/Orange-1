-- 用户通知关联表 (Inbox & Read Status)
DROP TABLE IF EXISTS "user_notifications";
CREATE TABLE IF NOT EXISTS "user_notifications" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "notification_id" INTEGER NOT NULL,
    "is_read" INTEGER DEFAULT 0, -- 0:unread, 1:read
    "read_time" DATETIME
);

-- 联合唯一索引，确保一个用户对一条通知只有一条记录
CREATE UNIQUE INDEX IF NOT EXISTS "idx_user_notification" ON "user_notifications"("user_id", "notification_id");
