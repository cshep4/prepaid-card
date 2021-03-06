package service

import . "simple-prepaid-card/dao"

func NewMerchantService() MerchantService {
	merchantService := MerchantService{}
	merchantService.transactionDao = TransactionDAO{}
	merchantService.accountDao = AccountDAO{}

	return merchantService
}

func NewUserService() UserService {
	userService := UserService{}
	userService.accountDao = AccountDAO{}
	userService.transactionDao = TransactionDAO{}

	return userService
}
