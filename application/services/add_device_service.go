package services

import (
	"crypto/rand"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/pipe-network/signaling-server/application/ports"
	"github.com/pipe-network/signaling-server/domain/values"
	"github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack/v5"
	"io"
)

var (
	NoMatchingMessage  = errors.New("could not find a message that match")
	NoUUIDForPublicKey = errors.New("no uuid for public key found")
	UUIDsDoesNotMatch  = errors.New("the given uuid does not match the stored one")
)

type AddDeviceService interface {
	OnAddDeviceMessage(connection *websocket.Conn, message []byte) error
}

type AddDeviceServiceImpl struct {
	publicKeyToUUID       map[values.Key]uuid.UUID
	keyPairStorage        ports.KeyPairStorage
	deviceTokenRepository ports.DeviceTokenRepository
}

func NewAddDeviceServiceImpl(
	keyPairStorage ports.KeyPairStorage,
	deviceTokenRepository ports.DeviceTokenRepository,
) AddDeviceService {
	return &AddDeviceServiceImpl{
		publicKeyToUUID:       map[values.Key]uuid.UUID{},
		keyPairStorage:        keyPairStorage,
		deviceTokenRepository: deviceTokenRepository,
	}
}

func (a *AddDeviceServiceImpl) OnAddDeviceMessage(connection *websocket.Conn, message []byte) error {
	addDeviceMessage, err := values.AddDeviceMessageFromBytes(message)
	if err != nil {
		return err
	}

	addDeviceRequest := values.AddDeviceRequestMessage{}
	err = msgpack.Unmarshal(addDeviceMessage.Data, &addDeviceRequest)
	if err == nil {
		err = a.onAddDeviceRequestMessage(connection, addDeviceMessage.PublicKey)
		if err != nil {
			return err
		}
		return nil
	}

	// AddDeviceRequestMessage was not successful, so we try AddDeviceSolvedMessage
	decryptedAddDeviceSolvedMessage, err := values.DecryptMessage(
		addDeviceMessage.Data,
		addDeviceMessage.Nonce,
		addDeviceMessage.PublicKey,
		a.keyPairStorage.PrivateKey(),
	)
	if err != nil {
		return err
	}

	addDeviceSolvedMessage := values.AddDeviceSolvedMessage{}
	err = msgpack.Unmarshal(decryptedAddDeviceSolvedMessage, &addDeviceSolvedMessage)
	if err != nil {
		return NoMatchingMessage
	}

	err = a.onAddDeviceSolvedMessage(addDeviceMessage.PublicKey, addDeviceSolvedMessage)
	if err != nil {
		return err
	}

	err = connection.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "device token was updated"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddDeviceServiceImpl) generateUUID(publicKey values.Key) uuid.UUID {
	v4uuid := uuid.NewV4()
	a.publicKeyToUUID[publicKey] = v4uuid
	return v4uuid
}

func (a *AddDeviceServiceImpl) onAddDeviceRequestMessage(
	connection *websocket.Conn,
	devicePublicKey values.Key,
) error {
	v4uuid := a.generateUUID(devicePublicKey)
	addDeviceControlMessage := values.NewAddDeviceControlMessage(v4uuid.String())
	addDeviceControlMessagePack, err := msgpack.Marshal(addDeviceControlMessage)
	if err != nil {
		return err
	}

	nonce := [values.NonceByteLength]byte{}
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	encryptedDeviceControlMessagePack := values.EncryptMessage(
		addDeviceControlMessagePack,
		nonce,
		devicePublicKey,
		a.keyPairStorage.PrivateKey(),
	)
	addDeviceMessage := values.AddDeviceMessage{
		Data:      encryptedDeviceControlMessagePack,
		PublicKey: a.keyPairStorage.PublicKey().Bytes(),
		Nonce:     nonce,
	}
	packedAddDeviceMessage := addDeviceMessage.ToBytes()
	err = connection.WriteMessage(websocket.BinaryMessage, packedAddDeviceMessage)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddDeviceServiceImpl) onAddDeviceSolvedMessage(
	devicePublicKey values.Key,
	addDeviceSolvedMessage values.AddDeviceSolvedMessage,
) error {
	storedUUID, ok := a.publicKeyToUUID[devicePublicKey]
	if !ok {
		return NoUUIDForPublicKey
	}

	if storedUUID.String() != addDeviceSolvedMessage.UUID {
		return UUIDsDoesNotMatch
	}

	err := a.deviceTokenRepository.CreateOrUpdateToken(values.Device{
		Token:     addDeviceSolvedMessage.DeviceToken,
		PublicKey: devicePublicKey.HexString(),
	})
	if err != nil {
		return err
	}

	return nil
}
