package utils

import (
	"fmt"
	"os"

	"github.com/wneessen/go-mail"
)

func SendOTP(otp string, to string) error {

	email_addr := os.Getenv("GMAIL_ADDRESS")
	email_pass := os.Getenv("GMAIL_APP_PASS")

	msg := mail.NewMsg()
	msg.From(email_addr)
	msg.To(to)
	msg.Subject("Unibox OTP")

	template := fmt.Sprintf(`
		<div style="font-family: Arial, Helvetica, sans-serif; background-color:#f6f7fb; padding:40px 0;">
		<div style="max-width:520px; margin:0 auto; background:#ffffff; border-radius:12px; padding:32px; box-shadow:0 4px 20px rgba(0,0,0,0.08);">

			<h2 style="margin:0; font-size:22px; color:#111827;">Unibox</h2>

			<p style="margin-top:10px; font-size:14px; color:#6b7280;">Secure One-Time Password Verification</p>

			<hr style="border:none; border-top:1px solid #e5e7eb; margin:20px 0;" />

			<p style="font-size:16px; color:#374151; margin-bottom:10px;">Hello,</p>

			<p style="font-size:15px; color:#4b5563; line-height:1.6;">
			Use the following OTP to complete your verification process. 
			Please do not share this code with anyone.
			</p>

			<div style="text-align:center; margin:30px 0;">
			<div style="display:inline-block; padding:14px 28px; font-size:28px; letter-spacing:6px; background:#f3f4f6; border-radius:8px; font-weight:bold; color:#111827;">
				%s
			</div>
			</div>

			<p style="font-size:14px; color:#6b7280; text-align:center;">This OTP will expire in <b>7 minutes</b>.</p>
			<p style="font-size:12px; color:#9ca3af; text-align:center; margin-top:30px;">If you didn’t request this, you can safely ignore this email.</p>

		</div>
		</div>
	`, otp)

	msg.SetBodyString(mail.TypeTextHTML, template)

	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(email_addr),
		mail.WithPassword(email_pass),
	)

	if err != nil {
		return err
	}

	return client.DialAndSend(msg)
}
