package mysql

import "GoForum/models"

func GetCommunity() (data []*models.Community, err error) {
	sqlstr := "select community_id,community_name from community "
	err = db.Select(&data, sqlstr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetCommunityByID(communityID int64) (data *models.Community, err error) {
	data = new(models.Community)
	sqlstr := "select community_id,community_name from community where community_id = ?"
	err = db.Get(data, sqlstr, communityID)
	if err != nil {
		return nil, err
	}
	return data, nil
}
