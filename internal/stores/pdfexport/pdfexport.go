package pdfexport

import (
	"brotchaschte2/internal/database"
	query "brotchaschte2/internal/database/generated"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io"
	"strings"

	"codeberg.org/go-pdf/fpdf"
)

//go:embed "fonts"
var fontFiles embed.FS

type PDFExportStore struct {
	db      *sql.DB
	queries *query.Queries
}

func New() PDFExportStore {
	db := database.New()
	return PDFExportStore{
		db:      db.DB(),
		queries: db.Queries(),
	}
}

func (p *PDFExportStore) DeliveryRuns(ctx context.Context, w io.Writer) error {
	orders, err := p.queries.GetOrderedDeliveries(ctx)
	if err != nil {
		return fmt.Errorf("get ordered orders: %w", err)
	}

	header := [...]string{"x", "300g", "500g", "1kg", "Zeit", "Bestellung", "Nachname", "Vorname", "Adresse", "Bemerkung"}
	column_widths := [...]float64{6, 10, 10, 10, 13, 23, 30, 25, 80, 80}
	wSum := 0.0
	for _, v := range column_widths {
		wSum += v
	}
	left := (297 - wSum) / 2

	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(5, 10, 5)

	regular, err := fontFiles.ReadFile("fonts/GoogleSansCode-Regular.ttf")
	if err != nil {
		fmt.Println(err)
	}
	pdf.AddUTF8FontFromBytes("GoogleSansCode", "", regular)
	bold, err := fontFiles.ReadFile("fonts/GoogleSansCode-Bold.ttf")
	if err != nil {
		fmt.Println(err)
	}
	pdf.AddUTF8FontFromBytes("GoogleSansCode", "B", bold)

	delivery := struct {
		valid bool
		id    int64
		until string
		label string
	}{valid: true}

	pdf.SetFooterFunc(func() {
		pdf.SetFont("GoogleSansCode", "B", 9)
		pdf.SetTextColor(160, 160, 160)
		pdf.SetXY(left, -10)
		if delivery.valid {
			pdf.CellFormat(20, 7, fmt.Sprintf("Liefernr: %d", delivery.id), "", 0, "", false, 0, "")
		} else {
			pdf.SetTextColor(255, 80, 80)
			pdf.CellFormat(20, 7, "Nicht eingeteilte Bestellungen", "", 0, "", false, 0, "")
		}
	})

	for i, order := range orders {
		if order.DeliveryID.Int64 != delivery.id || order.DeliveryID.Valid != delivery.valid {
			pdf.AddPage()

			delivery.valid = order.DeliveryID.Valid
			delivery.id = order.DeliveryID.Int64
			delivery.until = order.DeliveryUntil.String
			delivery.label = order.DeliveryLabel.String

			pdf.SetX(left)
			pdf.SetFont("GoogleSansCode", "B", 9)
			pdf.SetTextColor(30, 30, 30)
			if delivery.valid {
				pdf.CellFormat(60, 7, fmt.Sprintf("Liefernr: %d", delivery.id), "", 0, "", false, 0, "")
				pdf.CellFormat(60, 7, fmt.Sprintf("Lieferzeit: %s", delivery.until), "", 0, "", false, 0, "")
				pdf.CellFormat(60, 7, fmt.Sprintf("Notiz: %s", delivery.label), "", 0, "", false, 0, "")
			} else {
				pdf.SetTextColor(255, 80, 80)
				pdf.CellFormat(60, 7, "Nicht eingeteilte Bestellungen", "", 0, "", false, 0, "")
			}
			pdf.Ln(-1)
			pdf.Ln(-1)

			pdf.SetFont("GoogleSansCode", "B", 9)
			pdf.SetX(left)
			pdf.SetFillColor(80, 80, 100)
			pdf.SetTextColor(255, 255, 255)
			for i, col := range header {
				pdf.CellFormat(column_widths[i], 7, col, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)

		}
		pdf.SetX(left)
		pdf.SetFont("GoogleSansCode", "", 9)
		pdf.SetFillColor(220, 220, 230)
		pdf.SetTextColor(30, 30, 30)
		fill := i%2 == 0

		// address := pkg.MaybeFromSQL[string](order.GeocodedAddress)

		pdf.CellFormat(column_widths[0], 6, "", "BL", 0, "", fill, 0, "")
		pdf.CellFormat(column_widths[1], 6, formatAmount(order.Amountsizes), "BL", 0, "", fill, 0, "")
		pdf.CellFormat(column_widths[2], 6, formatAmount(order.Amountsizem), "BL", 0, "", fill, 0, "")
		pdf.CellFormat(column_widths[3], 6, formatAmount(order.Amountsizel), "BL", 0, "", fill, 0, "")
		pdf.CellFormat(column_widths[4], 6, order.Deliverymethod, "BL", 0, "R", fill, 0, "")
		pdf.CellFormat(column_widths[5], 6, order.Ordermethod.String, "BL", 0, "", fill, 0, "")
		pdf.CellFormat(column_widths[6], 6, order.Lastname, "BL", 0, "", fill, 0, "")
		pdf.CellFormat(column_widths[7], 6, order.Firstname, "BL", 0, "", fill, 0, "")
		pdf.CellFormat(column_widths[8], 6, fmt.Sprintf("%s %s %s", order.Address, order.Postalcode, order.City), "BL", 0, "", fill, 0, "")
		pdf.MultiCell(column_widths[9], 6, strings.TrimRight(order.Remarks.String, "\r\n"), "BLR", "", fill)

		//pdf.Ln(-1)
	}

	if err := pdf.Output(w); err != nil {
		return fmt.Errorf("creating pdf: %w", err)
	}
	return nil
}

func formatAmount(amount int64) string {
	if amount == 0 {
		return "-"
	}
	return fmt.Sprintf("%d", amount)
}
