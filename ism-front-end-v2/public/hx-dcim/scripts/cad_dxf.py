# -*- coding: utf-8 -*-
"""LibreDWG DXF 组码读取（只读）。ezdxf 对这份 DXF 会崩，所以自己解析。"""
import re

U_RE = re.compile(r"\\U\+([0-9A-Fa-f]{4})")
MTEXT_CODE = re.compile(r"\\[A-Za-z][^;\\]*;")


def decode_text(s):
    s2 = U_RE.sub(lambda m: chr(int(m.group(1), 16)), s)
    raw = bytearray()
    out = []

    def flush():
        if not raw:
            return
        for enc in ("gbk", "gb18030"):
            try:
                out.append(raw.decode(enc))
                raw.clear()
                return
            except Exception:
                pass
        out.append(raw.decode("latin1", "replace"))
        raw.clear()

    for ch in s2:
        o = ord(ch)
        if o < 128:
            flush()
            out.append(ch)
        elif o <= 255:
            raw.append(o)
        else:
            flush()
            out.append(ch)
    flush()
    return "".join(out)


def strip_mtext(s):
    s = s.replace("\\P", "\n").replace("\\~", " ")
    s = MTEXT_CODE.sub("", s)
    s = s.replace("{", "").replace("}", "")
    return re.sub(r"[ \t]+", " ", s).strip()


def _pairs(path):
    with open(path, "r", encoding="latin-1") as f:
        while True:
            c = f.readline()
            if not c:
                return
            v = f.readline()
            if not v:
                return
            cs = c.strip()
            if not cs or not (cs[0].isdigit() or (cs[0] == "-" and cs[1:].isdigit())):
                continue
            yield int(cs), v.rstrip("\r\n")


def load_model(path):
    """模型空间 LINE / LWPOLYLINE / INSERT / TEXT / MTEXT。"""
    section = None
    ent = None
    rec = {}
    lines = []
    inserts = []
    texts = []

    def commit():
        nonlocal ent, rec
        if section != "ENTITIES" or not ent:
            ent = None
            rec = {}
            return
        if ent == "LINE" and all(k in rec for k in ("x1", "y1", "x2", "y2")):
            lines.append(rec)
        elif ent == "LWPOLYLINE" and len(rec.get("pts") or []) >= 2:
            pts = rec["pts"]
            col = rec.get("color")
            for i in range(1, len(pts)):
                x1, y1 = pts[i - 1]
                x2, y2 = pts[i]
                if x1 is None or y1 is None or x2 is None or y2 is None:
                    continue
                lines.append({"x1": x1, "y1": y1, "x2": x2, "y2": y2, "color": col, "lw": 1})
        elif ent == "INSERT" and "x" in rec and rec.get("name"):
            inserts.append(rec)
        elif ent in ("TEXT", "MTEXT") and rec.get("_t"):
            s = strip_mtext(decode_text("".join(rec["_t"])))
            if s:
                texts.append({
                    "x": rec.get("x"), "y": rec.get("y"), "h": rec.get("h") or 0,
                    "s": s, "kind": ent,
                })
        ent = None
        rec = {}

    for code, val in _pairs(path):
        if code == 0:
            commit()
            name = val.strip()
            if name == "SECTION":
                ent = "SECTION"
                rec = {}
                continue
            if name == "ENDSEC":
                section = None
                ent = None
                rec = {}
                continue
            ent = name
            rec = {"_t": [], "pts": []}
            continue
        if ent == "SECTION" and code == 2:
            section = val.strip()
            ent = None
            continue
        if not ent or ent == "SECTION":
            continue
        try:
            if code == 10:
                fv = float(val)
                if ent == "LINE":
                    rec["x1"] = fv
                elif ent == "LWPOLYLINE":
                    rec["pts"].append([fv, None])
                elif "x" not in rec:
                    rec["x"] = fv
            elif code == 20:
                fv = float(val)
                if ent == "LINE":
                    rec["y1"] = fv
                elif ent == "LWPOLYLINE" and rec["pts"]:
                    rec["pts"][-1][1] = fv
                elif "y" not in rec:
                    rec["y"] = fv
            elif code == 11 and ent == "LINE":
                rec["x2"] = float(val)
            elif code == 21 and ent == "LINE":
                rec["y2"] = float(val)
            elif code == 62:
                rec["color"] = int(float(val))
            elif code == 2 and ent == "INSERT":
                rec["name"] = val.strip()
            elif code in (1, 3) and ent in ("TEXT", "MTEXT"):
                rec["_t"].append(val)
            elif code == 40 and ent in ("TEXT", "MTEXT"):
                rec["h"] = float(val)
            elif code == 41 and ent == "INSERT":
                rec["sx"] = float(val)
        except Exception:
            pass
    commit()
    return {"lines": lines, "inserts": inserts, "texts": texts}
