package mail

import "fmt"

// ResetPasswordBody generates the HTML body for a password reset email.
func ResetPasswordBody(resetLink string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Reset Your Password</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 0;
			padding: 0;
			background-color: #f7f7f7;
			color: #333;
			text-align: center;
		}
		.container {
			width: 100%%;
			max-width: 600px;
			margin: 0 auto;
			background-color: #fff;
			padding: 20px;
			border-radius: 8px;
			box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		}
		h1 { color: #18181B; }
		.reset-button {
			display: inline-block;
			background-color: #18181B;
			color: white;
			padding: 12px 20px;
			text-decoration: none;
			border-radius: 5px;
			font-weight: bold;
			margin-top: 10px;
		}
		.footer {
			margin-top: 20px;
			font-size: 12px;
			color: #888;
		}
		p { font-size: 14px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Reset Password Akun</h1>
		<p>Kami menerima permintaan untuk mereset kata sandi akun Anda. Silakan gunakan link berikut untuk melanjutkan proses reset password:</p>
		<a href="%s" style="color:white!important" class="reset-button">Reset Password</a>
		<p>Link ini hanya berlaku selama 1 jam. Demi keamanan informasi Anda, mohon jangan membagikan link ini kepada siapa pun.</p>
		<p>Apabila Anda tidak merasa melakukan permintaan ini, silakan abaikan email ini.</p>
		<div class="footer">
			<p><strong>Hormat Kami,</strong></p>
			<p><strong>Visea Hive</strong></p>
		</div>
	</div>
</body>
</html>`, resetLink)
}

// NotificationBody generates a generic notification email body.
func NotificationBody(name, title, message string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s - Visea Hive</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			line-height: 1.6;
			margin: 0;
			padding: 20px;
			background-color: #ffffff;
			color: #333;
		}
		.email-content {
			max-width: 600px;
			margin: 0 auto;
		}
		h1 {
			color: #8e44ad;
			font-size: 24px;
			margin-bottom: 20px;
		}
		p {
			margin-bottom: 15px;
			font-size: 14px;
		}
		.greeting { font-weight: bold; margin-bottom: 20px; }
		.signature { margin-top: 30px; }
	</style>
</head>
<body>
	<div class="email-content">
		<h1>%s</h1>
		<div class="greeting">
			Kepada Yth:<br>
			%s<br>
			Di Tempat
		</div>
		<p>%s</p>
		<div class="signature">
			<p><strong>Hormat Kami,</strong></p>
			<p><strong>Visea Hive</strong></p>
		</div>
	</div>
</body>
</html>`, title, title, name, message)
}

// VerifyEmailBody generates the HTML body for an email verification email.
func VerifyEmailBody(name, verifyLink string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Verifikasi Email Anda</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 0;
			padding: 0;
			background-color: #f7f7f7;
			color: #333;
			text-align: center;
		}
		.container {
			width: 100%%;
			max-width: 600px;
			margin: 0 auto;
			background-color: #fff;
			padding: 20px;
			border-radius: 8px;
			box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		}
		h1 { color: #18181B; }
		.verify-button {
			display: inline-block;
			background-color: #18181B;
			color: white;
			padding: 12px 20px;
			text-decoration: none;
			border-radius: 5px;
			font-weight: bold;
			margin-top: 10px;
		}
		.footer {
			margin-top: 20px;
			font-size: 12px;
			color: #888;
		}
		p { font-size: 14px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Verifikasi Email</h1>
		<p>Halo %s,</p>
		<p>Terima kasih telah mendaftar di Visea Hive. Silakan klik tombol di bawah ini untuk memverifikasi alamat email Anda:</p>
		<a href="%s" style="color:white!important" class="verify-button">Verifikasi Email</a>
		<p>Link ini berlaku selama 24 jam.</p>
		<p>Jika Anda tidak membuat akun ini, Anda dapat mengabaikan email ini.</p>
		<div class="footer">
			<p><strong>Hormat Kami,</strong></p>
			<p><strong>Visea Hive</strong></p>
		</div>
	</div>
</body>
</html>`, name, verifyLink)
}
