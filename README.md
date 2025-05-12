# API 线上笔试任务说明

## 📌 任务目标

你需要建立一组 RESTful API，包含以下功能：
1. 用户登录模拟（POST /login）
2. 查询该用户的报告列表（GET /reports?user_id=xxx）
3. 查询单份报告内容（GET /report/:report_id）
4. 生成报告的 AI 摘要（POST /report/:id/summary）

## 🧱 JSON 模拟数据
请参考本文件夹中的 `users.json` 和 `reports.json` 两个文件作为数据源。

## 🛠️ 技术建议
- 使用 Node.js + Express 或 Python Flask 开发均可
- 可将数据存储在内存中，或使用 SQLite 管理
- 若无法接入 GPT，模拟摘要结果即可

## ✅ 提交内容
- GitHub 仓库或 Zip 压缩包（附上可执行说明）
- API 结构清晰、功能逻辑完整
- 加分项：Postman 测试文件 / Swagger 文档 / 线上部署网址

祝你开发顺利 🚀


## 您好

感谢您申请我们的AI工具开发实习生（AI-Enhanced Full Stack Intern）职位！
我们对您的背景非常感兴趣，为了更进一步了解您的实践能力，我们诚挚邀请您参与接下来的线上技术挑战任务。
本次任务将模拟您未来在团队中可能参与的开发模块，包含API设计、数据结构管理与AI工具整合的基本逻辑。
🔧 任务内容简介：
• 实现登录、报告查询与摘要生成的简易API系统
• 使用您熟悉的技术（Node.js/Flask/SQLite/JSON皆可）
• 所需数据与任务说明已包含在附件ZIP文件中

📌 请注意：本任务要求候选人自行确保项目能够在所选择的环境中正确安装与运行。我们也鼓励您（非必要但加分）将项目部署至临时性的生产环境，方便我们直接浏览与测试成果。
📅 完成期限：建议在收到任务后3天内回复GitHub链接或完整压缩包
📨 回复方式：
请将完成的项目上传至GitHub或云盘（Google Drive/WeTransfer），
并回复本邮件附上链接即可。我们将尽快安排后续技术面试。

若您有任何问题，也欢迎随时来信咨询！

祝您开发顺利 🍀

敬祝
顺心如意，

丰鑫团队


curl 下载pdf文件：
curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTgzMTgwNjM0NTEzMjYzODIwOCIsImVtYWlsIjoiMTQ0QHFxLmNvbSIsImlzcyI6ImFpLWVuaGFuY2UtdGVzdCIsInN1YiI6IjE4MzE4MDYzNDUxMzI2MzgyMDgiLCJleHAiOjE3NDY5NTkxOTYsIm5iZiI6MTc0Njk1NTU5NiwiaWF0IjoxNzQ2OTU1NTk2fQ.U1pjFSb2azQaQBb3hyPleJuudi72vD1ZtmHGtYVjcWFVMwENd4Em6eUCoMG2zf7KMmdIZlTDsIKtkviXm_yhHYVJwoEFkWgks-v5s8AfaBiGPztpysoAaDJPLzi_ICwb3shBXOQNs9YhBKCwk3Dqpj9PhgL3NviQRqwPleTrt3LO31mvvrceWbhX2B2v-QUtcKyOkG9hz02yAjUaKvq_NzEL9RT6JfvU07ZBAyA4vBfdeffKDF6wvf3nbjSR0md1jSxdjSrNg_KojtZSXlO4QOui2ybBU2IVMTgNpjsv9YVeAuysCA1E4s1xm6iXXDrMiiSDPiFUm8EX0NQf42XxQA"      -o download.pdf      "http://localhost:8080/api/v1/report/513205d4-8705-4b61-83ed-8ac88c31bd58/pdf"
