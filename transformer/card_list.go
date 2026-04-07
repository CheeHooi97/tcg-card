package transformer

import (
	"pkm/model"

	"github.com/ivpusic/grpool"
)

type CardList struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Set            string `json:"set"`
	Rarity         string `json:"rarity"`
	Price          string `json:"price"`
	IncreasedPrice string `json:"increasedPrice"`
	PercentChange  string `json:"percentChange"`
	PhotoUrl       string `json:"photoUrl"`
	model.BaseModel
}

func ToCardList(c *model.Card, rarity, increasedPrice, percentChange string) *CardList {
	card := CardList{
		Id:             c.Id,
		Name:           c.Name,
		Set:            c.SetName,
		Rarity:         rarity,
		Price:          c.Ungrade,
		IncreasedPrice: increasedPrice,
		PercentChange:  percentChange,
		PhotoUrl:       c.PhotoUrl,
		BaseModel:      c.BaseModel,
	}

	return &card
}

func ToCardLists(d []*model.Card, rarityMap, increasedPriceMap, percentChangeMap map[string]string) []*CardList {
	size := len(d)
	o := make([]*CardList, size)
	pool := grpool.NewPool(20, 20)
	pool.WaitCount(size)
	defer pool.Release()
	for n, item := range d {
		pool.JobQueue <- func(index int, val *model.Card) func() {
			return func() {
				defer pool.JobDone()
				o[index] = ToCardList(val, rarityMap[val.Id], increasedPriceMap[val.Id], percentChangeMap[val.Id])
			}
		}(n, item)
	}
	pool.WaitAll()
	return o
}
