
package middleware

import "ISMServer/utils/errmsg"

// CheckApiToken API Token 认证（用于 OpenAPI）
// 简单验证：检查 JWT 是否有效
func CheckApiToken(token string) (int, string, string) {
	code, username, role, _, _ := JwtToken(token)
	if code == errmsg.SUCCSE {
		return errmsg.SUCCSE, username, role
	}
	return errmsg.ERROR_TOKEN_NOT_EXIST, "", ""
}

// GenerateAPIToken 为用户生成长期 API Token（复用 JWT 机制）
func GenerateAPIToken(username, role, name, user_uuid string) (string, int64, int) {
	return SetToken(username, role, name, user_uuid)
}
