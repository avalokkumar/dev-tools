#!/usr/bin/env python3
"""End-to-end MCP smoke harness.

Spawns `devforge mcp`, performs the MCP handshake, lists all tools, then
calls each one with a valid argument set and asserts the response shape.
Includes both reference plugins (Go example-hello, Python example-python)
when DEVFORGE_PLUGIN_DIR points at a directory containing them.

Exits non-zero if any tool fails.
"""
from __future__ import annotations

import json
import os
import subprocess
import sys
import time
from typing import Any, Dict


# ---------- expected tool inventory ----------

BUILTINS = [
    # Phase A-D
    "cron_next", "cron_parse",
    "csv_format", "csv_validate",
    "diff_compare",
    "faker_generate", "faker_kinds",
    "json_format", "json_validate",
    "jwt_decode", "jwt_verify",
    "regex_explain", "regex_test",
    "tz_convert", "tz_list",
    "uuid_generate", "uuid_hash",
    "yaml_convert", "yaml_format", "yaml_validate",
    # Phase E1 — encoding
    "enc_base64_encode", "enc_base64_decode",
    "enc_url_encode", "enc_url_decode",
    "enc_html_encode", "enc_html_decode",
    "enc_hex_encode", "enc_hex_decode",
    # Phase E2 — strings + time
    "str_case", "str_diff", "str_stats", "str_sort_unique", "str_replace",
    "time_convert", "time_relative", "time_duration",
    # Phase E3 — sql + md + id + color
    "sql_format", "sql_validate",
    "md_to_html", "md_table_from_csv",
    "id_ulid", "id_slug",
    "color_convert",
    # Phase E4 — data
    "data_json_to_csv", "data_csv_to_json",
    "data_json_to_xml", "data_xml_to_json",
    "data_flatten", "data_unflatten",
    "data_key_rename",
    # Phase E5 — git + devops
    "git_patch", "git_commit_format", "git_ignore_gen",
    "dockerfile_lint", "env_parse", "env_diff", "k8s_validate",
    # Phase E6 + E7 — net
    "url_parse", "headers_analyze", "dns_lookup", "http_request",
    # Phase E8 — crypto + totp
    "crypto_aes_encrypt", "crypto_aes_decrypt",
    "crypto_rsa_keygen", "crypto_hmac",
    "crypto_password_hash", "crypto_password_strength",
    "totp_generate", "totp_verify",
    # Phase E9 — code + math + ip
    "code_fmt_go", "code_fmt_xml", "code_fmt_html",
    "math_eval", "math_unit",
    "ip_calc",
]
PLUGINS = ["hello_say", "py_reverse"]
ALL_TOOLS = sorted(BUILTINS + PLUGINS)


# ---------- per-tool argument fixtures + result probes ----------

# Each entry: (args, predicate over decoded result text -> bool, label)
HMAC_TOKEN = (
    # An HS256 token signed with secret "test-secret-do-not-use",
    # payload {"sub":"alok","exp":4102444800} (Jan 1 2100).
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."
    "eyJzdWIiOiJhbG9rIiwiZXhwIjo0MTAyNDQ0ODAwfQ."
    "zmSWUZIoyEKOSgEFh3hMHAhlAMGMY9UHJ4xQkVTWTA0"
)

CASES: Dict[str, Dict[str, Any]] = {
    "uuid_generate": {
        "args": {"version": 7, "count": 3},
        "check": lambda r: isinstance(r.get("values"), list) and len(r["values"]) == 3,
    },
    "uuid_hash": {
        "args": {"input": "devforge", "algos": ["sha256"]},
        "check": lambda r: r.get("digests", {}).get("sha256")
        == "1508762a1758bc783abc2d0fb5189d3de02c6d6919bc91cb26ff838c70d57e8e",
    },
    "json_format": {
        "args": {"input": '{"b":1,"a":2}', "indent": 2, "sortKeys": True},
        "check": lambda r: '"a": 2' in r["output"] and '"b": 1' in r["output"],
    },
    "json_validate": {
        "args": {"input": '{"ok":true}'},
        "check": lambda r: r.get("valid") is True,
    },
    "yaml_format": {
        "args": {"input": "a: 1\nb: 2\n"},
        "check": lambda r: "a: 1" in r["output"] and "b: 2" in r["output"],
    },
    "yaml_validate": {
        "args": {"input": "a: 1\nb: 2\n"},
        "check": lambda r: r.get("valid") is True,
    },
    "yaml_convert": {
        "args": {"input": "name: alok", "to": "json"},
        "check": lambda r: '"name"' in r["output"] and "alok" in r["output"],
    },
    "csv_format": {
        "args": {"input": "a,b\n1,2\n", "alignColumns": True},
        "check": lambda r: "a, b" in r["output"],
    },
    "csv_validate": {
        "args": {"input": "a,b\n1,2\n", "expectedColumns": ["a", "b"]},
        "check": lambda r: r.get("valid") is True,
    },
    "diff_compare": {
        "args": {
            "left": '{"a":1,"b":2}',
            "right": '{"a":99,"c":3}',
            "mode": "json",
        },
        "check": lambda r: r.get("summary", {}).get("changes") == 1
        and r["summary"]["adds"] == 1
        and r["summary"]["removes"] == 1,
    },
    "regex_test": {
        "args": {"pattern": r"\d+", "input": "build 42 release 7"},
        "check": lambda r: len(r.get("matches", [])) == 2
        and r["matches"][0]["value"] == "42",
    },
    "regex_explain": {
        "args": {"pattern": r"^\d{4}-\d{2}$"},
        "check": lambda r: any(n["token"] == "^" for n in r.get("tree", [])),
    },
    "cron_parse": {
        "args": {"expression": "*/5 * * * *"},
        "check": lambda r: len(r.get("fields", [])) == 5
        and r["fields"][0]["value"] == "*/5",
    },
    "cron_next": {
        "args": {
            "expression": "0 12 * * *",
            "n": 2,
            "from": "2026-05-09T00:00:00Z",
            "tz": "UTC",
        },
        "check": lambda r: len(r.get("runs", [])) == 2,
    },
    "jwt_decode": {
        "args": {"token": HMAC_TOKEN},
        "check": lambda r: r.get("header", {}).get("alg") == "HS256"
        and r["payload"]["sub"] == "alok",
    },
    "jwt_verify": {
        "args": {
            "token": HMAC_TOKEN,
            "key": "test-secret-do-not-use",
            "expectedAlgs": ["HS256"],
        },
        "check": lambda r: r.get("valid") is True,
    },
    "tz_convert": {
        "args": {
            "time": "2026-05-09T12:00:00Z",
            "fromTZ": "America/New_York",
            "toTZ": "Asia/Tokyo",
        },
        # NY May 12:00 (UTC-4 EDT) → Tokyo +13h → 01:00 next day.
        "check": lambda r: r.get("converted", "").startswith("2026-05-10T01:00"),
    },
    "tz_list": {
        "args": {"filter": "tokyo"},
        "check": lambda r: isinstance(r, list)
        and any(z["name"] == "Asia/Tokyo" for z in r),
    },
    "faker_generate": {
        "args": {
            "spec": {
                "fields": [
                    {"name": "id", "kind": "sequence"},
                    {"name": "email", "kind": "email"},
                ]
            },
            "count": 3,
            "seed": 42,
            "format": "csv",
        },
        "check": lambda r: r.get("output", "").startswith("id,email")
        and r["output"].count("\n") == 4,  # header + 3 rows
    },
    "faker_kinds": {
        "args": {},
        "check": lambda r: len(r.get("kinds", [])) >= 15
        and len(r.get("locales", [])) >= 5,
    },
    "hello_say": {
        "args": {"name": "alok"},
        "check": lambda r: "hello, alok" in r.get("message", ""),
    },
    "py_reverse": {
        "args": {"input": "DevForge"},
        "check": lambda r: r.get("output") == "egroFveD",
    },
    # ---- Phase E1 ----
    "enc_base64_encode": {"args": {"input": "foo"}, "check": lambda r: r.get("output") == "Zm9v"},
    "enc_base64_decode": {"args": {"input": "Zm9v"}, "check": lambda r: r.get("output") == "foo"},
    "enc_url_encode":    {"args": {"input": "a b&c"}, "check": lambda r: r.get("output") == "a+b%26c"},
    "enc_url_decode":    {"args": {"input": "a+b%26c"}, "check": lambda r: r.get("output") == "a b&c"},
    "enc_html_encode":   {"args": {"input": "<a>"}, "check": lambda r: "&lt;" in r.get("output", "")},
    "enc_html_decode":   {"args": {"input": "&lt;a&gt;"}, "check": lambda r: r.get("output") == "<a>"},
    "enc_hex_encode":    {"args": {"input": "abc"}, "check": lambda r: r.get("output") == "616263"},
    "enc_hex_decode":    {"args": {"input": "616263"}, "check": lambda r: r.get("output") == "abc"},
    # ---- Phase E2 ----
    "str_case":        {"args": {"input": "hello world", "mode": "snake"}, "check": lambda r: r.get("output") == "hello_world"},
    "str_diff":        {"args": {"left": "a\nb", "right": "a\nB"}, "check": lambda r: r.get("summary", {}).get("adds") == 1},
    "str_stats":       {"args": {"input": "hello world"}, "check": lambda r: r.get("words") == 2},
    "str_sort_unique": {"args": {"input": "c\na\nb", "unique": True}, "check": lambda r: r.get("output") == "a\nb\nc"},
    "str_replace":     {"args": {"input": "foo", "pattern": "o", "replacement": "0"}, "check": lambda r: r.get("output") == "f00"},
    "time_convert":    {"args": {"input": "1715184000"}, "check": lambda r: r.get("epochS") == 1715184000},
    "time_relative":   {"args": {"from": "2026-05-09T12:00:00Z", "to": "2026-05-09T12:05:00Z"}, "check": lambda r: "from now" in r.get("phrase", "")},
    "time_duration":   {"args": {"input": "1h30m"}, "check": lambda r: r.get("hours") == 1 and r.get("minutes") == 30},
    # ---- Phase E3 ----
    "sql_format":         {"args": {"input": "select * from t"}, "check": lambda r: "SELECT" in r.get("output", "")},
    "sql_validate":       {"args": {"input": "SELECT * FROM users"}, "check": lambda r: any(d.get("code") == "SQL.LINT.SELECT_STAR" for d in r.get("diagnostics", []))},
    "md_to_html":         {"args": {"input": "# Hi"}, "check": lambda r: "<h1" in r.get("output", "")},
    "md_table_from_csv":  {"args": {"input": "a,b\n1,2"}, "check": lambda r: "| a | b |" in r.get("output", "")},
    "id_ulid":            {"args": {"count": 3}, "check": lambda r: len(r.get("values", [])) == 3},
    "id_slug":            {"args": {"input": "Hello, World!"}, "check": lambda r: r.get("output") == "hello-world"},
    "color_convert":      {"args": {"input": "#ff0000"}, "check": lambda r: r.get("hex") == "#ff0000"},
    # ---- Phase E4 ----
    "data_json_to_csv":  {"args": {"input": '[{"a":1,"b":2}]'}, "check": lambda r: "a,b" in r.get("output", "")},
    "data_csv_to_json":  {"args": {"input": "a,b\n1,2"}, "check": lambda r: '"a"' in r.get("output", "")},
    "data_json_to_xml":  {"args": {"input": '{"x":1}'}, "check": lambda r: "<x>1</x>" in r.get("output", "")},
    "data_xml_to_json":  {"args": {"input": "<root><x>1</x></root>"}, "check": lambda r: '"x"' in r.get("output", "")},
    "data_flatten":      {"args": {"input": '{"a":{"b":1}}'}, "check": lambda r: '"a.b"' in r.get("output", "")},
    "data_unflatten":    {"args": {"input": '{"a.b":1}'}, "check": lambda r: '"b"' in r.get("output", "")},
    "data_key_rename":   {"args": {"input": '{"old":1}', "rules": [{"from": "old", "to": "new"}]}, "check": lambda r: '"new"' in r.get("output", "")},
    # ---- Phase E5 ----
    "git_patch":         {"args": {"left": "a\n", "right": "b\n"}, "check": lambda r: "+b" in r.get("output", "")},
    "git_commit_format": {"args": {"input": "feat(x): hi"}, "check": lambda r: r.get("type") == "feat"},
    "git_ignore_gen":    {"args": {"templates": ["go"]}, "check": lambda r: "*.exe" in r.get("output", "")},
    "dockerfile_lint":   {"args": {"input": "FROM nginx:latest\n"}, "check": lambda r: any(d.get("code") == "DOCKER.LINT.LATEST_TAG" for d in r.get("diagnostics", []))},
    "env_parse":         {"args": {"input": "FOO=1"}, "check": lambda r: r.get("values", {}).get("FOO") == "1"},
    "env_diff":          {"args": {"left": "A=1", "right": "A=2"}, "check": lambda r: any(c.get("key") == "A" for c in r.get("changed", []))},
    "k8s_validate":      {"args": {"input": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: x\n"}, "check": lambda r: r.get("valid") is True},
    # ---- Phase E6 + E7 ----
    "url_parse":       {"args": {"input": "https://example.com/p?a=1"}, "check": lambda r: r.get("hostname") == "example.com"},
    "headers_analyze": {"args": {"headers": {"X-Content-Type-Options": "nosniff"}}, "check": lambda r: any(f.get("ok") for f in r.get("findings", []))},
    "dns_lookup":      {"args": {"host": "127.0.0.1", "type": "ptr", "allowPrivate": True}, "check": lambda r: r.get("host") == "127.0.0.1"},
    "http_request":    {"args": {"url": "http://localhost:9/x"}, "check": lambda r: any(d.get("code") == "HTTP.PRIVATE_BLOCKED" for d in r.get("diagnostics", []))},
    # ---- Phase E8 ----
    "crypto_aes_encrypt":      {"args": {"plaintext": "hi", "passphrase": "k", "pbkdf2Iters": 1000}, "check": lambda r: bool(r.get("output"))},
    "crypto_aes_decrypt":      {"args": {"input": "xxx", "passphrase": "k"}, "check": lambda r: any(d.get("code", "").startswith("CRYPTO.AES") for d in r.get("diagnostics", []))},
    "crypto_rsa_keygen":       {"args": {"bits": 2048}, "check": lambda r: "PRIVATE KEY" in r.get("privatePem", "")},
    "crypto_hmac":             {"args": {"input": "Hi There", "key": "0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b", "keyEncoding": "hex", "algorithm": "sha256"}, "check": lambda r: r.get("output") == "b0344c61d8db38535ca8afceaf0bf12b881dc200c9833da726e9376c2e32cff7"},
    "crypto_password_hash":    {"args": {"password": "hunter2", "bcryptCost": 4}, "check": lambda r: r.get("hash", "").startswith("$2")},
    "crypto_password_strength":{"args": {"password": "Abcdefgh1!"}, "check": lambda r: r.get("score", 0) >= 2},
    "totp_generate":           {"args": {"secret": "12345678901234567890", "secretEncoding": "raw", "digits": 8, "at": 59}, "check": lambda r: r.get("code") == "94287082"},
    "totp_verify":             {"args": {"code": "94287082", "secret": "12345678901234567890", "secretEncoding": "raw", "digits": 8, "at": 59}, "check": lambda r: r.get("valid") is True},
    # ---- Phase E9 ----
    "code_fmt_go":   {"args": {"input": "package x\nfunc f( ){}\n"}, "check": lambda r: "func f()" in r.get("output", "")},
    "code_fmt_xml":  {"args": {"input": "<a><b>1</b></a>"}, "check": lambda r: "  <b>" in r.get("output", "")},
    "code_fmt_html": {"args": {"input": "<p>x</p>"}, "check": lambda r: "<p>x</p>" in r.get("output", "")},
    "math_eval":     {"args": {"expression": "(2+3)*4"}, "check": lambda r: r.get("value") == 20},
    "math_unit":     {"args": {"value": 1, "from": "gib", "to": "mib"}, "check": lambda r: r.get("value") == 1024},
    "ip_calc":       {"args": {"cidr": "192.168.1.0/24"}, "check": lambda r: r.get("usableHosts") == "254"},
}


# ---------- MCP driver ----------

class MCP:
    def __init__(self, binary: str, plugin_dir: str | None) -> None:
        env = os.environ.copy()
        if plugin_dir:
            env["DEVFORGE_PLUGIN_DIR"] = plugin_dir
        self.proc = subprocess.Popen(
            [binary, "mcp"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=env,
            bufsize=0,
        )
        self.next_id = 0

    def _send(self, method: str, params: Any = None) -> Dict[str, Any]:
        self.next_id += 1
        msg: Dict[str, Any] = {
            "jsonrpc": "2.0",
            "id": self.next_id,
            "method": method,
        }
        if params is not None:
            msg["params"] = params
        line = (json.dumps(msg) + "\n").encode()
        assert self.proc.stdin
        self.proc.stdin.write(line)
        self.proc.stdin.flush()
        return self._recv()

    def _recv(self) -> Dict[str, Any]:
        assert self.proc.stdout
        deadline = time.time() + 10
        while time.time() < deadline:
            line = self.proc.stdout.readline()
            if not line:
                break
            try:
                return json.loads(line.decode())
            except json.JSONDecodeError:
                continue
        raise RuntimeError("no MCP response (timeout or EOF)")

    def initialize(self) -> Dict[str, Any]:
        return self._send(
            "initialize",
            {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {"name": "mcp-smoke", "version": "0"},
            },
        )

    def list_tools(self) -> list[str]:
        r = self._send("tools/list")
        return sorted(t["name"] for t in r["result"]["tools"])

    def call(self, name: str, args: Dict[str, Any]) -> Any:
        r = self._send("tools/call", {"name": name, "arguments": args})
        if "error" in r:
            raise RuntimeError(f"{name}: {r['error']}")
        content = r["result"]["content"]
        # mcp-go returns an array of {type:'text', text:<json>}
        text = content[0]["text"]
        try:
            return json.loads(text)
        except json.JSONDecodeError:
            return text

    def close(self) -> None:
        try:
            assert self.proc.stdin
            self.proc.stdin.close()
        except Exception:
            pass
        try:
            self.proc.wait(timeout=2)
        except subprocess.TimeoutExpired:
            self.proc.kill()


# ---------- main ----------

def main() -> int:
    binary = "./bin/devforge"
    plugin_dir = sys.argv[1] if len(sys.argv) > 1 else None

    print(f"binary    : {binary}")
    print(f"plugin dir: {plugin_dir or '(none)'}")
    print()

    mcp = MCP(binary, plugin_dir)
    try:
        # 1. Handshake
        init = mcp.initialize()
        srv = init["result"]["serverInfo"]
        print(f"[handshake] serverInfo = {srv}")
        if srv["name"] != "devforge":
            print("FAIL: serverInfo.name != devforge")
            return 1

        # 2. List tools
        listed = mcp.list_tools()
        print(f"[tools/list] {len(listed)} tools registered")
        missing = [t for t in ALL_TOOLS if t not in listed]
        extra = [t for t in listed if t not in ALL_TOOLS]
        if missing:
            print(f"FAIL: missing tools: {missing}")
            return 1
        if extra:
            print(f"NOTE: extra tools (probably new builtins): {extra}")
        print()

        # 3. Call each
        passed = 0
        failed = 0
        for name in ALL_TOOLS:
            if name not in CASES:
                print(f"[skip] {name} — no fixture")
                continue
            spec = CASES[name]
            try:
                res = mcp.call(name, spec["args"])
                ok = bool(spec["check"](res))
            except Exception as e:  # noqa: BLE001
                print(f"[FAIL] {name:18s} → exception: {e}")
                failed += 1
                continue
            if ok:
                print(f"[PASS] {name:18s}")
                passed += 1
            else:
                print(f"[FAIL] {name:18s} → unexpected result: {json.dumps(res)[:160]}")
                failed += 1

        print()
        print(f"summary: {passed} passed, {failed} failed, "
              f"{len(ALL_TOOLS)} total")
        return 0 if failed == 0 else 1
    finally:
        mcp.close()


if __name__ == "__main__":
    sys.exit(main())
