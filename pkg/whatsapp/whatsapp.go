package whatsapp

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	qrCode "github.com/skip2/go-qrcode"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	waproto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
)

var WhatsAppDatastore *sqlstore.Container
var WhatsAppClient = make(map[string]*whatsmeow.Client)

func init() {
	var err error

	dbType, err := env.GetEnvString("WHATSAPP_DATASTORE_TYPE")
	if err != nil {
		log.Print(nil).Fatal("Error Parse Environment Variable for WhatsApp Client Datastore Type")
	}

	dbURI, err := env.GetEnvString("WHATSAPP_DATASTORE_URI")
	if err != nil {
		log.Print(nil).Fatal("Error Parse Environment Variable for WhatsApp Client Datastore URI")
	}

	datastore, err := sqlstore.New(dbType, dbURI, nil)
	if err != nil {
		log.Print(nil).Fatal("Error Connect WhatsApp Client Datastore")
	}

	WhatsAppDatastore = datastore
}

func WhatsAppInitClient(device *store.Device, jid string) {
	var err error

	if WhatsAppClient[jid] == nil {
		if device == nil {
			// Initialize New WhatsApp Client Device in Datastore
			device = WhatsAppDatastore.NewDevice()
		}

		// Set Client Properties
		store.CompanionProps.Os = proto.String("Go WhatsApp Multi-Device REST")
		store.CompanionProps.PlatformType = waproto.CompanionProps_DESKTOP.Enum()

		// Set Client Versions
		version.Major, err = env.GetEnvInt("WHATSAPP_VERSION_MAJOR")
		if err == nil {
			store.CompanionProps.Version.Primary = proto.Uint32(uint32(version.Major))
		}
		version.Minor, err = env.GetEnvInt("WHATSAPP_VERSION_MINOR")
		if err == nil {
			store.CompanionProps.Version.Secondary = proto.Uint32(uint32(version.Minor))
		}
		version.Patch, err = env.GetEnvInt("WHATSAPP_VERSION_PATCH")
		if err == nil {
			store.CompanionProps.Version.Tertiary = proto.Uint32(uint32(version.Patch))
		}

		// Initialize New WhatsApp Client
		// And Save it to The Map
		WhatsAppClient[jid] = whatsmeow.NewClient(device, nil)
	}
}

func WhatsAppGenerateQR(qrChan <-chan whatsmeow.QRChannelItem) (string, int) {
	qrChanCode := make(chan string)
	qrChanTimeout := make(chan int)
	qrChanBase64 := make(chan string)

	// Get QR Code Data and Timeout
	go func() {
		for evt := range qrChan {
			if evt.Event == "code" {
				qrChanCode <- evt.Code
				qrChanTimeout <- int(evt.Timeout.Seconds())
			}
		}
	}()

	// Generate QR Code Data to PNG Base64 Format
	go func() {
		select {
		case tmp := <-qrChanCode:
			png, _ := qrCode.Encode(tmp, qrCode.Medium, 256)
			qrChanBase64 <- base64.StdEncoding.EncodeToString(png)
		}
	}()

	// Return QR Code and Timeout Information
	return <-qrChanBase64, <-qrChanTimeout
}

func WhatsAppLogin(jid string) (string, int, error) {
	if WhatsAppClient[jid] != nil {
		// Make Sure WebSocket Connection is Disconnected
		WhatsAppClient[jid].Disconnect()

		if WhatsAppClient[jid].Store.ID == nil {
			// Device ID is not Exist
			// Generate QR Code
			qrChanGenerate, _ := WhatsAppClient[jid].GetQRChannel(context.Background())

			// Connect WebSocket while Initialize QR Code Data to be Sent
			err := WhatsAppClient[jid].Connect()
			if err != nil {
				return "", 0, err
			}

			// Get Generated QR Code and Timeout Information
			qrImage, qrTimeout := WhatsAppGenerateQR(qrChanGenerate)

			// Return QR Code in Base64 Format and Timeout Information
			return "data:image/png;base64," + qrImage, qrTimeout, nil
		} else {
			// Device ID is Exist
			// Reconnect WebSocket
			err := WhatsAppReconnect(jid)
			if err != nil {
				return "", 0, err
			}

			return "WhatsApp Client is Reconnected", 0, nil
		}
	}

	// Return Error WhatsApp Client is not Valid
	return "", 0, errors.New("WhatsApp Client is not Valid")
}

func WhatsAppReconnect(jid string) error {
	if WhatsAppClient[jid] != nil {
		// Make Sure WebSocket Connection is Disconnected
		WhatsAppClient[jid].Disconnect()

		// Make Sure Store ID is not Empty
		// To do Reconnection
		if WhatsAppClient[jid].Store.ID != nil {
			err := WhatsAppClient[jid].Connect()
			if err != nil {
				return err
			}

			return nil
		}

		return errors.New("WhatsApp Client Store ID is Empty, Please Re-Login and Scan QR Code Again")
	}

	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppLogout(jid string) error {
	if WhatsAppClient[jid] != nil {
		// Logout WhatsApp Client and Disconnect from WebSocket
		err := WhatsAppClient[jid].Logout()
		if err != nil {
			// Force Disconnect
			WhatsAppClient[jid].Disconnect()

			// Manually Delete Device from Datastore Store
			err = WhatsAppClient[jid].Store.Delete()
			if err != nil {
				return err
			}
		}

		// Free WhatsApp Client Map
		WhatsAppClient[jid] = nil
		delete(WhatsAppClient, jid)

		return nil
	}

	// Return Error WhatsApp Client is not Valid
	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppClientIsOK(jid string) error {
	// Make Sure WhatsApp Client is Connected
	if !WhatsAppClient[jid].IsConnected() {
		return errors.New("WhatsApp Client is not Connected")
	}

	// Make Sure WhatsApp Client is Logged In
	if !WhatsAppClient[jid].IsLoggedIn() {
		return errors.New("WhatsApp Client is not Logged In")
	}

	return nil
}

func WhatsAppComposeJID(jid string) types.JID {
	// Decompose JID First Before Recomposing
	jid = WhatsAppDecomposeJID(jid)

	// Check if JID Contains '-' Symbol
	if strings.ContainsRune(jid, '-') {
		// Check if the JID is a Group ID
		if len(strings.SplitN(jid, "-", 2)) == 2 {
			// Return JID as Group Server (@g.us)
			return types.NewJID(jid, types.GroupServer)
		}
	}

	// Return JID as Default User Server (@s.whatsapp.net)
	return types.NewJID(jid, types.DefaultUserServer)
}

func WhatsAppDecomposeJID(jid string) string {
	// Check if JID Contains '@' Symbol
	if strings.ContainsRune(jid, '@') {
		// Split JID Based on '@' Symbol
		// and Get Only The First Section Before The Symbol
		buffers := strings.Split(jid, "@")
		jid = buffers[0]
	}

	// Check if JID First Character is '+' Symbol
	if jid[0] == '+' {
		// Remove '+' Symbol from JID
		jid = jid[1:]
	}

	return jid
}

func WhatsAppSendText(jid string, rjid string, message string) error {
	if WhatsAppClient[jid] != nil {
		// Make Sure WhatsApp Client is OK
		err := WhatsAppClientIsOK(jid)
		if err != nil {
			return err
		}

		// Compose WhatsApp Proto
		content := &waproto.Message{
			Conversation: proto.String(message),
		}

		// Send WhatsApp Message Proto
		_, err = WhatsAppClient[jid].SendMessage(WhatsAppComposeJID(rjid), "", content)
		if err != nil {
			return err
		}

		return nil
	}

	// Return Error WhatsApp Client is not Valid
	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppSendLocation(jid string, rjid string, latitude float64, longitude float64) error {
	if WhatsAppClient[jid] != nil {
		// Make Sure WhatsApp Client is OK
		err := WhatsAppClientIsOK(jid)
		if err != nil {
			return err
		}

		// Compose WhatsApp Proto
		content := &waproto.Message{
			LocationMessage: &waproto.LocationMessage{
				DegreesLatitude:  proto.Float64(latitude),
				DegreesLongitude: proto.Float64(longitude),
			},
		}

		// Send WhatsApp Message Proto
		_, err = WhatsAppClient[jid].SendMessage(WhatsAppComposeJID(rjid), "", content)
		if err != nil {
			return err
		}

		return nil
	}

	// Return Error WhatsApp Client is not Valid
	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppSendDocument(jid string, rjid string, fileBytes []byte, fileType string, fileName string) error {
	if WhatsAppClient[jid] != nil {
		// Make Sure WhatsApp Client is OK
		err := WhatsAppClientIsOK(jid)
		if err != nil {
			return err
		}

		// Upload File to WhatsApp Storage Server
		fileUploaded, err := WhatsAppClient[jid].Upload(context.Background(), fileBytes, whatsmeow.MediaDocument)

		// Compose WhatsApp Proto
		content := &waproto.Message{
			DocumentMessage: &waproto.DocumentMessage{
				Url:           proto.String(fileUploaded.URL),
				DirectPath:    proto.String(fileUploaded.DirectPath),
				Mimetype:      proto.String(fileType),
				Title:         proto.String(fileName),
				FileName:      proto.String(fileName),
				FileLength:    proto.Uint64(fileUploaded.FileLength),
				FileSha256:    fileUploaded.FileSHA256,
				FileEncSha256: fileUploaded.FileEncSHA256,
				MediaKey:      fileUploaded.MediaKey,
			},
		}

		// Send WhatsApp Message Proto
		_, err = WhatsAppClient[jid].SendMessage(WhatsAppComposeJID(rjid), "", content)
		if err != nil {
			return err
		}

		return nil
	}

	// Return Error WhatsApp Client is not Valid
	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppSendImage(jid string, rjid string, imageBytes []byte, imageType string, imageCaption string) error {
	if WhatsAppClient[jid] != nil {
		// Make Sure WhatsApp Client is OK
		err := WhatsAppClientIsOK(jid)
		if err != nil {
			return err
		}

		// Upload Image to WhatsApp Storage Server
		imageUploaded, err := WhatsAppClient[jid].Upload(context.Background(), imageBytes, whatsmeow.MediaImage)

		// Compose WhatsApp Proto
		content := &waproto.Message{
			ImageMessage: &waproto.ImageMessage{
				Url:           proto.String(imageUploaded.URL),
				DirectPath:    proto.String(imageUploaded.DirectPath),
				Mimetype:      proto.String(imageType),
				Caption:       proto.String(imageCaption),
				FileLength:    proto.Uint64(imageUploaded.FileLength),
				FileSha256:    imageUploaded.FileSHA256,
				FileEncSha256: imageUploaded.FileEncSHA256,
				MediaKey:      imageUploaded.MediaKey,
			},
		}

		// Send WhatsApp Message Proto
		_, err = WhatsAppClient[jid].SendMessage(WhatsAppComposeJID(rjid), "", content)
		if err != nil {
			return err
		}

		return nil
	}

	// Return Error WhatsApp Client is not Valid
	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppSendAudio(jid string, rjid string, audioBytes []byte, audioType string) error {
	if WhatsAppClient[jid] != nil {
		// Make Sure WhatsApp Client is OK
		err := WhatsAppClientIsOK(jid)
		if err != nil {
			return err
		}

		// Upload Audio to WhatsApp Storage Server
		audioUploaded, err := WhatsAppClient[jid].Upload(context.Background(), audioBytes, whatsmeow.MediaAudio)

		// Compose WhatsApp Proto
		content := &waproto.Message{
			AudioMessage: &waproto.AudioMessage{
				Url:           proto.String(audioUploaded.URL),
				DirectPath:    proto.String(audioUploaded.DirectPath),
				Mimetype:      proto.String(audioType),
				FileLength:    proto.Uint64(audioUploaded.FileLength),
				FileSha256:    audioUploaded.FileSHA256,
				FileEncSha256: audioUploaded.FileEncSHA256,
				MediaKey:      audioUploaded.MediaKey,
			},
		}

		// Send WhatsApp Message Proto
		_, err = WhatsAppClient[jid].SendMessage(WhatsAppComposeJID(rjid), "", content)
		if err != nil {
			return err
		}

		return nil
	}

	// Return Error WhatsApp Client is not Valid
	return errors.New("WhatsApp Client is not Valid")
}

func WhatsAppSendVideo(jid string, rjid string, videoBytes []byte, videoType string, videoCaption string) error {
	if WhatsAppClient[jid] != nil {
		// Make Sure WhatsApp Client is OK
		err := WhatsAppClientIsOK(jid)
		if err != nil {
			return err
		}

		// Upload Video to WhatsApp Storage Server
		videoUploaded, err := WhatsAppClient[jid].Upload(context.Background(), videoBytes, whatsmeow.MediaVideo)

		// Compose WhatsApp Proto
		content := &waproto.Message{
			VideoMessage: &waproto.VideoMessage{
				Url:           proto.String(videoUploaded.URL),
				DirectPath:    proto.String(videoUploaded.DirectPath),
				Mimetype:      proto.String(videoType),
				Caption:       proto.String(videoCaption),
				FileLength:    proto.Uint64(videoUploaded.FileLength),
				FileSha256:    videoUploaded.FileSHA256,
				FileEncSha256: videoUploaded.FileEncSHA256,
				MediaKey:      videoUploaded.MediaKey,
			},
		}

		// Send WhatsApp Message Proto
		_, err = WhatsAppClient[jid].SendMessage(WhatsAppComposeJID(rjid), "", content)
		if err != nil {
			return err
		}

		return nil
	}

	// Return Error WhatsApp Client is not Valid
	return errors.New("WhatsApp Client is not Valid")
}
