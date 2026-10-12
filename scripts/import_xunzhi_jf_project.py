#!/usr/bin/env python3
"""POST /ImportProject 导入寻智 xunzhi-jf 竖切包（含 alarmTriggers）。

项目已存在（code 4001）时不会更新触发器/点表：请删项目后重导，或运行 scripts/xunzhi_jf_sync_alarm_db.py。

重导后 project_uuid / 默认大屏 display_uid 会变（勿写死旧 uuid）。前端以 ProjectList 名称
「寻智机房环境-xunzhi-jf」解析；浏览器可清 ProjectUuid Cookie 或重新进入该项目。
导入成功后可用 sqlite3 核对: SELECT uuid,name FROM project_lists WHERE name LIKE '%寻智%';
"""
import json
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PKG = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_complete_project_package.json"
BASE = "http://127.0.0.1:8081"


def curl_json(method_path, body=None, token=None, project_uuid=None):
    cmd = ["curl", "-s", "-X", "POST", f"{BASE}{method_path}"]
    if token:
        cmd += ["-H", f"Authorization:{token}"]
    if project_uuid:
        cmd += ["-H", f"ProjectUuid:{project_uuid}"]
    if body is not None:
        cmd += ["-H", "Content-Type:application/json", "-d", json.dumps(body)]
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=120)
    return json.loads(r.stdout or "{}")


def login():
    import hashlib

    pwd_md5 = hashlib.md5(b"123456").hexdigest()
    data = curl_json("/login", {"Username": "admin", "password": pwd_md5})
    if data.get("code") != 1000:
        raise SystemExit(f"login failed: {data}")
    return data["data"]["token"]


def main():
    if not PKG.is_file():
        raise SystemExit(f"missing package: {PKG}")
    token = login()
    with open(PKG, "r", encoding="utf-8") as f:
        pkg = json.load(f)
    # ImportProject 接受完整包（去掉 _meta）
    pkg.pop("_meta", None)
    resp = curl_json("/ImportProject", pkg, token=token)
    print(json.dumps(resp, ensure_ascii=False, indent=2))
    if resp.get("code") not in (0, 1000, 2002):
        sys.exit(1)


if __name__ == "__main__":
    main()
