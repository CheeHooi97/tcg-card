package service

import "pkm/repository"

type Services struct {
	AdminService             *AdminService
	BGSService               *BGSService
	BGSUrlService            *BGSUrlService
	BGSLoggingService        *BGSLoggingService
	CGCService               *CGCService
	CGCUrlService            *CGCUrlService
	CGCLoggingService        *CGCLoggingService
	CardService              *CardService
	CardPriceService         *CardPriceService
	PSAService               *PSAService
	PSAUrlService            *PSAUrlService
	PSALoggingService        *PSALoggingService
	SetService               *SetService
	TAGService               *TAGService
	TAGUrlService            *TAGUrlService
	TAGLoggingService        *TAGLoggingService
	TokenService             *TokenService
	UserService              *UserService
	UserCardService          *UserCardService
	UserCardSearchLogService *UserCardSearchLogService
	UserDeviceService        *UserDeviceService
}

func InitializeService(repos *repository.Repositories) *Services {
	return &Services{
		AdminService:             NewAdminService(repos.AdminRepo),
		BGSService:               NewBGSService(repos.BGSRepo),
		BGSUrlService:            NewBGSUrlService(repos.BGSUrlRepo),
		BGSLoggingService:        NewBGSLoggingService(repos.BGSLoggingRepo),
		CGCService:               NewCGCService(repos.CGCRepo),
		CGCUrlService:            NewCGCUrlService(repos.CGCUrlRepo),
		CGCLoggingService:        NewCGCLoggingService(repos.CGCLoggingRepo),
		CardService:              NewCardService(repos.CardRepo),
		CardPriceService:         NewCardPriceService(repos.CardPriceRepo),
		PSAService:               NewPSAService(repos.PSARepo),
		PSAUrlService:            NewPSAUrlService(repos.PSAUrlRepo),
		PSALoggingService:        NewPSALoggingService(repos.PSALoggingRepo),
		SetService:               NewSetService(repos.SetRepo),
		TAGService:               NewTAGService(repos.TAGRepo),
		TAGUrlService:            NewTAGUrlService(repos.TAGUrlRepo),
		TAGLoggingService:        NewTAGLoggingService(repos.TAGLoggingRepo),
		TokenService:             NewTokenService(repos.TokenRepo),
		UserService:              NewUserService(repos.UserRepo),
		UserCardService:          NewUserCardService(repos.UserCardRepo),
		UserCardSearchLogService: NewUserCardSearchLogService(repos.UserCardSearchLogRepo),
		UserDeviceService:        NewUserDeviceService(repos.UserDeviceRepo),
	}
}
