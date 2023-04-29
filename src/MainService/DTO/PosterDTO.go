package DTO

type PosterDTO struct {
	Title       string  `json:"title" binding:"required,min=5,max=255"`
	Description string  `json:"description" binding:"min=5,max=1000"`
	Status      string  `json:"status" binding:"required,oneof=lost found"`
	TelID       string  `json:"tel_id" binding:"min=5,max=255"`
	UserPhone   string  `json:"user_phone" binding:"min=11,max=13"`
	Alert       bool    `json:"alert" binding:"required"`
	Chat        bool    `json:"chat" binding:"required"`
	Award       float64 `json:"award"`
	UserID      uint    `json:"user_id" binding:"required,min=1"`
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
}
