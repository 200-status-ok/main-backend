package DTO

type CreatePosterDTO struct {
	Title       string  `json:"title" binding:"required,max=255"`
	Description string  `json:"description" binding:"max=1000"`
	Status      string  `json:"status" binding:"required,oneof=lost found"`
	TelID       string  `json:"tel_id" binding:"max=255"`
	UserPhone   string  `json:"user_phone" binding:"min=11,max=13"`
	Alert       bool    `json:"alert" binding:"required"`
	Chat        bool    `json:"chat" binding:"required"`
	Award       float64 `json:"award"`
	UserID      uint    `json:"user_id" binding:"required,min=1"`
	State       string  `json:"state"`
}

type UpdatePosterDTO struct {
	Title       string   `json:"title" binding:"max=255"`
	Description string   `json:"description" binding:"max=1000"`
	Status      string   `json:"status" binding:"oneof=lost found ''"`
	TelID       string   `json:"tel_id" binding:"max=255"`
	UserPhone   string   `json:"user_phone" binding:"max=13"`
	Alert       string   `json:"alert" binding:"oneof=true false ''"`
	Chat        string   `json:"chat" binding:"oneof=true false ''"`
	Award       float64  `json:"award"` //todo if you want to update reward to 0, set it to -1
	UserID      uint     `json:"user_id"`
	State       string   `json:"state" binding:"oneof=pending accepted rejected ''"`
	ImgUrls     []string `json:"img_urls"`
	TagIds      []int    `json:"tag_ids"`
}

type FilterObject struct { //todo move this to another file
	Status       string
	SearchPhrase string
	TimeStart    int64
	TimeEnd      int64
	OnlyRewards  bool
	Lat          float64
	Lon          float64
	TagIds       []int
	State        string
}

type GeneratedPosterTags struct { //todo modar move maybe, reGenerate
	Result Result `json:"result"`
	Status Status `json:"status"`
}
type TagName struct {
	Fa string `json:"fa"`
}
type Tag struct {
	Confidence float64 `json:"confidence"`
	Tag        TagName `json:"tag"`
}
type Result struct {
	Tags []Tag `json:"tags"`
}
type Status struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type GeneratedPosterColors struct {
	Result struct {
		Colors struct {
			BackgroundColors []struct {
				B                           int     `json:"b"`
				ClosestPaletteColor         string  `json:"closest_palette_color"`
				ClosestPaletteColorHTMLCode string  `json:"closest_palette_color_html_code"`
				ClosestPaletteColorParent   string  `json:"closest_palette_color_parent"`
				ClosestPaletteDistance      float64 `json:"closest_palette_distance"`
				G                           int     `json:"g"`
				HTMLCode                    string  `json:"html_code"`
				Percent                     float64 `json:"percent"`
				R                           int     `json:"r"`
			} `json:"background_colors"`
			ColorPercentThreshold float64 `json:"color_percent_threshold"`
			ColorVariance         float64 `json:"color_variance"`
			ForegroundColors      []struct {
				B                           int     `json:"b"`
				ClosestPaletteColor         string  `json:"closest_palette_color"`
				ClosestPaletteColorHTMLCode string  `json:"closest_palette_color_html_code"`
				ClosestPaletteColorParent   string  `json:"closest_palette_color_parent"`
				ClosestPaletteDistance      float64 `json:"closest_palette_distance"`
				G                           int     `json:"g"`
				HTMLCode                    string  `json:"html_code"`
				Percent                     float64 `json:"percent"`
				R                           int     `json:"r"`
			} `json:"foreground_colors"`
			ImageColors []struct {
				B                           int     `json:"b"`
				ClosestPaletteColor         string  `json:"closest_palette_color"`
				ClosestPaletteColorHTMLCode string  `json:"closest_palette_color_html_code"`
				ClosestPaletteColorParent   string  `json:"closest_palette_color_parent"`
				ClosestPaletteDistance      float64 `json:"closest_palette_distance"`
				G                           int     `json:"g"`
				HTMLCode                    string  `json:"html_code"`
				Percent                     float64 `json:"percent"`
				R                           int     `json:"r"`
			} `json:"image_colors"`
			ObjectPercentage float64 `json:"object_percentage"`
		} `json:"colors"`
	} `json:"result"`
	Status struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"status"`
}
