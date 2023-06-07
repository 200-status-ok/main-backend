package Repository

import (
	"gorm.io/gorm"
)

type PosterRepository struct {
	db *gorm.DB
}

func NewPosterRepository(db *gorm.DB) *PosterRepository {
	return &PosterRepository{db: db}
}

type PosterResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (r *PosterRepository) GetImagesTextsPosterByID(posterID uint) ([]string, *PosterResult, error) {
	var images []string
	var posterResult PosterResult

	err := r.db.Table("images").
		Select("images.url").
		Where("images.poster_id = ?", posterID).
		Scan(&images).Error

	if err != nil {
		return nil, &PosterResult{}, err
	}

	err = r.db.Table("posters").
		Select("posters.title, posters.description").
		Where("posters.id = ?", posterID).
		Scan(&posterResult).Error

	if err != nil {
		return nil, &PosterResult{}, err
	}

	return images, &posterResult, nil
}

func (r *PosterRepository) UpdatePosterState(posterID uint, state string) error {
	err := r.db.Table("posters").
		Where("id = ?", posterID).
		Update("state", state).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *PosterRepository) UpdateTags(result map[string]string) error {
	for key, value := range result {
		err := r.db.Table("tags").
			Where("name = ?", key).
			Update("state", value).Error

		if err != nil {
			return err
		}
	}
	return nil
}
