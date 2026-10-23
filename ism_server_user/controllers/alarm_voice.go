package controllers

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// 报警播报走穹盾项目里单宇辰的声纹（F5-TTS 零样本克隆），不再用浏览器把音调拉高的系统音。

type AlarmVoiceController struct {
	beego.Controller
}

type alarmVoiceReq struct {
	Text string `json:"text"`
}

type voiceLineReader struct {
	ch chan string
}

func newVoiceLineReader(r io.Reader) *voiceLineReader {
	ch := make(chan string, 8)
	go func() {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			ch <- sc.Text()
		}
		close(ch)
	}()
	return &voiceLineReader{ch: ch}
}

func (r *voiceLineReader) next(d time.Duration) (string, error) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case s, ok := <-r.ch:
		if !ok {
			return "", io.EOF
		}
		return s, nil
	case <-timer.C:
		return "", errors.New("语音合成超时")
	}
}

type danyuchenWorker struct {
	mu   sync.Mutex
	cmd  *exec.Cmd
	in   io.WriteCloser
	out  *voiceLineReader
	root string
}

var alarmSpeaker danyuchenWorker

func init() {
	go func() {
		time.Sleep(8 * time.Second)
		if _, err := synthDanyuchen("请注意，设备报警。"); err != nil {
			logs.Error("单宇辰报警语音预热失败: %v", err)
			return
		}
		logs.Info("单宇辰报警语音已就绪")
	}()
}

func alarmVoiceRoot() string {
	root, err := config.String("alarmvoiceroot")
	if err != nil || strings.TrimSpace(root) == "" {
		return "/Users/yunfanleo/cursorProjects/ 低空管控智能体/web"
	}
	return strings.TrimSpace(root)
}

func (w *danyuchenWorker) ensure(root string) error {
	if w.cmd != nil && w.cmd.Process != nil {
		return nil
	}
	python := filepath.Join(root, ".venv-stt", "bin", "python")
	script := filepath.Join(root, "scripts", "voice_clone.py")
	if _, err := os.Stat(python); err != nil {
		return errors.New("找不到单宇辰语音环境")
	}
	cmd := exec.Command(python, script)
	cmd.Env = append(os.Environ(),
		"HF_HUB_OFFLINE=1",
		"TRANSFORMERS_OFFLINE=1",
		"HF_HUB_DISABLE_TELEMETRY=1",
		"PYTHONUNBUFFERED=1",
		"PYTHONHASHSEED=1",
	)
	cmd.Stderr = io.Discard
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	reader := newVoiceLineReader(stdout)
	deadline := time.Now().Add(120 * time.Second)
	for {
		line, err := reader.next(time.Until(deadline))
		if err != nil {
			_ = cmd.Process.Kill()
			return err
		}
		var msg struct {
			Ready bool   `json:"ready"`
			Error string `json:"error"`
		}
		if json.Unmarshal([]byte(line), &msg) != nil {
			continue
		}
		if msg.Error != "" {
			_ = cmd.Process.Kill()
			return errors.New(msg.Error)
		}
		if msg.Ready {
			break
		}
	}
	w.cmd = cmd
	w.in = stdin
	w.out = reader
	w.root = root
	return nil
}

func (w *danyuchenWorker) reset() {
	if w.cmd != nil && w.cmd.Process != nil {
		_ = w.cmd.Process.Kill()
		_ = w.cmd.Wait()
	}
	w.cmd = nil
	w.in = nil
	w.out = nil
}

func danyuchenPrompt(root string) string {
	raw, err := os.ReadFile(filepath.Join(root, "data", "voiceprint", "meta.json"))
	if err != nil {
		return "各位好。"
	}
	var meta struct {
		Prompt string `json:"prompt"`
	}
	if json.Unmarshal(raw, &meta) != nil || strings.TrimSpace(meta.Prompt) == "" {
		return "各位好。"
	}
	return strings.TrimSpace(meta.Prompt)
}

func clipAlarmText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Map(func(r rune) rune {
		if r < 32 {
			return -1
		}
		return r
	}, text)
	if utf8.RuneCountInString(text) > 80 {
		runes := []rune(text)
		text = string(runes[:80])
	}
	return text
}

func synthDanyuchen(text string) (string, error) {
	text = clipAlarmText(text)
	if text == "" {
		return "", errors.New("没有要播的内容")
	}
	root := alarmVoiceRoot()
	ref := filepath.Join(root, "data", "voiceprint", "ref.wav")
	if _, err := os.Stat(ref); err != nil {
		return "", errors.New("找不到单宇辰声纹")
	}
	const speed = 0.48
	const nfe = 16
	sum := sha1.Sum([]byte(fmt.Sprintf("%.2f\n%s", speed, text)))
	name := hex.EncodeToString(sum[:]) + ".wav"
	dir := filepath.Join("data", "alarm-voice")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	out := filepath.Join(dir, name)
	if info, err := os.Stat(out); err == nil && info.Size() > 44 {
		return out, nil
	}

	alarmSpeaker.mu.Lock()
	defer alarmSpeaker.mu.Unlock()
	if info, err := os.Stat(out); err == nil && info.Size() > 44 {
		return out, nil
	}
	if err := alarmSpeaker.ensure(root); err != nil {
		alarmSpeaker.reset()
		return "", err
	}
	req, _ := json.Marshal(map[string]interface{}{
		"ref": ref, "refText": danyuchenPrompt(root), "text": text,
		"out": out, "speed": speed, "nfe": nfe,
	})
	if _, err := io.WriteString(alarmSpeaker.in, string(req)+"\n"); err != nil {
		alarmSpeaker.reset()
		return "", err
	}
	line, err := alarmSpeaker.out.next(90 * time.Second)
	if err != nil {
		alarmSpeaker.reset()
		return "", err
	}
	var reply struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if json.Unmarshal([]byte(line), &reply) != nil || !reply.Ok {
		msg := reply.Error
		if msg == "" {
			msg = "克隆失败"
		}
		alarmSpeaker.reset()
		return "", errors.New(msg)
	}
	return out, nil
}

func (c *AlarmVoiceController) Speak() {
	text := strings.TrimSpace(c.GetString("text"))
	if text == "" && len(c.Ctx.Input.RequestBody) > 0 {
		var req alarmVoiceReq
		_ = json.Unmarshal(c.Ctx.Input.RequestBody, &req)
		text = req.Text
	}
	path, err := synthDanyuchen(text)
	if err != nil {
		c.Ctx.Output.SetStatus(503)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	wav, err := os.ReadFile(path)
	if err != nil {
		c.Ctx.Output.SetStatus(503)
		c.Data["json"] = map[string]string{"error": "语音文件读不出来"}
		c.ServeJSON()
		return
	}
	c.Ctx.Output.Header("Content-Type", "audio/wav")
	c.Ctx.Output.Header("Cache-Control", "no-store")
	_ = c.Ctx.Output.Body(wav)
}
