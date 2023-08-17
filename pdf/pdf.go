package pdf

import (
	"bytes"

	"github.com/go-pdf/fpdf"
	"smark.freecoop.net/grafana-email/config"
)

var pdfWidthMargin = 20.0

func InitPDF() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pageWidth, _ := pdf.GetPageSize()

	// Calculate the center position for the text
	text := "Text in the Center of the Page"
	textWidth := pdf.GetStringWidth(text)
	x := (pageWidth - textWidth) / 2
	switch config.Instance.PdfTitleAlign {
	case "left":
		x = pdfWidthMargin
	case "right":
		x = pageWidth - textWidth - pdfWidthMargin
	}
	// Draw the centered text
	pdf.Text(x, 40, text)
	pdf.Ln(-1)
	return pdf
}

func InsertImageInNewPage(imageName string, pdf *fpdf.Fpdf, image []byte) {
	pdf.AddPage()
	InsertImage(imageName, pdf, image, 50)
}

func InsertImage(imageName string, pdf *fpdf.Fpdf, image []byte, yStart float64) {
	pdf.RegisterImageReader(imageName, "png", bytes.NewReader(image))
	imageWidth, _ := pdf.GetPageSize()
	imageHeight := 0.0 // Set to 0 for automatic height calculation based on the aspect ratio
	imageY := yStart
	// Insert the image into the PDF
	pdf.ImageOptions(imageName, pdfWidthMargin, imageY, imageWidth-pdfWidthMargin*2, imageHeight, false, fpdf.ImageOptions{}, 0, "")
}
