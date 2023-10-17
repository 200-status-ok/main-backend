package Repository

import (
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) GetPaymentById(id int) (Model.Payment, error) {
	var payment Model.Payment
	result := r.db.First(&payment, id)
	if result.Error != nil {
		return Model.Payment{}, result.Error
	}

	return payment, nil
}

func (r *PaymentRepository) CreatePayment(payment Model.Payment) (Model.Payment, error) {
	result := r.db.Create(&payment)
	if result.Error != nil {
		return Model.Payment{}, result.Error
	}

	return payment, nil
}

func (r *PaymentRepository) DeletePayment(id uint) error {
	var paymentModel Model.Payment
	result := r.db.Find(&paymentModel, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	result = r.db.Delete(&paymentModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PaymentRepository) UpdatePayment(id uint, payment Model.Payment) (Model.Payment, error) {
	var paymentModel Model.Payment
	result := r.db.First(&paymentModel, id)
	if result.Error != nil {
		return Model.Payment{}, result.Error
	}
	paymentModel.SetStatus(payment.GetStatus())
	paymentModel.SetAmount(payment.GetAmount())
	result = r.db.Save(&paymentModel)
	if result.Error != nil {
		return Model.Payment{}, result.Error
	}

	return paymentModel, nil
}

func (r *PaymentRepository) GetPaymentsOfUser(userID uint) ([]Model.Payment, error) {
	var payments []Model.Payment
	result := r.db.Find(&payments, "user_id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return payments, nil
}

func (r *PaymentRepository) GetPaymentByTrackID(trackID string) (Model.Payment, error) {
	var payment Model.Payment
	result := r.db.First(&payment, "track_id = ?", trackID)
	if result.Error != nil {
		return Model.Payment{}, result.Error
	}

	return payment, nil
}
