# Liveguard Panel

Liveguard Panel 是 Liveguard 前端界面。

## 📄 页面

- **index**：主页，侧边栏会展示口袋token的登录状态
- **Live**：获取直播列表，创建总结任务
- **Task**：管理任务
  - Restore：在任务失败、服务重启时，手动将任务回退到上一个状态重新进行
  - Delete：删除任务，任务表以live_id作为唯一键，任务无法被重复创建
- **Story**：查看生成的结果

## 🚀 技术栈
- **Python 3.12**：使用最新版本的 Python
- **Streamlit**：轻量且简单的 Web UI 框架
- **Poetry**：依赖管理和软件包管理工具
- **API 客户端**：`api_client.py` 负责与后端通信

## 🎯 选择 Streamlit 的理由
1. **快速构建 UI**
   - 无需掌握 HTML/CSS/JS，仅使用 Python 即可搭建 Web 界面
2. **数据可视化简单**
   - 原生支持图表和表格展示
3. **开发效率高**
   - 代码量少，适用于 PoC（概念验证）和 MVP（最小可行产品）开发
4. **易于集成 API 客户端**
   - 通过 `requests` 等库获取 API 数据，并实时展示
