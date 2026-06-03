package services

import (
	"back-end/models"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var meses = []string{
	"janeiro", "fevereiro", "março", "abril", "maio", "junho",
	"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
}

func logoPath() string {
	if path := os.Getenv("LOGO_PATH"); path != "" {
		return path
	}
	return "/root/arruda_Clara/img/logo.png"
}

func novoCabecalho(titulo string) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(20, 20, 20)

	logoFile := logoPath()
	if _, err := os.Stat(logoFile); os.IsNotExist(err) {
		log.Printf("Atenção: Logo não encontrado em %s", logoFile)
	} else {
		pdf.ImageOptions(logoFile, 20, 15, 170, 30, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	}

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(60, 50)
	pdf.Cell(90, 10, tr("GABINETE DO VEREADOR MARCOS ARRUDA"))
	pdf.SetXY(85, 60)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, tr(titulo))
	pdf.Ln(30)
	return pdf
}

func dataExtenso() string {
	now := time.Now()
	return fmt.Sprintf("São Roque, %d de %s de %d", now.Day(), meses[now.Month()-1], now.Year())
}

func GeneratePDF(requestData models.Inquerito) error {
	pdf := novoCabecalho("INDICAÇÃO N.º ___ / 2025")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Excelentíssimo Senhor Presidente,"))
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr("Este Vereador, submentendo este documento ao Chefe do Poder Executivo, diretamente ou através de departamento ou divisão competente,"))
	pdf.SetFont("Arial", "B", 11)
	pdf.Write(10, tr(" INDICA"))
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr(" ao Senhor Prefeito Municipal, que "))
	pdf.Write(10, tr(requestData.Reclamacao+"."))
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Nestes termos,"))
	pdf.Ln(5)
	pdf.Cell(0, 10, tr("Aguarda deferimento."))
	pdf.Ln(30)

	pdf.Cell(0, 10, tr(dataExtenso()))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 10, tr("VEREADOR MARCOS ARRUDA"), "", 1, "C", false, 0, "")

	return pdf.OutputFileAndClose("indicacao.pdf")
}

func GenerateRequerimento(requestData models.Inquerito) error {
	if !strings.Contains(requestData.Reclamacao, "$$") {
		requestData.Reclamacao = "$$" + requestData.Reclamacao
	}
	args := strings.Split(requestData.Reclamacao, "$$")

	pdf := novoCabecalho("REQUERIMENTO N.º ___ / 2025")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Senhor Presidente,"))
	pdf.Ln(15)
	pdf.Cell(0, 10, tr("Senhores Vereadores,"))
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr("Este Vereador, submentendo este documento ao Chefe do Poder Executivo, diretamente ou através de departamento ou divisão competente,"))
	pdf.SetFont("Arial", "B", 11)
	pdf.Write(10, tr(" REQUER"))
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr(" que, seja oficiado ao Prefeito Municipal para que nos sejam prestadas as seguintes informações:"))
	pdf.Ln(15)

	requests := strings.Split(args[1], "|")
	for i, request := range requests {
		pdf.SetFont("Arial", "", 11)
		pdf.Write(10, tr(fmt.Sprintf("%d. %s", i+1, request)))
		pdf.Ln(15)
	}

	pdf.Cell(0, 10, tr(dataExtenso()))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 10, tr("VEREADOR MARCOS ARRUDA"), "", 1, "C", false, 0, "")

	return pdf.OutputFileAndClose("requerimento.pdf")
}

func GenerateOficio(requestData models.Inquerito) error {
	now := time.Now()
	pdf := novoCabecalho("OFÍCIO N.º ___ / " + now.Format("2006"))
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Senhor Presidente,"))
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr("Este Vereador, submentendo este documento ao Chefe do Poder Executivo, diretamente ou através de departamento ou divisão competente,"))
	pdf.SetFont("Arial", "B", 11)
	pdf.Write(10, tr(" OFICIA"))
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr(" ao Senhor Prefeito Municipal, que "))
	pdf.Write(10, tr(requestData.Reclamacao+"."))
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Nestes termos,"))
	pdf.Ln(5)
	pdf.Cell(0, 10, tr("Aguarda deferimento."))
	pdf.Ln(30)

	pdf.Cell(0, 10, tr(dataExtenso()))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 10, tr("VEREADOR MARCOS ARRUDA"), "", 1, "C", false, 0, "")

	return pdf.OutputFileAndClose("oficio.pdf")
}
