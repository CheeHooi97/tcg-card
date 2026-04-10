package repository

import (
	"errors"
	"fmt"
	"pkm/model"
	"pkm/utils"
	"regexp"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type CardRepository interface {
	Create(card *model.Card) error
	GetById(id string) (*model.Card, error)
	GetByCardNameAndSet(cardName, set string) bool
	SearchCardKeywords(keyword string) ([]*model.Card, error)
	SearchCardBySort(sortType, userId string) ([]*model.Card, error)
	GetAllCards() ([]*model.Card, error)
	Update(card *model.Card) error
	Delete(id string) error
}

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) Create(card *model.Card) error {
	result := r.db.Create(card)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cardRepository) GetById(id string) (*model.Card, error) {
	var card model.Card
	result := r.db.First(&card, id)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return &card, nil
}

func (r *cardRepository) GetByCardNameAndSet(cardName, set string) bool {
	var auth model.Card
	result := r.db.
		Where("name = ? AND set_name = ?", cardName, set).
		First(&auth)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func (r *cardRepository) SearchCardKeywords(keyword string) ([]*model.Card, error) {
	var cards []*model.Card
	query := r.db

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return cards, nil
	}

	lower := strings.ToLower(keyword)

	// Regex patterns
	fullSetRegex := regexp.MustCompile(`\b(\d{1,3})\/(\d{1,3})\b`)
	singleNumberRegex := regexp.MustCompile(`(?i)#?\b\d+\b`)

	// 1️⃣ If contains "pokemon" → search set name
	if strings.Contains(lower, "pokemon") {
		query = query.Where(
			"set_name LIKE ?",
			"%"+utils.CapitalizeFirst(keyword)+"%",
		)
	}

	// 2️⃣ Handle full card format: 010/120
	if fullSetRegex.MatchString(keyword) {
		matches := fullSetRegex.FindStringSubmatch(keyword)

		// Normalize card number: 010 → 10
		cardNo := strings.TrimLeft(matches[1], "0")
		if cardNo == "" {
			cardNo = "0"
		}

		setTotal := matches[2]

		query = query.
			Where("name REGEXP ?", fmt.Sprintf("#%s(\\b|$)", cardNo)).
			Where("set_number = ?", setTotal)

		// Remove "010/120" from keyword so name search stays clean
		keyword = fullSetRegex.ReplaceAllString(keyword, "")
	}

	// 3️⃣ Handle single card number: #20 or 20
	if singleNumberRegex.MatchString(keyword) {
		number := singleNumberRegex.FindString(keyword)
		number = strings.TrimPrefix(number, "#")
		number = strings.TrimLeft(number, "0")
		if number == "" {
			number = "0"
		}

		query = query.Where(
			"name REGEXP ?",
			fmt.Sprintf("#%s(\\b|$)", number),
		)

		// Remove number so it doesn't interfere with name search
		keyword = singleNumberRegex.ReplaceAllString(keyword, "")
	}

	// 4️⃣ Remaining text → card name
	keyword = strings.TrimSpace(keyword)
	if keyword != "" && !strings.Contains(strings.ToLower(keyword), "pokemon") {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// Execute query
	result := query.Find(&cards)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return cards, nil
}

func (r *cardRepository) SearchCardBySort(sortType, userId string) ([]*model.Card, error) {
	var cards []*model.Card
	query := r.db

	switch sortType {
	case "price_asc":
		query = query.
			Where("ungraded IS NOT NULL AND ungraded != ''").
			Order("CAST(REPLACE(ungraded,'$','') AS DECIMAL(10,2)) ASC")
	case "price_desc":
		query = query.
			Where("ungraded IS NOT NULL AND ungraded != ''").
			Order("CAST(REPLACE(ungraded,'$','') AS DECIMAL(10,2)) DESC")
	case "alphabet_asc":
		query = query.Order("name ASC")
	case "alphabet_desc":
		query = query.Order("name DESC")
	case "popular":
		case "momentum_up", "momentum_down":
			type MomentumPrice struct {
				CardID string
				Price  string
			}

			var rows []MomentumPrice

			err := r.db.
				Table("card_price").
				Select("card_id, price").
				Where("created_date_time <= NOW() - INTERVAL '7 days'").
				Order("created_date_time DESC").
				Find(&rows).Error

		if err != nil {
			return nil, err
		}

		weekAgoPrice := make(map[string]string)
			for _, row := range rows {
				if _, exists := weekAgoPrice[row.CardID]; !exists {
					weekAgoPrice[row.CardID] = row.Price
				}
			}

		if len(weekAgoPrice) == 0 {
			return cards, nil
		}

		var allCards []model.Card
		err = r.db.
			Where("ungraded IS NOT NULL AND ungraded != ''").
			Find(&allCards).Error

		if err != nil {
			return nil, err
		}

			type Momentum struct {
				CardID string
				Delta  float64
			}

		momentums := make([]Momentum, 0)

		for _, c := range allCards {
			oldPrice, ok := weekAgoPrice[c.Id]
			if !ok {
				continue
			}

			old := utils.ParsePrice(oldPrice)

			current := utils.ParsePrice(c.Ungrade)
			if current <= 0 {
				continue
			}

				momentums = append(momentums, Momentum{
					CardID: c.Id,
					Delta:  current - old,
				})
		}

		if len(momentums) == 0 {
			return cards, nil
		}

		sort.Slice(momentums, func(i, j int) bool {
			if sortType == "momentum_down" {
				return momentums[i].Delta < momentums[j].Delta
			}
			return momentums[i].Delta > momentums[j].Delta
		})

	ids := make([]string, 0, 100)
	for _, m := range momentums {
		ids = append(ids, m.CardID)
	}

		err = r.db.
			Where("id IN ?", ids).
			Order(utils.OrderByField("id", ids)).
			Find(&cards).Error

		return cards, err
	case "relevance":
		// relevance
		var ids []string

		err := r.db.
			Table("user_card_search_log").
			Select("card_id").
			Where("user_id = ?", userId).
			Order("created_date_time DESC").
			Limit(100).
			Pluck("card_id", &ids).Error

		if err != nil {
			return nil, err
		}

		if len(ids) == 0 {
			return cards, nil
		}

		err = query.
			Where("id IN ?", ids).
			Order(utils.OrderByField("id", ids)).
			Find(&cards).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	default:
	}

	// Execute query
	result := query.Find(&cards)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return cards, nil
}

// func (r *cardRepository) SearchCardList(keyword, set, order string) ([]*model.Card, error) {
// 	var cards []*model.Card

// 	query := r.db

// 	if keyword != "" {
// 		if strings.Contains(keyword, "Pokemon") {
// 			keyword = utils.CapitalizeFirst(keyword)
// 			query = query.Where("set_name LIKE ?", "%"+keyword+"%")
// 		} else {
// 			re := regexp.MustCompile(`\d+`)
// 			numbers := re.FindAllString(keyword, -1)

// 			namePart := keyword
// 			for _, num := range numbers {
// 				namePart = strings.ReplaceAll(namePart, num, "")
// 			}
// 			namePart = strings.TrimSpace(namePart)

// 			if namePart != "" {
// 				query = query.Where("name LIKE ?", "%"+namePart+"%")
// 			}

// 			for _, num := range numbers {
// 				query = query.Where("name REGEXP ?", fmt.Sprintf(`#\w*%s\b`, num))
// 			}
// 		}
// 	}

// 	if set != "" {
// 		query = query.Where("set_name = ?", set)
// 	}

// 	if order != "" {
// 		switch order {
// 		case "alphabet_asc":
// 			query = query.Where("name LIKE ?", "%#%").
// 				Order("CAST(SUBSTRING_INDEX(name, '#', -1) AS UNSIGNED) ASC")
// 		case "alphabet_desc":
// 			query = query.Where("name LIKE ?", "%#%").
// 				Order("CAST(SUBSTRING_INDEX(name, '#', -1) AS UNSIGNED) DESC")
// 		case "price_asc":
// 			query = query.Order("CAST(REPLACE(ungrade, '$', '') AS DECIMAL(10,2)) ASC")
// 		case "price_desc":
// 			query = query.Order("CAST(REPLACE(ungrade, '$', '') AS DECIMAL(10,2)) DESC")
// 		case "popular":
// 			query = query.
// 				Joins(`
// 			JOIN (
// 				SELECT card_id, COUNT(*) AS search_count
// 				FROM usercardsearchlog
// 				WHERE created_date_time >= DATE_SUB(NOW(), INTERVAL 7 DAY)
// 				GROUP BY card_id
// 			) ucs ON ucs.card_id = cards.id
// 		`).
// 				Order("ucs.search_count DESC")
// 		case "momentum_up":
// 			query = query.
// 				Joins(`
// 			JOIN (
// 				SELECT
// 					card_id,
// 					MAX(CASE WHEN created_date_time >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN price END) -
// 					MIN(CASE WHEN created_date_time < DATE_SUB(NOW(), INTERVAL 7 DAY) THEN price END) AS price_change
// 				FROM ungrade
// 				GROUP BY card_id
// 			) cp ON cp.card_id = cards.id
// 		`).
// 				Order("cp.price_change DESC")
// 		case "momentum_down":
// 		default:
// 			// "relevance"
// 		}

// 	}

// 	result := query.Find(&cards)
// 	if result.Error != nil {
// 		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return nil, result.Error
// 		}
// 	}
// 	return cards, nil
// }

func (r *cardRepository) GetAllCards() ([]*model.Card, error) {
	var card []*model.Card
	result := r.db.
		Order("created_date_time ASC").
		Find(&card)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	return card, nil
}

func (r *cardRepository) Update(card *model.Card) error {
	result := r.db.Save(card)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *cardRepository) Delete(id string) error {
	result := r.db.Model(&model.Card{}).Where("id = ?", id).Update("status", false)
	return result.Error
}

