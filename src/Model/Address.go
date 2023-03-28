package Model

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Province      string  `gorm:"type:varchar(100);" json:"province"`
	City          string  `gorm:"type:varchar(100);" json:"city"`
	AddressDetail string  `gorm:"type:text;" json:"address_detail"`
	Latitude      float64 `gorm:"type:decimal;" json:"latitude"`
	Longitude     float64 `gorm:"type:decimal;" json:"longitude"`
	PosterId      uint    `gorm:"type:int;unique" json:"address_id"`
}

func (a *Address) GetProvince() string {
	return a.Province
}

func (a *Address) SetProvince(province string) {
	a.Province = province
}

func (a *Address) GetCity() string {
	return a.City
}

func (a *Address) SetCity(city string) {
	a.City = city
}

func (a *Address) GetAddressDetail() string {
	return a.AddressDetail
}

func (a *Address) SetAddressDetail(addressDetail string) {
	a.AddressDetail = addressDetail
}

func (a *Address) GetLatitude() float64 {
	return a.Latitude
}

func (a *Address) SetLatitude(latitude float64) {
	a.Latitude = latitude
}

func (a *Address) GetLongitude() float64 {
	return a.Longitude
}

func (a *Address) SetLongitude(longitude float64) {
	a.Longitude = longitude
}

func (a *Address) GetPosterId() uint {
	return a.PosterId
}

func (a *Address) SetPosterId(posterId uint) {
	a.PosterId = posterId
}
