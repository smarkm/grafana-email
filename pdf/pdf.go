package pdf

import (
	"bytes"

	"github.com/go-pdf/fpdf"
)

func InitPDF() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetTitle("Hello----", true)
	pdf.SetFont("Arial", "B", 16)

	pdf.Ln(-1)
	return pdf
}

func InsertImageInNewPage(imageName string, pdf *fpdf.Fpdf, image []byte) {
	pdf.AddPage()
	InsertImage(imageName, pdf, image)
}

func InsertImage(imageName string, pdf *fpdf.Fpdf, image []byte) {
	pdf.RegisterImageReader(imageName, "png", bytes.NewReader(image))
	imageWidth, _ := pdf.GetPageSize()
	imageHeight := 0.0 // Set to 0 for automatic height calculation based on the aspect ratio
	imageX := 20.0
	imageY := 50.0

	// Insert the image into the PDF
	pdf.ImageOptions(imageName, imageX, imageY, imageWidth-40, imageHeight, false, fpdf.ImageOptions{}, 0, "")
}
