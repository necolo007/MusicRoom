package music_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/necolo007/MusicRoom/core/gin/dbs"
	"github.com/necolo007/MusicRoom/core/libx"
	"github.com/necolo007/MusicRoom/internal/app/music/music_dto"
	"github.com/necolo007/MusicRoom/internal/app/music/music_entity"
	"github.com/necolo007/MusicRoom/pkg/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	MUSIC_DIR   = "./assets/music"   // 修改为新的音乐存储路径
	PICTURE_DIR = "./assets/picture" // 修改为新的图片存储路径
	LYRIC_DIR   = "./assets/lyric"   // 新增歌词存储路径
)

// 初始化上传目录
func init() {
	var err error
	err = os.MkdirAll(MUSIC_DIR, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(PICTURE_DIR, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(LYRIC_DIR, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func UploadMusic(c *gin.Context) {
	// 直接从表单获取基本信息
	name := c.PostForm("name")
	artist := c.PostForm("artist")
	album := c.PostForm("album")
	tags := c.PostForm("tags")

	// 校验必填字段
	if name == "" || artist == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "音乐名称和艺术家为必填项"})
		return
	}

	// 获取用户ID
	userIDStr := libx.Uid(c)
	userID, err := strconv.ParseUint(strconv.Itoa(int(userIDStr)), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 处理音频文件
	musicFile, err := c.FormFile("music_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "音乐文件缺失"})
		return
	}

	// 检查文件类型
	fileExt := strings.ToLower(filepath.Ext(musicFile.Filename))
	if !utils.IsAllowedAudioType(fileExt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的音乐文件格式，仅支持mp3, wav, ogg, flac"})
		return
	}

	// 生成唯一文件名
	musicFileName := uuid.New().String() + fileExt
	musicFilePath := filepath.Join(MUSIC_DIR, musicFileName)

	// 保存音乐文件
	if err := c.SaveUploadedFile(musicFile, musicFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存音乐文件"})
		return
	}

	// 初始化音乐实体
	music := music_entity.Music{
		Name:     name,
		Artist:   artist,
		Album:    album,
		FilePath: musicFilePath,
		FileType: fileExt[1:], // 去掉点号
		FileSize: int64(musicFile.Size),
		UserID:   uint(userID),
		Tags:     tags,
	}

	// 处理封面图片（如果有）
	coverFile, err := c.FormFile("cover_image")
	if err == nil {
		coverExt := strings.ToLower(filepath.Ext(coverFile.Filename))
		if !utils.IsAllowedImageType(coverExt) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式，仅支持jpg, jpeg, png"})
			return
		}

		coverFileName := uuid.New().String() + coverExt
		coverFilePath := filepath.Join(PICTURE_DIR, coverFileName)

		if err := c.SaveUploadedFile(coverFile, coverFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存封面图片"})
			return
		}

		music.CoverPath = coverFilePath
	}

	// 处理歌词文件（如果有）
	lyricFile, err := c.FormFile("lyric_file")
	if err == nil {
		lyricExt := strings.ToLower(filepath.Ext(lyricFile.Filename))
		if lyricExt != ".lrc" && lyricExt != ".txt" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的歌词文件格式，仅支持lrc, txt"})
			return
		}

		lyricFileName := uuid.New().String() + lyricExt
		lyricFilePath := filepath.Join(LYRIC_DIR, lyricFileName)

		if err := c.SaveUploadedFile(lyricFile, lyricFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存歌词文件"})
			return
		}

		// 读取歌词内容存入数据库
		lyricContent, err := os.ReadFile(lyricFilePath)
		if err == nil {
			music.Lyrics = string(lyricContent)
		}
	} else {
		// 检查是否直接从表单提交了歌词文本
		lyrics := c.PostForm("lyrics")
		if lyrics != "" {
			// 保存歌词到文件
			lyricFileName := uuid.New().String() + ".lrc"
			lyricFilePath := filepath.Join(LYRIC_DIR, lyricFileName)
			if err := os.WriteFile(lyricFilePath, []byte(lyrics), 0644); err == nil {
				music.Lyrics = lyrics
			}
		}
	}

	// 保存到数据库
	if err := dbs.DB.Create(&music).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存音乐信息到数据库"})
		// 清理已上传的文件
		os.Remove(musicFilePath)
		if music.CoverPath != "" {
			os.Remove(music.CoverPath)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "音乐上传成功",
		"id":      music.ID,
	})
}

// GetMusic 获取单首音乐的信息
func GetMusic(c *gin.Context) {
	//获取用户id
	id := libx.Uid(c)
	musicID, err := strconv.Atoi(strconv.Itoa(int(id)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var music music_entity.Music
	if err := dbs.DB.First(&music, musicID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该音乐"})
		return
	}

	// 转换为DTO
	dto := music_entity.MusicDTO{
		ID:        music.ID,
		Name:      music.Name,
		Artist:    music.Artist,
		Album:     music.Album,
		CoverPath: music.CoverPath,
		Duration:  music.Duration,
		FileType:  music.FileType,
		CreatedAt: music.CreatedAt,
		UserID:    music.UserID,
		Tags:      music.Tags,
	}

	c.JSON(http.StatusOK, dto)
}

// ListMusic 列出音乐列表
func ListMusic(c *gin.Context) {
	var params music_dto.MusicQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	} else if params.PageSize > 100 {
		params.PageSize = 100
	}

	query := dbs.DB.Model(&music_entity.Music{})

	// 应用筛选条件
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ? OR artist LIKE ? OR album LIKE ?", keyword, keyword, keyword)
	}

	if params.UserID > 0 {
		query = query.Where("user_id = ?", params.UserID)
	}

	if params.Tags != "" {
		// 拆分标签并构建查询
		tags := strings.Split(params.Tags, ",")
		for _, tag := range tags {
			query = query.Where("FIND_IN_SET(?, tags) > 0", strings.TrimSpace(tag))
		}
	}

	// 应用排序
	if params.SortBy != "" {
		order := "DESC"
		if strings.ToLower(params.Order) == "asc" {
			order = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s", params.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页
	offset := (params.Page - 1) * params.PageSize
	var musics []music_entity.Music
	query.Limit(params.PageSize).Offset(offset).Find(&musics)

	// 转换为DTO
	dtos := make([]music_entity.MusicDTO, len(musics))
	for i, music := range musics {
		dtos[i] = music_entity.MusicDTO{
			ID:        music.ID,
			Name:      music.Name,
			Artist:    music.Artist,
			Album:     music.Album,
			CoverPath: music.CoverPath,
			Duration:  music.Duration,
			FileType:  music.FileType,
			CreatedAt: music.CreatedAt,
			UserID:    music.UserID,
			Tags:      music.Tags,
		}
	}

	c.JSON(http.StatusOK, music_dto.MusicListResponse{
		Total: total,
		Items: dtos,
	})
}

// DeleteMusic 删除音乐
func DeleteMusic(c *gin.Context) {
	id := c.Query("id")
	musicID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 获取用户ID（假设从JWT或会话中获取）
	userID := libx.Uid(c)

	var music music_entity.Music
	if err := dbs.DB.First(&music, musicID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该音乐"})
		return
	}

	role := libx.GetRole(c)
	// 验证权限（只有上传者或管理员可以删除）
	if music.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此音乐"})
		return
	}

	// 删除相关文件
	if music.FilePath != "" {
		os.Remove(music.FilePath)
	}
	if music.CoverPath != "" {
		os.Remove(music.CoverPath)
	}

	// 从数据库删除
	if err := dbs.DB.Delete(&music).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "音乐已成功删除"})
}
