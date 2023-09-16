package service

func Init() {
	sessInit()
	storageObj.New()
	paymentObj.New()
}
