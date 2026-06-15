/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-20 15:48:32
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	xlst "github.com/ivahaev/go-xlsx-templater"
)

type ReportController struct {
	beego.Controller
}

func buildDifferencePairsAndSeries(deviceList, dataList []string, bucketSize int) (map[string][]int, []models.HourDataPair, [][]any) {
	pairIndexMap := make(map[string][]int)
	pairs := make([]models.HourDataPair, len(deviceList))
	series := make([][]any, len(deviceList))

	for i := range deviceList {
		deviceName := deviceList[i]
		dataName := dataList[i]
		key := deviceName + "\x00" + dataName
		pairIndexMap[key] = append(pairIndexMap[key], i)
		pairs[i] = models.HourDataPair{DeviceName: deviceName, DataName: dataName}
		series[i] = make([]any, bucketSize)
		for j := range series[i] {
			series[i][j] = "-"
		}
	}

	return pairIndexMap, pairs, series
}

func fetchDifferenceBatch(reportType string, pairs []models.HourDataPair, startTime, endTime string) (int, []models.DevicesHistoryDataList, []models.DevicesCHHistoryData) {
	var records []models.DevicesHistoryDataList
	var recordsCK []models.DevicesCHHistoryData
	var code int

	switch reportType {
	case "day":
		if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
			code, records = models.GetMysqlDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 2 {
			code, records = models.GetTsDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
		} else if protocol_common.HistoryRecordDbType == 4 {
			code, records = models.GetInfluxDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 3 {
			code, recordsCK = models.GetClickHouseDifferenceDataBatch(pairs, startTime, endTime)
		}
	case "week":
		if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
			code, records = models.GetMysqlWeeklyDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 2 {
			code, records = models.GetTsWeeklyDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
		} else if protocol_common.HistoryRecordDbType == 4 {
			code, records = models.GetInfluxWeeklyDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 3 {
			code, recordsCK = models.GetClickHouseWeeklyDifferenceDataBatch(pairs, startTime, endTime)
		}
	case "month":
		if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
			code, records = models.GetMysqlMonthlyDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 2 {
			code, records = models.GetTsMonthlyDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
		} else if protocol_common.HistoryRecordDbType == 4 {
			code, records = models.GetInfluxMonthlyDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 3 {
			code, recordsCK = models.GetClickHouseMonthlyDifferenceDataBatch(pairs, startTime, endTime)
		}
	case "year":
		if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
			code, records = models.GetMysqlYearlyDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 2 {
			code, records = models.GetTsYearlyDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
		} else if protocol_common.HistoryRecordDbType == 4 {
			code, records = models.GetInfluxYearlyDifferenceDataBatch(pairs, startTime, endTime)
		} else if protocol_common.HistoryRecordDbType == 3 {
			code, recordsCK = models.GetClickHouseYearlyDifferenceDataBatch(pairs, startTime, endTime)
		}
	}

	return code, records, recordsCK
}

func fillSeriesFromHistoryRecords(records []models.DevicesHistoryDataList, pairIndexMap map[string][]int, series [][]any, indexFn func(time.Time) (int, bool)) {
	for _, r := range records {
		key := r.DeviceName + "\x00" + r.DataName
		idxs, ok := pairIndexMap[key]
		if !ok {
			continue
		}

		bucketIndex, ok := indexFn(r.RecordTime)
		if !ok {
			continue
		}

		for _, idx := range idxs {
			series[idx][bucketIndex] = r.DataValue
		}
	}
}

func fillSeriesFromClickHouseRecords(records []models.DevicesCHHistoryData, pairIndexMap map[string][]int, series [][]any, indexFn func(time.Time) (int, bool)) {
	for _, r := range records {
		key := r.DeviceName + "\x00" + r.DataName
		idxs, ok := pairIndexMap[key]
		if !ok {
			continue
		}

		bucketIndex, ok := indexFn(r.RecordTime)
		if !ok {
			continue
		}

		for _, idx := range idxs {
			series[idx][bucketIndex] = r.DataValue
		}
	}
}

func (c *ReportController) GetAlarmHistoryList() {

	var data interface{}
	var params = make(map[string]interface{})
	var code int

	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	//json数据封装到user对象中
	err := json.Unmarshal(rawData, &params)
	if err != nil {
		code = -1
	} else {
		data, code = models.GetAlarmHistoryList(ProjectUuid, params)
	}
	result := map[string]interface{}{
		"code": code,
		"list": data,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *ReportController) GetDataHistoryList() {

	var data interface{}
	var params = make(map[string]interface{})
	var code int

	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	//json数据封装到user对象中
	err := json.Unmarshal(rawData, &params)
	if err != nil {
		code = -1
	} else {
		if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
			data, code = models.GetDataHistoryList(ProjectUuid, params)
		} else if protocol_common.HistoryRecordDbType == 2 {
			data, code = models.GetDataTsHistoryList(ProjectUuid, params, protocol_common.HistoryRecordTsDb)
		} else if protocol_common.HistoryRecordDbType == 4 {
			data, code = models.GetDataInfluxHistoryList(ProjectUuid, params)
		} else if protocol_common.HistoryRecordDbType == 3 {
			data, code = models.GetDataClickHouseHistoryList(ProjectUuid, params)
		}

	}
	result := map[string]interface{}{
		"code": code,
		"list": data,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

type AlignedRow struct {
	MainTime time.Time     // 主轴时间
	Items    []AlignedItem // 所有 UUID 对应的对齐项（顺序一致即可，无需 key）
}

// AlignedItem 表示一个对齐后的数据项（对应一个 DataUuid）
type AlignedItem struct {
	DataName   string // 数据名称（若未匹配则为空）
	DeviceName string // 设备名称
	DataValue  string // 值
}

// AutoAlignToSlice 自动对齐为 slice 形式（支持通用 HistoryDataItem 接口）
func AutoAlignToSlice(data []HistoryDataItem, timeWindow time.Duration) []struct {
	MainTime time.Time
	Items    []AlignedItem
} {
	if len(data) == 0 || timeWindow < 0 {
		return nil
	}

	// 1. 按 DataUuid 分组
	groups := make(map[string][]HistoryDataItem)
	var earliestPt HistoryDataItem

	for i, pt := range data {
		groups[pt.GetDataUuid()] = append(groups[pt.GetDataUuid()], pt)
		if i == 0 || pt.GetRecordTime().Before(earliestPt.GetRecordTime()) {
			earliestPt = pt
		}
	}

	if earliestPt == nil { // 防御性检查
		return nil
	}

	mainUUID := earliestPt.GetDataUuid()
	mainPoints := groups[mainUUID]

	// 2. 对所有分组按时间排序
	sortedGroups := make(map[string][]HistoryDataItem)
	for uuid, pts := range groups {
		sorted := make([]HistoryDataItem, len(pts))
		copy(sorted, pts)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].GetRecordTime().Before(sorted[j].GetRecordTime())
		})
		sortedGroups[uuid] = sorted
	}

	// 3. 固定 UUID 顺序（避免每次输出顺序不同）
	allUUIDs := make([]string, 0, len(groups))
	for uuid := range groups {
		allUUIDs = append(allUUIDs, uuid)
	}
	sort.Strings(allUUIDs)

	// 4. 对每个主轴点进行对齐
	var result []struct {
		MainTime time.Time
		Items    []AlignedItem
	}

	for _, mainPt := range mainPoints {
		var items []AlignedItem

		for _, uuid := range allUUIDs {
			points := sortedGroups[uuid]
			var matched HistoryDataItem = nil
			minDiff := time.Duration(1<<63 - 1)

			if len(points) > 0 {
				i := sort.Search(len(points), func(j int) bool {
					return !points[j].GetRecordTime().Before(mainPt.GetRecordTime())
				})

				candidates := []HistoryDataItem{}
				if i < len(points) {
					candidates = append(candidates, points[i])
				}
				if i > 0 {
					candidates = append(candidates, points[i-1])
				}

				for j := range candidates {
					diff := candidates[j].GetRecordTime().Sub(mainPt.GetRecordTime())
					if diff < 0 {
						diff = -diff
					}
					if diff < minDiff {
						minDiff = diff
						matched = candidates[j]
					}
				}

				if minDiff > timeWindow {
					matched = nil
				}
			}

			if matched != nil {
				items = append(items, AlignedItem{
					DataName:   matched.GetDataName(),
					DeviceName: matched.GetDeviceName(),
					DataValue:  matched.GetDataValue(),
				})
			}
		}

		result = append(result, struct {
			MainTime time.Time
			Items    []AlignedItem
		}{
			MainTime: mainPt.GetRecordTime(),
			Items:    items,
		})
	}

	return result
}

// 豆包API配置
type DoubaoConfig struct {
	APIKey    string // 火山引擎AK（原生调用用AK作为API Key）
	Model     string // doubao-lite/doubao-pro
	BaseURL   string // 固定为https://ark.cn-beijing.volces.com/api/v3
	Timeout   int64  // 超时时间（秒）
	UseProxy  bool   // 是否使用代理
	ProxyAddr string // 代理地址
}

// 消息结构体
type ChatMessage struct {
	Role    string `json:"role"`    // user/assistant/system
	Content string `json:"content"` // 消息内容
}

// 非流式请求结构体
type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// 非流式响应结构体
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// 流式响应分片结构体
type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// 创建HTTP客户端（含代理）
func newHTTPClient(cfg DoubaoConfig) *http.Client {
	transport := &http.Transport{}
	// 配置代理
	if cfg.UseProxy {
		proxyURL, err := url.Parse(cfg.ProxyAddr)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}
	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}
}

// 基础对话：一次性获取完整响应
func ChatWithDoubao(ctx context.Context, cfg DoubaoConfig, prompt string) (string, error) {
	client := newHTTPClient(cfg)
	// 构造请求体
	reqBody := ChatCompletionRequest{
		Model: cfg.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败：%v", err)
	}
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", cfg.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败：%v", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey) // 鉴权：Bearer + AK
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败：%v", err)
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败：%v", err)
	}
	// 解析响应
	var chatResp ChatCompletionResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", fmt.Errorf("解析响应失败：%v，原始响应：%s", err, string(body))
	}
	// 处理错误
	if chatResp.Error.Code != 0 {
		return "", fmt.Errorf("API错误：%d - %s", chatResp.Error.Code, chatResp.Error.Message)
	}
	// 提取结果
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("无响应结果")
	}
	return chatResp.Choices[0].Message.Content, nil
}

// 流式对话：实时返回响应
func StreamChatWithDoubao(ctx context.Context, cfg DoubaoConfig, prompt string) error {
	client := newHTTPClient(cfg)
	// 构造请求体
	reqBody := ChatCompletionRequest{
		Model: cfg.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		Stream: true,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求体失败：%v", err)
	}
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", cfg.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("创建请求失败：%v", err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败：%v", err)
	}
	defer resp.Body.Close()
	// 流式读取响应
	fmt.Println("===== 流式回复 =====")
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // 读取完毕
		}
		// 处理SSE格式数据
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}
			// 解析分片
			var chunk StreamChunk
			err := json.Unmarshal([]byte(data), &chunk)
			if err != nil {
				continue
			}
			// 实时输出
			if len(chunk.Choices) > 0 {
				fmt.Print(chunk.Choices[0].Delta.Content)
			}
		}
	}
	fmt.Println("\n")
	return nil
}

// 数据分析：传入结构化数据让豆包分析
func AnalyzeDataWithDoubao(ctx context.Context, cfg DoubaoConfig, data interface{}, analyzePrompt string) (string, error) {
	// 格式化数据为JSON
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("格式化数据失败：%v", err)
	}
	// 构建完整提示词
	fullPrompt := fmt.Sprintf(`请分析以下业务数据：
%s
分析要求：
%s`, string(dataJSON), analyzePrompt)
	// 调用基础对话接口
	return ChatWithDoubao(ctx, cfg, fullPrompt)
}

// 1. 先定义通用接口（字段名需与models包中结构体实际字段一致）
type HistoryDataItem interface {
	GetDataUuid() string
	GetRecordTime() time.Time
	GetDataName() string
	GetDeviceName() string
	GetDataValue() string
}

// 2. 为models.DevicesHistoryDataList创建本地包装类型
type LocalDevicesHistoryData models.DevicesHistoryDataList

// 3. 为本地包装类型实现HistoryDataItem接口
func (d LocalDevicesHistoryData) GetDataUuid() string {
	return d.DataUuid // 确保models包中结构体有DataUuid字段
}
func (d LocalDevicesHistoryData) GetRecordTime() time.Time {
	return d.RecordTime // 确保有RecordTime字段
}
func (d LocalDevicesHistoryData) GetDataName() string {
	return d.DataName
}
func (d LocalDevicesHistoryData) GetDeviceName() string {
	return d.DeviceName
}
func (d LocalDevicesHistoryData) GetDataValue() string {
	return d.DataValue
}

// 4. 为models.DevicesCHHistoryData创建本地包装类型
type LocalDevicesCHHistoryData models.DevicesCHHistoryData

// 5. 为ClickHouse版本的本地类型实现接口
func (d LocalDevicesCHHistoryData) GetDataUuid() string {
	return d.DataUuid // 按models包中实际字段名调整
}
func (d LocalDevicesCHHistoryData) GetRecordTime() time.Time {
	return d.RecordTime
}
func (d LocalDevicesCHHistoryData) GetDataName() string {
	return d.DataName
}
func (d LocalDevicesCHHistoryData) GetDeviceName() string {
	return d.DeviceName
}
func (d LocalDevicesCHHistoryData) GetDataValue() string {
	return d.DataValue
}
func (c *ReportController) GetDataHistoryReport() {

	var data interface{}
	var params = make(map[string]interface{})
	var code int

	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	//json数据封装到user对象中
	err := json.Unmarshal(rawData, &params)
	if err != nil {
		code = -1
	} else {
		if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
			data, code = models.GetDataHistoryReport(ProjectUuid, params)
		} else if protocol_common.HistoryRecordDbType == 2 {
			data, code = models.GetDataTsHistoryList(ProjectUuid, params, protocol_common.HistoryRecordTsDb)
		} else if protocol_common.HistoryRecordDbType == 4 {
			data, code = models.GetDataInfluxHistoryList(ProjectUuid, params)
		} else if protocol_common.HistoryRecordDbType == 3 {
			data, code = models.GetDataClickHouseHistoryList(ProjectUuid, params)
		}
		if code == 0 {

			var timeWindows int = 30
			dataWindows, ok := params["dataWindows"]
			if ok {
				// 安全转换 float64 → int（JSON 解析的数字默认是 float64）
				if val, ok := dataWindows.(float64); ok {
					timeWindows = int(val)
				}
			}

			// 定义通用接口切片
			var historyItems []HistoryDataItem
			var convertOk bool

			// 安全类型断言 + 转换为本地包装类型
			if protocol_common.HistoryRecordDbType != 3 {
				// MySQL/TimescaleDB/InfluxDB 类型：转换为LocalDevicesHistoryData
				var dataList []models.DevicesHistoryDataList
				dataList, convertOk = data.([]models.DevicesHistoryDataList)
				if convertOk {
					// 把models类型切片转换为本地包装类型切片（实现接口）
					historyItems = make([]HistoryDataItem, len(dataList))
					for i, item := range dataList {
						historyItems[i] = LocalDevicesHistoryData(item) // 类型转换
					}
				}
			} else {
				// ClickHouse 类型：转换为LocalDevicesCHHistoryData
				var chDataList []models.DevicesCHHistoryData
				chDataList, convertOk = data.([]models.DevicesCHHistoryData)
				if convertOk {
					historyItems = make([]HistoryDataItem, len(chDataList))
					for i, item := range chDataList {
						historyItems[i] = LocalDevicesCHHistoryData(item) // 类型转换
					}
				}
			}

			// 类型转换成功才执行对齐
			if convertOk && len(historyItems) > 0 {
				rows := AutoAlignToSlice(historyItems, time.Duration(timeWindows)*time.Second)
				data = rows
			} else {
				// 类型转换失败时的错误处理（避免返回错误数据）
				code = -2 // 自定义错误码：数据类型不支持
				data = nil
			}
		}
	}
	result := map[string]interface{}{
		"code": code,
		"list": data,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *ReportController) GetDiyDataHistoryList() {

	var data string
	var params = make(map[string]interface{})
	var code int
	rawData := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	//json数据封装到user对象中
	err := json.Unmarshal(rawData, &params)
	if err != nil {
		code = -1
	} else {
		if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
			data, code, _ = models.GetDiyDataHistoryList(ProjectUuid, params)
		} else if protocol_common.HistoryRecordDbType == 2 {
			data, code, _ = models.GetDiyDataTsHistoryList(ProjectUuid, params, protocol_common.HistoryRecordTsDb)
		} else if protocol_common.HistoryRecordDbType == 4 {
			data, code, _ = models.GetDiyDataInfluxHistoryList(ProjectUuid, params)
		} else if protocol_common.HistoryRecordDbType == 3 {
			data, code, _ = models.GetDiyDataClickHouseHistoryList(ProjectUuid, params)
		}

	}
	result := map[string]interface{}{
		"code": code,
		"path": data,
	}
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *ReportController) GetChartDataHistoryList() {

	rawData := c.Ctx.Input.RequestBody
	var result map[string]interface{}

	if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
		code, data := models.GetChartDataHistoryList(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 2 {
		code, data := models.GetChartDataTsHistoryList(rawData, protocol_common.HistoryRecordTsDb)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 4 {
		code, data := models.GetChartDataInfluxHistoryList(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 3 {
		code, data := models.GetChartDataClickHouseHistoryList(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func (c *ReportController) GetTrendChartData() {

	rawData := c.Ctx.Input.RequestBody
	var result map[string]interface{}

	if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
		code, data := models.GetMysqlTrendChartData(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 2 {
		code, data := models.GetTsTrendChartData(rawData, protocol_common.HistoryRecordTsDb)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 4 {
		code, data := models.GetInfluxTrendChartData(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 3 {
		code, data := models.GetClickHouseTrendChartData(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func (c *ReportController) TestChartData() {

	var err error
	var code int = 0
	var results []map[string]interface{}

	ctx := map[string]interface{}{
		"DeviceName": "Item name",
	}
	doc := xlst.New()
	err4 := doc.ReadTemplate("tt/ty.xlsx")

	if err4 == nil {
		models.Db.Raw("select * from devices_history_data_list").Scan(&results)
		ctx["DataModel"] = results
		err = doc.Render(ctx)
		if err != nil {
			code = -1
		}
		err = doc.Save("tt/ty2.xlsx")
		if err != nil {
			code = -2
		}
	} else {
		code = -3
	}
	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *ReportController) GetHistoryHour() {
	rawData := c.Ctx.Input.RequestBody
	var code int = 0
	var result map[string]interface{}

	// 前端传入参数
	type FrontendParams struct {
		QueryDate  string   `json:"QueryDate"`  // 格式：2026-04（年月）
		DeviceList []string `json:"DeviceList"` // 设备名
		DataList   []string `json:"DataList"`   // 数据名
	}

	var fp FrontendParams
	var totalData [][]any // 二维数组：[设备数][7天]

	// 1. 校验ProjectUuid
	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result = map[string]interface{}{"code": -6, "realData": nil}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 2. 解析前端参数
	err := json.Unmarshal(rawData, &fp)
	if err != nil {
		code = -1
	} else {
		// 3. 校验参数长度
		if len(fp.DeviceList) == 0 || len(fp.DataList) == 0 || len(fp.DeviceList) != len(fp.DataList) {
			code = -7
		} else {
			// 4. 解析目标日期（带时区）
			queryDate, err := time.Parse("2006-01-02", fp.QueryDate)
			if err != nil {
				code = -8 // 日期格式错误
			} else {
				// 当天起止时间
				startTime := queryDate.Format("2006-01-02 00:00:00")
				endTime := queryDate.AddDate(0, 0, 1).Format("2006-01-02 00:00:00")

				// 5. 一次性批量获取所有 (DeviceName, DataName) 组合的数据
				pairIndexMap := make(map[string][]int)
				pairs := make([]models.HourDataPair, len(fp.DeviceList))
				hourSeries := make([][]any, len(fp.DeviceList))
				for i := range fp.DeviceList {
					deviceName := fp.DeviceList[i]
					dataName := fp.DataList[i]
					key := deviceName + "\x00" + dataName
					pairIndexMap[key] = append(pairIndexMap[key], i)
					pairs[i] = models.HourDataPair{DeviceName: deviceName, DataName: dataName}
					hourSeries[i] = make([]any, 24) // 24小时
					for j := range hourSeries[i] {
						hourSeries[i][j] = "-"
					}
				}

				var records []models.DevicesHistoryDataList
				var recordsCK []models.DevicesCHHistoryData
				if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
					code, records = models.GetMysqlHourDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 2 {
					code, records = models.GetTsHourDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
				} else if protocol_common.HistoryRecordDbType == 4 {
					code, records = models.GetInfluxHourDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 3 {
					code, recordsCK = models.GetClickHouseHourDataBatch(pairs, startTime, endTime)
				}

				if code == 0 {
					if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 || protocol_common.HistoryRecordDbType == 2 || protocol_common.HistoryRecordDbType == 4 {
						for _, r := range records {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 计算该记录所属的小时（0-23）
							hourIndex := r.RecordTime.Hour()
							if hourIndex < 0 || hourIndex >= 24 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][hourIndex] = r.DataValue
							}
						}
					} else {
						for _, r := range recordsCK {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 计算该记录所属的小时（0-23）
							hourIndex := r.RecordTime.Hour()
							if hourIndex < 0 || hourIndex >= 24 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][hourIndex] = r.DataValue
							}
						}
					}
				}

				totalData = hourSeries
			}
		}
	}

	// 返回
	result = map[string]interface{}{
		"code":     code,
		"realData": totalData,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()
	_ = json.NewEncoder(gw).Encode(result)
}

// GetHistoryYearDifferenceReport
// 获取年度差值报告（当年每个月的差值）
func (c *ReportController) GetHistoryYearDifferenceReport() {
	rawData := c.Ctx.Input.RequestBody
	var code int = 0
	var result map[string]interface{}

	// 前端传入参数
	type FrontendParams struct {
		QueryDate  string   `json:"QueryDate"`  // 格式：2026（年份）
		DeviceList []string `json:"DeviceList"` // 设备名
		DataList   []string `json:"DataList"`   // 数据名
	}

	var fp FrontendParams
	var totalData [][]any // 二维数组：[设备数][12个月]

	// 1. 校验ProjectUuid
	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result = map[string]interface{}{"code": -6, "realData": nil}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 2. 解析前端参数
	err := json.Unmarshal(rawData, &fp)
	if err != nil {
		code = -1
	} else {
		// 3. 校验参数长度
		if len(fp.DeviceList) == 0 || len(fp.DataList) == 0 || len(fp.DeviceList) != len(fp.DataList) {
			code = -7
		} else {
			// 4. 解析目标年份
			queryYear, err := time.Parse("2006", fp.QueryDate)
			if err != nil {
				code = -8 // 年份格式错误
			} else {
				// 计算当年第一天和最后一天
				startTime := time.Date(queryYear.Year(), 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
				endTime := time.Date(queryYear.Year(), 12, 31, 23, 59, 59, 0, time.Local).Format("2006-01-02 15:04:05")

				// 5. 一次性批量获取所有 (DeviceName, DataName) 组合的数据
				pairIndexMap := make(map[string][]int)
				pairs := make([]models.HourDataPair, len(fp.DeviceList))
				hourSeries := make([][]any, len(fp.DeviceList))
				for i := range fp.DeviceList {
					deviceName := fp.DeviceList[i]
					dataName := fp.DataList[i]
					key := deviceName + "\x00" + dataName
					pairIndexMap[key] = append(pairIndexMap[key], i)
					pairs[i] = models.HourDataPair{DeviceName: deviceName, DataName: dataName}
					hourSeries[i] = make([]any, 12) // 12个月
					for j := range hourSeries[i] {
						hourSeries[i][j] = "-"
					}
				}

				var records []models.DevicesHistoryDataList
				var recordsCK []models.DevicesCHHistoryData
				if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
					code, records = models.GetMysqlYearlyDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 2 {
					code, records = models.GetTsYearlyDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
				} else if protocol_common.HistoryRecordDbType == 4 {
					code, records = models.GetInfluxYearlyDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 3 {
					code, recordsCK = models.GetClickHouseYearlyDifferenceDataBatch(pairs, startTime, endTime)
				}

				if code == 0 {
					if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 || protocol_common.HistoryRecordDbType == 2 || protocol_common.HistoryRecordDbType == 4 {
						for _, r := range records {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 确保记录时间在查询年份内
							if r.RecordTime.Year() != queryYear.Year() {
								continue
							}
							// 计算月份索引（0-11）
							monthIndex := int(r.RecordTime.Month()) - 1
							if monthIndex < 0 || monthIndex >= 12 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][monthIndex] = r.DataValue
							}
						}
					} else {
						for _, r := range recordsCK {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 确保记录时间在查询年份内
							if r.RecordTime.Year() != queryYear.Year() {
								continue
							}
							// 计算月份索引（0-11）
							monthIndex := int(r.RecordTime.Month()) - 1
							if monthIndex < 0 || monthIndex >= 12 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][monthIndex] = r.DataValue
							}
						}
					}
				}

				totalData = hourSeries
			}
		}
	}

	// 返回
	result = map[string]interface{}{
		"code":     code,
		"realData": totalData,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()
	_ = json.NewEncoder(gw).Encode(result)
}
func (c *ReportController) GetHistoryMonthDifferenceReport() {
	rawData := c.Ctx.Input.RequestBody
	var code int = 0
	var result map[string]interface{}

	// 前端传入参数
	type FrontendParams struct {
		QueryDate  string   `json:"QueryDate"`  // 格式：2026-04-04
		DeviceList []string `json:"DeviceList"` // 设备名
		DataList   []string `json:"DataList"`   // 数据名
	}

	var fp FrontendParams
	var totalData [][]any // 二维数组：[设备数][当月天数]

	// 1. 校验ProjectUuid
	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result = map[string]interface{}{"code": -6, "realData": nil}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 2. 解析前端参数
	err := json.Unmarshal(rawData, &fp)
	if err != nil {
		code = -1
	} else {
		// 3. 校验参数长度
		if len(fp.DeviceList) == 0 || len(fp.DataList) == 0 || len(fp.DeviceList) != len(fp.DataList) {
			code = -7
		} else {
			// 4. 解析目标日期（年月格式：YYYY-MM）
			queryDate, err := time.Parse("2006-01", fp.QueryDate)
			if err != nil {
				code = -8 // 日期格式错误
			} else {
				// 计算当月第一天（00:00:00）
				startTime := queryDate.Format("2006-01-02 00:00:00")
				// 计算当月最后一天（23:59:59）
				lastDay := queryDate.AddDate(0, 1, -1)
				endTime := lastDay.Format("2006-01-02") + " 23:59:59"

				// 5. 一次性批量获取所有 (DeviceName, DataName) 组合的数据
				pairIndexMap := make(map[string][]int)
				pairs := make([]models.HourDataPair, len(fp.DeviceList))
				hourSeries := make([][]any, len(fp.DeviceList))
				// 计算当月天数：下个月第一天减一天，然后取日期值
				lastDayOfMonth := queryDate.AddDate(0, 1, -1)
				monthDays := lastDayOfMonth.Day()
				for i := range fp.DeviceList {
					deviceName := fp.DeviceList[i]
					dataName := fp.DataList[i]
					key := deviceName + "\x00" + dataName
					pairIndexMap[key] = append(pairIndexMap[key], i)
					pairs[i] = models.HourDataPair{DeviceName: deviceName, DataName: dataName}
					hourSeries[i] = make([]any, monthDays)
					for j := range hourSeries[i] {
						hourSeries[i][j] = "-"
					}
				}

				var records []models.DevicesHistoryDataList
				var recordsCK []models.DevicesCHHistoryData
				if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
					code, records = models.GetMysqlMonthlyDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 2 {
					code, records = models.GetTsMonthlyDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
				} else if protocol_common.HistoryRecordDbType == 4 {
					code, records = models.GetInfluxMonthlyDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 3 {
					code, recordsCK = models.GetClickHouseMonthlyDifferenceDataBatch(pairs, startTime, endTime)
				}

				if code == 0 {
					if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 || protocol_common.HistoryRecordDbType == 2 || protocol_common.HistoryRecordDbType == 4 {
						for _, r := range records {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 确保记录时间在查询月份内
							if r.RecordTime.Year() != queryDate.Year() || r.RecordTime.Month() != queryDate.Month() {
								continue
							}
							// 计算日期是当月第几天（1-31）
							dayOfMonth := r.RecordTime.Day()
							dayIndex := dayOfMonth - 1 // 转换为0-indexed
							if dayIndex < 0 || dayIndex >= len(hourSeries[0]) {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][dayIndex] = r.DataValue
							}
						}
					} else {
						for _, r := range recordsCK {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 计算日期是当月第几天（1-31）
							dayOfMonth := r.RecordTime.Day()
							dayIndex := dayOfMonth - 1 // 转换为0-indexed
							if dayIndex < 0 || dayIndex >= len(hourSeries[0]) {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][dayIndex] = r.DataValue
							}
						}
					}
				}

				totalData = hourSeries
			}
		}
	}

	// 返回
	result = map[string]interface{}{
		"code":     code,
		"realData": totalData,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()
	_ = json.NewEncoder(gw).Encode(result)
}
func (c *ReportController) GetHistoryDayDifferenceReport() {
	rawData := c.Ctx.Input.RequestBody
	var code int = 0
	var result map[string]interface{}

	// 前端传入参数
	type FrontendParams struct {
		QueryDate  string   `json:"QueryDate"`  // 格式：2026-04-04
		DeviceList []string `json:"DeviceList"` // 设备名
		DataList   []string `json:"DataList"`   // 数据名
	}

	var fp FrontendParams
	var totalData [][]any // 二维数组：[设备数][7天]

	// 1. 校验ProjectUuid
	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result = map[string]interface{}{"code": -6, "realData": nil}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 2. 解析前端参数
	err := json.Unmarshal(rawData, &fp)
	if err != nil {
		code = -1
	} else {
		// 3. 校验参数长度
		if len(fp.DeviceList) == 0 || len(fp.DataList) == 0 || len(fp.DeviceList) != len(fp.DataList) {
			code = -7
		} else {
			// 4. 解析目标日期（带时区）
			queryDate, err := time.Parse("2006-01-02", fp.QueryDate)
			if err != nil {
				code = -8 // 日期格式错误
			} else {
				// 当天起止时间
				startTime := queryDate.Format("2006-01-02 00:00:00")
				endTime := queryDate.AddDate(0, 0, 1).Format("2006-01-02 00:00:00")
				// 5. 一次性批量获取所有 (DeviceName, DataName) 组合的数据
				pairIndexMap := make(map[string][]int)
				pairs := make([]models.HourDataPair, len(fp.DeviceList))
				hourSeries := make([][]any, len(fp.DeviceList))
				for i := range fp.DeviceList {
					deviceName := fp.DeviceList[i]
					dataName := fp.DataList[i]
					key := deviceName + "\x00" + dataName
					pairIndexMap[key] = append(pairIndexMap[key], i)
					pairs[i] = models.HourDataPair{DeviceName: deviceName, DataName: dataName}
					hourSeries[i] = make([]any, 24) // 24小时
					for j := range hourSeries[i] {
						hourSeries[i][j] = "-"
					}
				}

				var records []models.DevicesHistoryDataList
				var recordsCK []models.DevicesCHHistoryData
				if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
					code, records = models.GetMysqlDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 2 {
					code, records = models.GetTsDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
				} else if protocol_common.HistoryRecordDbType == 4 {
					code, records = models.GetInfluxDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 3 {
					code, recordsCK = models.GetClickHouseDifferenceDataBatch(pairs, startTime, endTime)
				}

				if code == 0 {
					if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 || protocol_common.HistoryRecordDbType == 2 || protocol_common.HistoryRecordDbType == 4 {
						for _, r := range records {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 计算该记录所属的小时（0-23）
							hourIndex := r.RecordTime.Hour()
							if hourIndex < 0 || hourIndex >= 24 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][hourIndex] = r.DataValue
							}
						}
					} else {
						for _, r := range recordsCK {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 计算该记录所属的小时（0-23）
							hourIndex := r.RecordTime.Hour()
							if hourIndex < 0 || hourIndex >= 24 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][hourIndex] = r.DataValue
							}
						}
					}
				}

				totalData = hourSeries
			}
		}
	}

	// 返回
	result = map[string]interface{}{
		"code":     code,
		"realData": totalData,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()
	_ = json.NewEncoder(gw).Encode(result)
}
func (c *ReportController) GetHistoryWeeklyDifferenceReport() {
	rawData := c.Ctx.Input.RequestBody
	var code int = 0
	var result map[string]interface{}

	// 前端传入参数
	type FrontendParams struct {
		QueryDate  []string `json:"QueryDate"`  // 格式：2026-04-04
		DeviceList []string `json:"DeviceList"` // 设备名
		DataList   []string `json:"DataList"`   // 数据名
	}

	var fp FrontendParams
	var totalData [][]any // 二维数组：[设备数][24小时]

	// 1. 校验ProjectUuid
	SProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if SProjectUuid == "" {
		result = map[string]interface{}{"code": -6, "realData": nil}
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 2. 解析前端参数
	err := json.Unmarshal(rawData, &fp)
	if err != nil {
		code = -1
	} else {
		// 3. 校验参数长度
		if len(fp.DeviceList) == 0 || len(fp.DataList) == 0 || len(fp.DeviceList) != len(fp.DataList) {
			code = -7
		} else if len(fp.QueryDate) < 2 {
			code = -8
		} else {
			// 4. 解析目标日期（带时区）
			queryDate, err := time.Parse("2006-01-02", fp.QueryDate[0])
			queryDate2, err2 := time.Parse("2006-01-02", fp.QueryDate[1])
			if err != nil || err2 != nil {
				code = -8 // 日期格式错误
			} else {
				// 当天起止时间
				startTime := queryDate.Format("2006-01-02 00:00:00")
				endTime := queryDate2.AddDate(0, 0, 1).Format("2006-01-02 00:00:00")

				// 5. 一次性批量获取所有 (DeviceName, DataName) 组合的数据
				pairIndexMap := make(map[string][]int)
				pairs := make([]models.HourDataPair, len(fp.DeviceList))
				hourSeries := make([][]any, len(fp.DeviceList))
				for i := range fp.DeviceList {
					deviceName := fp.DeviceList[i]
					dataName := fp.DataList[i]
					key := deviceName + "\x00" + dataName
					pairIndexMap[key] = append(pairIndexMap[key], i)
					pairs[i] = models.HourDataPair{DeviceName: deviceName, DataName: dataName}
					hourSeries[i] = make([]any, 7) // 7天
					for j := range hourSeries[i] {
						hourSeries[i][j] = "-"
					}
				}

				var records []models.DevicesHistoryDataList
				var recordsCK []models.DevicesCHHistoryData
				if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
					code, records = models.GetMysqlWeeklyDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 2 {
					code, records = models.GetTsWeeklyDifferenceDataBatch(pairs, startTime, endTime, protocol_common.HistoryRecordTsDb)
				} else if protocol_common.HistoryRecordDbType == 4 {
					code, records = models.GetInfluxWeeklyDifferenceDataBatch(pairs, startTime, endTime)
				} else if protocol_common.HistoryRecordDbType == 3 {
					code, recordsCK = models.GetClickHouseWeeklyDifferenceDataBatch(pairs, startTime, endTime)
				}

				if code == 0 {
					if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 || protocol_common.HistoryRecordDbType == 2 || protocol_common.HistoryRecordDbType == 4 {
						for _, r := range records {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 计算该日期是星期几（0=周日, 1=周一, ..., 6=周六）
							dayOfWeek := int(r.RecordTime.Weekday())
							if dayOfWeek == 0 {
								dayOfWeek = 7 // 调整为1-7（周一到周日）
							}
							dayIndex := dayOfWeek - 1 // 转换为0-6索引
							if dayIndex < 0 || dayIndex >= 7 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][dayIndex] = r.DataValue
							}
						}
					} else {
						for _, r := range recordsCK {
							key := r.DeviceName + "\x00" + r.DataName
							idxs, ok := pairIndexMap[key]
							if !ok {
								continue
							}
							// 计算该日期是星期几（0=周日, 1=周一, ..., 6=周六）
							dayOfWeek := int(r.RecordTime.Weekday())
							if dayOfWeek == 0 {
								dayOfWeek = 7 // 调整为1-7（周一到周日）
							}
							dayIndex := dayOfWeek - 1 // 转换为0-6索引
							if dayIndex < 0 || dayIndex >= 7 {
								continue
							}
							for _, idx := range idxs {
								hourSeries[idx][dayIndex] = r.DataValue
							}
						}
					}
				}

				totalData = hourSeries
			}
		}
	}

	// 返回
	result = map[string]interface{}{
		"code":     code,
		"realData": totalData,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()
	_ = json.NewEncoder(gw).Encode(result)
}

func (c *ReportController) GetTrendChartDataByDate() {

	rawData := c.Ctx.Input.RequestBody
	var result map[string]interface{}

	if protocol_common.HistoryRecordDbType == 1 || protocol_common.HistoryRecordDbType == 5 {
		code, data := models.GetMysqlTrendChartDataByDate(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 2 {
		code, data := models.GetTsTrendChartDataByDate(rawData, protocol_common.HistoryRecordTsDb)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 4 {
		code, data := models.GetInfluxTrendChartDataByDate(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	} else if protocol_common.HistoryRecordDbType == 3 {
		code, data := models.GetClickHouseTrendChartDataByDate(rawData)
		result = map[string]interface{}{
			"code": code,
			"data": data,
		}
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}
