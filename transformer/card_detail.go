package transformer

import (
	"pkm/model"
)

type CardDetail struct {
	Id             string     `json:"id"`
	Name           string     `json:"name"`
	Set            string     `json:"set"`
	Rarity         string     `json:"rarity"`
	Price          string     `json:"price"`
	IncreasedPrice string     `json:"increasedPrice"`
	PercentChange  string     `json:"percentChange"`
	PhotoUrl       string     `json:"photoUrl"`
	PSAGrade       *model.PSA `json:"psaGrade"`
	BGSGrade       *model.BGS `json:"bgsGrade"`
	TAGGrade       *model.TAG `json:"tagGrade"`
	CGCGrade       *model.CGC `json:"cgcGrade"`
	model.BaseModel
}

func ToCardDetail(c *model.Card, rarity, increasedPrice, percentChange string, psaGrade *model.PSA,
	bgsGrade *model.BGS, tagGrade *model.TAG, cgcGrade *model.CGC) *CardDetail {
	card := CardDetail{
		Id:             c.Id,
		Name:           c.Name,
		Set:            c.SetName,
		Rarity:         rarity,
		Price:          c.Ungrade,
		IncreasedPrice: increasedPrice,
		PercentChange:  percentChange,
		PhotoUrl:       c.PhotoUrl,
		PSAGrade:       psaGrade,
		BGSGrade:       bgsGrade,
		TAGGrade:       tagGrade,
		CGCGrade:       cgcGrade,
		BaseModel:      c.BaseModel,
	}

	return &card
}
