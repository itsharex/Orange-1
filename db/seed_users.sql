-- 初始化用户数据
-- 密码: qq123456

INSERT INTO users (id, username, password, name, email, phone, role, department, position, status, create_time) VALUES 
(1, 'admin', '$2a$10$sp/NLWYRQjt9zVCq6HeOieFaFNBl79RoBXdePqxg9UhwQyT1/C7vu', '管理员', 'admin@orange.com', '13800000000', 'admin', '技术部', '系统管理员', 1, CURRENT_TIMESTAMP);

-- 测试用户 (密码: admin123)
INSERT INTO users (id, username, password, name, email, phone, role, department, position, status, create_time) VALUES 
(2, 'xu', '$2a$10$sp/NLWYRQjt9zVCq6HeOieFaFNBl79RoBXdePqxg9UhwQyT1/C7vu', '郑旭', 'xu@company.com', '13800000001', 'user', '技术部', '项目经理', 1, CURRENT_TIMESTAMP);
