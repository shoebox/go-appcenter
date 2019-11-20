package appcenter

type AppCenterRepository interface {
	Upload(name string, name2 string) (string, error)
}

type AppCenterService struct {
	AppCenterRepository
}

func (service *AppCenterService) Upload(name string, name2 string) (string, error) {
	return "", nil
}
