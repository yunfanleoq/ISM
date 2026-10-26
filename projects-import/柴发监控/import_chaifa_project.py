#!/usr/bin/env python3
"""
将「后沙峪改造-柴发部分」ISM 项目包导入到正在运行的 ISM 系统。

会调用 POST /ImportProject，新建一个独立项目（与现有项目互不影响）。

用法:
  python3 import_chaifa_project.py
  python3 import_chaifa_project.py --base-url http://127.0.0.1:8081 \\
      --user admin --password 123456 \\
      --package ./后沙峪改造-柴发部分_ISM项目包.json
"""

from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
import sys
from pathlib import Path


def md5_hex(s: str) -> str:
    return hashlib.md5(s.encode("utf-8")).hexdigest()


def curl_json(method: str, url: str, headers: dict, body: bytes | None = None, timeout: int = 600) -> dict:
    cmd = ["curl", "-sS", "-X", method, url, "--max-time", str(timeout)]
    for k, v in headers.items():
        cmd += ["-H", f"{k}: {v}"]
    if body is not None:
        cmd += ["--data-binary", "@-"]
    r = subprocess.run(cmd, input=body, capture_output=True)
    if r.returncode != 0:
        raise RuntimeError(r.stderr.decode("utf-8", errors="replace") or "curl failed")
    text = r.stdout.decode("utf-8", errors="replace")
    try:
        return json.loads(text)
    except json.JSONDecodeError as e:
        raise RuntimeError(f"响应非 JSON: {text[:500]}") from e


def login(base: str, user: str, password: str) -> tuple[str, str]:
    """返回 (token, user_uuid)。password 为原始密码，内部 MD5。"""
    pwd_md5 = md5_hex(password)
    resp = curl_json(
        "POST",
        f"{base}/login",
        {"Content-Type": "application/json"},
        json.dumps({"Username": user, "password": pwd_md5}, ensure_ascii=False).encode("utf-8"),
        timeout=30,
    )
    code = resp.get("code")
    if code != 1000:
        raise RuntimeError(f"登录失败 code={code} resp={resp}")
    data = resp.get("data") or {}
    token = data.get("token") or ""
    uuid = data.get("uuid") or data.get("Uuid") or ""
    if not token:
        raise RuntimeError(f"登录成功但 token 为空: {resp}")
    return token, uuid


def main() -> int:
    ap = argparse.ArgumentParser(description="导入柴发 ISM 项目包（新建独立项目）")
    ap.add_argument("--base-url", default="http://127.0.0.1:8081")
    ap.add_argument("--user", default="admin")
    ap.add_argument("--password", default="123456")
    ap.add_argument(
        "--package",
        default=str(Path(__file__).with_name("后沙峪改造-柴发部分_ISM项目包.json")),
    )
    args = ap.parse_args()

    pkg_path = Path(args.package).expanduser().resolve()
    if not pkg_path.exists():
        print(f"找不到项目包: {pkg_path}", file=sys.stderr)
        return 1

    with open(pkg_path, "rb") as f:
        raw = f.read()
    try:
        pkg = json.loads(raw)
    except json.JSONDecodeError as e:
        print(f"项目包 JSON 无效: {e}", file=sys.stderr)
        return 1

    st = pkg.get("statistics") or {}
    print(f"项目包: {pkg_path.name}")
    print(f"  项目名: {pkg.get('project', {}).get('name')}")
    print(
        f"  模型={st.get('deviceModels')} 组={st.get('registerGroups')} "
        f"点={st.get('registerPoints')} 设备={st.get('devices')}"
    )
    print(f"登录 {args.base_url} 用户={args.user} ...")
    token, _ = login(args.base_url.rstrip("/"), args.user, args.password)
    print("登录成功，开始导入（数据点较多，可能需要数分钟）...")

    # 注意：Authorization 不要加 Bearer 前缀
    resp = curl_json(
        "POST",
        f"{args.base_url.rstrip('/')}/ImportProject",
        {
            "Authorization": token,
            "Content-Type": "application/json",
        },
        raw,
        timeout=1800,
    )
    print(json.dumps(resp, ensure_ascii=False, indent=2))
    if resp.get("code") != 0:
        print("导入失败", file=sys.stderr)
        return 2

    data = resp.get("data") or {}
    print()
    print("导入成功（已新建独立项目）")
    print(f"  project_uuid = {data.get('project_uuid')}")
    print(f"  模型={data.get('model_count')} 数据点={data.get('point_count')} "
          f"设备={data.get('device_count')} 树节点={data.get('tree_count')}")
    print("请用 admin 刷新「项目列表」，进入「后沙峪改造-柴发部分」即可。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
