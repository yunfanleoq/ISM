#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""解压数据点位转发表三个压缩包(GBK 文件名解码)，并把 .xls 批量转 .xlsx。

产物布局:
  数据点位转发表/_unzipped/A/<...>.xlsx        (excel(1).zip, A 系列生产楼)
  数据点位转发表/_unzipped/模块/<...>.xlsx     (20260613.zip, B 系列机房模块)
  数据点位转发表/_unzipped/配电室/<...>.xlsx   (TB-配电室.zip, B 系列配电室)
"""
import os, sys, zipfile, shutil, subprocess, glob

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, "数据点位转发表")
OUT = os.path.join(SRC, "_unzipped")

ZIPS = {
    "A": "excel(1).zip",
    "模块": "20260613.zip",
    "配电室": "TB-配电室.zip",
}
SOFFICE = "/opt/homebrew/bin/soffice"


def decode_name(raw: str) -> str:
    try:
        return raw.encode("cp437").decode("gbk")
    except Exception:
        return raw


def extract():
    for sub, zf_name in ZIPS.items():
        dest = os.path.join(OUT, sub)
        os.makedirs(dest, exist_ok=True)
        zpath = os.path.join(SRC, zf_name)
        n = 0
        with zipfile.ZipFile(zpath) as zf:
            for info in zf.infolist():
                name = decode_name(info.filename)
                base = os.path.basename(name)
                if not base or info.is_dir():
                    continue
                if base.startswith("~$"):  # excel 临时锁文件
                    continue
                low = base.lower()
                if not (low.endswith(".xls") or low.endswith(".xlsx")):
                    continue
                with zf.open(info) as f, open(os.path.join(dest, base), "wb") as o:
                    o.write(f.read())
                n += 1
        print(f"[extract] {sub:6s} <- {zf_name}: {n} files")


def convert_xls():
    """把所有 .xls 转成 .xlsx(soffice headless)。"""
    xls_files = []
    for sub in ZIPS:
        xls_files += glob.glob(os.path.join(OUT, sub, "*.xls"))
    if not xls_files:
        print("[convert] no .xls to convert")
        return
    print(f"[convert] converting {len(xls_files)} .xls files via soffice ...")
    for x in xls_files:
        d = os.path.dirname(x)
        r = subprocess.run([SOFFICE, "--headless", "--convert-to", "xlsx",
                            "--outdir", d, x],
                           capture_output=True, text=True)
        ok = os.path.exists(os.path.splitext(x)[0] + ".xlsx")
        print(f"   {'OK ' if ok else 'ERR'} {os.path.basename(x)}")
        if not ok:
            print("       ", r.stdout.strip(), r.stderr.strip())


def summary():
    print("\n[summary]")
    for sub in ZIPS:
        d = os.path.join(OUT, sub)
        xlsx = sorted(glob.glob(os.path.join(d, "*.xlsx")))
        xls = sorted(glob.glob(os.path.join(d, "*.xls")))
        print(f"  {sub:6s}: {len(xlsx)} xlsx, {len(xls)} xls(原始)")


if __name__ == "__main__":
    if "--clean" in sys.argv:
        shutil.rmtree(OUT, ignore_errors=True)
    extract()
    convert_xls()
    summary()
    print("\n[done] ->", OUT)
