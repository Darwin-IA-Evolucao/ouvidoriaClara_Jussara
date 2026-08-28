import { apiGet, apiPostForm } from '../utils/api'
import type { Contato, ImportContatosResult } from '../types'

export async function getAllContatos(limit: number = 0, offset: number = 0): Promise<{ contatos: Contato[]; total: number }> {
  const res = await apiGet<{ contatos: Contato[]; total: number } | Contato[] | null>(`/contatos?limit=${limit}&offset=${offset}`)
  if (Array.isArray(res)) return { contatos: res, total: res.length }
  if (res && Array.isArray(res.contatos)) return { contatos: res.contatos, total: res.total ?? res.contatos.length }
  return { contatos: [], total: 0 }
}

export async function importContatos(file: File, campanha?: string): Promise<ImportContatosResult> {
  const formData = new FormData()
  formData.append('file', file)
  if (campanha) formData.append('campanha', campanha)
  return apiPostForm<ImportContatosResult>('/contatos/import', formData)
}
