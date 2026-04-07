package transformer

import (
	"pkm/model"

	"github.com/ivpusic/grpool"
)

type SetList struct {
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

func ToSetList(c *model.Card, rarity, increasedPrice, percentChange string) *SetList {
	card := SetList{
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

func ToSetLists(d []*model.Card, rarityMap, increasedPriceMap, percentChangeMap map[string]string) []*SetList {
	size := len(d)
	o := make([]*SetList, size)
	pool := grpool.NewPool(20, 20)
	pool.WaitCount(size)
	defer pool.Release()
	for n, item := range d {
		pool.JobQueue <- func(index int, val *model.Card) func() {
			return func() {
				defer pool.JobDone()
				o[index] = ToSetList(val, rarityMap[val.Id], increasedPriceMap[val.Id], percentChangeMap[val.Id])
			}
		}(n, item)
	}
	pool.WaitAll()
	return o
}
