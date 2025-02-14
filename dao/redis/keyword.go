package redis

import (
	"blog/dao/mysql"
	"blog/models"
	"fmt"
	"strings"
)

const (
	year                 = "year"
	month                = "month"
	week                 = "week"
	defaultIncreaseCount = 1
)

func SetEssayKeyword(essayKeyword *models.EssayIdAndKeyword) (err error) {
	var eid int64
	if eid, err = mysql.GetEssaySnowflakeID((*essayKeyword).EssayId); err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", getRedisKey(KeyEssayKeyword), eid)

	// 创建 Redis 管道
	pipe := client.Pipeline()

	// 首先删除现有的所有关键词
	pipe.Del(key)

	// 如果有新的关键词，则设置它们
	if len((*essayKeyword).Keywords) > 0 {
		// 使用 SADD 命令添加到集合
		for _, keyword := range (*essayKeyword).Keywords {
			pipe.SAdd(key, strings.ToLower(strings.TrimSpace(keyword)))
		}
	}

	// 执行管道命令
	_, err = pipe.Exec()
	if err != nil {
		return fmt.Errorf("failed to set essay keywords: %w", err)
	}
	return nil
}

func IncreaseSearchKeyword(preKey string, keyword string) (err error) {
	preKey = getRedisKey(preKey)
	return SetYearMonthWeekTimesZoneForZset(preKey, keyword, defaultIncreaseCount)
}

func GetEssayKeywordsForIndex(e *[]models.DataAboutEssay) (err error) {
	keyPre := getRedisKey(KeyEssayKeyword)
	for i := range *e {
		ids, err := mysql.GetEssaySnowflakeID((*e)[i].ID)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s%d", keyPre, ids)
		keywords, err := client.SMembers(key).Result()
		if err != nil {
			return err
		}
		(*e)[i].Keywords = append(keywords, (*e)[i].Name)
	}
	return err
}

func GetEssayKeywordsForOne(eid int64) (keywords []string, err error) {
	keyPre := getRedisKey(KeyEssayKeyword)
	key := fmt.Sprintf("%s%d", keyPre, eid)
	keywords, err = client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	return keywords, nil
}

func GetSearchKeywordRank(rankKind *models.RankKindForZset) (err error) {
	preKey := getRedisKey(KeySearchKeyWordTimes)
	return GetYearMonthWeekTimesZoneForZsetRank(rankKind, preKey)
}

func CleanLowerKeywordsZsetEveryMonth() error {
	return CleanLowerZsetEveryMonth()
}
