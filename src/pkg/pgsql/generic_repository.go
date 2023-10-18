package pgsql

import (
	"context"
	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{
		db: db,
	}
}

func (r *GenericRepository[T]) Add(entity *T, ctx context.Context) error {
	return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *GenericRepository[T]) AddAll(entity *[]T, ctx context.Context) error {
	return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *GenericRepository[T]) GetById(id int, ctx context.Context) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).Where("id = ? AND is_active = ?", id, true).FirstOrInit(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *GenericRepository[T]) Get(params *T, ctx context.Context) *T {
	var entity T
	r.db.WithContext(ctx).Where(&params).FirstOrInit(&entity)
	return &entity
}

func (r *GenericRepository[T]) GetAll(ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *GenericRepository[T]) Where(params *T, ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Where(&params).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *GenericRepository[T]) Update(entity *T, ctx context.Context) error {
	return r.db.WithContext(ctx).Save(&entity).Error
}

func (r GenericRepository[T]) UpdateAll(entities *[]T, ctx context.Context) error {
	return r.db.WithContext(ctx).Save(&entities).Error
}

func (r *GenericRepository[T]) Delete(id int, ctx context.Context) error {
	var entity T
	return r.db.WithContext(ctx).FirstOrInit(&entity).UpdateColumn("is_active", false).Error
}

func (r *GenericRepository[T]) SkipTake(skip int, take int, ctx context.Context) (*[]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Offset(skip).Limit(take).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return &entities, nil
}

func (r *GenericRepository[T]) Count(ctx context.Context) int64 {
	var entity T
	var count int64
	r.db.WithContext(ctx).Model(&entity).Count(&count)
	return count
}

func (r *GenericRepository[T]) CountWhere(params *T, ctx context.Context) int64 {
	var entity T
	var count int64
	r.db.WithContext(ctx).Model(&entity).Where(&params).Count(&count)
	return count
}
