package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

func (bc *BotController) handleText(c telebot.Context) error {
	return bc.processInput(c, c.Text(), false)
}

func (bc *BotController) handleDocument(c telebot.Context) error {
	doc := c.Message().Document
	if doc == nil {
		return nil
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
	return bc.processInput(c, finalInput, false)
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
	return bc.processInput(c, transcribedText, true)
}

func isSupportedDocument(filename, mimeType string) bool {
	return strings.HasSuffix(filename, ".md") || mimeType == "application/pdf"
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
