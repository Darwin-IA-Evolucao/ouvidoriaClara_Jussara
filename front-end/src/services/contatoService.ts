import { apiGet, apiPost, apiPostForm } from '../utils/api'
import type { Contato, ImportContatosResult } from '../types'

export async function getAllContatos(): Promise<Contato[]> {
  return apiGet<Contato[]>('/contatos')
}

export async function getContatoByTelefone(telefone: string): Promise<Contato> {
  return apiGet<Contato>(`/contatos/${telefone}`)
}

export async function createContato(contato: {
  telefone: string
  nome: string
  campanha: string
}): Promise<{ message: string }> {
  return apiPost<{ message: string }>('/contatos', contato)
}

export async function importContatos(file: File, campanha: string): Promise<ImportContatosResult> {
  const form = new FormData()
  form.append('file', file)
  form.append('campanha', campanha)
  return apiPostForm<ImportContatosResult>('/contatos/import', form)
}
