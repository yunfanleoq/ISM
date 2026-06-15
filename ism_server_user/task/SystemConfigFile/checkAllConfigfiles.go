/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-15 10:33:48
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ISMConfig

import (
	"fmt"
	"os"
)

// 检查文件是否存在，如果不存在则创建文件
func checkAndCreateConfigFile(filePath string) error {
	// 判断文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，创建文件
		_, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("创建文件失败: %s", err)
		}
		fmt.Printf("文件 %s 创建成功\n", filePath)
	}
	return nil
}

func CheckAllConfigFiles() {
	checkAndCreateConfigFile("conf/ISMNodeConfig.conf")
}
