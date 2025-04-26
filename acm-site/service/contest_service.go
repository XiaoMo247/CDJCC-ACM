package service

import (
	"acm-site/database"
	"acm-site/model"
)

func UploadContest(data *model.Contest) error {
	return database.DB.Create(&data).Error
}

func DeleteContest(id uint) error {
	return database.DB.Delete(&model.Contest{}, id).Error
}

func GetAllContests() ([]model.Contest, error) {
	var contests []model.Contest
	err := database.DB.Order("time desc").Find(&contests).Error
	return contests, err
}
