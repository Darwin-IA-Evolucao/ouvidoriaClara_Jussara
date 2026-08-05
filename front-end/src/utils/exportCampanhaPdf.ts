import { jsPDF } from 'jspdf'
import autoTable from 'jspdf-autotable'
import { formatPhoneDisplay } from './phone'

export type ContatoPdfRow = {
  nome: string
  telefone: string
  campanha?: string
}

function dataGeracao(): string {
  return new Date().toLocaleString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function safeFileName(s: string): string {
  return s.replace(/[^\w\-]+/g, '_').replace(/_+/g, '_').slice(0, 60)
}

export function exportCampanhaPdf(titulo: string, contatos: ContatoPdfRow[], comCampanha = false) {
  const doc = new jsPDF()
  const geradoEm = dataGeracao()

  doc.setFontSize(16)
  doc.setFont('helvetica', 'bold')
  doc.text(titulo, 14, 18)

  doc.setFontSize(10)
  doc.setFont('helvetica', 'normal')
  doc.setTextColor(100)
  doc.text(`Gerado em: ${geradoEm}`, 14, 26)
  doc.text(`Total de contatos: ${contatos.length}`, 14, 32)
  doc.setTextColor(0)

  const head = comCampanha
    ? [['#', 'Tipo Liderança', 'Nome', 'Telefone']]
    : [['#', 'Nome', 'Telefone']]

  const body = contatos.map((c, i) => {
    const nome = c.nome || 'Sem nome'
    const telefone = formatPhoneDisplay(c.telefone) || c.telefone
    if (comCampanha) return [String(i + 1), c.campanha || '-', nome, telefone]
    return [String(i + 1), nome, telefone]
  })

  autoTable(doc, {
    startY: 38,
    head,
    body,
    styles: { fontSize: 9, cellPadding: 3 },
    headStyles: { fillColor: [65, 102, 156], textColor: 255, fontStyle: 'bold' },
    alternateRowStyles: { fillColor: [245, 247, 250] },
    columnStyles: comCampanha
      ? { 0: { cellWidth: 22 }, 1: { cellWidth: 40 }, 2: { cellWidth: 50 }, 3: { cellWidth: 45 } }
      : { 0: { cellWidth: 22 }, 1: { cellWidth: 75 }, 2: { cellWidth: 55 } },
  })

  const file = `lideranca_${safeFileName(titulo)}_${new Date().toISOString().slice(0, 10)}.pdf`
  doc.save(file)
}
