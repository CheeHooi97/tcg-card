package transformer

import (
	"pkm/model"

	"github.com/ivpusic/grpool"
)

type PokemonList struct {
	Data struct {
		Items []struct {
			CardName    string           `json:"cardName"`
			CardNumber  string           `json:"cardNumber"`
			Variation   string           `json:"variation"`
			Grades      map[string]int64 `json:"grades"`
			GradesCount int64            `json:"count,omitempty"`
		} `json:"items"`
		Total int64 `json:"total"`
		Limit int64 `json:"limit"`
	} `json:"data"`
}

func ToPokemonList(d *model.PSA, count int64) *PokemonList {
	ticket := PokemonList{}

	return &ticket
}

func ToPokemonLists(d []*model.PSA, countMap map[string]int64) []*PokemonList {
	size := len(d)
	o := make([]*PokemonList, size)
	pool := grpool.NewPool(20, 20)
	pool.WaitCount(size)
	defer pool.Release()
	for n, item := range d {
		pool.JobQueue <- func(index int, val *model.PSA) func() {
			return func() {
				defer pool.JobDone()
				o[index] = ToPokemonList(val, countMap[val.Id])
			}
		}(n, item)
	}
	pool.WaitAll()
	return o
}
