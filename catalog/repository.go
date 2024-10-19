package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	elastic "github.com/olivere/elastic/v7"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	CreateProduct(context.Context, Product) error
	GetProductById(context.Context, string) (*Product, error)
	GetProducts(context.Context, uint64, uint64) ([]Product, error)
	GetProductsByIds(context.Context, []string) ([]Product, error)
	SearchProduct(c context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &elasticRepository{client: client}, err
}

func (r *elasticRepository) Close() {

}

func (r *elasticRepository) CreateProduct(c context.Context, p Product) error {
	_, err := r.client.Index().
		Index("catalog.product").
		Id(p.Id).
		BodyJson(
			productDocument{
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		).
		Do(c)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (r *elasticRepository) GetProductById(c context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog.product").
		Id(id).
		Do(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}

	p := productDocument{}
	if err := json.Unmarshal(res.Source, &p); err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		Id:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *elasticRepository) GetProducts(c context.Context, skip, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog.product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
		Do(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []Product{}

	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, Product{
				Id:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}
func (r *elasticRepository) GetProductsByIds(c context.Context, ids []string) ([]Product, error) {

	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().
			Index("catalog.product").
			Id(id))
	}

	res, err := r.client.MultiGet().
		Add(items...).
		Do(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []Product{}

	for _, hit := range res.Docs {
		p := productDocument{}
		if err := json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, Product{
				Id:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}
func (r *elasticRepository) SearchProduct(c context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog.product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(take)).
		Do(c)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []Product{}

	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, Product{
				Id:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}
