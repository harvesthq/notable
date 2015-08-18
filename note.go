package notable

type Note struct {
	Author    string `json:"author"`
	AvatarURL string `json:"avatar_url"`
	Category  string `json:"category"`
	Text      string `json:"text"`
}
