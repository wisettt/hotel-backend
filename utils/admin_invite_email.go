package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// SendAdminInviteEmail sends an account setup invite email for admins.
func SendAdminInviteEmail(recipientEmail, inviteLink, name, role string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("SMTP_FROM_NAME")

	if smtpUser == "" || smtpPass == "" || smtpHost == "" || smtpPort == "" {
		log.Printf("[MOCK EMAIL] invite to:%s role:%s link:%s", recipientEmail, role, inviteLink)
		return nil
	}

	safe := func(s string) string {
		return strings.ReplaceAll(strings.TrimSpace(s), "\r\n", " ")
	}

	name = safe(name)
	role = safe(role)
	inviteLink = safe(inviteLink)

	if !(strings.HasPrefix(inviteLink, "http://") || strings.HasPrefix(inviteLink, "https://")) {
		inviteLink = "https://" + strings.TrimLeft(inviteLink, "/")
	}

	from := fmt.Sprintf("%s <%s>", fromName, smtpUser)
	to := []string{recipientEmail}
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	subject := "You're invited to Horizon Hotel System"
	boundary := "----=_INVITE_EMAIL_BOUNDARY"

	plainBody := fmt.Sprintf(
		"Hi %s,\n\n"+
		"You have been invited to join Horizon Hotel as a %s.\n"+
		"Please set your password using the link below:\n%s\n\n"+
		"If you did not expect this invitation, you can ignore this email.\n",
		name, role, inviteLink,
	)

	htmlBody := fmt.Sprintf(`<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Invitation</title>
<style>
body { background:#f5f7fb; font-family:Arial, Helvetica, sans-serif; color:#222; }
.container { max-width:640px; margin:20px auto; }
.card { background:#fff; border:1px solid #e6eef6; padding:24px; border-radius:8px; }
.btn { display:inline-block; padding:12px 20px; background:#0b74ff; color:#fff; text-decoration:none; border-radius:6px; margin-top:16px; }
</style>
</head>
<body>
<div class="container">
  <div class="card">
    <h2>You're invited</h2>
    <p>Hi %s,</p>
    <p>You have been invited to join Horizon Hotel as a <strong>%s</strong>.</p>
    <p>Click the button below to set your password.</p>
    <a class="btn" href="%s" target="_blank">Set up my account</a>
    <p>If you did not expect this invitation, you can ignore this email.</p>
  </div>
</div>
</body>
</html>`,
		name, role, inviteLink,
	)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", recipientEmail))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundary))

	sb.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	sb.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	sb.WriteString(plainBody + "\r\n")

	sb.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	sb.WriteString("Content-Type: text/html; charset=utf-8\r\n\r\n")
	sb.WriteString(htmlBody + "\r\n")

	sb.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	if err := smtp.SendMail(addr, auth, smtpUser, to, []byte(sb.String())); err != nil {
		log.Printf("Failed to send invite email to %s: %v", recipientEmail, err)
		return err
	}

	log.Printf("Invite email sent to %s", recipientEmail)
	return nil
}
