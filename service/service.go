package service

import . "prepaid-card/dao"

func NewMerchantService() MerchantService {
	merchantService := MerchantService{}
	merchantService.transactionDao = TransactionDAO{}
	merchantService.accountDao = AccountDAO{}

	return merchantService
}

func NewUserService() UserService {
	userService := UserService{}
	userService.dao = AccountDAO{}

	return userService
}
