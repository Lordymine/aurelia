package telegram

import (
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

func (bc *BotController) handleText(c telebot.Context) error {
	return bc.processInput(c, c.Text(), nil, false)
}

func (bc *BotController) handlePhoto(c telebot.Context) error {
	photo := c.Message().Photo
	if photo == nil {
		return nil
	}

	if c.Message().AlbumID != "" {
		return bc.handlePhotoAlbum(c, photo)
	}

	return bc.processPhotoInput(c, strings.TrimSpace(c.Message().Caption), []albumPhoto{{
		messageID: c.Message().ID,
		photo:     *photo,
	}})
}

func (bc *BotController) handlePhotoAlbum(c telebot.Context, photo *telebot.Photo) error {
	albumID := c.Message().AlbumID
	isOwner := bc.storeAlbumPhoto(albumID, c.Message().ID, strings.TrimSpace(c.Message().Caption), *photo)
	if !isOwner {
		return nil
	}

	time.Sleep(900 * time.Millisecond)

	caption, photos, ok := bc.flushAlbumPhotos(albumID)
	if !ok {
		return nil
	}
	return bc.processPhotoInput(c, caption, photos)
}

func (bc *BotController) processPhotoInput(c telebot.Context, caption string, photos []albumPhoto) error {
	if len(photos) == 0 {
		return nil
	}

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.UploadingPhoto, 4*time.Second)
	defer stopTyping()

	text := strings.TrimSpace(caption)

	var photoPaths []string
	for _, p := range photos {
		filePath, err := bc.downloadTelegramFile(&p.photo.File, fmt.Sprintf("photo_%d.jpg", p.messageID))
		if err != nil {
			log.Printf("Failed to download photo: %v", err)
			continue
		}
		photoPaths = append(photoPaths, filePath)
	}

	if text == "" {
		if len(photoPaths) > 1 {
			text = "Analise estas imagens."
		} else {
			text = "Analise esta imagem."
		}
	}

	for _, path := range photoPaths {
		text += fmt.Sprintf("\n\nImagem: %s", path)
	}

	return bc.processInput(c, text, nil, false)
}

func (bc *BotController) handleDocument(c telebot.Context) error {
	doc := c.Message().Document
	if doc == nil {
		return nil
	}

	if isSupportedImageDocument(doc.FileName, doc.MIME) {
		return bc.handleImageDocument(c, doc)
	}

	if !isSupportedDocument(doc.FileName, doc.MIME) {
		log.Println("Unsupported document type:", doc.MIME)
		return SendContextText(c, unsupportedDocumentMessage)
	}

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	filePath, err := bc.downloadTelegramFile(&doc.File, doc.FileID+"_"+doc.FileName)
	if err != nil {
		log.Println("Failed to download file:", err)
		return SendContextText(c, downloadFailureMessage)
	}
	defer func() { _ = os.Remove(filePath) }()

	finalInput := buildDocumentInput(c.Message().Caption, doc.FileName, doc.MIME, filePath)
	return bc.processInput(c, finalInput, nil, false)
}

func (bc *BotController) handleImageDocument(c telebot.Context, doc *telebot.Document) error {
	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.UploadingPhoto, 4*time.Second)
	defer stopTyping()

	text := strings.TrimSpace(c.Message().Caption)
	if text == "" {
		text = "Analise esta imagem."
	}

	filePath, err := bc.downloadTelegramFile(&doc.File, doc.FileID+"_"+doc.FileName)
	if err != nil {
		log.Printf("Failed to download image document: %v", err)
		return bc.processInput(c, text, nil, false)
	}

	text += fmt.Sprintf("\n\nImagem: %s", filePath)
	return bc.processInput(c, text, nil, false)
}

func (bc *BotController) handleVoice(c telebot.Context) error {
	fileID, filename, ok := resolveAudioAttachment(c)
	if !ok {
		return nil
	}

	stopRecording := startChatActionLoop(bc.bot, c.Chat(), telebot.RecordingAudio, 4*time.Second)
	defer stopRecording()

	filePath, err := bc.downloadTelegramFile(&telebot.File{FileID: fileID}, fileID+"_"+filename)
	if err != nil {
		log.Println("Failed to download audio:", err)
		return SendContextText(c, downloadFailureMessage)
	}
	defer func() { _ = os.Remove(filePath) }()

	transcribedText, err := bc.transcribeAudioFile(filePath)
	if err != nil {
		var msgErr sendContextTextError
		if ok := errorAs(err, &msgErr); ok {
			return SendContextText(c, msgErr.Error())
		}
		return SendContextText(c, audioProcessingFailureMessage)
	}
	return bc.processInput(c, transcribedText, nil, true)
}

func isSupportedDocument(filename, mimeType string) bool {
	return strings.HasSuffix(filename, ".md") || mimeType == "application/pdf"
}

func isSupportedImageDocument(filename, mimeType string) bool {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(mimeType)), "image/") {
		return true
	}
	guessed := mime.TypeByExtension(strings.ToLower(filepath.Ext(filename)))
	return strings.HasPrefix(guessed, "image/")
}

func (bc *BotController) downloadTelegramFile(file *telebot.File, filename string) (string, error) {
	filePath := filepath.Join(os.TempDir(), filename)
	if err := bc.bot.Download(file, filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

func buildDocumentInput(caption, filename, mimeType, filePath string) string {
	var extractedText string
	if strings.HasSuffix(filename, ".md") {
		content, err := os.ReadFile(filePath)
		if err == nil {
			extractedText = string(content)
		}
	} else if mimeType == "application/pdf" {
		extractedText = fmt.Sprintf("[Parsed content of PDF %s]", filename)
	}

	return fmt.Sprintf("%s\n\n[Analise o anexo %s]:\n%s", caption, filename, extractedText)
}

func (bc *BotController) storeAlbumPhoto(albumID string, messageID int, caption string, photo telebot.Photo) bool {
	bc.albumMu.Lock()
	defer bc.albumMu.Unlock()

	album, ok := bc.pendingAlbums[albumID]
	if !ok {
		album = &pendingAlbum{ownerMessageID: messageID}
		bc.pendingAlbums[albumID] = album
	}
	if caption != "" && album.caption == "" {
		album.caption = caption
	}
	album.photos = append(album.photos, albumPhoto{messageID: messageID, photo: photo})
	return album.ownerMessageID == messageID
}

func (bc *BotController) flushAlbumPhotos(albumID string) (string, []albumPhoto, bool) {
	bc.albumMu.Lock()
	defer bc.albumMu.Unlock()

	album, ok := bc.pendingAlbums[albumID]
	if !ok {
		return "", nil, false
	}
	delete(bc.pendingAlbums, albumID)

	photos := append([]albumPhoto(nil), album.photos...)
	sort.SliceStable(photos, func(i, j int) bool {
		return photos[i].messageID < photos[j].messageID
	})
	return album.caption, photos, true
}

func resolveAudioAttachment(c telebot.Context) (string, string, bool) {
	switch {
	case c.Message().Voice != nil:
		return c.Message().Voice.FileID, "voice.ogg", true
	case c.Message().Audio != nil:
		return c.Message().Audio.FileID, "audio.mp3", true
	default:
		return "", "", false
	}
}

func (bc *BotController) transcribeAudioFile(filePath string) (string, error) {
	if bc.stt == nil || !bc.stt.IsAvailable() {
		return "", SendContextTextError(audioNotConfiguredMessage)
	}

	log.Printf("Enviando audio [%s] para transcricao via Groq API...", filePath)
	transcribedText, err := bc.stt.Transcribe(context.Background(), filePath)
	if err != nil {
		log.Printf("Groq STT error: %v\n", err)
		return "", SendContextTextError(audioProcessingFailureMessage)
	}
	if strings.TrimSpace(transcribedText) == "" {
		return "", SendContextTextError(emptyAudioMessage)
	}
	return transcribedText, nil
}

type sendContextTextError string

// SendContextTextError creates a sendContextTextError.
func SendContextTextError(message string) error {
	return sendContextTextError(message)
}

func (e sendContextTextError) Error() string {
	return string(e)
}

func errorAs(err error, target *sendContextTextError) bool {
	if err == nil {
		return false
	}
	value, ok := err.(sendContextTextError)
	if !ok {
		return false
	}
	*target = value
	return true
}
