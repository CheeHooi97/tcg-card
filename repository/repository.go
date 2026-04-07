package repository

import "gorm.io/gorm"

type Repositories struct {
	AdminRepo             AdminRepository
	BGSRepo               BGSRepository
	BGSUrlRepo            BGSUrlRepository
	BGSLoggingRepo        BGSLoggingRepository
	CGCRepo               CGCRepository
	CGCUrlRepo            CGCUrlRepository
	CGCLoggingRepo        CGCLoggingRepository
	CardRepo              CardRepository
	CardPriceRepo         CardPriceRepository
	PSARepo               PSARepository
	PSAUrlRepo            PSAUrlRepository
	PSALoggingRepo        PSALoggingRepository
	SetRepo               SetRepository
	TAGRepo               TAGRepository
	TAGUrlRepo            TAGUrlRepository
	TAGLoggingRepo        TAGLoggingRepository
	TokenRepo             TokenRepository
	UserRepo              UserRepository
	UserCardRepo          UserCardRepository
	UserCardSearchLogRepo UserCardSearchLogRepository
	UserDeviceRepo        UserDeviceRepository
}

func InitializeRepository(db *gorm.DB) *Repositories {
	return &Repositories{
		AdminRepo:             NewAdminRepository(db),
		BGSRepo:               NewBGSRepository(db),
		BGSUrlRepo:            NewBGSUrlRepository(db),
		BGSLoggingRepo:        NewBGSLoggingRepository(db),
		CGCRepo:               NewCGCRepository(db),
		CGCUrlRepo:            NewCGCUrlRepository(db),
		CGCLoggingRepo:        NewCGCLoggingRepository(db),
		CardRepo:              NewCardRepository(db),
		CardPriceRepo:         NewCardPriceRepository(db),
		PSARepo:               NewPSARepository(db),
		PSAUrlRepo:            NewPSAUrlRepository(db),
		PSALoggingRepo:        NewPSALoggingRepository(db),
		SetRepo:               NewSetRepository(db),
		TAGRepo:               NewTAGRepository(db),
		TAGUrlRepo:            NewTAGUrlRepository(db),
		TAGLoggingRepo:        NewTAGLoggingRepository(db),
		TokenRepo:             NewTokenRepository(db),
		UserCardSearchLogRepo: NewUserCardSearchLogRepository(db),
		UserCardRepo:          NewUserCardRepository(db),
		UserRepo:              NewUserRepository(db),
		UserDeviceRepo:        NewUserDeviceRepository(db),
	}
}
