package music_dto

import "github.com/necolo007/MusicRoom/internal/app/music/music_entity"

// MusicQueryParams 表示列出音乐的查询参数
type MusicQueryParams struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"` // 通过名称、艺术家、专辑搜索
	UserID   uint   `form:"user_id" json:"user_id"` // 按上传者过滤
	SortBy   string `form:"sort_by" json:"sort_by"` // created_at, plays等
	Order    string `form:"order" json:"order"`     // asc, desc
	Tags     string `form:"tags" json:"tags"`       // 按标签过滤
}

// MusicListResponse 表示列出音乐的响应
type MusicListResponse struct {
	Total int64                   `json:"total"`
	Items []music_entity.MusicDTO `json:"items"`
}
