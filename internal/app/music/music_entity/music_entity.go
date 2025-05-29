package music_entity

import (
	"gorm.io/gorm"
	"time"
)

// Music 表示系统中的音乐曲目
type Music struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name   string `gorm:"size:255;not null" json:"name"`   // 音乐名称
	Artist string `gorm:"size:255;not null" json:"artist"` // 艺术家/歌手
	Album  string `gorm:"size:255" json:"album"`           // 所属专辑

	CoverPath string `gorm:"size:1024" json:"cover_path"`         // 封面图片路径
	FilePath  string `gorm:"size:1024;not null" json:"file_path"` // 音乐文件路径

	Duration int    `gorm:"default:0" json:"duration"`  // 音乐时长(秒)
	FileSize int64  `gorm:"default:0" json:"file_size"` // 文件大小(字节)
	FileType string `gorm:"size:50" json:"file_type"`   // 文件类型(mp3, flac等)

	Lyrics string `gorm:"type:text" json:"lyrics"`       // 歌词
	UserID uint   `gorm:"index;not null" json:"user_id"` // 上传用户ID

	Tags string `gorm:"size:255" json:"tags"` // 标签，以逗号分隔
}

// MusicDTO 是音乐信息的数据传输对象
type MusicDTO struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Artist       string    `json:"artist"`
	Album        string    `json:"album"`
	CoverPath    string    `json:"cover_path"`
	ResourcePath string    `gorm:"size:1024" json:"resource_path"`
	Duration     int       `json:"duration"`
	FileType     string    `json:"file_type"`
	CreatedAt    time.Time `json:"created_at"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name,omitempty"` // 可从用户服务填充
	Tags         string    `json:"tags"`
}
