package dao

import (
	"context"
	"encoding/json"
	"github.com/ecodeclub/ekit/slice"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const ArticleIndexName = "article_index"
const TagIndexName = "tags_index"

type Article struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Status  int32    `json:"status"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type ArticleElasticDAO struct {
	client *elastic.Client
}

func NewArticleElasticDAO(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{client: client}
}

func (h *ArticleElasticDAO) Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error) {
	queryString := strings.Join(keywords, " ")
	// 文章，标题或者内容任何一个匹配上
	// 并且状态 status 必须是已发表的状态

	// status 精确查找
	statusTerm := elastic.NewTermQuery("status", 2)

	// 标签命中
	tagArtIdAnys := slice.Map(tagArtIds, func(idx int, src int64) any {
		return src
	})

	// 内容或者标题，模糊查找（match）
	title := elastic.NewMatchQuery("title", queryString)
	content := elastic.NewMatchQuery("content", queryString)
	or := elastic.NewBoolQuery().Should(title, content)
	if len(tagArtIds) > 0 {
		tag := elastic.NewTermsQuery("id", tagArtIdAnys...).
			Boost(2.0)
		or = or.Should(tag)
	}

	and := elastic.NewBoolQuery().Must(statusTerm, or)

	//return NewSearcher[Article](h.client, ArticleIndexName).
	//	Query(and).Do(ctx)
	resp, err := h.client.Search(ArticleIndexName).Query(and).Do(ctx)
	if err != nil {
		return nil, err
	}
	var res []Article
	for _, hit := range resp.Hits.Hits {
		var art Article
		err = json.Unmarshal(hit.Source, &art)
		if err != nil {
			return nil, err
		}
		res = append(res, art)
	}
	return res, nil
}

func (h *ArticleElasticDAO) InputArticle(ctx context.Context, art Article) error {
	_, err := h.client.Index().Index(ArticleIndexName).
		// 为什么要指定 ID？
		// 确保后面文章更新的时候，我们这里产生类似的两条数据，而是更新了数据
		Id(strconv.FormatInt(art.Id, 10)).
		BodyJson(art).Do(ctx)
	return err
}

func NewArticleRepository(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{
		client: client,
	}
}

type Searcher[T any] struct {
	client  *elastic.Client
	idxName []string
	query   elastic.Query
}

func NewSearcher[T any](client *elastic.Client, idxName ...string) *Searcher[T] {
	return &Searcher[T]{
		client:  client,
		idxName: idxName,
	}
}

func (s *Searcher[T]) Query(q elastic.Query) *Searcher[T] {
	s.query = q
	return s
}

//
//func (s *Searcher[T]) Do1(ctx context.Context) (T, error) {
//
//}

func (s *Searcher[T]) Do(ctx context.Context) ([]T, error) {
	resp, err := s.client.Search(s.idxName...).Do(ctx)
	res := make([]T, 0, resp.Hits.TotalHits.Value)
	for _, hit := range resp.Hits.Hits {
		var t T
		err = json.Unmarshal(hit.Source, &t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

//func (s *Searcher[T]) Resp() *elastic.SearchResult {
//
//}
