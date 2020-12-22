package bigcache

import (
	"encoding/json"
	"fmt"
	"github.com/liuhongdi/digv10/global"
	"github.com/liuhongdi/digv10/model"
	"strconv"
)

//bigcache中索引的名字
func getArticleCacheName(articleId uint64) (string) {
	return "article_"+strconv.FormatUint(articleId,10)
}

//从bigcache得到一篇文章
func GetOneArticleBigCache(articleId uint64) (*model.Article,error) {
	fmt.Println("bigcache:GetOneArticleBigCache")
	key := getArticleCacheName(articleId);
	val,err := global.BigCache.Get(key)
	if (err != nil) {
		return nil,err
	} else {
		article := model.Article{}
		if err := json.Unmarshal([]byte(val), &article); err != nil {
			return nil,err
		}
		return &article,nil
	}
}
//向bigcache保存一篇文章
func SetOneArticleBigCache(articleId uint64,article *model.Article) (error) {
	key := getArticleCacheName(articleId);
	content,err := json.Marshal(article)
	if (err != nil){
		fmt.Println(err)
		return err;
	}
	errSet := global.BigCache.Set(key,[]byte(content))
	if (errSet != nil) {
		return errSet
	}
	fmt.Println("len:",global.BigCache.Len())
	fmt.Println("Capacity:",global.BigCache.Capacity())
	//global.BigCache.
	return nil
}
