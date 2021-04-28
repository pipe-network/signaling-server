package services

import (
	"fmt"
	"github.com/NaySoftware/go-fcm"
	"github.com/pipe-network/signaling-server/application/ports"
	"github.com/pipe-network/signaling-server/application/services"
)

type FCMNotificationService struct {
	ServerKey string
}

var _ ports.NotificationService = (*FCMNotificationService)(nil)

func NewFCMNotificationService(
	flagService services.FlagService,
) ports.NotificationService {
	fcmNotificationService := &FCMNotificationService{
		ServerKey: flagService.String(services.FCMServerKey),
	}
	return fcmNotificationService
}

func (f *FCMNotificationService) Notify(title string, message string, data interface{}, deviceId string) error {
	client := fcm.NewFcmClient(f.ServerKey)
	//payload := &fcm.NotificationPayload{Title: title, Body: message}
	client.SetMsgData(data)
	//client.SetNotificationPayload(payload)
	client.AppendDevices([]string{deviceId})

	status, err := client.Send()
	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}
	return err
}
