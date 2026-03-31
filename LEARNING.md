# Claude Code 学习计划

## 学习进度

| 模块 | 状态 | 完成日期 | 备注 |
|------|------|---------|------|
| MCP（Model Context Protocol） | ✅ 已完成 | 2026-03-30 | 连接外部工具/服务 |
| Hooks 系统 | ✅ 已完成 | 2026-03-30 | 事件触发钩子，Stop/PreToolUse/PostToolUse/UserPromptSubmit |
| Settings / 权限管理 | ✅ 已完成 | 2026-03-30 | 三层配置文件，allow/deny/ask 权限规则，通配符语法 |
| Slash Commands | ✅ 已完成 | 2026-03-31 | Skills 系统，frontmatter，$ARGUMENTS，动态注入，go-review/commit/daily 三个实战 skill |
| Claude API | ✅ 已完成 | 2026-03-31 | Messages API、Vision API、并行调用、优雅降级，实战：图片转 ASCII 同时生成描述 |
| Agent SDK | 🔄 学习中 | - | - |
| Subagents（子代理） | 🔲 待学习 | - | - |

---

## 各模块详情

### ✅ MCP（Model Context Protocol）
- **是什么**：模型上下文协议，让 Claude 连接外部工具和数据源
- **能做什么**：连接 Jira、Slack、数据库、文件系统等外部服务
- **学习资料**：[MCP 官方文档](https://docs.anthropic.com/en/docs/claude-code/mcp)

---

### 🔲 Hooks 系统（下一步）
- **是什么**：事件触发钩子，在 Claude 执行前后自动触发 shell 命令
- **能做什么**：
  - Claude 停止后自动运行测试
  - 文件修改前自动备份
  - 出错时发送通知
  - 提交前自动格式化代码
- **配置位置**：`settings.json`
- **学习资料**：[Hooks 文档](https://docs.anthropic.com/en/docs/claude-code/hooks)

---

### 🔲 Settings / 权限管理
- **是什么**：通过 `settings.json` 控制 Claude 的行为和权限
- **能做什么**：
  - 允许/禁止特定 shell 命令
  - 设置环境变量
  - 管理文件读写权限
  - 配置 git 操作权限
- **配置层级**：全局 `~/.claude/settings.json` → 项目级 `settings.json`

---

### 🔲 Slash Commands（自定义命令）
- **是什么**：`/` 开头的自定义命令，封装常用提示词
- **能做什么**：
  - 把常用操作封装成一个命令
  - 复用提示词模板
  - 快速触发标准化流程
- **存放位置**：`.claude/commands/` 目录下的 `.md` 文件

---

### 🔲 Claude API
- **是什么**：直接调用 Anthropic Claude API
- **能做什么**：
  - 构建自己的应用集成 Claude 能力
  - Messages API（对话）
  - Vision API（图像分析）
  - Tool Use API（工具调用）
  - Batch API（批量处理）
- **学习资料**：[API 参考](https://docs.anthropic.com/en/reference/messages)

---

### 🔲 Agent SDK
- **是什么**：构建自定义 AI Agent 的框架
- **能做什么**：
  - 创建有工具、有记忆、有决策逻辑的 Agent
  - 实现代码审查 Agent、文档生成 Agent 等
  - 构建复杂自动化工作流
- **学习资料**：[Agent SDK 文档](https://docs.anthropic.com/en/docs/claude-code/sdk)

---

### 🔲 Subagents（子代理）
- **是什么**：多个专职小 Agent 协作完成任务
- **能做什么**：
  - 任务并行化处理
  - 角色专业分工
  - 主 Agent 调度多个子 Agent

---

## 推荐学习顺序

```
第 1 步：Hooks 系统           ← ✅ 已完成
第 2 步：Settings / 权限管理  ← 当前位置
第 3 步：Slash Commands       ← 提升工作效率
第 4 步：Claude API           ← 构建自己的应用
第 5 步：Agent SDK            ← 进阶，构建自定义 Agent
第 6 步：Subagents            ← 高级，多 Agent 协作
```

---

## 学习笔记

### 2026-03-30
- 梳理了 Claude Code 完整学习路线
- MCP 已掌握，确定下一步学习 Hooks 系统
- Hooks 系统学习完成：Stop/PreToolUse/PostToolUse/UserPromptSubmit 四种事件，command 类型 hook，exit code 控制阻断，systemMessage 输出格式，三个实战 hook 全部验证通过
