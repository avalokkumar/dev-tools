#!/usr/bin/env python3
"""DevForge example plugin (Python).

Demonstrates the language-agnostic plugin transport: newline-delimited
JSON-RPC 2.0 over stdio. The host calls `initialize` once, then `invoke`
per call. This plugin contributes one Tool ("py") with one Operation
("reverse") that reverses an input string.
"""

import json
import sys
from typing import Any, Dict


SDK = 1

OPERATIONS = [
    {
        "tool": "py",
        "op": "reverse",
        "description": "Reverse the characters of an input string.",
        "inputSchema": {
            "type": "object",
            "required": ["input"],
            "properties": {"input": {"type": "string"}},
        },
    },
]


def handle_invoke(params: Dict[str, Any]) -> Dict[str, Any]:
    name = params.get("name", "")
    args = params.get("args", {}) or {}
    if name == "py_reverse":
        s = str(args.get("input", ""))
        return {"output": s[::-1]}
    raise ValueError(f"unknown operation: {name}")


def reply_result(req_id: Any, result: Any) -> None:
    sys.stdout.write(
        json.dumps({"jsonrpc": "2.0", "id": req_id, "result": result}) + "\n"
    )
    sys.stdout.flush()


def reply_error(req_id: Any, code: int, message: str) -> None:
    sys.stdout.write(
        json.dumps(
            {
                "jsonrpc": "2.0",
                "id": req_id,
                "error": {"code": code, "message": message},
            }
        )
        + "\n"
    )
    sys.stdout.flush()


def main() -> None:
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        try:
            frame = json.loads(line)
        except json.JSONDecodeError:
            continue

        method = frame.get("method")
        req_id = frame.get("id")
        params = frame.get("params") or {}

        if method == "initialize":
            reply_result(
                req_id,
                {"plugin": "example-python", "operations": OPERATIONS},
            )
        elif method == "invoke":
            try:
                reply_result(req_id, handle_invoke(params))
            except Exception as e:  # noqa: BLE001
                reply_error(req_id, -32000, str(e))
        elif method == "shutdown":
            reply_result(req_id, None)
            return
        else:
            if req_id is not None:
                reply_error(req_id, -32601, f"method not found: {method}")


if __name__ == "__main__":
    main()
