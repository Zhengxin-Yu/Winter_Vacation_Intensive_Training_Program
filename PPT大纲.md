1. 封面
项目名称：酒店行李寄存系统（Go + Gin）
学校/课程/姓名/日期
2. 项目背景与目标
酒店行李寄存流程痛点（手工登记、易错、难追溯）
系统目标：标准化寄存流程、可追溯、易查询
3. 系统架构
前后端分离
技术栈：Gin + GORM + MySQL + JWT + Redis（可选）
模块结构（handlers / services / repositories）
4. 数据库设计
users
luggage_items
luggage_history
luggage_storerooms
luggage_migrations
行李修改日志表（更新前后快照）
5. 核心功能流程
登录鉴权（JWT）
寄存：生成 6 位取件码
查询：取件码 / 客人姓名
取件：自动写历史
6. 行李寄存单管理
新增 / 修改寄存信息
修改记录可追溯（old_data / new_data）
7. 寄存室管理
创建 / 停用
剩余容量与已存数量统计
8. 图片上传与展示
/api/upload 获取图片 URL
photo_urls 多图、photo_url 兼容
9. API 设计规范
统一前缀 /api
统一响应格式 {message, error}
关键接口列表（简要展示）
10. 测试与演示流程
登录 → 上传 → 寄存 → 查询 → 取件
工具：Apifox / curl
11. 遇到的问题与解决
JWT token 格式错误
数据库字段调整
取件码冲突处理
12. 总结与展望
已完成：寄存、取件、查询、日志、上传
可扩展：前端页面、统计报表、二维码展示