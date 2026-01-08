-- 初始化字典数据

-- 款项阶段
INSERT INTO dictionaries (id, code, name, status, create_time) VALUES 
(1, 'payment_stage', '款项阶段', 1, CURRENT_TIMESTAMP);

INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(1, '首付款', 'deposit', 1, 1, CURRENT_TIMESTAMP),
(1, '进度款', 'progress', 2, 1, CURRENT_TIMESTAMP),
(1, '尾款', 'final', 3, 1, CURRENT_TIMESTAMP),
(1, '全款', 'all', 4, 1, CURRENT_TIMESTAMP);

-- 支付方式
INSERT INTO dictionaries (id, code, name, status, create_time) VALUES 
(2, 'payment_method', '支付方式', 1, CURRENT_TIMESTAMP);

INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(2, '银行转账', 'bank_transfer', 1, 1, CURRENT_TIMESTAMP),
(2, '支付宝', 'alipay', 2, 1, CURRENT_TIMESTAMP),
(2, '微信支付', 'wechat', 3, 1, CURRENT_TIMESTAMP),
(2, '现金', 'cash', 4, 1, CURRENT_TIMESTAMP);

-- 项目状态
INSERT INTO dictionaries (id, code, name, status, create_time) VALUES 
(3, 'project_status', '项目状态', 1, CURRENT_TIMESTAMP);

INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(3, '未开始', 'notstarted', 1, 1, CURRENT_TIMESTAMP),
(3, '进行中', 'active', 2, 1, CURRENT_TIMESTAMP),
(3, '已完成', 'completed', 3, 1, CURRENT_TIMESTAMP),
(3, '已逾期', 'overdue', 4, 1, CURRENT_TIMESTAMP),
(3, '已归档', 'archived', 5, 1, CURRENT_TIMESTAMP);

-- 项目类型
INSERT INTO dictionaries (id, code, name, status, create_time) VALUES 
(4, 'project_type', '项目类型', 1, CURRENT_TIMESTAMP);

INSERT INTO dictionary_item (dictionary_id, label, value, sort, status, create_time) VALUES 
(4, 'Web开发', 'web', 1, 1, CURRENT_TIMESTAMP),
(4, '移动应用', 'mobile', 2, 1, CURRENT_TIMESTAMP),
(4, 'UI设计', 'design', 3, 1, CURRENT_TIMESTAMP),
(4, 'SaaS系统', 'saas', 4, 1, CURRENT_TIMESTAMP),
(4, '其他', 'other', 5, 1, CURRENT_TIMESTAMP);
