package ports

type NotificationService interface {
	Notify(title string, message string, data interface{}, deviceId string) error
}
