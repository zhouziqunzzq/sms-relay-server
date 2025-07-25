package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/zhouziqunzzq/sms-relay-server/common"
	"github.com/zhouziqunzzq/sms-relay-server/models"
)

func forwardSMSByEmail(ctx context.Context, smsRelayRequest models.SMSRelayRequest) error {
	if smsRelayRequest.PhoneNumber.ForwardDestinations.Email.IsEmpty() {
		return nil // No email to forward to
	}

	toAddr := smsRelayRequest.PhoneNumber.ForwardDestinations.Email.Email
	logger.Printf("Forwarding SMS to email: %s", toAddr)

	// Fetch SMTP credentials from Secrets Manager
	username, password, err := getSMTPCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch SMTP credentials: %w", err)
	}

	// Compose the email message
	fromHeader := fmt.Sprintf("From: %s\r\n", username)
	toHeader := fmt.Sprintf("To: %s\r\n", toAddr)
	subject := fmt.Sprintf("Subject: SMS Relay for %s - %s\r\n",
		smsRelayRequest.DeviceName, smsRelayRequest.PhoneNumber.Name)
	body := fmt.Sprintf("Device: %s (%s)\nPhone Number: %s (%s)\nFrom: %s\nMessage: %s",
		smsRelayRequest.DeviceName, smsRelayRequest.Device.ID,
		smsRelayRequest.PhoneNumber.Name, smsRelayRequest.PhoneNumber.PhoneNumber,
		smsRelayRequest.SMS.From, smsRelayRequest.SMS.Body)
	msg := []byte(fromHeader + toHeader + subject + "\r\n" + body)

	if useSSL {
		// Establish a TLS connection
		serverAddr := fmt.Sprintf("%s:%s", smtpServer, smtpPort)
		conn, err := tls.Dial("tcp", serverAddr, &tls.Config{
			InsecureSkipVerify: false,
		})
		if err != nil {
			return fmt.Errorf("failed to establish TLS connection: %w", err)
		}
		defer conn.Close()

		// Create an SMTP client over the TLS connection
		client, err := smtp.NewClient(conn, smtpServer)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Quit()

		// Authenticate and send the email
		if err := client.Auth(smtp.PlainAuth("", username, password, smtpServer)); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		if err := client.Mail(username); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}
		if err := client.Rcpt(toAddr); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get writer: %w", err)
		}
		if _, err := writer.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
		if err := writer.Close(); err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}
	} else {
		// Send the email without SSL
		err = smtp.SendMail(
			fmt.Sprintf("%s:%s", smtpServer, smtpPort),
			smtp.PlainAuth("", username, password, smtpServer),
			username,
			[]string{toAddr},
			msg,
		)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	logger.Printf("Successfully forwarded SMS to email: %s", toAddr)
	return nil
}

func getSMTPCredentials(ctx context.Context) (username string, password string, err error) {
	// Fetch SMTP username
	username, err = common.GetSecretValue(ctx, secretsClient, smtpUsernameSecretName, "username")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch SMTP username: %w", err)
	}

	// Fetch SMTP password
	password, err = common.GetSecretValue(ctx, secretsClient, smtpPasswordSecretName, "password")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch SMTP password: %w", err)
	}

	return username, password, nil
}
