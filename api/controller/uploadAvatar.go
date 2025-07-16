package controller

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadAvatarResponse 头像上传响应结构
type UploadAvatarResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AvatarUrl string `json:"avatarUrl"`
		FileName  string `json:"fileName"`
	} `json:"data"`
}

// UploadAvatar 处理头像上传
func UploadAvatar(c *gin.Context) {
	var response UploadAvatarResponse

	// 获取用户账号
	accountNum := c.PostForm("account_num")
	if accountNum == "" {
		response.Code = 400
		response.Msg = "账号不能为空"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.Code = 400
		response.Msg = "获取上传文件失败"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	// 验证文件类型
	if !isValidImageType(header.Filename) {
		response.Code = 400
		response.Msg = "只支持jpg、jpeg、png、gif格式的图片"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 验证文件大小 (5MB)
	if header.Size > 5*1024*1024 {
		response.Code = 400
		response.Msg = "文件大小不能超过5MB"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 创建上传目录
	uploadDir := "./static/uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Code = 500
		response.Msg = "创建上传目录失败"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 生成唯一文件名
	fileExt := filepath.Ext(header.Filename)
	fileName := generateFileName(accountNum, fileExt)
	filePath := filepath.Join(uploadDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		response.Code = 500
		response.Msg = "创建文件失败"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		response.Code = 500
		response.Msg = "保存文件失败"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 生成访问URL
	avatarUrl := fmt.Sprintf("/api/static/uploads/avatars/%s", fileName)

	// 更新数据库中的头像URL (这里需要调用你的用户服务)
	if err := updateUserAvatar(accountNum, avatarUrl); err != nil {
		response.Code = 500
		response.Msg = "更新用户头像失败"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 返回成功响应
	response.Code = 200
	response.Msg = "头像上传成功"
	response.Data.AvatarUrl = avatarUrl
	response.Data.FileName = fileName
	c.JSON(http.StatusOK, response)
}

// isValidImageType 验证是否为有效的图片类型
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// generateFileName 生成唯一文件名
func generateFileName(accountNum, ext string) string {
	// 使用时间戳和UUID生成唯一文件名
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	uuidStr := uuid.New().String()[:8]

	// 也可以使用MD5生成文件名
	hash := md5.Sum([]byte(accountNum + timestamp + uuidStr))
	hashStr := fmt.Sprintf("%x", hash)[:16]

	return fmt.Sprintf("avatar_%s_%s%s", accountNum, hashStr, ext)
}

// updateUserAvatar 更新用户头像URL (需要根据你的实际数据库结构实现)
func updateUserAvatar(accountNum, avatarUrl string) error {
	// 这里需要调用你的用户服务来更新数据库
	// 示例代码，需要根据你的实际架构调整

	// 1. 如果是直接数据库操作：
	// db.Exec("UPDATE users SET avatar_url = ? WHERE account_num = ?", avatarUrl, accountNum)

	// 2. 如果是通过RPC调用用户服务：
	// userClient.UpdateAvatar(context.Background(), &pb.UpdateAvatarRequest{
	//     AccountNum: accountNum,
	//     AvatarUrl:  avatarUrl,
	// })

	// 3. 如果是通过HTTP调用用户服务：
	// 这里可以调用你的用户服务的更新接口

	// 暂时返回nil，你需要根据实际情况实现
	return nil
}
