import { apiGet, apiPost, apiPut, apiDelete } from '../utils/api'
import type { Endereco, Logradouro } from '../types'

export async function getRegiao(rua: string, bairro?: string): Promise<Logradouro> {
  const params = new URLSearchParams()
  params.set('rua', rua)
  if (bairro) params.set('bairro', bairro)
  return apiGet<Logradouro>(`/getRegiao?${params.toString()}`)
}

export async function getAllEnderecos(limit: number = 0, offset: number = 0, regiao?: string): Promise<{ enderecos: Logradouro[]; total: number }> {
  const params = new URLSearchParams()
  params.set('limit', String(limit))
  params.set('offset', String(offset))
  if (regiao) params.set('regiao', regiao)
  const res = await apiGet<{ enderecos: Logradouro[]; total?: number } | null>(`/enderecos?${params.toString()}`)
  if (res && Array.isArray(res.enderecos)) return { enderecos: res.enderecos, total: res.total ?? res.enderecos.length }
  return { enderecos: [], total: 0 }
}

export async function getAllRegioes(): Promise<string[]> {
  const res = await apiGet<string[] | null>('/endereco/regiao')
  if (Array.isArray(res)) return res
  return []
}

export async function cadastrarEnderecos(enderecos: Endereco[]): Promise<void> {
  await apiPost<void>('/cadastrarEnderecos', enderecos)
}

export async function createEndereco(data: { logradouro: string; bairro: string; regiao: string }): Promise<void> {
  await apiPost<void>('/cadastrarEnderecos', [{ Logradouro: data.logradouro, Bairro: data.bairro, 'Região': data.regiao }])
}

export async function updateEndereco(id: number, data: { logradouro: string; bairro: string; regiao: string }): Promise<void> {
  await apiPut<void>(`/endereco/${id}`, data)
}

export async function deleteEndereco(id: number): Promise<void> {
  await apiDelete<void>(`/endereco/${id}`)
}
