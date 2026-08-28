import { apiGet } from '../utils/api'
import type { Stat } from '../types'

export async function getStats(inicio?: string, fim?: string): Promise<Stat> {
  const params = new URLSearchParams()
  if (inicio) params.set('inicio', inicio)
  if (fim) params.set('fim', fim)
  const qs = params.toString()
  const res = await apiGet<Partial<Stat> | null>(`/stats${qs ? `?${qs}` : ''}`)
  return {
    numPessoas: res?.numPessoas ?? 0,
    numReclamacoes: res?.numReclamacoes ?? 0,
    percIndicacao: res?.percIndicacao ?? 0,
    percRequerimento: res?.percRequerimento ?? 0,
    regioes: res?.regioes ?? [],
    categorias: res?.categorias ?? [],
    tipos: res?.tipos ?? [],
    indicacoesAprovadas: res?.indicacoesAprovadas ?? 0,
    totalIndicacoes: res?.totalIndicacoes ?? 0,
    requerimentosAprovados: res?.requerimentosAprovados ?? 0,
    totalRequerimentos: res?.totalRequerimentos ?? 0,
  }
}
