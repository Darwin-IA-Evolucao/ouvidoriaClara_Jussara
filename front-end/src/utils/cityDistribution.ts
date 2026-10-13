import type { OcorrenciaPorCidade } from '../types'

export function normalizeCidades(data: OcorrenciaPorCidade[] | null | undefined): OcorrenciaPorCidade[] {
  const map = new Map<string, OcorrenciaPorCidade>()
  for (const c of data ?? []) {
    const nome = (c.cidade || '').trim() || 'Sem cidade'
    const key = nome.toLowerCase()
    const prev = map.get(key)
    if (!prev) {
      map.set(key, { cidade: nome, categoria: { ...(c.categoria || {}) }, total: c.total || 0 })
    } else {
      for (const [cat, n] of Object.entries(c.categoria || {})) {
        prev.categoria[cat] = (prev.categoria[cat] || 0) + n
      }
      prev.total += c.total || 0
    }
  }
  return Array.from(map.values()).sort((a, b) => b.total - a.total)
}
