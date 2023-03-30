package Repository

import (
	"github.com/403-access-denied/main-backend/src/DTO"
	"github.com/403-access-denied/main-backend/src/Model"
	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) GetAddressById(id int) (Model.Address, error) {
	var address Model.Address
	result := r.db.First(&address, id)
	if result.Error != nil {
		return Model.Address{}, result.Error
	}

	return address, nil
}

func (r *AddressRepository) CreateAddress(address DTO.AddressDTO, posterID uint) (Model.Address, error) {
	var addressModel Model.Address
	addressModel.SetPosterId(uint(posterID))
	addressModel.SetProvince(address.Province)
	addressModel.SetCity(address.City)
	addressModel.SetAddressDetail(address.AddressDetail)
	addressModel.SetLatitude(address.Latitude)
	addressModel.SetLongitude(address.Longitude)
	result := r.db.Create(&addressModel)
	if result.Error != nil {
		return Model.Address{}, result.Error
	}

	return addressModel, nil
}

func (r *AddressRepository) UpdateAddress(address DTO.AddressDTO, posterID uint) (Model.Address, error) {
	var addressModel Model.Address
	result := r.db.First(&addressModel, "poster_id = ?", posterID)
	if result.Error != nil {
		return Model.Address{}, result.Error
	}
	addressModel.SetProvince(address.Province)
	addressModel.SetCity(address.City)
	addressModel.SetAddressDetail(address.AddressDetail)
	addressModel.SetLatitude(address.Latitude)
	addressModel.SetLongitude(address.Longitude)
	result = r.db.Save(&addressModel)
	if result.Error != nil {
		return Model.Address{}, result.Error
	}

	return addressModel, nil
}
