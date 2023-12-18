package services

import (
	"goback/models"
	"goback/utils"
)

func AllInfo() ([]models.Info, error) {
	var infos []models.Info
	utils.DB.Find(&infos)

	if infos == nil || len(infos) == 0 {
		return make([]models.Info, 0), nil
	}
	return infos, nil

}
