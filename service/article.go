package service

import (
	"fmt"
	"github.com/liuhongdi/digv10/bigcache"
	"github.com/liuhongdi/digv10/dao"
	"github.com/liuhongdi/digv10/global"
	"github.com/liuhongdi/digv10/model"
	"github.com/liuhongdi/digv10/rediscache"
	"strconv"
)

//得到一篇文章的详情
func GetOneArticle(articleId uint64) (*model.Article, error) {
	//get from bigcache
	article,err := bigcache.GetOneArticleBigCache(articleId);
	if ( err != nil) {
		//get from redis
		article,errSel := rediscache.GetOneArticleRedisCache(articleId)
		if (errSel != nil) {
			//get from mysql
			article,errSel := dao.SelectOneArticle(articleId);
			if (errSel != nil) {
				return nil,errSel
			} else {
				//set redis cache
				errSetR := rediscache.SetOneArticleRedisCache(articleId,article)
				if (errSetR != nil){
					fmt.Println(errSetR)
				}
				//set bigcache
				errSetB := bigcache.SetOneArticleBigCache(articleId,article)
				if (errSetB != nil){
					fmt.Println(errSetB)
				}
				return article,nil
			}
			//return nil,errSel
		} else {
			//set bigcache
			errSet := bigcache.SetOneArticleBigCache(articleId,article)
			if (errSet != nil){
				return nil,errSet
			} else {
				return article,errSel
			}
		}

	} else {
		return article,err
	}
}

func GetArticleSum() (int, error) {
	return dao.SelectcountAll()
}

//得到多篇文章，按分页返回
func GetArticleList(page int ,pageSize int) ([]*model.Article,error) {
	articles, err := dao.SelectAllArticle(page,pageSize)
	if err != nil {
		return nil,err
	} else {
		return articles,nil
	}
}

//从redis更新bigcache
func UpdateArticleBigcache(articleId uint64) (error) {
	//get from redis
	article,errSel := rediscache.GetOneArticleRedisCache(articleId)
	if (errSel != nil) {

		return errSel
	} else {
		errSetB := bigcache.SetOneArticleBigCache(articleId, article)
		if (errSetB != nil) {
			fmt.Println(errSetB)
			return errSetB
		}

		return nil
	}
}

//订阅redis消息
func SubMessage(channel string) {
	pubsub := global.RedisDb.Subscribe(channel)
	fmt.Println("subscribe begin Receive")
	_, err := pubsub.Receive()
	if err != nil {
		return
	}
	fmt.Println("subscribe begin channel")
	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println("message:")
		fmt.Println( msg.Channel, msg.Payload, "\r\n")
		//把字符串转articleid
		articleId,errUint := strconv.ParseUint(msg.Payload, 0, 64)
		if (errUint != nil) {
			fmt.Println(errUint)
		} else {
			//更新bigcache
			errB := UpdateArticleBigcache(articleId)
			if (errB != nil){
				fmt.Println(errB)
			}
		}

	}
}