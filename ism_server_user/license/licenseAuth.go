/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:58:05
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package licenseAuth

import (
	protocolCommon "ISMServer/protocol/common"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/denisbrodbeck/machineid"
	"github.com/forgoer/openssl"
	"github.com/shirou/gopsutil/cpu"
	"github.com/tjfoc/gmsm/sm4"
)

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
	//去掉明文后面的填充数据
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
	if upper {
		result = strings.ToUpper(result)
	}
	if half {
		result = result[8:24]
	}
	return result
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

//高级加密标准（Adevanced Encryption Standard ,AES）

// 16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
// key不能泄露
var PwdKey = []byte("QQ3n!JCy@N&QEEOQ8wUqO2U$3aRjqpeK")

// PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充的反向操作，删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}

// 实现加密
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 实现解密
func AesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// 加密base64
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEcrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

// 解密
func DePwdCode(pwd string) ([]byte, error) {
	//解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	//执行AES解密
	return AesDeCrypt(pwdByte, PwdKey)

}

func GetPhysicalID() string {
	id, _ := machineid.ProtectedID("myAppName")
	return GetMd5String(id, true, true)
}

func CheckLicense() bool {
	var code int = 0
	var lisceseAuth map[string]interface{}

	var ISMProtectedID string
	var getError error

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
	m := getMd5String(id, true, true)
	protocolCommon.ISMProtectedID = m
	content, err := os.ReadFile("data/auth/active.dat")
	if err != nil {
		code = -1
	} else {
		var lisceseAuth map[string]interface{}
		authContent, autherr := SM4Decrypt(string(content))
		if autherr != nil {
			code = -5
		} else {
			jsonErr := json.Unmarshal([]byte(authContent), &lisceseAuth)
			if jsonErr == nil {
				IsBatch := int(lisceseAuth["IsBatch"].(float64))
				if IsBatch == 1 {
					code = 0
				} else {
					if lisceseAuth["ISMProtectedID"] != m {
						code = -6
					} else {
						code = 0
					}
				}
			} else {
				code = -4
			}
		}
	}
	if code != 0 {
		protocolCommon.IsLicense = false
	} else {
		protocolCommon.IsLicense = true
	}

	key := []byte("abcdsxyzhkj12345")
	protocolCommon.IsLicense = true
	content, errread := os.ReadFile("static/company/license.lic")
	if errread != nil {
		protocolCommon.ConfigPageCount = 2
		protocolCommon.ConfigAppCount = 1
		protocolCommon.IsOem = false
		return protocolCommon.IsLicense
	}
	base64Content, base64err := base64.StdEncoding.DecodeString(string(content))
	if base64err != nil {
		protocolCommon.ConfigPageCount = 2
		protocolCommon.ConfigAppCount = 1
		protocolCommon.IsOem = false
		return protocolCommon.IsLicense
	}
	deLiscese, deErr := openssl.AesECBDecrypt(base64Content, key, openssl.PKCS7_PADDING)
	if deErr != nil {
		protocolCommon.ConfigPageCount = 2
		protocolCommon.ConfigAppCount = 1
		protocolCommon.IsOem = false
		return protocolCommon.IsLicense
	}
	jsonErr := json.Unmarshal(deLiscese, &lisceseAuth)
	if jsonErr != nil {
		protocolCommon.ConfigPageCount = 2
		protocolCommon.ConfigAppCount = 1
		protocolCommon.IsOem = false
		return protocolCommon.IsLicense
	}
	if lisceseAuth["ConfigPageCount"] != nil {
		protocolCommon.ConfigPageCount = int(lisceseAuth["ConfigPageCount"].(float64))
	} else if lisceseAuth["systemUrl"] != "www.ismctl.com" {
		protocolCommon.ConfigPageCount = 1000
	} else {
		protocolCommon.ConfigPageCount = 2
	}
	if lisceseAuth["ConfigAppCount"] != nil {
		protocolCommon.ConfigAppCount = int(lisceseAuth["ConfigAppCount"].(float64))
	} else if lisceseAuth["systemUrl"] != "www.ismctl.com" {
		protocolCommon.ConfigAppCount = 1000
	} else {
		protocolCommon.ConfigAppCount = 1
	}
	if lisceseAuth["systemUrl"] != "www.ismctl.com" {
		protocolCommon.IsOem = true
		return protocolCommon.IsLicense
	} else if (lisceseAuth["isAuth"] != nil) && int(lisceseAuth["isAuth"].(float64)) == 1 {
		protocolCommon.IsOem = true
		return protocolCommon.IsLicense
	}

	return protocolCommon.IsLicense
}

// 生成32位md5字串
func GetMd5String(s string, upper bool, half bool) string {
	h := md5.New()
	h.Write([]byte(s))
	result := hex.EncodeToString(h.Sum(nil))
	if upper {
		result = strings.ToUpper(result)
	}
	if half {
		result = result[8:24]
	}
	return result
}

// 利用随机数生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b), true, false)
}
