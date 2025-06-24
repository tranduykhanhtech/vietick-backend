package email

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
	"vietick-backend/internal/config"
)

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

func (e *EmailService) SendEmailVerification(toEmail, toName, verificationToken string) error {
	verificationURL := fmt.Sprintf("https://vietick.com/verify-email?token=%s", verificationToken)
	
	subject := "Verify Your VietTick Account"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Verification</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #1DA1F2;">Welcome to VietTick!</h1>
        <p>Hi %s,</p>
        <p>Thank you for signing up for VietTick. To complete your registration, please verify your email address by clicking the button below:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #1DA1F2; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">Verify Email Address</a>
        </div>
        <p>If the button doesn't work, you can also copy and paste the following link into your browser:</p>
        <p style="word-break: break-all; color: #1DA1F2;">%s</p>
        <p>This verification link will expire in 24 hours.</p>
        <p>If you didn't create an account with VietTick, please ignore this email.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="font-size: 12px; color: #666;">This is an automated message, please do not reply to this email.</p>
    </div>
</body>
</html>
	`, toName, verificationURL, verificationURL)

	return e.sendEmail(toEmail, toName, subject, body)
}

func (e *EmailService) SendPasswordReset(toEmail, toName, resetToken string) error {
	resetURL := fmt.Sprintf("https://vietick.com/reset-password?token=%s", resetToken)
	
	subject := "Reset Your VietTick Password"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #1DA1F2;">Password Reset Request</h1>
        <p>Hi %s,</p>
        <p>We received a request to reset your password for your VietTick account. Click the button below to reset your password:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #1DA1F2; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
        </div>
        <p>If the button doesn't work, you can also copy and paste the following link into your browser:</p>
        <p style="word-break: break-all; color: #1DA1F2;">%s</p>
        <p>This password reset link will expire in 1 hour.</p>
        <p>If you didn't request a password reset, please ignore this email. Your password will not be changed.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="font-size: 12px; color: #666;">This is an automated message, please do not reply to this email.</p>
    </div>
</body>
</html>
	`, toName, resetURL, resetURL)

	return e.sendEmail(toEmail, toName, subject, body)
}

func (e *EmailService) SendVerificationApproval(toEmail, toName string) error {
	subject := "Your VietTick Verification Has Been Approved!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verification Approved</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #1DA1F2;">Congratulations! You're Verified!</h1>
        <p>Hi %s,</p>
        <p>Great news! Your identity verification has been approved and you now have a verified blue checkmark on your VietTick profile.</p>
        <p>Your verified status helps other users know that your account is authentic and trustworthy.</p>
        <div style="text-align: center; margin: 30px 0;">
            <div style="background-color: #f0f8ff; padding: 20px; border-radius: 10px;">
                <h2 style="color: #1DA1F2; margin: 0;">✓ Verified Account</h2>
                <p style="margin: 10px 0 0 0; color: #666;">You now have the blue checkmark!</p>
            </div>
        </div>
        <p>Thank you for being a valued member of the VietTick community!</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="font-size: 12px; color: #666;">This is an automated message, please do not reply to this email.</p>
    </div>
</body>
</html>
	`, toName)

	return e.sendEmail(toEmail, toName, subject, body)
}

func (e *EmailService) sendEmail(toEmail, toName, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.config.FromEmail, e.config.FromName))
	m.SetHeader("To", m.FormatAddress(toEmail, toName))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.config.SMTPHost, e.config.SMTPPort, e.config.SMTPUser, e.config.SMTPPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
