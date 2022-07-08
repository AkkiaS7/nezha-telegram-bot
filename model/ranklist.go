package model

type RankList struct {
	UserID        int64 `gorm:"primary_key"`
	ServerCount   int
	OnlineCount   int
	MemTotal      int64
	MemUsedTotal  int64
	DiskTotal     int64
	DiskUsedTotal int64
	Load15Total   float64
}

func GetServerCountTop(limit int) []*RankList {
	var rankList []*RankList
	DB.Table("rank_lists").Select("user_id,server_count").Order("server_count desc").Limit(limit).Scan(&rankList)
	return rankList
}

func GetOnlineCountTop(limit int) []*RankList {
	var rankList []*RankList
	DB.Table("rank_lists").Select("user_id,online_count").Order("online_count desc").Limit(limit).Scan(&rankList)
	return rankList
}

func GetMemTotalTop(limit int) []*RankList {
	var rankList []*RankList
	DB.Table("rank_lists").Select("user_id,mem_total").Order("mem_total desc").Limit(limit).Scan(&rankList)
	return rankList
}

func GetMemUsedTop(limit int) []*RankList {
	var rankList []*RankList
	DB.Table("rank_lists").Select("user_id,mem_used_total").Order("mem_used_total desc").Limit(limit).Scan(&rankList)
	return rankList
}

func GetDiskTotalTop(limit int) []*RankList {
	var rankList []*RankList
	DB.Table("rank_lists").Select("user_id,disk_total").Order("disk_total desc").Limit(limit).Scan(&rankList)
	return rankList
}

func GetDiskUsedTop(limit int) []*RankList {
	var rankList []*RankList
	DB.Table("rank_lists").Select("user_id,disk_used_total").Order("disk_used_total desc").Limit(limit).Scan(&rankList)
	return rankList
}

func GetLoad15Top(limit int) []*RankList {
	var rankList []*RankList
	DB.Table("rank_lists").Select("user_id,load15_total").Order("load15_total desc").Limit(limit).Scan(&rankList)
	return rankList
}

func (rl *RankList) Save() {
	DB.Save(rl)
}

func (rl *RankList) Delete() {
	DB.Delete(rl)
}

func DeleteRankByUserID(userID int64) error {
	return DB.Delete(&RankList{}, "user_id = ?", userID).Error
}
