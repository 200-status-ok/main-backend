package Repository

import (
	"github.com/403-access-denied/main-backend/src/MainService/DTO"
	"github.com/403-access-denied/main-backend/src/MainService/Model"
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

func (r *AddressRepository) CreateAddress(addresses []DTO.CreateAddressDTO, posterID uint) ([]Model.Address, error) {
	var addressesModel []Model.Address
	for _, address := range addresses {
		addressesModel = append(addressesModel, Model.Address{
			Province:      address.Province,
			City:          address.City,
			AddressDetail: address.AddressDetail,
			Latitude:      address.Latitude,
			Longitude:     address.Longitude,
			PosterId:      posterID,
		})
	}
	result := r.db.Create(&addressesModel)
	if result.Error != nil {
		return []Model.Address{}, result.Error
	}

	return addressesModel, nil
}

func (r *AddressRepository) UpdateAddress(addresses []DTO.CreateAddressDTO, posterID uint) ([]Model.Address, error) {
	var addressesModel []Model.Address
	for _, address := range addresses {
		addressesModel = append(addressesModel, Model.Address{
			Province:      address.Province,
			City:          address.City,
			AddressDetail: address.AddressDetail,
			Latitude:      address.Latitude,
			Longitude:     address.Longitude,
			PosterId:      posterID,
		})
	}
	result := r.db.Save(&addressesModel)
	if result.Error != nil {
		return []Model.Address{}, result.Error
	}

	return addressesModel, nil
}
