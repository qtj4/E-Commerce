package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailConfig struct {
	SenderEmail    string
	SenderPassword string
	SMTPHost       string
	SMTPPort       string
}

type OrderEmailData struct {
	OrderID     string
	UserEmail   string
	Items       []OrderItemData
	TotalAmount float64
	OrderStatus string
}

type OrderItemData struct {
	ProductName string
	Quantity    int
	Price       float64
	Subtotal    float64
}

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Order Confirmation</title>
</head>
<body>
    <h2>Order Confirmation</h2>
    <p>Thank you for your order!</p>
    <p><strong>Order ID:</strong> {{.OrderID}}</p>
    <h3>Order Details:</h3>
    <table border="1" cellpadding="5" cellspacing="0">
        <tr>
            <th>Product</th>
            <th>Quantity</th>
            <th>Price</th>
            <th>Subtotal</th>
        </tr>
        {{range .Items}}
        <tr>
            <td>{{.ProductName}}</td>
            <td>{{.Quantity}}</td>
            <td>${{printf "%.2f" .Price}}</td>
            <td>${{printf "%.2f" .Subtotal}}</td>
        </tr>
        {{end}}
        <tr>
            <td colspan="3" align="right"><strong>Total:</strong></td>
            <td><strong>${{printf "%.2f" .TotalAmount}}</strong></td>
        </tr>
    </table>
    <p>Order Status: {{.OrderStatus}}</p>
    <p>Thank you for shopping with us!</p>
</body>
</html>
`

func SendOrderConfirmationEmail(config EmailConfig, data OrderEmailData) error {
	// Parse template
	tmpl, err := template.New("orderEmail").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	// Execute template with data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	// Set up email headers
	headers := make(map[string]string)
	headers["From"] = config.SenderEmail
	headers["To"] = data.UserEmail
	headers["Subject"] = "Order Confirmation - Order #" + data.OrderID
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Compose message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body.String()

	// Authentication
	auth := smtp.PlainAuth("", config.SenderEmail, config.SenderPassword, config.SMTPHost)

	// Send email
	err = smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort,
		auth,
		config.SenderEmail,
		[]string{data.UserEmail},
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
