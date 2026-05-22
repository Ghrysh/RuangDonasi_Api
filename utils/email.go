package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

func SendResetPasswordEmail(toEmail string, resetToken string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := "465"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	frontendURL := os.Getenv("FRONTEND_URL")

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken)

	subject := "Subject: Reset Kata Sandi Anda - Ruang Donasi\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<h2>Permintaan Reset Kata Sandi</h2>
		<p>Kami menerima permintaan untuk mereset kata sandi akun Ruang Donasi Anda.</p>
		<p>Silakan klik tautan di bawah ini untuk membuat kata sandi baru. Tautan ini hanya berlaku selama 15 menit.</p>
		<a href="%s" style="display:inline-block; padding:10px 20px; background-color:#16a34a; color:#fff; text-decoration:none; border-radius:5px;">Reset Kata Sandi</a>
		<br><br>
		<p>Jika Anda tidak merasa meminta reset kata sandi, abaikan saja email ini.</p>
	`, resetLink)

	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		return fmt.Errorf("gagal koneksi TLS ke SMTP: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("gagal membuat SMTP client: %v", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("autentikasi email gagal (cek App Password): %v", err)
	}

	if err = client.Mail(smtpUser); err != nil {
		return err
	}
	if err = client.Rcpt(toEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}