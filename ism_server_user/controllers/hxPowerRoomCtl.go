package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// 一次图编辑的保存：把这一间的位移写回 public/hx-dcim/data/hx-power-rooms/<室>.json。
// 只改动到的数字和字符串，不重排整份文件。

type HxPowerRoomController struct {
	beego.Controller
}

type hxNumPatch struct {
	I int      `json:"i"`
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
	S      *string  `json:"s"`
	W      *float64 `json:"w"`
	H      *float64 `json:"h"`
	Fs     *float64 `json:"fs"`
	C      *string  `json:"c"`
	CLock  *bool    `json:"cLock"`
	FsLock *bool    `json:"fsLock"`
}

type hxWirePatch struct {
	I   int         `json:"i"`
	Pts [][]float64 `json:"pts"`
}

type hxRoomPatch struct {
	Room   string        `json:"room"`
	Texts  []hxNumPatch  `json:"texts"`
	Knives []hxNumPatch  `json:"knives"`
	Boxes  []hxNumPatch  `json:"boxes"`
	Wires  []hxWirePatch `json:"wires"`
	Remove struct {
		Texts  []int `json:"texts"`
		Knives []int `json:"knives"`
		Boxes  []int `json:"boxes"`
		Wires  []int `json:"wires"`
	} `json:"remove"`
	Add struct {
		Texts  []json.RawMessage `json:"texts"`
		Knives []json.RawMessage `json:"knives"`
		Boxes  []json.RawMessage `json:"boxes"`
		Wires  []json.RawMessage `json:"wires"`
	} `json:"add"`
}

var hxRoomCodeRe = regexp.MustCompile(`^[0-9A-Za-z]+$`)

func (c *HxPowerRoomController) Save() {
	var patch hxRoomPatch
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &patch); err != nil {
		c.hxFail("请求不是一次图补丁")
		return
	}
	path, err := hxRoomFile(patch.Room)
	if err != nil {
		c.hxFail(err.Error())
		return
	}
	src, err := os.ReadFile(path)
	if err != nil {
		c.hxFail("读不到这间一次图")
		return
	}
	out, err := applyHxRoomPatch(src, patch)
	if err != nil {
		c.hxFail(err.Error())
		return
	}
	if err := writeHxRoomAtomic(path, out); err != nil {
		c.hxFail(err.Error())
		return
	}
	c.Data["json"] = map[string]interface{}{"code": 200, "room": patch.Room}
	c.ServeJSON()
}

// 整份布局一次换上：先在同目录写成临时文件，校验是完整 JSON，fsync 后再改名。
// 失败只删临时文件，正式文件保持原样。
func writeHxRoomAtomic(path string, out []byte) error {
	if !json.Valid(out) {
		return errors.New("要写的不是完整 JSON")
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("写入失败")
	}
	_, werr := f.Write(out)
	if werr == nil {
		werr = f.Sync()
	}
	cerr := f.Close()
	if werr != nil || cerr != nil {
		os.Remove(tmp)
		return errors.New("写入失败")
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return errors.New("写入失败")
	}
	return nil
}

func (c *HxPowerRoomController) hxFail(msg string) {
	c.Data["json"] = map[string]interface{}{"code": -1, "msg": msg}
	c.ServeJSON()
}

func hxRoomFile(code string) (string, error) {
	if !hxRoomCodeRe.MatchString(code) {
		return "", errors.New("室号不对")
	}
	// 5A1 不在 hx-power-rooms 目录，画面在已核对过的这份布局里。
	if code == "5A1" {
		dir, err := filepath.Abs(filepath.Join("..", "ism-front-end-v2", "public", "hx-dcim", "data"))
		if err != nil {
			return "", errors.New("路径不对")
		}
		path := filepath.Join(dir, "hx-5a1-cad-layout.json")
		if !hxInside(dir, path) {
			return "", errors.New("路径不对")
		}
		if st, err := os.Stat(path); err != nil || st.IsDir() {
			return "", errors.New("没有这间一次图")
		}
		return path, nil
	}
	dir, err := filepath.Abs(filepath.Join("..", "ism-front-end-v2", "public", "hx-dcim", "data", "hx-power-rooms"))
	if err != nil {
		return "", errors.New("路径不对")
	}
	path := filepath.Join(dir, code+".json")
	if !hxInside(dir, path) {
		return "", errors.New("路径不对")
	}
	if st, err := os.Stat(path); err != nil || st.IsDir() {
		return "", errors.New("没有这间一次图")
	}
	return path, nil
}

func hxInside(dir, path string) bool {
	rel, err := filepath.Rel(dir, path)
	return err == nil && rel != ".." && !bytes.HasPrefix([]byte(rel), []byte(".."+string(filepath.Separator)))
}

func applyHxRoomPatch(src []byte, patch hxRoomPatch) ([]byte, error) {
	var err error
	if src, err = patchNums(src, "texts", patch.Texts); err != nil {
		return nil, err
	}
	if src, err = patchNums(src, "knives", patch.Knives); err != nil {
		return nil, err
	}
	if src, err = patchNums(src, "boxes", patch.Boxes); err != nil {
		return nil, err
	}
	if src, err = patchWires(src, patch.Wires); err != nil {
		return nil, err
	}
	if src, err = removeIdx(src, "wires", patch.Remove.Wires); err != nil {
		return nil, err
	}
	if src, err = removeIdx(src, "boxes", patch.Remove.Boxes); err != nil {
		return nil, err
	}
	if src, err = removeIdx(src, "knives", patch.Remove.Knives); err != nil {
		return nil, err
	}
	if src, err = removeIdx(src, "texts", patch.Remove.Texts); err != nil {
		return nil, err
	}
	if src, err = appendRaw(src, "texts", patch.Add.Texts); err != nil {
		return nil, err
	}
	if src, err = appendRaw(src, "knives", patch.Add.Knives); err != nil {
		return nil, err
	}
	if src, err = appendRaw(src, "boxes", patch.Add.Boxes); err != nil {
		return nil, err
	}
	if src, err = appendRaw(src, "wires", patch.Add.Wires); err != nil {
		return nil, err
	}
	return src, nil
}

type hxRep struct {
	a, b int
	buf  []byte
}

func patchNums(src []byte, key string, items []hxNumPatch) ([]byte, error) {
	if len(items) == 0 {
		return src, nil
	}
	spans, err := objectSpans(src, key)
	if err != nil {
		return nil, err
	}
	var reps []hxRep
	for _, it := range items {
		if it.I < 0 || it.I >= len(spans) {
			return nil, errors.New(key + " 下标超出")
		}
		obj := append([]byte(nil), src[spans[it.I][0]:spans[it.I][1]]...)
		if it.X != nil {
			obj, err = replaceJSONNumber(obj, "x", *it.X)
			if err != nil {
				return nil, err
			}
		}
		if it.Y != nil {
			obj, err = replaceJSONNumber(obj, "y", *it.Y)
			if err != nil {
				return nil, err
			}
		}
		if it.S != nil {
			obj, err = replaceJSONString(obj, "s", *it.S)
			if err != nil {
				return nil, err
			}
		}
		if it.W != nil {
			obj, err = replaceJSONNumber(obj, "w", *it.W)
			if err != nil {
				return nil, err
			}
		}
		if it.H != nil {
			obj, err = replaceJSONNumber(obj, "h", *it.H)
			if err != nil {
				return nil, err
			}
		}
		if it.Fs != nil {
			obj, err = replaceJSONNumber(obj, "fs", *it.Fs)
			if err != nil {
				return nil, err
			}
		}
		if it.C != nil {
			obj, err = replaceJSONString(obj, "c", *it.C)
			if err != nil {
				return nil, err
			}
		}
		if it.CLock != nil {
			obj, err = setJSONBool(obj, "cLock", *it.CLock)
			if err != nil {
				return nil, err
			}
		}
		if it.FsLock != nil {
			obj, err = setJSONBool(obj, "fsLock", *it.FsLock)
			if err != nil {
				return nil, err
			}
		}
		reps = append(reps, hxRep{spans[it.I][0], spans[it.I][1], obj})
	}
	return applyReps(src, reps), nil
}

func patchWires(src []byte, items []hxWirePatch) ([]byte, error) {
	if len(items) == 0 {
		return src, nil
	}
	spans, err := objectSpans(src, "wires")
	if err != nil {
		return nil, err
	}
	var reps []hxRep
	for _, it := range items {
		if it.I < 0 || it.I >= len(spans) {
			return nil, errors.New("wires 下标超出")
		}
		obj := append([]byte(nil), src[spans[it.I][0]:spans[it.I][1]]...)
		flat := make([]float64, 0, len(it.Pts)*2)
		for _, p := range it.Pts {
			if len(p) < 2 {
				return nil, errors.New("线段坐标不全")
			}
			flat = append(flat, p[0], p[1])
		}
		obj, err = replacePts(obj, flat)
		if err != nil {
			return nil, err
		}
		reps = append(reps, hxRep{spans[it.I][0], spans[it.I][1], obj})
	}
	return applyReps(src, reps), nil
}

func removeIdx(src []byte, key string, idx []int) ([]byte, error) {
	if len(idx) == 0 {
		return src, nil
	}
	spans, err := objectSpans(src, key)
	if err != nil {
		return nil, err
	}
	seen := map[int]bool{}
	var ordered []int
	for _, i := range idx {
		if i < 0 || i >= len(spans) || seen[i] {
			continue
		}
		seen[i] = true
		ordered = append(ordered, i)
	}
	for a := 0; a < len(ordered); a++ {
		for b := a + 1; b < len(ordered); b++ {
			if ordered[b] > ordered[a] {
				ordered[a], ordered[b] = ordered[b], ordered[a]
			}
		}
	}
	for _, i := range ordered {
		src = dropSpan(src, spans[i][0], spans[i][1])
		spans, err = objectSpans(src, key)
		if err != nil {
			return nil, err
		}
	}
	return src, nil
}

func appendRaw(src []byte, key string, objs []json.RawMessage) ([]byte, error) {
	if len(objs) == 0 {
		return src, nil
	}
	a, b, err := arrayBounds(src, key)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	pretty := bytes.Contains(src[a:b], []byte("\n"))
	empty := true
	for _, c := range src[a+1 : b] {
		if c != ' ' && c != '\n' && c != '\r' && c != '\t' {
			empty = false
			break
		}
	}
	if !empty {
		buf.WriteByte(',')
	}
	for i, raw := range objs {
		if !json.Valid(raw) {
			return nil, errors.New(key + " 新增不是 JSON")
		}
		if pretty {
			buf.WriteByte('\n')
			var indented bytes.Buffer
			if err := json.Indent(&indented, raw, "    ", "  "); err != nil {
				return nil, err
			}
			buf.Write(indented.Bytes())
		} else {
			if i > 0 || !empty {
				buf.WriteByte(' ')
			}
			buf.Write(raw)
		}
	}
	if pretty {
		buf.WriteByte('\n')
		buf.WriteString("  ")
	}
	out := make([]byte, 0, len(src)+buf.Len())
	out = append(out, src[:b]...)
	out = append(out, buf.Bytes()...)
	out = append(out, src[b:]...)
	return out, nil
}

func applyReps(src []byte, reps []hxRep) []byte {
	for i := 0; i < len(reps); i++ {
		for j := i + 1; j < len(reps); j++ {
			if reps[j].a > reps[i].a {
				reps[i], reps[j] = reps[j], reps[i]
			}
		}
	}
	for _, r := range reps {
		out := make([]byte, 0, len(src)-(r.b-r.a)+len(r.buf))
		out = append(out, src[:r.a]...)
		out = append(out, r.buf...)
		out = append(out, src[r.b:]...)
		src = out
	}
	return src
}

func dropSpan(src []byte, start, end int) []byte {
	left := start
	for left > 0 && isJSONSpace(src[left-1]) {
		left--
	}
	if left > 0 && src[left-1] == ',' {
		return append(append([]byte{}, src[:left-1]...), src[end:]...)
	}
	right := end
	for right < len(src) && isJSONSpace(src[right]) {
		right++
	}
	if right < len(src) && src[right] == ',' {
		right++
	}
	return append(append([]byte{}, src[:start]...), src[right:]...)
}

func objectSpans(src []byte, key string) ([][2]int, error) {
	a, b, err := arrayBounds(src, key)
	if err != nil {
		return nil, err
	}
	arr := src[a : b+1]
	var spans [][2]int
	depth := 0
	inStr := false
	esc := false
	obj := -1
	for i := 0; i < len(arr); i++ {
		c := arr[i]
		if inStr {
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{', '[':
			if c == '{' && depth == 1 {
				obj = i
			}
			depth++
		case '}', ']':
			depth--
			if c == '}' && depth == 1 && obj >= 0 {
				spans = append(spans, [2]int{a + obj, a + i + 1})
				obj = -1
			}
		}
	}
	return spans, nil
}

func arrayBounds(src []byte, key string) (int, int, error) {
	pat := []byte(`"` + key + `"`)
	i := bytes.Index(src, pat)
	if i < 0 {
		return 0, 0, errors.New("没有 " + key)
	}
	j := i + len(pat)
	for j < len(src) && isJSONSpace(src[j]) {
		j++
	}
	if j >= len(src) || src[j] != ':' {
		return 0, 0, errors.New("没有 " + key)
	}
	j++
	for j < len(src) && isJSONSpace(src[j]) {
		j++
	}
	if j >= len(src) || src[j] != '[' {
		return 0, 0, errors.New("没有 " + key)
	}
	end := matchJSON(src, j)
	if end < 0 {
		return 0, 0, errors.New(key + " 没有收尾")
	}
	return j, end, nil
}

func matchJSON(src []byte, open int) int {
	depth := 0
	inStr := false
	esc := false
	for i := open; i < len(src); i++ {
		c := src[i]
		if inStr {
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{', '[':
			depth++
		case '}', ']':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func setJSONBool(obj []byte, field string, val bool) ([]byte, error) {
	key := []byte(`"` + field + `"`)
	lit := []byte("false")
	if val {
		lit = []byte("true")
	}
	i := bytes.Index(obj, key)
	if i < 0 {
		end := bytes.LastIndexByte(obj, '}')
		if end < 0 {
			return nil, errors.New("没有 " + field)
		}
		ins := append([]byte(`,"`+field+`":`), lit...)
		out := make([]byte, 0, len(obj)+len(ins))
		out = append(out, obj[:end]...)
		out = append(out, ins...)
		out = append(out, obj[end:]...)
		return out, nil
	}
	j := i + len(key)
	for j < len(obj) && isJSONSpace(obj[j]) {
		j++
	}
	if j >= len(obj) || obj[j] != ':' {
		return nil, errors.New("没有 " + field)
	}
	j++
	for j < len(obj) && isJSONSpace(obj[j]) {
		j++
	}
	k := j
	if bytes.HasPrefix(obj[j:], []byte("true")) {
		k = j + 4
	} else if bytes.HasPrefix(obj[j:], []byte("false")) {
		k = j + 5
	} else {
		return nil, errors.New("没有 " + field)
	}
	out := make([]byte, 0, len(obj)+len(lit))
	out = append(out, obj[:j]...)
	out = append(out, lit...)
	out = append(out, obj[k:]...)
	return out, nil
}

func replaceJSONNumber(obj []byte, field string, val float64) ([]byte, error) {
	i := bytes.Index(obj, []byte(`"`+field+`"`))
	if i < 0 {
		return nil, errors.New("没有 " + field)
	}
	j := i + len(field) + 2
	for j < len(obj) && isJSONSpace(obj[j]) {
		j++
	}
	if j >= len(obj) || obj[j] != ':' {
		return nil, errors.New("没有 " + field)
	}
	j++
	for j < len(obj) && isJSONSpace(obj[j]) {
		j++
	}
	k := j
	if k < len(obj) && obj[k] == '-' {
		k++
	}
	for k < len(obj) && ((obj[k] >= '0' && obj[k] <= '9') || obj[k] == '.') {
		k++
	}
	if k == j || (obj[j] == '-' && k == j+1) {
		return nil, errors.New("没有 " + field)
	}
	num := formatHxNum(val)
	out := make([]byte, 0, len(obj)+len(num))
	out = append(out, obj[:j]...)
	out = append(out, num...)
	out = append(out, obj[k:]...)
	return out, nil
}

func replaceJSONString(obj []byte, field, val string) ([]byte, error) {
	i := bytes.Index(obj, []byte(`"`+field+`"`))
	if i < 0 {
		return nil, errors.New("没有 " + field)
	}
	j := i + len(field) + 2
	for j < len(obj) && isJSONSpace(obj[j]) {
		j++
	}
	if j >= len(obj) || obj[j] != ':' {
		return nil, errors.New("没有 " + field)
	}
	j++
	for j < len(obj) && isJSONSpace(obj[j]) {
		j++
	}
	if j >= len(obj) || obj[j] != '"' {
		return nil, errors.New("没有 " + field)
	}
	k := j + 1
	esc := false
	for k < len(obj) {
		if esc {
			esc = false
			k++
			continue
		}
		if obj[k] == '\\' {
			esc = true
			k++
			continue
		}
		if obj[k] == '"' {
			k++
			break
		}
		k++
	}
	quoted, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(obj)+len(quoted))
	out = append(out, obj[:j]...)
	out = append(out, quoted...)
	out = append(out, obj[k:]...)
	return out, nil
}

func replacePts(obj []byte, flat []float64) ([]byte, error) {
	a, b, err := arrayBounds(obj, "pts")
	if err != nil {
		return nil, err
	}
	pts := obj[a : b+1]
	var buf bytes.Buffer
	vi := 0
	for i := 0; i < len(pts); i++ {
		c := pts[i]
		if c == '-' || (c >= '0' && c <= '9') {
			j := i
			if pts[j] == '-' {
				j++
			}
			for j < len(pts) && ((pts[j] >= '0' && pts[j] <= '9') || pts[j] == '.') {
				j++
			}
			if vi >= len(flat) {
				return nil, errors.New("线段点数对不上")
			}
			buf.WriteString(formatHxNum(flat[vi]))
			vi++
			i = j - 1
			continue
		}
		buf.WriteByte(c)
	}
	if vi != len(flat) {
		return nil, errors.New("线段点数对不上")
	}
	out := make([]byte, 0, len(obj)-len(pts)+buf.Len())
	out = append(out, obj[:a]...)
	out = append(out, buf.Bytes()...)
	out = append(out, obj[b+1:]...)
	return out, nil
}

func formatHxNum(v float64) string {
	r := math.Round(v*100) / 100
	return strconv.FormatFloat(r, 'f', -1, 64)
}

func isJSONSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}
