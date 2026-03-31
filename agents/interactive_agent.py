"""
方案 C：带 Session 的交互式 CLI Agent
多轮对话，Agent 记住上下文，支持追问

演示知识点：
- session_id 从 init 消息中获取
- resume=session_id 恢复上下文
- 每轮对话都保留之前的文件读取、分析结果
"""

import asyncio
from claude_agent_sdk import query, ClaudeAgentOptions

SYSTEM_PROMPT = """你是这个 Go 项目的专属助手（ascii_convert_go）。
你熟悉项目结构：main.go（HTTP服务）、converter/（核心转换）、describer/（Claude Vision）、mcp/（MCP Server）。
每次回答要简洁，中文回答。"""


async def chat_turn(prompt: str, session_id: str | None) -> str:
    """执行一轮对话，返回 session_id"""
    new_session_id = session_id
    result_text = ""

    options = ClaudeAgentOptions(
        allowed_tools=["Read", "Glob", "Grep", "Bash"],
        permission_mode="acceptEdits",
        cwd="/Users/koopli/workplace/self_proj/ascii_convert_go",
        system_prompt=SYSTEM_PROMPT,
    )

    # 有 session_id 时恢复上下文，否则新建
    if session_id:
        options = ClaudeAgentOptions(
            allowed_tools=["Read", "Glob", "Grep", "Bash"],
            permission_mode="acceptEdits",
            cwd="/Users/koopli/workplace/self_proj/ascii_convert_go",
            system_prompt=SYSTEM_PROMPT,
            resume=session_id,
        )

    async for message in query(prompt=prompt, options=options):
        cls = type(message).__name__

        # session_id 在 AssistantMessage 或 ResultMessage 上
        if hasattr(message, "session_id") and message.session_id:
            new_session_id = message.session_id

        if cls == "ResultMessage":
            result_text = getattr(message, "result", "")

    return new_session_id, result_text


async def main():
    print("🤖 ASCII Convert Go 交互式 Agent")
    print("   (输入 'quit' 退出，'new' 开始新会话，'session' 查看当前 session ID)\n")

    session_id = None
    turn = 0

    while True:
        try:
            user_input = input("你: ").strip()
        except (EOFError, KeyboardInterrupt):
            print("\n👋 再见！")
            break

        if not user_input:
            continue

        if user_input.lower() == "quit":
            print("👋 再见！")
            break

        if user_input.lower() == "new":
            session_id = None
            turn = 0
            print("🔄 已开始新会话\n")
            continue

        if user_input.lower() == "session":
            if session_id:
                print(f"📍 当前 Session: {session_id[:16]}...\n")
            else:
                print("📍 还没有 Session（发送第一条消息后会自动创建）\n")
            continue

        turn += 1
        status = f"[第{turn}轮]" + ("" if not session_id else " [续接上下文]")
        print(f"🔄 {status} 思考中...")

        session_id, result = await chat_turn(user_input, session_id)

        print(f"\nAgent: {result}\n")
        print("-" * 40)


if __name__ == "__main__":
    asyncio.run(main())
