/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-24 15:43:37
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/middleware"
	"ISMServer/models"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	protocolCommon "ISMServer/protocol/common"
	"ISMServer/utils/errmsg"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	license "ISMServer/license"

	sysfont "github.com/adrg/sysfont"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/denisbrodbeck/machineid"
	"github.com/forgoer/openssl"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/mattn/anko/vm"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	gopsutilnet "github.com/shirou/gopsutil/v3/net"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

type ISMSystem struct {
	beego.Controller
}

var upgradeDir string = "data/upgrade/"
var tempDir string = "data/tempDir/"
var dbbackup string = "data/dbbackup/"

// SM4加密
func SM4Encrypt(data string) (result string, err error) {
	//字符串转byte切片
	plainText := []byte(data)
	//建议从配置文件中读取秘钥，进行统一管理
	SM4Key := "WFQKl64#kv!@Jy*L"
	//todo 注意：iv需要是随机的，进一步保证加密的安全性，将iv的值和加密后的数据一起返回给外部
	SM4Iv := "Q0UTr$9YAztk&glp"
	iv := []byte(SM4Iv)
	key := []byte(SM4Key)
	//实例化sm4加密对象
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("加密失败")
	}
	//明文数据填充
	paddingData := paddingLastGroup(plainText, block.BlockSize())
	//声明SM4的加密工作模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//为填充后的数据进行加密处理
	cipherText := make([]byte, len(paddingData))
	//使用CryptBlocks这个核心方法，将paddingData进行加密处理，将加密处理后的值赋值到cipherText中
	blockMode.CryptBlocks(cipherText, paddingData)
	//加密结果使用hex转成字符串，方便外部调用
	cipherString := hex.EncodeToString(cipherText)
	return cipherString, nil
}

// SM4解密 传入string 输出string
func SM4Decrypt(data string) (res string, err error) {
	//秘钥
	SM4Key := "WFQKl64#kv!@Jy*L"
	//iv是Initialization Vector，初始向量，
	SM4Iv := "Q0UTr$9YAztk&glp"
	iv := []byte(SM4Iv)
	key := []byte(SM4Key)
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("解密失败")
	}
	//使用hex解码
	decodeString, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	//CBC模式 优点：具有较好的安全性，能够隐藏明文的模式和重复性。 缺点：加密过程是串行的，不适合并行处理。
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//下文有详解这段代码的含义
	blockMode.CryptBlocks(decodeString, decodeString)
	plainText, errRe := unPaddingLastGroup(decodeString)
	if errRe == 0 {
		//直接返回字符串类型，方便外部调用
		return string(plainText), nil
	} else {
		return string(plainText), fmt.Errorf("解密失败")
	}
}

// 明文数据填充
func paddingLastGroup(plainText []byte, blockSize int) []byte {
	//1.计算最后一个分组中明文后需要填充的字节数
	padNum := blockSize - len(plainText)%blockSize
	//2.将字节数转换为byte类型
	char := []byte{byte(padNum)}
	//3.创建切片并初始化
	newPlain := bytes.Repeat(char, padNum)
	//4.将填充数据追加到原始数据后
	newText := append(plainText, newPlain...)
	return newText
}

// 去掉明文后面的填充数据
func unPaddingLastGroup(plainText []byte) ([]byte, int8) {
	//1.拿到切片中的最后一个字节
	length := len(plainText)
	lastChar := plainText[length-1]
	//2.将最后一个数据转换为整数
	number := int(lastChar)
	if length-number < 0 {
		return plainText, -1
	}
	return plainText[:length-number], 0
}

// 生成32位md5字串
func getMd5String(s string, upper bool, half bool) string {
	h := md5.New()
	h.Write([]byte(s))
	result := hex.EncodeToString(h.Sum(nil))
	if upper == true {
		result = strings.ToUpper(result)
	}
	if half == true {
		result = result[8:24]
	}
	return result
}
func generateNetMachineCode() (string, error) {
	// 获取所有网卡接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// 选择第一个非Loopback接口（如果有的话）
	var iface *net.Interface
	for _, intf := range interfaces {
		if intf.Flags&net.FlagLoopback == 0 {
			iface = &intf
			break
		}
	}

	// 获取网卡的MAC地址
	hwAddr := iface.HardwareAddr.String()

	// 生成机器码（这里使用MD5哈希函数作为示例）
	hash := md5.New()
	hash.Write([]byte(hwAddr))
	hashValue := hash.Sum(nil)
	machineCode := hex.EncodeToString(hashValue)

	return machineCode, nil
}

func generateCpuMachineCode() (string, error) {
	var PhysicalID string
	cpus, _ := cpu.Info()
	for _, c := range cpus {
		PhysicalID = PhysicalID + c.PhysicalID
	}

	hash := md5.New()
	hash.Write([]byte(PhysicalID))
	hashValue := hash.Sum(nil)
	machineCode := hex.EncodeToString(hashValue)

	return machineCode, nil
}

func (c *ISMSystem) GetPhysicalIDCheck() {
	var code int = 0
	result := map[string]interface{}{
		"code": 0,
		"id":   "",
	}
	if protocolCommon.IsLicense {
		if protocolCommon.IsAuthLimit && protocolCommon.IsAuthTimeLimit {
			code = -3
		} else {
			code = 0
		}
	} else {
		code = -1
	}

	result["code"] = code
	result["id"] = protocolCommon.ISMProtectedID
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func isExpired(createDateStr string, authDays int) (bool, error, time.Duration) {
	beijingLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return false, fmt.Errorf("加载北京时间时区失败: %v", err), 0
	}
	layout := "2006-01-02_15-04-05"
	isExpiredVar := false

	createDate, err := time.Parse(layout, createDateStr)
	if err != nil {
		return false, fmt.Errorf("解析创建时间失败：%v", err), 0
	}

	// 2. 计算过期时间（创建时间 + 授权天数）
	expireDate := createDate.AddDate(0, 0, authDays)
	// 3. 对比当前时间（使用本地时间或UTC，根据需求调整）
	now := time.Now().In(beijingLoc)
	remaining := expireDate.Sub(now)
	if remaining <= 0 {
		isExpiredVar = true
	} else {
		isExpiredVar = false
	}
	// 判断当前时间是否超过过期时间
	return isExpiredVar, nil, remaining
}
func (c *ISMSystem) WitePhysicalID() {

	var err_code = -1

	var ISMProtectedID string
	var getError error

	type authStu struct {
		AuthCode string `json:"authCode"`
	}
	result := map[string]interface{}{
		"code": err_code,
	}
	var paramsJson authStu

	CheckMachineCodeType, checkErr := config.Int("CheckMachineCodeType")
	if checkErr != nil {
		CheckMachineCodeType = 0
	}
	if CheckMachineCodeType == 1 {
		ISMProtectedID, getError = generateCpuMachineCode()
	} else {
		ISMProtectedID, getError = generateNetMachineCode()
	}

	if getError != nil {
		ISMProtectedID = "ISMProtectedID"
	}
	id, _ := machineid.ProtectedID(ISMProtectedID)
	mac_code := getMd5String(id, true, true)
	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &paramsJson)
	if err != nil {
		err_code = -3
	} else {
		authContent, deError := SM4Decrypt(paramsJson.AuthCode)
		if deError != nil {
			err_code = -2
		} else {

			var lisceseAuth map[string]interface{}

			jsonErr := json.Unmarshal([]byte(authContent), &lisceseAuth)
			if jsonErr != nil {
				err_code = -9
			} else {
				if lisceseAuth["ISMProtectedID"] != mac_code {
					err_code = -10
				} else {
					IsBatch := int(lisceseAuth["IsBatch"].(float64))
					if IsBatch == 1 {
						protocolCommon.IsAuthLimit = false
						protocolCommon.IsAuthTimeLimit = false
					} else {
						IsAuthLimit, ok1 := lisceseAuth["IsAuthLimit"]
						CreateDate, okDate := lisceseAuth["CreateDate"]
						AuthorizationDays, ok2 := lisceseAuth["AuthorizationDays"]
						if okDate && ok1 && ok2 {
							limitDays := int(IsAuthLimit.(float64))
							AuthorizationDaysInt := int(AuthorizationDays.(float64))
							AuthorizationCreateDate := CreateDate.(string)
							if limitDays == 1 && AuthorizationDaysInt > 0 {
								protocolCommon.IsAuthLimit = true
								protocolCommon.AuthorizationDays = AuthorizationDaysInt
								protocolCommon.AuthorizationCreateDate = AuthorizationCreateDate
								isExpire, expireErr, remainTime := isExpired(AuthorizationCreateDate, AuthorizationDaysInt)

								days := int(remainTime.Hours() / 24)
								remainingAfterDays := remainTime - time.Duration(days)*24*time.Hour

								hours := int(remainingAfterDays.Hours())
								remainingAfterHours := remainingAfterDays - time.Duration(hours)*time.Hour

								minutes := int(remainingAfterHours.Minutes())

								protocolCommon.AuthRemainingTimeDays = days
								protocolCommon.AuthRemainingTimeHours = hours
								protocolCommon.AuthRemainingTimeMinutes = minutes

								if expireErr == nil && isExpire {
									protocolCommon.IsAuthTimeLimit = true
								} else {
									protocolCommon.IsAuthTimeLimit = false
								}
							} else {
								protocolCommon.IsAuthLimit = false
							}
						} else {
							protocolCommon.IsAuthLimit = false
						}
					}

					err2 := os.WriteFile("data/auth/active.dat", []byte(paramsJson.AuthCode), 0666) //写入文件(字节数组)
					if err2 != nil {
						err_code = -5
					} else {
						protocolCommon.IsLicense = true
						err_code = 0
					}
				}
			}

		}
	}
	result["code"] = err_code
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *ISMSystem) RolesList() {

	result := map[string]interface{}{
		"code": 0,
		"list": models.SystemRoleList(),
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *ISMSystem) GetSystemFonts() {

	result := map[string]interface{}{
		"code": 0,
	}
	finder := sysfont.NewFinder(nil)

	for _, font := range finder.List() {
		fmt.Println(font.Family, font.Name, font.Filename)
	}
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *ISMSystem) SetDebug() {

	result := map[string]interface{}{
		"code": 0,
	}
	type DebugStu struct {
		DebugEnable bool `json:"DebugEnable"`
	}
	var DebugParam DebugStu
	data := c.Ctx.Input.RequestBody
	//json数据封装到user对象中
	err := json.Unmarshal(data, &DebugParam)
	if err != nil {
		result := map[string]interface{}{
			"code": -1,
			"msg":  "参数错误",
		}
		c.Data["json"] = result
		c.ServeJSON()
		return
	} else {
		protocolCommon.ModbusDebug = DebugParam.DebugEnable
		protocolCommon.IEC104Debug = DebugParam.DebugEnable

		protocolCommonFunc.CloseChanel()
	}
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *ISMSystem) GetSystemData() {

	var getLists []models.SystemDataModel
	var code int64

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		getLists = models.GetSystemDataList(ProjectUuid)

	} else {
		code = -1
		getLists = nil
	}
	result := map[string]interface{}{
		"code": code,
		"list": getLists,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *ISMSystem) OnlineCheckUpgrade() {

	var getLists []models.SystemDataModel
	var code int64
	result := map[string]interface{}{
		"code":   code,
		"Result": getLists,
	}
	upgradeServer, upgradeServererr := config.String("upgradeServer")
	if upgradeServererr != nil || upgradeServer == "" {
		upgradeServer = "http://www.ismctl.com/ism"
	}
	url := upgradeServer + "/ismUpgrade_V3.json"
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		result["code"] = -2
		c.Data["json"] = result
		c.ServeJSON() //返回json格式
		return
	}
	//得到返回结果
	ResBody, _ := ioutil.ReadAll(res.Body)
	//对返回的json数据做解析
	var dataAttr map[string]interface{}

	if err := json.Unmarshal(ResBody, &dataAttr); err == nil {
		result["Result"] = dataAttr
	} else {
		result["code"] = -1
	}
	res.Body.Close()
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *ISMSystem) LocalUpgrade() {

	type UploadResult struct {
		Code int
	}
	var reponse_result UploadResult

	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".zip": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	//创建目录
	uploadDir := upgradeDir
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + "lastversion.zip"
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	reponse_result.Code = 0
	c.Data["json"] = reponse_result

	c.ServeJSON() //返回json格式
	if reponse_result.Code == 0 {
		os.Exit(0)
	}
}
func (c *ISMSystem) UpdateDataModel() {

	type UploadResult struct {
		Code int
	}
	var reponse_result UploadResult

	suuid := c.Ctx.Input.Param(":muid")
	if suuid == "" {
		reponse_result.Code = -6
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".xlsx": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	//创建目录
	uploadDir := tempDir
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + h.Filename
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	excelfile, err := excelize.OpenFile(fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	var DataModelType int
	var DataModelTypeIndex int = -1
	var DataModelTypeErr error
	var convErr error

	for _, sheetName := range excelfile.GetSheetList() {
		rows, err := excelfile.GetRows(sheetName)
		if err != nil {
			continue
		}
		for index, row := range rows {
			// 跳过第一行标题
			if index == 0 {

				for k, v := range row {
					if v == "模型类型(勿修改)" {
						DataModelTypeIndex = k
						break
					}
				}
				continue
			}
			if DataModelTypeIndex == -1 {
				reponse_result.Code = -6
				c.Data["json"] = reponse_result
				c.ServeJSON()
				return
			}
			if index == 1 {
				DataModelType, DataModelTypeErr = strconv.Atoi(row[DataModelTypeIndex])
				if DataModelTypeErr != nil {
					reponse_result.Code = -6
					c.Data["json"] = reponse_result
					c.ServeJSON()
					return
				}
			}
			//SNMP设备
			if DataModelType == 1 {

			} else if DataModelType == 2 { //Modbus设备
				var setparams models.ModbusDevicesDataModel
				setparams.Name = row[0]
				setparams.RegisterAddress, convErr = strconv.Atoi(row[1])
				if convErr != nil {
					continue
				}
				setparams.Auth = row[2]
				setparams.Type = row[3]
				setparams.ByteOrder = row[4]
				setparams.DataUnit = row[5]
				setparams.ConversionExpression = row[6]

				if row[7] == "是" {
					setparams.IsAlarm = 1
				} else if row[7] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[8] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[8] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[8] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[8] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[8] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[9]
				setparams.AlarmClearMessage = row[10]

				if row[11] == "是" {
					setparams.IsRecord = 1
				} else if row[11] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[12] == "定时存储" {
					setparams.RecordType = 1
				} else if row[12] == "即时存储" {
					setparams.RecordType = 2
				} else if row[12] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[13])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[14]
				setparams.FloatAccuracy = row[15]
				if len(row) > 17 {
					setparams.Uuid = row[17]
				} else {
					setparams.Uuid = ""
				}
				models.ModbusRegisterAddressUpdate(setparams)
			} else if DataModelType == 3 { //OPCUA设备
				var setparams models.OpcuaDevicesDataModel

				setparams.Name = row[0]
				setparams.Nodeid = row[1]

				setparams.Auth = row[2]
				if row[3] == "Boolean" {
					setparams.Type = "1"
				} else if row[3] == "SByte" {
					setparams.Type = "2"
				} else if row[3] == "Byte" {
					setparams.Type = "3"
				} else if row[3] == "Int16" {
					setparams.Type = "4"
				} else if row[3] == "UInt16" {
					setparams.Type = "5"
				} else if row[3] == "Int32" {
					setparams.Type = "6"
				} else if row[3] == "UInt32" {
					setparams.Type = "7"
				} else if row[3] == "Int64" {
					setparams.Type = "8"
				} else if row[3] == "UInt64" {
					setparams.Type = "9"
				} else if row[3] == "Float" {
					setparams.Type = "10"
				} else if row[3] == "Double" {
					setparams.Type = "11"
				} else if row[3] == "String" {
					setparams.Type = "12"
				} else {
					setparams.Type = "1"
				}

				setparams.DataUnit = row[4]
				setparams.ConversionExpression = row[5]

				if row[6] == "是" {
					setparams.IsAlarm = 1
				} else if row[6] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[7] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[7] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[7] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[7] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[7] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[8]
				setparams.AlarmClearMessage = row[9]

				if row[10] == "是" {
					setparams.IsRecord = 1
				} else if row[10] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[11] == "定时存储" {
					setparams.RecordType = 1
				} else if row[11] == "即时存储" {
					setparams.RecordType = 2
				} else if row[11] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[12])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[13]
				setparams.FloatAccuracy = row[14]
				setparams.ModelType = 3
				if len(row) > 16 {
					setparams.Uuid = row[16]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {
					var getDevicesDataModel models.OpcuaDevicesDataModel
					existData := models.Db.Model(&models.OpcuaDevicesDataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.OpcuaNodeAdd(setparams)
					} else {
						models.OpcuaNodeEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.OpcuaNodeAdd(setparams)
				}

			} else if DataModelType == 5 { //RESTFul设备

			} else if DataModelType == 6 { //静态数据

			} else if DataModelType == 7 { //自定义数据

			} else if DataModelType == 15 { //西门子PLC
				var setparams models.SimS7DataModel

				setparams.Name = row[0]
				if row[1] == "DB" {
					setparams.DataFromType = 0
				} else if row[1] == "I区" {
					setparams.DataFromType = 1
				} else if row[1] == "Q区" {
					setparams.DataFromType = 2
				} else if row[1] == "M区" {
					setparams.DataFromType = 3
				} else {
					setparams.DataFromType = 0
				}
				setparams.DBIndex, convErr = strconv.Atoi(row[2])
				if convErr != nil {
					setparams.DBIndex = 0
				}
				setparams.DBOffset = row[3]
				setparams.StringMaxLength, convErr = strconv.Atoi(row[4])
				if convErr != nil {
					setparams.StringMaxLength = 255
				}
				if row[5] == "有符号" {
					setparams.IsHaveUnsigned = 0
				} else if row[5] == "无符号" {
					setparams.IsHaveUnsigned = 1
				} else {
					setparams.IsHaveUnsigned = 0
				}

				setparams.Auth = row[6]

				if row[7] == "Bool" {
					setparams.Type = "0"
				} else if row[7] == "Byte" {
					setparams.Type = "1"
				} else if row[7] == "SINT" {
					setparams.Type = "2"
				} else if row[7] == "INT" {
					setparams.Type = "3"
				} else if row[7] == "DINT" {
					setparams.Type = "4"
				} else if row[7] == "USINT" {
					setparams.Type = "5"
				} else if row[7] == "UINT" {
					setparams.Type = "6"
				} else if row[7] == "UDINT" {
					setparams.Type = "7"
				} else if row[7] == "LINT" {
					setparams.Type = "8"
				} else if row[7] == "REAL" {
					setparams.Type = "9"
				} else if row[7] == "LREAL" {
					setparams.Type = "10"
				} else if row[7] == "TIME" {
					setparams.Type = "11"
				} else if row[7] == "LTIME" {
					setparams.Type = "12"
				} else if row[7] == "S5TIME" {
					setparams.Type = "13"
				} else if row[7] == "DATE" {
					setparams.Type = "14"
				} else if row[7] == "DATE_AND_TIME" {
					setparams.Type = "15"
				} else if row[7] == "WSTRING" {
					setparams.Type = "16"
				} else if row[7] == "String" {
					setparams.Type = "17"
				} else if row[7] == "ULINT" {
					setparams.Type = "18"
				} else if row[7] == "WORD" {
					setparams.Type = "19"
				} else if row[7] == "DWORD" {
					setparams.Type = "20"
				} else if row[7] == "LWORD" {
					setparams.Type = "21"
				} else if row[7] == "CHAR" {
					setparams.Type = "22"
				} else if row[7] == "WCHAR" {
					setparams.Type = "23"
				}

				setparams.DataUnit = row[8]
				setparams.ConversionExpression = row[9]

				if row[10] == "是" {
					setparams.IsAlarm = 1
				} else if row[10] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[11] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[11] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[11] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[11] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[11] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[12]
				setparams.AlarmClearMessage = row[13]

				if row[14] == "是" {
					setparams.IsRecord = 1
				} else if row[14] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[15] == "定时存储" {
					setparams.RecordType = 1
				} else if row[15] == "即时存储" {
					setparams.RecordType = 2
				} else if row[15] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[16])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[17]
				setparams.FloatAccuracy = row[18]
				setparams.ModelType = 15
				if len(row) > 20 {
					setparams.Uuid = row[20]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {
					var getDevicesDataModel models.SimS7DataModel
					existData := models.Db.Model(&models.SimS7DataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.SimS7ModelDataAdd(setparams)
					} else {
						models.SimS7ModelDataEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.SimS7ModelDataAdd(setparams)
				}
			} else if DataModelType == 20 { //MQTT设备
				var setparams models.MqttDevicesDataModel

				setparams.Name = row[0]
				setparams.Identifier = row[1]

				setparams.Auth = row[2]
				if row[3] == "Boolean" {
					setparams.Type = "1"
				} else if row[3] == "Int" {
					setparams.Type = "2"
				} else if row[3] == "Float" {
					setparams.Type = "3"
				} else if row[3] == "Double" {
					setparams.Type = "4"
				} else if row[3] == "String" {
					setparams.Type = "5"
				} else {
					setparams.Type = "1"
				}

				setparams.DataUnit = row[4]
				setparams.ConversionExpression = row[5]

				if row[6] == "是" {
					setparams.IsAlarm = 1
				} else if row[6] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[7] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[7] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[7] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[7] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[7] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[8]
				setparams.AlarmClearMessage = row[9]

				if row[10] == "是" {
					setparams.IsRecord = 1
				} else if row[10] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[11] == "定时存储" {
					setparams.RecordType = 1
				} else if row[11] == "即时存储" {
					setparams.RecordType = 2
				} else if row[11] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[12])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[13]
				setparams.FloatAccuracy = row[14]
				setparams.ModelType = 3
				if len(row) > 16 {
					setparams.Uuid = row[16]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {
					var getDevicesDataModel models.MqttDevicesDataModel
					existData := models.Db.Model(&models.MqttDevicesDataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.MqttNodeAdd(setparams)
					} else {
						models.MqttNodeEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.MqttNodeAdd(setparams)
				}
			} else if DataModelType == 30 { //DLT645电表
				var setparams models.Dlt645DevicesDataModel

				setparams.Name = row[0]

				setparams.DataIdentification = row[1]

				setparams.Auth = row[2]

				if row[3] == "Boolean" {
					setparams.Type = "1"
				} else if row[3] == "SByte" {
					setparams.Type = "2"
				} else if row[3] == "Byte" {
					setparams.Type = "3"
				} else if row[3] == "Int16" {
					setparams.Type = "4"
				} else if row[3] == "UInt16" {
					setparams.Type = "5"
				} else if row[3] == "Int32" {
					setparams.Type = "6"
				} else if row[3] == "UInt32" {
					setparams.Type = "7"
				} else if row[3] == "Int64" {
					setparams.Type = "8"
				} else if row[3] == "UInt64" {
					setparams.Type = "9"
				} else if row[3] == "Float" {
					setparams.Type = "10"
				} else if row[3] == "Double" {
					setparams.Type = "11"
				} else if row[3] == "String" {
					setparams.Type = "12"
				} else {
					continue
				}

				setparams.DataUnit = row[4]
				setparams.ConversionExpression = row[5]

				if row[6] == "是" {
					setparams.IsAlarm = 1
				} else if row[6] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[7] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[7] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[7] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[7] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[7] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[8]
				setparams.AlarmClearMessage = row[8]

				if row[10] == "是" {
					setparams.IsRecord = 1
				} else if row[10] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[11] == "定时存储" {
					setparams.RecordType = 1
				} else if row[11] == "即时存储" {
					setparams.RecordType = 2
				} else if row[11] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[12])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[13]
				// setparams.FloatAccuracy = row[13]
				setparams.ModelType = 30
				if len(row) > 16 {
					setparams.Uuid = row[16]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {
					var getDevicesDataModel models.Dlt645DevicesDataModel
					existData := models.Db.Model(&models.Dlt645DevicesDataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.DLT645NodeAdd(setparams)
					} else {
						models.DLT645NodeEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.DLT645NodeAdd(setparams)
				}
			} else if DataModelType == 40 { //IEC104电力规约
				var setparams models.IEC104DevicesDataModel

				setparams.Name = row[0]

				if row[1] == "遥信" {
					setparams.DataCategory = 1
				} else if row[1] == "遥测" {
					setparams.DataCategory = 2
				} else if row[1] == "脉冲(电量)" {
					setparams.DataCategory = 3
				} else if row[1] == "遥控" {
					setparams.DataCategory = 4
				} else if row[1] == "遥调(设点)" {
					setparams.DataCategory = 5
				} else {
					continue
				}
				setparams.DataPoint, convErr = strconv.Atoi(row[2])
				if convErr != nil {
					continue
				}

				setparams.Auth = row[3]

				if row[4] == "Float" {
					setparams.Type = "10"
				} else if row[4] == "Int" {
					setparams.Type = "8"
				} else {
					continue
				}

				setparams.DataUnit = row[5]
				setparams.ConversionExpression = row[6]

				if row[7] == "是" {
					setparams.IsAlarm = 1
				} else if row[7] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[8] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[8] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[8] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[8] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[8] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[9]
				setparams.AlarmClearMessage = row[10]

				if row[11] == "是" {
					setparams.IsRecord = 1
				} else if row[11] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[12] == "定时存储" {
					setparams.RecordType = 1
				} else if row[12] == "即时存储" {
					setparams.RecordType = 2
				} else if row[12] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[13])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[14]
				setparams.FloatAccuracy = row[15]
				setparams.ModelType = 40
				if len(row) > 17 {
					setparams.Uuid = row[17]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {
					var getDevicesDataModel models.IEC104DevicesDataModel
					existData := models.Db.Model(&models.IEC104DevicesDataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.IEC104NodeAdd(setparams)
					} else {
						models.IEC104NodeEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.IEC104NodeAdd(setparams)
				}
			} else if DataModelType == 350 { //IEC61850设备
				var setparams models.IEC61850DevicesDataModel

				setparams.Name = row[0]
				setparams.Nodeid = row[1]

				setparams.Auth = row[2]
				if row[3] == "Boolean" {
					setparams.Type = "1"
				} else if row[3] == "SByte" {
					setparams.Type = "2"
				} else if row[3] == "Byte" {
					setparams.Type = "3"
				} else if row[3] == "Int16" {
					setparams.Type = "4"
				} else if row[3] == "UInt16" {
					setparams.Type = "5"
				} else if row[3] == "Int32" {
					setparams.Type = "6"
				} else if row[3] == "UInt32" {
					setparams.Type = "7"
				} else if row[3] == "Int64" {
					setparams.Type = "8"
				} else if row[3] == "UInt64" {
					setparams.Type = "9"
				} else if row[3] == "Float" {
					setparams.Type = "10"
				} else if row[3] == "Double" {
					setparams.Type = "11"
				} else if row[3] == "String" {
					setparams.Type = "12"
				} else {
					setparams.Type = "1"
				}

				setparams.DataUnit = row[4]
				setparams.ConversionExpression = row[5]

				if row[6] == "是" {
					setparams.IsAlarm = 1
				} else if row[6] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[7] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[7] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[7] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[7] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[7] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[8]
				setparams.AlarmClearMessage = row[9]

				if row[10] == "是" {
					setparams.IsRecord = 1
				} else if row[10] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[11] == "定时存储" {
					setparams.RecordType = 1
				} else if row[11] == "即时存储" {
					setparams.RecordType = 2
				} else if row[11] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[12])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[13]
				setparams.FloatAccuracy = row[14]
				setparams.ModelType = 350
				if len(row) > 16 {
					setparams.Uuid = row[16]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {

					var getDevicesDataModel models.IEC61850DevicesDataModel
					existData := models.Db.Model(&models.IEC61850DevicesDataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.IEC61850NodeAdd(setparams)
					} else {
						models.IEC61850NodeEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.IEC61850NodeAdd(setparams)
				}

			} else if DataModelType == 470 { //HJ212设备
				var setparams models.HJ212DevicesDataModel

				setparams.Name = row[0]
				setparams.EncodeID = row[1]

				setparams.Auth = row[2]
				if row[3] == "Boolean" {
					setparams.Type = "1"
				} else if row[3] == "Int" {
					setparams.Type = "2"
				} else if row[3] == "Float" {
					setparams.Type = "3"
				} else if row[3] == "Double" {
					setparams.Type = "4"
				} else if row[3] == "String" {
					setparams.Type = "5"
				} else {
					setparams.Type = "1"
				}

				setparams.DataUnit = row[4]
				setparams.ConversionExpression = row[5]

				if row[6] == "是" {
					setparams.IsAlarm = 1
				} else if row[6] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[7] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[7] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[7] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[7] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[7] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[8]
				setparams.AlarmClearMessage = row[9]

				if row[10] == "是" {
					setparams.IsRecord = 1
				} else if row[10] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[11] == "定时存储" {
					setparams.RecordType = 1
				} else if row[11] == "即时存储" {
					setparams.RecordType = 2
				} else if row[11] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[12])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[13]
				setparams.FloatAccuracy = row[14]
				setparams.ModelType = 470
				if len(row) > 16 {
					setparams.Uuid = row[16]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {
					var getDevicesDataModel models.HJ212DevicesDataModel
					existData := models.Db.Model(&models.HJ212DevicesDataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.HJ212NodeAdd(setparams)
					} else {
						models.HJ212NodeEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.HJ212NodeAdd(setparams)
				}
			} else if DataModelType == 500 { //BACnet
				var setparams models.BacnetDevicesDataModel

				setparams.Name = row[0]
				BACnetZone := row[1]
				if BACnetZone == "AI" {
					setparams.BacnetZone = 0
				} else if BACnetZone == "AO" {
					setparams.BacnetZone = 1
				} else if BACnetZone == "AV" {
					setparams.BacnetZone = 2
				} else if BACnetZone == "BI" {
					setparams.BacnetZone = 3
				} else if BACnetZone == "BO" {
					setparams.BacnetZone = 4
				} else if BACnetZone == "BV" {
					setparams.BacnetZone = 5
				} else if BACnetZone == "MSI" {
					setparams.BacnetZone = 6
				} else if BACnetZone == "MSO" {
					setparams.BacnetZone = 7
				} else if BACnetZone == "MSV" {
					setparams.BacnetZone = 8
				} else if BACnetZone == "DEV" {
					setparams.BacnetZone = 9
				} else if BACnetZone == "ACC" {
					setparams.BacnetZone = 10
				}
				BACnetAddress, convErr := strconv.Atoi(row[2])
				if convErr != nil {
					continue
				} else {
					setparams.BacnetAddress = BACnetAddress
				}

				setparams.Auth = row[3]

				if row[4] == "Bit" {
					setparams.Type = "1"
				} else if row[4] == "Int" {
					setparams.Type = "2"
				} else if row[4] == "Float" {
					setparams.Type = "3"
				}

				setparams.DataUnit = row[5]
				setparams.ConversionExpression = row[6]

				if row[7] == "是" {
					setparams.IsAlarm = 1
				} else if row[7] == "否" {
					setparams.IsAlarm = 0
				} else {
					setparams.IsAlarm = 0
				}

				if row[8] == "提示" {
					setparams.AlarmLevel = 0
				} else if row[8] == "次要" {
					setparams.AlarmLevel = 1
				} else if row[8] == "重要" {
					setparams.AlarmLevel = 2
				} else if row[8] == "紧急" {
					setparams.AlarmLevel = 3
				} else if row[8] == "致命" {
					setparams.AlarmLevel = 4
				} else {
					setparams.AlarmLevel = 0
				}

				setparams.AlarmMessage = row[9]
				setparams.AlarmClearMessage = row[10]

				if row[11] == "是" {
					setparams.IsRecord = 1
				} else if row[11] == "否" {
					setparams.IsRecord = 0
				} else {
					setparams.IsRecord = 0
				}
				if row[12] == "定时存储" {
					setparams.RecordType = 1
				} else if row[12] == "即时存储" {
					setparams.RecordType = 2
				} else if row[12] == "变化存储" {
					setparams.RecordType = 0
				} else {
					setparams.RecordType = 0
				}

				setparams.RecordInterval, convErr = strconv.Atoi(row[13])
				if convErr != nil {
					setparams.RecordInterval = 60
				}
				setparams.RecordDataCharge = row[14]
				setparams.ModelType = 500
				if len(row) > 16 {
					setparams.Uuid = row[16]
				} else {
					setparams.Uuid = ""
				}
				setparams.Muid = suuid
				if setparams.Uuid != "" {
					var getDevicesDataModel models.BacnetDevicesDataModel
					existData := models.Db.Model(&models.BacnetDevicesDataModel{}).Where("uuid = ? and muid=?", setparams.Uuid, setparams.Muid).First(&getDevicesDataModel).Error
					if existData == gorm.ErrRecordNotFound {
						models.BACnetNodeAdd(setparams)
					} else {
						models.BACnetNodeEdit(setparams.Muid, setparams.Uuid, setparams)
					}
				} else if setparams.Muid != "" {
					models.BACnetNodeAdd(setparams)
				}
			}
		}
	}

	reponse_result.Code = 0
	c.Data["json"] = reponse_result

	c.ServeJSON() //返回json格式
}
func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// 判断文件或目录是否存在
func GetFileInfo(src string) os.FileInfo {
	if fileInfo, e := os.Stat(src); e != nil {
		if os.IsNotExist(e) {
			return nil
		}
		return nil
	} else {
		return fileInfo
	}
}

// 拷贝文件
func CopyFile(src, dst string) bool {
	if len(src) == 0 || len(dst) == 0 {
		return false
	}
	srcFile, e := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if e != nil {
		logs.Error("copyfile", e)
		return false
	}
	defer srcFile.Close()

	dst = strings.Replace(dst, "\\", "/", -1)
	dstPathArr := strings.Split(dst, "/")
	dstPathArr = dstPathArr[0 : len(dstPathArr)-1]
	dstPath := strings.Join(dstPathArr, "/")

	dstFileInfo := GetFileInfo(dstPath)
	if dstFileInfo == nil {
		if e := os.MkdirAll(dstPath, os.ModePerm); e != nil {
			logs.Error("copyfile", e)
			return false
		}
	}
	//这里要把O_TRUNC 加上，否则会出现新旧文件内容出现重叠现象
	dstFile, e := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if e != nil {
		logs.Error("copyfile", e)
		return false
	}
	defer dstFile.Close()
	//fileInfo, e := srcFile.Stat()
	//fileInfo.Size() > 1024
	//byteBuffer := make([]byte, 10)
	if _, e := io.Copy(dstFile, srcFile); e != nil {
		logs.Error("copyfile", e)
		return false
	} else {
		return true
	}

}
func CopyPath(src, dst string) bool {
	srcFileInfo := GetFileInfo(src)
	if srcFileInfo == nil || !srcFileInfo.IsDir() {
		return false
	}
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relationPath := strings.Replace(path, src, "/", -1)
		dstPath := strings.TrimRight(strings.TrimRight(dst, "/"), "\\") + relationPath
		if !info.IsDir() {
			if CopyFile(path, dstPath) {
				return nil
			} else {
				return errors.New(path + " copy fail")
			}
		} else {
			if _, err := os.Stat(dstPath); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
						return err
					} else {
						return nil
					}
				} else {
					return err
				}
			} else {
				return nil
			}
		}
	})

	if err != nil {
		return false
	}
	return true

}
func UnzipAndUpgrade(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// BeginUpgrade 开始升级ISMSystem系统
func (c *ISMSystem) BeginUpgrade() int {
	var code int = 0
	result := map[string]interface{}{
		"code": code,
	}
	if !protocolCommon.IsLicense {
		result["code"] = errmsg.LOGIN_NO_AUTH
		c.Data["json"] = result
		c.ServeJSON() //返回json格式
		return 0
	}
	_, err3 := os.Stat(upgradeDir)

	if os.IsNotExist(err3) {
		os.Mkdir(upgradeDir, os.ModePerm)
	}
	upgradeServer, upgradeServererr := config.String("upgradeServer")
	if upgradeServererr != nil || upgradeServer == "" {
		upgradeServer = "http://www.ismctl.com/ism"
	}
	var downloadFilePath string
	if (runtime.GOARCH != "arm64") && (runtime.GOARCH != "arm") {
		downloadFilePath = upgradeServer + "/lastversion_V3.zip"
	} else {
		downloadFilePath = upgradeServer + "/lastversion_arm_V3.zip"
	}
	logs.Info("升级下载路径", downloadFilePath)
	err := downloadFile(upgradeDir+"lastversion.zip", downloadFilePath)
	if err != nil {
		code = -1
		return -2
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "升级了系统", errmsg.JournalLevelInfo, c.Ctx.Input)
	result["code"] = code
	c.Data["json"] = result

	c.ServeJSON() //返回json格式
	time.Sleep(time.Second * 5)
	os.Exit(0)
	return 0
}
func (c *ISMSystem) ExecSysScript() {
	type ExecScriptList struct {
		Script []string `json:"Script"`
	}
	var code int = 0
	var message string = "成功"
	var recvList ExecScriptList
	var GetScriptList []models.ISMScript
	var isError int = 0
	var exceResult any
	data := c.Ctx.Input.RequestBody
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := json.Unmarshal(data, &recvList)
		if err != nil {
			code = -1
			message = "JSON格式错误"
		} else {
			e := protocolCommonFunc.ScriptDefine()
			err := models.Db.Model(&models.ISMScript{}).Where("script_uuid in ?", recvList.Script).Find(&GetScriptList).Error
			if err == nil && len(GetScriptList) > 0 {
				for _, v := range GetScriptList {
					tempComponents, deErr := base64.StdEncoding.DecodeString(v.ScriptContent)
					if deErr == nil {
						v.ScriptContent = string(tempComponents)
					}
					exceResult, err = vm.Execute(e, nil, v.ScriptContent)
					if err != nil {
						isError = 1
					}
				}
				if isError == 1 {
					code = -5
					message = "脚本执行有错误，没有完全执行!!!"
				}
			} else {
				code = -4
				message = "脚本不存在"
			}
		}
	} else {
		code = -1
		message = "缺少项目ID"
	}
	result := map[string]interface{}{
		"code":   code,
		"msg":    message,
		"result": exceResult,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}
func readDirFile(dir string) []string {

	var filesPath []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return filesPath
	}

	for _, file := range files {
		if !file.IsDir() {
			filesPath = append(filesPath, dir+"/"+file.Name())
		}
	}
	return filesPath
}

func (c *ISMSystem) GetCustomPel() {
	type CustomPel struct {
		DirName  string   `json:"DirName"`
		FilePath []string `json:"FilePath"`
	}
	var systemPelPath string = "static/customPel/"
	var code int = 0
	var GetCustomPelList []CustomPel
	dirs, err := ioutil.ReadDir(systemPelPath)
	if err != nil {
		code = -1
	} else {
		for _, dir := range dirs {
			var SingleCustomPel CustomPel
			SingleCustomPel.DirName = dir.Name()
			SingleCustomPel.FilePath = readDirFile(systemPelPath + dir.Name())
			GetCustomPelList = append(GetCustomPelList, SingleCustomPel)
		}
	}

	result := map[string]interface{}{
		"code": code,
		"list": GetCustomPelList,
	}
	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}
func (c *ISMSystem) GetSystemDeviceInfo() {

	var code int = 0
	var message string = "成功"
	type deviceInfo struct {
		Name       string ` json:"name" `
		Uuid       string ` json:"Uuid" `
		Status     int    ` json:"Status" `
		Longitude  string ` json:"longitude" `
		Latitude   string ` json:"latitude" `
		AlarmCount int64  `json:"AlarmCount"`
	}
	var getMonitorList []deviceInfo
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	if ProjectUuid != "" {
		err := models.Db.Model(&models.MonitorList{}).Select("name,status,longitude,latitude,uuid").Where("type = ? and project_uuid=?", 1, ProjectUuid).Find(&getMonitorList).Error
		if err != nil {
			code = -1
			message = err.Error()
		}
		for k, v := range getMonitorList {
			models.Db.Model(&models.DevicesAlarmList{}).Where("device_uuid = ? and project_uuid=? and clear_time < ?", v.Uuid, ProjectUuid, "2007-01-02 15:04:05").Count(&getMonitorList[k].AlarmCount)
		}

	} else {
		code = -1
		message = "缺少项目ID"
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
		"list": getMonitorList,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}

// GetHostInfo 获取主机信息
func GetHostInfo() (*host.InfoStat, any) {
	hostInfo, _ := host.Info()

	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)

	return hostInfo, t.Local().Format("2006-01-02 15:04:05")
}

// GetSysLoad 获取系统负载信息
func GetSysLoad() (data []byte) {
	loadInfo, _ := load.Avg()

	fmt.Println(loadInfo.Load1)
	return
}

type CpuInfo struct {
	Number      int       //cup个数
	Cores       int32     //核数
	UsedPercent []float64 //cpu使用率
	Load        any
}

// GetCpuInfo 获取CPU信息
func GetCpuInfo() (data any) {
	var CpuInfoData CpuInfo
	cpus, _ := cpu.Info()
	for _, c := range cpus {
		CpuInfoData.Cores = CpuInfoData.Cores + c.Cores
	}
	CpuInfoData.Number = len(cpus)
	percent, _ := cpu.Percent(0, false) //获取CPU使用率
	CpuInfoData.UsedPercent = percent

	CpuInfoData.Load, _ = load.Avg()

	return CpuInfoData
}

// GetMemInfo 获取内存信息
func GetMemInfo() (data any) {
	hostInfo, _ := mem.VirtualMemory()
	tmpData := gconv.Map(hostInfo)
	var gomem runtime.MemStats
	runtime.ReadMemStats(&gomem)
	tmpData["goUsed"] = gomem.Sys
	return tmpData
}

// GetDiskInfo 获取磁盘信息
func GetDiskInfo() (data any) {
	//diskPart, _ := disk.Partitions(false)
	//for _, dp := range diskPart {
	//	fmt.Println(dp)
	//	diskUsed, _ := disk.Usage(dp.Mountpoint)
	//	fmt.Printf("分区总大小: %d MB \n", diskUsed.Total/1024/1024)
	//	fmt.Printf("分区使用率: %.3f %% \n", diskUsed.UsedPercent)
	//	fmt.Printf("分区inode使用率: %.3f %% \n", diskUsed.InodesUsedPercent)
	//}
	diskUsed, _ := disk.Usage("/")

	return diskUsed
}

// NetWorkInfo 网速信息
type NetWorkInfo struct {
	Name         string
	Receive      uint64
	Sent         uint64
	ReceiveSpeed uint64
	SentSpeed    uint64
}

// GetNetStatusInfo 获取网络信息
func GetNetStatusInfo() (data any) {

	IOCountersStat, _ := gopsutilnet.IOCounters(true)
	var netWorkInfo NetWorkInfo
	for _, n := range IOCountersStat {
		netWorkInfo.Name = n.Name
		netWorkInfo.Receive = n.BytesRecv + netWorkInfo.Receive
		netWorkInfo.Sent = n.BytesSent + netWorkInfo.Sent
	}
	return netWorkInfo
}

func (c *ISMSystem) GetSystemAnalysis() {

	var code int = 0
	var message string = "成功"
	type analysisInfo struct {
		DeviceCount    int64          ` json:"DeviceCount" `
		DeviceOffCount int64          ` json:"DeviceOffCount" `
		AppCount       int64          ` json:"AppCount" `
		Longitude      string         ` json:"longitude" `
		Latitude       string         ` json:"latitude" `
		AlarmCount     int64          `json:"AlarmCount"`
		VideoCount     int64          `json:"VideoCount"`
		DataCount      int64          `json:"DataCount"`
		MemCount       int64          `json:"MemCount"`
		CpuCount       int64          `json:"CpuCount"`
		HostInfo       *host.InfoStat `json:"HostInfo"`
		MemInfo        any            `json:"MemInfo"`
		CpuInfo        any            `json:"CpuInfo"`
		DiskInfo       any            `json:"DiskInfo"`
		BootTime       any            `json:"BootTime"`
		Goroutine      any            `json:"Goroutine"`
		NetWork        any            `json:"NetWork"`
	}

	var getAnalysisInfo analysisInfo
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	
	// 系统级信息始终返回（CPU/内存/磁盘不依赖项目）
	percent, _ := cpu.Percent(0, false)
	memInfo, _ := mem.VirtualMemory()
	getAnalysisInfo.MemCount = (int64)(memInfo.UsedPercent)
	getAnalysisInfo.CpuCount = (int64)(percent[0])
	getAnalysisInfo.HostInfo, getAnalysisInfo.BootTime = GetHostInfo()
	getAnalysisInfo.MemInfo = GetMemInfo()
	getAnalysisInfo.CpuInfo = GetCpuInfo()
	getAnalysisInfo.DiskInfo = GetDiskInfo()
	getAnalysisInfo.Goroutine = runtime.NumGoroutine()
	getAnalysisInfo.NetWork = GetNetStatusInfo()

	// 项目级信息仅在指定项目时查询
	if ProjectUuid != "" {
		models.Db.Model(&models.MonitorList{}).Where("ID >0 and Type = 1 and project_uuid = ? ", ProjectUuid).Count(&getAnalysisInfo.DeviceCount)
		models.Db.Model(&models.MonitorList{}).Where("ID >0 and Type = 1 and project_uuid = ? and (status = 0 or status = 2)", ProjectUuid).Count(&getAnalysisInfo.DeviceOffCount)
		models.Db.Model(&models.DisplayModels{}).Where("ID >0 and project_uuid = ? ", ProjectUuid).Count(&getAnalysisInfo.AppCount)
		models.Db.Model(&models.DevicesAlarmList{}).Where("project_uuid=? and clear_time < ?", ProjectUuid, "2007-01-02 15:04:05").Count(&getAnalysisInfo.AlarmCount)
		models.Db.Model(&models.DeviceRealData{}).Where("project_uuid=?", ProjectUuid).Count(&getAnalysisInfo.DataCount)
		models.Db.Model(&models.ProjectVideoList{}).Where("project_uuid=?", ProjectUuid).Count(&getAnalysisInfo.VideoCount)
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
		"list": getAnalysisInfo,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}
func (c *ISMSystem) GetSystemParams() {

	type paramsInfo struct {
		WSPort                   uint16 ` json:"WSPort" `
		WSAddress                string ` json:"WSAddress" `
		IsLicense                bool   ` json:"IsLicense" `
		IsOEM                    bool   ` json:"IsOEM" `
		DefaultLang              string ` json:"DefaultLang" `
		UserData                 any    ` json:"UserData" `
		IsAuthTimeLimit          bool   ` json:"IsAuthTimeLimit" `
		IsAuthLimit              bool   ` json:"IsAuthLimit" `
		AuthRemainingTimeDays    int    ` json:"AuthRemainingTimeDays" `
		AuthRemainingTimeHours   int    ` json:"AuthRemainingTimeHours" `
		AuthRemainingTimeMinutes int    ` json:"AuthRemainingTimeMinutes" `
	}
	var MenuData any
	var getParamsInfo paramsInfo

	WSPort, wserr := config.Int("WSPort")
	if wserr != nil {
		getParamsInfo.WSPort = 10215
	} else {
		getParamsInfo.WSPort = uint16(WSPort)
	}

	WSAddress, wsaerr := config.String("WSAddress")
	if wsaerr != nil || WSAddress == "" {
		getParamsInfo.WSAddress = "local"
	} else {
		getParamsInfo.WSAddress = WSAddress
	}
	DefaultLang, langerrr := config.String("DefaultLanguage")
	if langerrr != nil || DefaultLang == "" {
		getParamsInfo.DefaultLang = "CN"
	} else {
		getParamsInfo.DefaultLang = DefaultLang
	}
	getParamsInfo.IsAuthLimit = protocolCommon.IsAuthLimit
	getParamsInfo.IsAuthTimeLimit = protocolCommon.IsAuthTimeLimit
	getParamsInfo.AuthRemainingTimeDays = protocolCommon.AuthRemainingTimeDays
	getParamsInfo.AuthRemainingTimeHours = protocolCommon.AuthRemainingTimeHours
	getParamsInfo.AuthRemainingTimeMinutes = protocolCommon.AuthRemainingTimeMinutes
	getParamsInfo.IsLicense = protocolCommon.IsLicense
	getParamsInfo.IsOEM = protocolCommon.IsOem
	MenuListDefault := `[{"router":"root","children":[{"router":"dashboard"},{"router":"DataWarehouse"},{"router":"DeviceLibraryConfig"},{"router":"DeviceModel","children":[{"router":"SnmpModel"},{"router":"SnmpAdd"},{"router":"SnmpDetail"},{"router":"ModbusModel"},{"router":"ModbusAdd"},{"router":"ModbusRegister"},{"router":"ModbusDetail"},{"router":"DLT645Model"},{"router":"DLT645ModelAdd"},{"router":"DLT645ModelData"},{"router":"DLT645ModelDetail"},{"router":"OPCUAModel"},{"router":"OpcuaAdd"},{"router":"OpcuaNodeid"},{"router":"OpcuaDetail"},{"router":"MqttModel"},{"router":"MqttAdd"},{"router":"MqttNodeid"},{"router":"MqttDetail"},{"router":"SiemensS7Model"},{"router":"SimS7Add"},{"router":"S7DataList"},{"router":"SimS7Detail"},{"router":"IEC104Model"},{"router":"IEC104ModelAdd"},{"router":"IEC104ModelData"},{"router":"IEC104ModelDetail"},{"router":"IEC61850Model"},{"router":"IEC61850Add"},{"router":"IEC61850Nodeid"},{"router":"IEC61850Detail"},{"router":"RestFulModel"},{"router":"RestFulData"},{"router":"HJ212Model"},{"router":"HJ212Add"},{"router":"HJ212Nodeid"},{"router":"HJ212Detail"},{"router":"VirtualDevice"},{"router":"VirtualDeviceData"},{"router":"DeviceCustomData"},{"router":"SystemData"},{"router":"StaticDataAdd"},{"router":"StaticDataDetail"}]},{"router":"Application"},{"router":"VideoManager","children":[{"router":"videoScreen"},{"router":"videoList"},{"router":"GB28281List"}]},{"router":"RealTimeAlarm"},{"router":"AlarmStrategy","children":[{"router":"ModelTrigger"},{"router":"AlarmRestoreMask"}]},{"router":"TaskPlan"},{"router":"ISMScripts"},{"router":"DataPush","children":[{"router":"DataTemplete"},{"router":"IEC104DataTemplete"},{"router":"IEC104TempleteData"},{"router":"ModbusDataTemplete"},{"router":"ModbusTempleteData"},{"router":"DataInterface"}]},{"router":"Reporting","children":[{"router":"AlarmHistory"},{"router":"DataHistory"},{"router":"DiyReport"},{"router":"DiyReportTemplete"},{"router":"ReportTempleteContent"}]},{"router":"Network","children":[{"router":"SystemNetwork"},{"router":"ISMNetwork"}]},{"router":"DataBase","children":[{"router":"DbManager"},{"router":"HistoryManager"}]},{"router":"Setting","children":[{"router":"Account"},{"router":"UserAdd"},{"router":"UserManager"},{"router":"AlarmTipsSetting"},{"router":"SystemParams"},{"router":"AccessToken"}]},{"router":"Journal","children":[{"router":"OperationJournal"}]},{"router":"Help","children":[{"router":"AboutSystem"},{"router":"SystemAuth"}]}]}]`
	MenuReadData, err2 := os.ReadFile("conf/MenuConfig.json")
	if err2 != nil {
		json.Unmarshal([]byte(MenuListDefault), &MenuData)
	} else {
		err2 = json.Unmarshal(MenuReadData, &MenuData)
		if err2 != nil {
			json.Unmarshal([]byte(MenuListDefault), &MenuData)
		}
	}
	CasAuthUrl, err := config.String("CasAuthUrl")
	if err == nil && len(CasAuthUrl) > 0 {
		Authorization := c.Ctx.GetCookie("Authorization")
		if len(Authorization) > 0 {
			var rolesArray []rolesStu
			var permissionsArray []permissionsStu
			tokenCode, username, Role, Name, _ := middleware.JwtToken(Authorization)
			if tokenCode == errmsg.SUCCSE {
				userData := map[string]interface{}{
					"name":        Name,
					"avatar":      "",
					"Job":         "",
					"Role":        Role,
					"Menu":        MenuData,
					"Uuid":        username,
					"ProjectUUID": "",
				}
				permission := models.FindUserPermission(Role)
				permissionArray := strings.Split(permission, ",")
				rolesArray = append(rolesArray, rolesStu{Id: Role, Operation: permissionArray})

				returnData := map[string]interface{}{
					"expireAt":    36000,
					"user":        userData,
					"roles":       rolesArray,
					"permissions": permissionsArray,
					"token":       Authorization,
				}
				getParamsInfo.UserData = returnData
			}
		}
	} else {
		getParamsInfo.UserData = nil
	}

	result := map[string]interface{}{
		"list": getParamsInfo,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}
func (c *ISMSystem) GetSystemParamsList() {

	var code int = 0
	var message string = "成功"
	type paramsInfo struct {
		WSPort      uint16 ` json:"WSPort" `
		WSAddress   string ` json:"WSAddress" `
		DefaultLang string ` json:"DefaultLang" `
		Httpport    string ` json:"Httpport" `
	}
	var getParamsInfo paramsInfo

	WSPort, wserr := config.Int("WSPort")
	if wserr != nil {
		getParamsInfo.WSPort = 10215
	} else {
		getParamsInfo.WSPort = uint16(WSPort)
	}

	WSAddress, wsaerr := config.String("WSAddress")
	if wsaerr != nil || WSAddress == "" {
		getParamsInfo.WSAddress = "local"
	} else {
		getParamsInfo.WSAddress = WSAddress
	}
	DefaultLang, langerrr := config.String("DefaultLanguage")
	if langerrr != nil || DefaultLang == "" {
		getParamsInfo.DefaultLang = "CN"
	} else {
		getParamsInfo.DefaultLang = DefaultLang
	}

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
		"list": getParamsInfo,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}

type NetworkInterfaceIP struct {
	IP   string
	Mask string
}
type NetworkInterface struct {
	Name       string
	MacAddress string
	GateWay    string
	IPv4       []NetworkInterfaceIP
	IPv6       []NetworkInterfaceIP
}

func GetNetworkInterfaces() ([]NetworkInterface, error) {
	// interfaces, err := net.Interfaces()
	// if err != nil {
	// 	return nil, err
	// }
	var networkInterfaces []NetworkInterface
	// for _, iface := range interfaces {
	// 	addrs, err := iface.Addrs()
	// 	if err != nil {
	// 		continue
	// 	}
	// 	for _, addr := range addrs {
	// 		ipNet, ok := addr.(*net.IPNet)
	// 		if !ok || ipNet.IP.IsLoopback() {
	// 			continue
	// 		}
	// 		gateways, err := net.InterfaceAddrs()
	// 		if err != nil {
	// 			continue
	// 		}
	// 		var gatewayIP net.IP
	// 		for _, gw := range gateways {
	// 			gwIPNet, ok := gw.(*net.IPNet)
	// 			if ok && !gwIPNet.IP.IsLoopback() && ipNet.Contains(gwIPNet.IP) {
	// 				gatewayIP = gwIPNet.IP.Mask(gwIPNet.Mask)
	// 				break
	// 			}
	// 		}
	// 		networkInterfaces = append(networkInterfaces, NetworkInterface{
	// 			Name:    iface.Name,
	// 			IP:      ipNet.IP,
	// 			Mask:    ipNet.Mask,
	// 			Gateway: gatewayIP,
	// 		})
	// 	}
	// }
	return networkInterfaces, nil
}
func getNetworkInfo() ([]NetworkInterface, error) {

	var NetworkInterfaceArray = make([]NetworkInterface, 0)

	intf, err := net.Interfaces()
	if err != nil {
		logs.Error("get network info failed: %v", err)
		return NetworkInterfaceArray, err
	}
	var saveSystemNetwork []NetworkInterface

	filename := "conf/NetWorker.json"
	netcontent, err := os.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(netcontent, &saveSystemNetwork)
		if err != nil {
			fmt.Println("err")
		}
	}

	for _, v := range intf {
		var tempInter NetworkInterface
		if v.Name == "" {
			continue
		}
		ips, err := v.Addrs()
		if err != nil {
			logs.Error("get network addr failed: %v", err)
			return NetworkInterfaceArray, err
		}

		//此处过滤loopback（本地回环）和isatap（isatap隧道）
		if !strings.Contains(v.Name, "sit") && !strings.Contains(v.Name, "docker") && !strings.Contains(v.Name, "lo") && !strings.Contains(v.Name, "Loopback") && !strings.Contains(v.Name, "isatap") {

			for _, item := range saveSystemNetwork {
				if item.Name == v.Name {
					if item.GateWay != "" {
						tempInter.GateWay = item.GateWay
					}
					break
				}
			}
			tempInter.Name = v.Name
			tempInter.MacAddress = v.HardwareAddr.String()
			var IpInfoArray []NetworkInterfaceIP

			for _, ip := range ips {
				var IpInfo NetworkInterfaceIP
				ipNet, ok := ip.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						IpInfo.IP = ipNet.IP.To4().String()
						IpInfo.Mask = net.IP(ipNet.Mask).String()
						IpInfoArray = append(IpInfoArray, IpInfo)
					}
				}
			}
			tempInter.IPv4 = IpInfoArray
			NetworkInterfaceArray = append(NetworkInterfaceArray, tempInter)
		}
	}
	return NetworkInterfaceArray, nil
}
func (c *ISMSystem) GetSystemNetworkInfo() {

	var code int = 0
	var message string = "成功"

	interfaces, err := getNetworkInfo()
	if err != nil {
		fmt.Println("Failed to get network interfaces:", err)
		return
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
		"list": interfaces,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}
func (c *ISMSystem) SaveSystemNetworkInfo() {

	var code int = -2
	var message string = "成功"
	var recvList []NetworkInterface

	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &recvList)
	if err != nil {
		code = -1
		message = "JSON格式错误"
	} else {
		code = 0
	}
	if code == 0 {
		if (runtime.GOARCH != "arm64") && (runtime.GOARCH != "arm") {
			code = -4
			message = "不支持此芯片"
		} else {
			NetworkBody, err1 := json.MarshalIndent(recvList, "", "\t")
			if err1 != nil {
				code = -1
				message = "JSON格式错误"
			} else {
				filename := "conf/NetWorker.json"
				err = os.WriteFile(filename, NetworkBody, 0644)
				if err != nil {
					logs.Error("Writing to file failed: %s", err)
					code = -3
					message = "文件写入错误"
				}
			}
		}
	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}
func (c *ISMSystem) RebootSystem() {

	var code int = -2
	var message string = "成功"

	if (runtime.GOARCH != "arm64") && (runtime.GOARCH != "arm") {
		code = -4
		message = "不支持此芯片"
	} else {

		// 如果无法使用shutdown命令，尝试使用reboot命令
		cmd := exec.Command("reboot")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			code = -5
			message = err.Error()
		} else {
			code = 0
			message = "成功"
		}

	}
	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}
func (c *ISMSystem) RebootISMSystem() {

	var code int = -2
	var message string = "成功"

	os.Exit(0)

	result := map[string]interface{}{
		"code": code,
		"msg":  message,
	}
	c.Data["json"] = result
	c.ServeJSON() //返回json格式

}

func (c *ISMSystem) GetAuthLicenseInfo() {

	var lisceseAuth map[string]interface{}
	key := []byte("abcdsxyzhkj12345")
	var code int = 0
	result := map[string]interface{}{
		"code": code,
	}
	content, errread := os.ReadFile("static/company/license.lic")
	if errread != nil {
		code = -1
	}
	base64Content, base64err := base64.StdEncoding.DecodeString(string(content))
	if base64err != nil {
		code = -2
	}
	deLiscese, deErr := openssl.AesECBDecrypt(base64Content, key, openssl.PKCS7_PADDING)
	if deErr != nil {
		code = -3

	}

	jsonErr := json.Unmarshal(deLiscese, &lisceseAuth)
	if jsonErr != nil {
		code = -4
	}
	result["code"] = code
	result["systemName"] = lisceseAuth["systemName"]
	result["systemUrl"] = lisceseAuth["systemUrl"]
	result["SystemAPPName"] = lisceseAuth["SystemAPPName"]
	result["systemCompany"] = lisceseAuth["systemCompany"]
	result["SystemLogo"] = lisceseAuth["SystemLogo"]
	result["systemLoginBg"] = lisceseAuth["systemLoginBg"]
	result["systemBg"] = lisceseAuth["systemBg"]

	c.Ctx.Output.Header("Content-Encoding", "gzip")
	gw := gzip.NewWriter(c.Ctx.ResponseWriter)
	defer gw.Close()

	err1 := json.NewEncoder(gw).Encode(result)
	if err1 != nil {
		fmt.Println("encode err: ", err1)
	}
}

func afterReboot() {
	tChannel := time.After(5 * time.Second) // 其内部其实是生成了一个Timer对象
	select {
	case <-tChannel:
		os.Exit(0)
	}
}
func (c *ISMSystem) AuthUpload() {

	type UploadResult struct {
		Code int
		Oem  int
	}
	var reponse_result UploadResult

	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".lic": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	//创建目录
	uploadDir := tempDir
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + "license.lic"
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	reponse_result.Code = 0
	key := []byte("abcdsxyzhkj12345")

	content, errread := os.ReadFile(fpath)
	if errread != nil {
		reponse_result.Code = -1
	}
	base64Content, base64err := base64.StdEncoding.DecodeString(string(content))
	if base64err != nil {
		reponse_result.Code = -2
	}
	_, deErr := openssl.AesECBDecrypt(base64Content, key, openssl.PKCS7_PADDING)
	if deErr != nil {
		reponse_result.Code = -3

	}
	if reponse_result.Code == 0 {
		CopyFile(fpath, "static/company/license.lic")
		// go afterReboot()
		//检查是否授权
		license.CheckLicense()

		if !protocolCommon.IsOem {
			reponse_result.Oem = 0
		} else {
			reponse_result.Oem = 1
		}
	}
	c.Data["json"] = reponse_result
	c.ServeJSON() //返回json格式
}
func (c *ISMSystem) BackupUpload() {

	type UploadResult struct {
		Code int
		Oem  int
	}
	var reponse_result UploadResult

	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".sql": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	//创建目录
	uploadDir := dbbackup
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + "BackupUpload_" + time.Now().Format(("2006-01-02_15-04-05")) + ".sql"
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	reponse_result.Code = 0
	c.Data["json"] = reponse_result
	c.ServeJSON() //返回json格式
}
func (c *ISMSystem) DiyUpload() {

	type UploadResult struct {
		Code int
	}
	var reponse_result UploadResult

	f, h, _ := c.GetFile("file") //获取上传的文件

	if h == nil || f == nil {
		reponse_result.Code = -1
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".zip": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		reponse_result.Code = -2
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	//创建目录
	uploadDir := tempDir
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		reponse_result.Code = -3
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	fpath := uploadDir + "diy.zip"
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = c.SaveToFile("file", fpath)
	if err != nil {
		reponse_result.Code = -4
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}

	reader, err := zip.OpenReader(fpath)
	if err != nil {
		reponse_result.Code = -5
		c.Data["json"] = reponse_result
		c.ServeJSON()
		return
	}
	defer reader.Close()
	for _, file := range reader.File {
		var decodeName string
		if file.Flags == 0 {
			// GBK/GB18030编码转换
			i := bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			// UTF-8编码
			decodeName = file.Name
		}

		filePath := path.Join("static/customPel/", decodeName)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				reponse_result.Code = -6
				c.Data["json"] = reponse_result
				c.ServeJSON()
				return
			}
			inFile, err := file.Open()
			if err != nil {
				reponse_result.Code = -7
				c.Data["json"] = reponse_result
				c.ServeJSON()
				return
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				reponse_result.Code = -8
				c.Data["json"] = reponse_result
				c.ServeJSON()
				return
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				reponse_result.Code = -9
				c.Data["json"] = reponse_result
				c.ServeJSON()
				return
			}
		}
	}
	reponse_result.Code = 0
	c.Data["json"] = reponse_result

	c.ServeJSON() //返回json格式
}
