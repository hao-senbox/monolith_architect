package order

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func formatVND(val float64) string {
	n := int64(math.Round(val))
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := strconv.FormatInt(n, 10)

	var out []byte
	c := 0
	for i := len(s) - 1; i >= 0; i-- {
		out = append(out, s[i])
		c++
		if c%3 == 0 && i != 0 {
			out = append(out, '.')
		}
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return sign + string(out) + " ₫"
}

func paymentLabel(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	switch t {
	case "cod":
		return "Thanh toán khi nhận hàng (COD)"
	case "online", "vnpay":
		return "Thanh toán trực tuyến"
	default:
		return strings.ToUpper(t)
	}
}

func BuildOrderEmailHTML(order Order, brandName string) string {
	// Tính tiền
	var subtotal float64
	for _, it := range order.OrderItems {
		if it.TotalPrice > 0 {
			subtotal += it.TotalPrice
		} else {
			subtotal += float64(it.Quantity) * it.Price
		}
	}

	discountAmount := 0.0
	if order.Discount != nil {
		discountAmount = subtotal - order.TotalPrice
		if discountAmount < 0 {
			discountAmount = 0
		}
	}

	grandTotal := order.TotalPrice
	if grandTotal == 0 {
		grandTotal = subtotal - discountAmount
	}

	// Dòng sản phẩm
	var itemsHTML strings.Builder
	for _, it := range order.OrderItems {
		lineTotal := it.TotalPrice
		if lineTotal == 0 {
			lineTotal = float64(it.Quantity) * it.Price
		}
		img := it.ProductImage
		if img == "" {
			img = "https://via.placeholder.com/64x64?text=IMG"
		}

		itemsHTML.WriteString(fmt.Sprintf(`
			<tr>
				<td style="padding:12px 0; border-bottom:1px solid #eee;">
					<table role="presentation" width="100%%" cellspacing="0" cellpadding="0">
						<tr>
							<td width="64" valign="top">
								<img src="%s" width="64" height="64" style="display:block; border-radius:8px; object-fit:cover;" alt="">
							</td>
							<td style="padding-left:12px;">
								<div style="font-weight:600; color:#111">%s</div>
								<div style="font-size:12px; color:#666; margin-top:2px;">Size: %s &nbsp;•&nbsp; SL: %d</div>
							</td>
							<td align="right" valign="top" style="white-space:nowrap;">
								<div style="font-size:13px; color:#666;">%s</div>
								<div style="font-weight:700; color:#111;">%s</div>
							</td>
						</tr>
					</table>
				</td>
			</tr>
		`, img, htmlEscape(it.ProductName), htmlEscape(it.Size), it.Quantity, formatVND(it.Price), formatVND(lineTotal)))
	}

	// Template KHÔNG còn logo → chỉ hiển thị brand + mã đơn
	// Trật tự tham số bên dưới (theo thứ tự xuất hiện %s):
	// 1 brandName  | 2 order.OrderCode
	// 3 order.ShippingAddress.Name | 4 brandName
	// 5 order.ShippingAddress.Name | 6 order.ShippingAddress.Address | 7 order.ShippingAddress.Phone
	// 8 paymentLabel(order.Type)
	// 9 itemsHTML
	// 10 formatVND(subtotal)
	// 11 discountRowHTML(discountAmount)
	// 12 formatVND(grandTotal)
	return fmt.Sprintf(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>%s - Đơn hàng #%s</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body style="margin:0; background:#f6f7f9; font-family:system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;">
  <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#f6f7f9; padding:24px 0;">
    <tr>
      <td align="center">
        <table role="presentation" width="640" cellspacing="0" cellpadding="0" style="max-width:640px; width:100%%; background:#ffffff; border-radius:16px; overflow:hidden; box-shadow:0 4px 24px rgba(0,0,0,0.06);">
          <tr>
            <td style="background:linear-gradient(135deg,#0ea5e9,#6366f1); padding:20px;" align="left">
              <table role="presentation" width="100%%">
                <tr>
                  <td style="color:#eaf2ff; font-weight:800; font-size:18px;">%s</td>
                  <td align="right" style="color:#eaf2ff; font-size:12px;">Mã đơn: <strong>#%s</strong></td>
                </tr>
              </table>
            </td>
          </tr>

          <tr>
            <td style="padding:24px 24px 8px 24px;">
              <div style="font-size:18px; font-weight:700; color:#111; margin-bottom:6px;">Cảm ơn bạn đã đặt hàng 👟</div>
              <div style="font-size:14px; color:#444;">Xin chào %s, đơn hàng của bạn đã được tạo thành công tại <strong>%s</strong>.</div>
            </td>
          </tr>

          <tr>
            <td style="padding:0 24px 8px 24px;">
              <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background:#f8fafc; border:1px solid #eef2f7; border-radius:12px;">
                <tr>
                  <td style="padding:16px;">
                    <table role="presentation" width="100%%">
                      <tr>
                        <td>
                          <div style="font-size:12px; color:#64748b;">ĐỊA CHỈ GIAO HÀNG</div>
                          <div style="font-weight:600; color:#111; margin-top:4px;">%s</div>
                          <div style="font-size:13px; color:#444; margin-top:2px;">%s</div>
                          <div style="font-size:13px; color:#444; margin-top:2px;">SĐT: %s</div>
                        </td>
                        <td align="right">
                          <div style="font-size:12px; color:#64748b;">PHƯƠNG THỨC THANH TOÁN</div>
                          <div style="font-weight:600; color:#111; margin-top:4px;">%s</div>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <tr><td style="height:8px;"></td></tr>

          <tr>
            <td style="padding:0 24px 8px 24px;">
              <div style="font-size:14px; font-weight:700; color:#111; margin:8px 0;">Sản phẩm</div>
              <table role="presentation" width="100%%" cellspacing="0" cellpadding="0">
                %s
              </table>
            </td>
          </tr>

          <tr>
            <td style="padding:8px 24px 4px 24px;">
              <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="border-top:1px dashed #e5e7eb;">
                <tr>
                  <td style="padding-top:12px; color:#555;">Tạm tính</td>
                  <td align="right" style="padding-top:12px; font-weight:600; color:#111;">%s</td>
                </tr>
                %s
                <tr>
                  <td style="padding-top:8px; font-size:16px; font-weight:800; color:#111;">TỔNG CỘNG</td>
                  <td align="right" style="padding-top:8px; font-size:16px; font-weight:800; color:#111;">%s</td>
                </tr>
              </table>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
		brandName, order.OrderCode,
		brandName, order.OrderCode,
		htmlEscape(order.ShippingAddress.Name), brandName,
		htmlEscape(order.ShippingAddress.Name),
		htmlEscape(order.ShippingAddress.Address),
		htmlEscape(order.ShippingAddress.Phone),
		paymentLabel(order.Type),
		itemsHTML.String(),
		formatVND(subtotal),
		discountRowHTML(discountAmount),
		formatVND(grandTotal),
	)
}

func discountRowHTML(discount float64) string {
	if discount <= 0 {
		return ""
	}
	return fmt.Sprintf(`
		<tr>
		  <td style="padding-top:6px; color:#16a34a;">Giảm giá</td>
		  <td align="right" style="padding-top:6px; font-weight:600; color:#16a34a;">- %s</td>
		</tr>`, formatVND(discount))
}

func htmlEscape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;",
	)
	return r.Replace(s)
}
