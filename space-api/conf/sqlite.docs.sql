CREATE VIRTUAL TABLE IF NOT EXISTS `sqlite3_keyword_docs` USING fts5 (
    -- 文章 ID, 后续返回依靠此字段
    `post_id`,
    -- 文章标题分词
    `title_split`,
    -- 文章分词内容
    `content_split`,
    -- 权重
    `weight`,
    -- 创建日期
    `post_updated_at`,
    -- 缓存创建的日期
    `record_created_at`,
    -- 缓存更新的日期
    `record_updated_at`
);
