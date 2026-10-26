#!/usr/bin/env python3
"""Upload a local file to Baidu Pan via cookie API, then create a share link."""
from __future__ import annotations

import hashlib
import json
import math
import os
import ssl
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

SLICE = 4 * 1024 * 1024  # 4MB
UA = "netdisk;11.4.0;BaiduPCS-Go"
CTX = ssl.create_default_context()


def env(name: str) -> str:
    v = os.environ.get(name, "").strip()
    if not v:
        raise SystemExit(f"missing env {name}")
    return v


class Pan:
    def __init__(self, bduss: str, stoken: str, bdstoken: str):
        self.cookie = f"BDUSS={bduss}; STOKEN={stoken}"
        self.bdstoken = bdstoken

    def _req(self, url: str, data: dict | None = None, raw: bytes | None = None, headers: dict | None = None, timeout: int = 120):
        hdrs = {
            "Cookie": self.cookie,
            "User-Agent": UA,
        }
        if headers:
            hdrs.update(headers)
        body = None
        if raw is not None:
            body = raw
        elif data is not None:
            body = urllib.parse.urlencode(data).encode()
            hdrs.setdefault("Content-Type", "application/x-www-form-urlencoded")
        req = urllib.request.Request(url, data=body, headers=hdrs)
        try:
            with urllib.request.urlopen(req, timeout=timeout, context=CTX) as r:
                return json.loads(r.read().decode())
        except urllib.error.HTTPError as e:
            txt = e.read().decode("utf-8", "ignore")
            raise RuntimeError(f"HTTP {e.code} {url}: {txt[:500]}") from e

    def mkdir(self, path: str) -> None:
        url = f"https://pan.baidu.com/api/create?a=commit&bdstoken={self.bdstoken}"
        j = self._req(url, {"path": path, "isdir": "1", "block_list": "[]"})
        # errno 0 ok, -8 already exists
        if j.get("errno") not in (0, -8):
            print("mkdir warn:", j)

    def precreate(self, path: str, size: int, block_list: list[str], content_md5: str, slice_md5: str) -> dict:
        url = f"https://pan.baidu.com/api/precreate?bdstoken={self.bdstoken}"
        data = {
            "path": path,
            "size": str(size),
            "isdir": "0",
            "autoinit": "1",
            "rtype": "1",
            "block_list": json.dumps(block_list),
            "content-md5": content_md5,
            "slice-md5": slice_md5,
        }
        return self._req(url, data)

    def upload_part(self, path: str, uploadid: str, partseq: int, chunk: bytes) -> dict:
        q = urllib.parse.urlencode(
            {
                "method": "upload",
                "type": "tmpfile",
                "app_id": "250528",
                "path": path,
                "uploadid": uploadid,
                "partseq": str(partseq),
            }
        )
        # multipart
        boundary = f"----BaiduPCS{int(time.time()*1000)}"
        body = (
            f"--{boundary}\r\n"
            f'Content-Disposition: form-data; name="file"; filename="blob"\r\n'
            f"Content-Type: application/octet-stream\r\n\r\n"
        ).encode() + chunk + f"\r\n--{boundary}--\r\n".encode()
        url = f"https://c3.pcs.baidu.com/rest/2.0/pcs/superfile2?{q}"
        return self._req(
            url,
            raw=body,
            headers={"Content-Type": f"multipart/form-data; boundary={boundary}"},
            timeout=300,
        )

    def create(self, path: str, size: int, uploadid: str, block_list: list[str]) -> dict:
        url = f"https://pan.baidu.com/api/create?isdir=0&bdstoken={self.bdstoken}"
        data = {
            "path": path,
            "size": str(size),
            "uploadid": uploadid,
            "block_list": json.dumps(block_list),
            "rtype": "1",
            "isdir": "0",
        }
        return self._req(url, data)

    def share(self, fs_id: int, pwd: str = "ism1", period: int = 7) -> dict:
        url = f"https://pan.baidu.com/share/set?channel=chunlei&clienttype=0&web=1&bdstoken={self.bdstoken}"
        data = {
            "schannel": "4",
            "channel_list": "[]",
            "period": str(period),
            "pwd": pwd,
            "fid_list": json.dumps([fs_id]),
        }
        return self._req(url, data)


def file_hashes(path: Path) -> tuple[list[str], str, str, int]:
    size = path.stat().st_size
    blocks: list[str] = []
    h_all = hashlib.md5()
    h_slice = hashlib.md5()
    slice_left = SLICE
    with path.open("rb") as f:
        while True:
            chunk = f.read(SLICE)
            if not chunk:
                break
            blocks.append(hashlib.md5(chunk).hexdigest())
            h_all.update(chunk)
            if slice_left > 0:
                take = chunk[:slice_left]
                h_slice.update(take)
                slice_left -= len(take)
    return blocks, h_all.hexdigest(), h_slice.hexdigest(), size


def main() -> int:
    local = Path(sys.argv[1]).resolve()
    remote_dir = sys.argv[2] if len(sys.argv) > 2 else "/来自：本地电脑/cursorProjects/ISM源码/交付物"
    pwd = sys.argv[3] if len(sys.argv) > 3 else "ism1"
    if not local.is_file():
        print("file not found", local)
        return 1

    pan = Pan(env("BAIDU_BDUSS"), env("BAIDU_STOKEN"), Path("/tmp/baidu_bdstoken.txt").read_text().strip())
    remote_path = remote_dir.rstrip("/") + "/" + local.name
    print(f"local={local} size={local.stat().st_size}")
    print(f"remote={remote_path}")

    pan.mkdir(remote_dir)
    print("computing md5 blocks...")
    t0 = time.time()
    blocks, content_md5, slice_md5, size = file_hashes(local)
    print(f"blocks={len(blocks)} content_md5={content_md5} slice_md5={slice_md5} took={time.time()-t0:.1f}s")

    print("precreate...")
    pre = pan.precreate(remote_path, size, blocks, content_md5, slice_md5)
    print("precreate:", {k: pre.get(k) for k in ("errno", "return_type", "uploadid", "block_list")})
    if pre.get("errno") != 0:
        print("precreate failed", pre)
        return 2

    # return_type 2 = rapid upload / already exists
    if pre.get("return_type") == 2:
        print("rapid upload ok")
        fs_id = pre.get("info", {}).get("fs_id") or pre.get("fs_id")
        if not fs_id:
            # fall through to create without uploadid sometimes returns path info
            info = pre.get("info") or {}
            fs_id = info.get("fs_id")
        if fs_id:
            sh = pan.share(int(fs_id), pwd=pwd, period=7)
            print("SHARE_RESULT", json.dumps(sh, ensure_ascii=False))
            link = sh.get("link") or sh.get("shorturl")
            print("SHARE_LINK", link)
            print("SHARE_PWD", pwd)
            Path("/tmp/baidu_share_result.json").write_text(json.dumps({"link": link, "pwd": pwd, "raw": sh, "path": remote_path}, ensure_ascii=False, indent=2))
            return 0

    uploadid = pre["uploadid"]
    need = pre.get("block_list")
    if need is None or need == []:
        need = list(range(len(blocks)))
    print(f"need upload parts: {len(need)}")

    with local.open("rb") as f:
        for i, partseq in enumerate(need):
            partseq = int(partseq)
            f.seek(partseq * SLICE)
            chunk = f.read(SLICE)
            for attempt in range(5):
                try:
                    r = pan.upload_part(remote_path, uploadid, partseq, chunk)
                    if r.get("md5") or r.get("errno") in (0, None):
                        break
                    print(f"part {partseq} unexpected {r}, retry {attempt}")
                except Exception as e:
                    print(f"part {partseq} err {e}, retry {attempt}")
                    time.sleep(1.5 * (attempt + 1))
            else:
                print("part failed permanently", partseq)
                return 3
            if (i + 1) % 10 == 0 or i == 0 or i + 1 == len(need):
                print(f"uploaded {i+1}/{len(need)} parts ({(i+1)*SLICE/1024/1024:.0f}MB approx)")

    print("create...")
    created = pan.create(remote_path, size, uploadid, blocks)
    print("create:", created)
    if created.get("errno") != 0:
        return 4
    fs_id = created.get("fs_id")
    sh = pan.share(int(fs_id), pwd=pwd, period=7)
    print("SHARE_RESULT", json.dumps(sh, ensure_ascii=False))
    link = sh.get("link") or sh.get("shorturl")
    print("SHARE_LINK", link)
    print("SHARE_PWD", pwd)
    Path("/tmp/baidu_share_result.json").write_text(
        json.dumps({"link": link, "pwd": pwd, "raw": sh, "path": remote_path, "fs_id": fs_id, "size": size}, ensure_ascii=False, indent=2)
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
