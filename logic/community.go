package logic

import (
	"GoForum/dao/mysql"
	"GoForum/models"
)

func GetCommunity() (data []*models.Community, err error) {
	return mysql.GetCommunity()
}

func GetDetailCommunity(communityID int64) (data *models.Community, err error) {
	return mysql.GetCommunityByID(communityID)
}
