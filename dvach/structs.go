package dvach

import (
	"encoding/json"
)

type AttachedFile struct {
	DisplayName   string `json:"displayname"`
	FullName      string `json:"fullname"`
	Name          string `json:"name"`
	MD5           string `json:"md5"`
	Path          string `json:"path"`
	ThumbnailPath string `json:"thumbnail"`
	Type          int    `json:"type"`
	NSFW          int    `json:"nsfw"`
}

type Post struct {
	Comment   string         `json:"comment"`
	Timestamp uint64         `json:"timestamp"`
	Files     []AttachedFile `json:"files"`
}

type thread struct {
	Posts []Post `json:"posts"`
}

type board struct {
	Thread []thread `json:"threads"`
}

func UnmarshalPosts(threadData []byte) ([]Post, error) {
	var b board

	err := json.Unmarshal(threadData, &b)
	if err != nil {
		return nil, err
	}

	return b.Thread[0].Posts, nil
}
