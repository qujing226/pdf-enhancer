-- 创建数据库
CREATE DATABASE IF NOT EXISTS reports_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE reports_db;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` varchar(64) NOT NULL COMMENT '用户ID',
  `name` varchar(100) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
  `salt` varchar(64) NOT NULL COMMENT '盐值',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建报告表
CREATE TABLE IF NOT EXISTS `reports` (
  `id` varchar(64) NOT NULL COMMENT '报告ID',
  `user_id` varchar(64) NOT NULL COMMENT '用户ID',
  `title` varchar(255) NOT NULL COMMENT '报告标题',
  `content` longtext COMMENT '报告内容',
  `summary` text COMMENT 'AI生成的摘要',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `pdf_path` varchar(255) DEFAULT NULL COMMENT 'PDF文件路径',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_reports_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报告表';

-- 插入示例用户数据（密码为测试密码的哈希值，实际应用中应该使用Argon2id生成）
INSERT INTO `users` (`id`, `name`, `email`, `password_hash`, `salt`) VALUES
('u123', 'Alice', 'alice@example.com', '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0c2FsdA$bgGZ6ZNWcszOhcLPPb8/8QcBGCDOm1wBnBNCy+fGqhg', 'saltsaltsaltsalt'),
('u456', 'Bob', 'bob@example.com', '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0c2FsdA$bgGZ6ZNWcszOhcLPPb8/8QcBGCDOm1wBnBNCy+fGqhg', 'saltsaltsaltsalt');

-- 插入示例报告数据
INSERT INTO `reports` (`id`, `user_id`, `title`, `content`, `summary`, `pdf_path`) VALUES
('rpt001', 'u123', 'Q1 2025 Investment Report', '这是完整报告内容，可以很长很长...', '这是由 GPT 模拟产生的摘要内容。', '/reports/rpt001.pdf'),
('rpt002', 'u456', 'Q4 2024 Investment Report', '这是 Q4 报告内容，描述资产配置与绩效。', '', '/reports/rpt002.pdf');

-- 创建MinIO存储桶（这部分需要在应用程序中实现，这里只是注释说明）
-- 应用启动时需要检查并创建名为'reports'的存储桶
-- 并设置适当的访问策略