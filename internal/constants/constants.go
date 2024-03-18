package constants

import "github.com/ffdb42/vk_trainee_task/internal/models"

const (
	AdminRole = "admin"
	UserRole  = "user"
)

const (
	SortDesc models.SortOrder = "DESC"
	SortAsc  models.SortOrder = "ASC"
)

const (
	SortByName        models.SortBy = "name"
	SortByRating      models.SortBy = "rating"
	SortByReleaseDate models.SortBy = "release_date"
)
