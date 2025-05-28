package utils

// IsAllowedAudioType 辅助函数：检查音频类型
func IsAllowedAudioType(ext string) bool {
	allowedTypes := []string{".mp3", ".wav", ".ogg", ".flac"}
	for _, t := range allowedTypes {
		if ext == t {
			return true
		}
	}
	return false
}

// IsAllowedImageType 辅助函数：检查图片类型
func IsAllowedImageType(ext string) bool {
	allowedTypes := []string{".jpg", ".jpeg", ".png"}
	for _, t := range allowedTypes {
		if ext == t {
			return true
		}
	}
	return false
}
