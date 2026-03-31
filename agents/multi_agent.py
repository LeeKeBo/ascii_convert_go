"""
方案 B：多 Agent 协作
主 Agent 接收图片路径 → 派子 Agent 做 ASCII 转换 → 派子 Agent 做 Vision 描述 → 汇总

演示知识点：
- AgentDefinition 定义专职子 Agent
- allowed_tools 包含 "Agent" 才能派子 Agent
- parent_tool_use_id 追踪消息属于哪个子 Agent
"""

import asyncio
import sys
from claude_agent_sdk import query, ClaudeAgentOptions, AgentDefinition


async def main(image_path: str):
    print(f"🖼️  处理图片：{image_path}\n")
    print("=" * 50)

    # 定义两个专职子 Agent
    agents = {
        "ascii-converter": AgentDefinition(
            description="专职把图片转为 ASCII 字符画，使用 go run 调用项目工具",
            prompt="""你是 ASCII 转换专家。
接收图片路径，调用项目的 Go 程序转换为 ASCII 字符画。
使用命令：go run main.go <图片路径>
输出转换结果或错误信息。""",
            tools=["Bash"],
        ),
        "image-describer": AgentDefinition(
            description="专职用 Claude Vision API 描述图片内容",
            prompt="""你是图片描述专家。
接收图片路径，用 Python 调用 Claude Vision API 生成一句话中文描述。
使用以下代码：
```python
import anthropic, base64, sys
client = anthropic.Anthropic()
with open(sys.argv[1], "rb") as f:
    b64 = base64.b64encode(f.read()).decode()
msg = client.messages.create(
    model="claude-haiku-4-5-20251001",
    max_tokens=100,
    messages=[{"role": "user", "content": [
        {"type": "image", "source": {"type": "base64", "media_type": "image/png", "data": b64}},
        {"type": "text", "text": "用一句话描述这张图片，中文，不超过30字"}
    ]}]
)
print(msg.content[0].text)
```
将代码保存为临时文件并执行，输出描述结果。""",
            tools=["Bash", "Write"],
        ),
    }

    results = {"ascii": "", "description": "", "subagent_msgs": []}

    async for message in query(
        prompt=f"""对图片 {image_path} 并行完成两件事：
1. 用 ascii-converter 子 Agent 将图片转为 ASCII 字符画（width=80）
2. 用 image-describer 子 Agent 生成图片的一句话中文描述

最后汇总两个结果，格式：
---
📝 图片描述：<描述>
🎨 ASCII 预览：（告知用户 ASCII 已生成，文件在哪）
---""",
        options=ClaudeAgentOptions(
            allowed_tools=["Bash", "Write", "Agent"],
            agents=agents,
            permission_mode="acceptEdits",
            cwd="/Users/koopli/workplace/self_proj/ascii_convert_go",
        ),
    ):
        cls = type(message).__name__
        parent_id = getattr(message, "parent_tool_use_id", None)

        if cls == "SystemMessage":
            data = getattr(message, "data", {})
            sid = data.get("session_id", "")
            if sid:
                print(f"🚀 Session 启动: {sid[:8]}...")

        elif cls == "AssistantMessage":
            prefix = "  [子Agent] " if parent_id else ""
            for block in getattr(message, "content", []):
                if hasattr(block, "text") and block.text.strip():
                    print(f"{prefix}💬 {block.text[:200]}")

        elif cls == "ResultMessage":
            print("\n" + "=" * 50)
            print("✅ 最终结果：")
            print(getattr(message, "result", ""))

    print("\n多 Agent 协作完成！")


if __name__ == "__main__":
    path = sys.argv[1] if len(sys.argv) > 1 else "test.png"
    asyncio.run(main(path))
