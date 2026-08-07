import * as XLSX from 'xlsx'
import { formatPhoneDisplay } from './phone'
import type { ContatoPdfRow } from './exportCampanhaPdf'

function safeFileName(s: string): string {
  return s.replace(/[^\w\-]+/g, '_').replace(/_+/g, '_').slice(0, 60)
}

export function exportCampanhaXlsx(titulo: string, contatos: ContatoPdfRow[], comCampanha = false) {
  const rows = contatos.map((c, i) => {
    const nome = c.nome || 'Sem nome'
    const telefone = formatPhoneDisplay(c.telefone) || c.telefone
    if (comCampanha) {
      return {
        '#': i + 1,
        'Tipo Liderança': c.campanha || '-',
        Nome: nome,
        Telefone: telefone,
      }
    }
    return {
      '#': i + 1,
      Nome: nome,
      Telefone: telefone,
    }
  })

  const ws = XLSX.utils.json_to_sheet(rows)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Contatos')
  XLSX.writeFile(wb, `lideranca_${safeFileName(titulo)}_${new Date().toISOString().slice(0, 10)}.xlsx`)
}
