/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:22
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	beego "github.com/beego/beego/v2/server/web"
)

type SystemImageController struct {
	beego.Controller
}

// CleanFileName 清理文件名中的特殊字符，适配 Beego 上传场景
func CleanFileName(filename string) string {
	if filename == "" {
		return "unknown_file"
	}

	// 过滤操作系统非法文件名字符，替换为下划线
	invalidChars := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|', '\t', '\n', '\r'}
	replaceMap := make(map[rune]string)
	for _, char := range invalidChars {
		replaceMap[char] = "_"
	}

	var cleanedBuilder strings.Builder
	for _, char := range filename {
		if replaceStr, exists := replaceMap[char]; exists {
			cleanedBuilder.WriteString(replaceStr)
		} else if unicode.IsControl(char) {
			cleanedBuilder.WriteString("_")
		} else {
			cleanedBuilder.WriteRune(char)
		}
	}

	cleanedName := cleanedBuilder.String()
	cleanedName = strings.Trim(cleanedName, " .") // 去除首尾空格/点
	if cleanedName == "" {
		return "unknown_file"
	}

	return cleanedName
}
func (c *SystemImageController) ImageList() {

	result := map[string]interface{}{
		"code": 0,
		"list": models.SystemImageList(),
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *SystemImageController) ImageUpload() {

	var imgType int = 1
	// 定义上传结果结构体
	type UploadResult struct {
		Code int    `json:"Code"`
		Path string `json:"Path"`
	}

	var responseResult UploadResult

	// 1. Beego 方式获取上传文件，正确处理错误
	f, h, err := c.GetFile("file")
	if err != nil || h == nil || f == nil {
		responseResult.Code = -1 // 文件获取失败
		c.Data["json"] = responseResult
		c.ServeJSON()
		return
	}

	// 2. 限制文件大小（示例：最大 100MB）
	const maxFileSize = 100 * 1024 * 1024 // 100MB
	if h.Size > maxFileSize {
		responseResult.Code = -5 // 文件过大
		c.Data["json"] = responseResult
		c.ServeJSON()
		return
	}

	// 3. 标准化文件后缀（转小写，避免大小写不一致）
	ext := strings.ToLower(path.Ext(h.Filename))

	// 4. 优化后缀验证逻辑（合并为一个 map，简化代码）
	allowExtMap := map[string]int{
		// 图片类型 (1)
		".png": 1, ".gif": 1, ".jpg": 1, ".jpeg": 1, ".svg": 1, ".bmp": 1,
		// 3D 文件 (2)
		".dae": 2, ".fbx": 2, ".ply": 2, ".gltf": 2, ".glb": 2, ".obj": 2, ".stl": 2, ".json": 2, ".mtl": 2, ".max": 2,
		// 文档 (3)
		".pdf": 3, ".docx": 3, ".xlsx": 3,
		// 视频 (4)
		".mp4": 4, ".flv": 4,
	}

	// 5. 验证文件后缀合法性
	fileType, ok := allowExtMap[ext]
	if !ok {
		responseResult.Code = -2 // 后缀非法
		c.Data["json"] = responseResult
		c.ServeJSON()
		return
	}
	imgType = fileType

	// 6. 清理文件名（核心优化：过滤特殊字符）
	fileNameOnly := CleanFileName(strings.TrimSuffix(h.Filename, ext))
	cleanedFileName := fileNameOnly + ext // 拼接清理后的完整文件名

	// 7. 创建上传目录（Beego 适配：权限改为 0755，更安全）
	uploadDir := models.SystemImagePath
	err = os.MkdirAll(uploadDir, 0755) // 避免 777 过大权限
	if err != nil {
		responseResult.Code = -3 // 目录创建失败
		c.Data["json"] = responseResult
		c.ServeJSON()
		return
	}

	// 8. Beego 方式保存文件（规范路径拼接）
	fpath := filepath.Join(uploadDir, cleanedFileName) // 跨系统路径拼接
	defer f.Close()                                    // 确保文件句柄关闭
	err = c.SaveToFile("file", fpath)
	if err != nil {
		responseResult.Code = -4 // 文件保存失败
		c.Data["json"] = responseResult
		c.ServeJSON()
		return
	}

	// 9. 插入数据库并返回结果
	responseResult.Path = fpath
	responseResult.Code = models.SystemImageInsert(fileNameOnly, fpath, imgType)
	c.Data["json"] = responseResult
	c.ServeJSON()
}

func (c *SystemImageController) ImageDel() {

	var delImageJson models.SystemImge
	var code int

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &delImageJson)
	if err != nil {
		code = -1
	} else {
		code = models.SystemImageDel(delImageJson.Path)
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
