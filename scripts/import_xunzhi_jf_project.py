#!/usr/bin/env python3
"""POST /ImportProject 导入寻智 xunzhi-jf 竖切包（含 alarmTriggers）。

项目已存在（code 4001）时不会更新组态/点表：请 ``--force``（ProjectDel 后重导）或删项目后重导。

重导后 project_uuid / 默认大屏 display_uid 会变（勿写死旧 uuid）。前端以 ProjectList 名称
「寻智机房环境-xunzhi-jf」解析；浏览器可清 ProjectUuid Cookie 或重新进入该项目。
"""
import argparse
import hashlib
import json
import os
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PKG = ROOT / "projects-import/xunzhi-jf/xunzhi-jf_complete_project_package.json"
PROJECT_NAME = "寻智机房环境-xunzhi-jf"
_HTTP_PORT = os.environ.get("ISM_HTTP_PORT", "8081")
BASE = os.environ.get("BASE_URL", f"http://127.0.0.1:{_HTTP_PORT}").rstrip("/")


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
    pwd_md5 = hashlib.md5(b"123456").hexdigest()
    data = curl_json("/login", {"Username": "admin", "password": pwd_md5})
    if data.get("code") != 1000:
        raise SystemExit(f"login failed: {data}")
    return data["data"]["token"]


def find_project_uuid(token: str) -> str:
    data = curl_json("/ProjectList", {}, token)
    for p in data.get("list") or []:
        pi = p.get("ProjectInfo", p)
        if pi.get("name") == PROJECT_NAME:
            return pi.get("uuid", "") or ""
    return ""


def project_del(token: str, proj_uuid: str) -> dict:
    return curl_json("/ProjectDel", {"uuid": proj_uuid}, token)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument(
        "--force",
        action="store_true",
        help="若已存在同名项目则 ProjectDel 后重导（更新平面 W2/组态绑定必需）",
    )
    args = ap.parse_args()

    if not PKG.is_file():
        raise SystemExit(f"missing package: {PKG}")
    token = login()
    existing = find_project_uuid(token)
    if existing and args.force:
        del_resp = project_del(token, existing)
        print(json.dumps({"project_del": del_resp, "deleted_uuid": existing}, ensure_ascii=False, indent=2))
        if del_resp.get("code") not in (0, 1000, 200):
            raise SystemExit(f"ProjectDel failed: {del_resp}")
        existing = find_project_uuid(token)
        if existing:
            raise SystemExit("ProjectDel 后仍存在同名项目，请手工清理后重试")

    with open(PKG, "r", encoding="utf-8") as f:
        pkg = json.load(f)
    pkg.pop("_meta", None)
    resp = curl_json("/ImportProject", pkg, token=token)
    print(json.dumps(resp, ensure_ascii=False, indent=2))
    if resp.get("code") == 4001:
        print(
            f"提示：项目「{PROJECT_NAME}」已存在，组态不会更新。请: "
            f"python3 scripts/import_xunzhi_jf_project.py --force",
            file=sys.stderr,
        )
    if resp.get("code") not in (0, 1000, 2002):
        sys.exit(1)


if __name__ == "__main__":
    main()
