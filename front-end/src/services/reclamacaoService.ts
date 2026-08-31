import { apiGet, apiPost, apiPut, apiDelete } from '../utils/api'
import type { Ocorrencia, OcorrenciaRequest, OcorrenciaUpdateRequest } from '../types'

export async function aprovarInquerito(id: number, mensagem: string, idUsuario: number): Promise<void> {
  await apiPost<void>(`/aprovar/${id}`, { mensagem, idUsuario })
}

export async function aprovarRequerimento(id: number, mensagem: string, idUsuario: number): Promise<void> {
  await apiPost<void>(`/aprovar/requerimento/${id}`, { mensagem, idUsuario })
}

export async function aprovarCausaAnimal(id: number, mensagem: string, idUsuario: number): Promise<void> {
  await apiPost<void>(`/aprovar/causa-animal/${id}`, { mensagem, idUsuario })
}

export async function aprovarComoAmbos(id: number, mensagem: string, idUsuario: number): Promise<void> {
  await apiPost<void>(`/indicreq/${id}`, { mensagem, idUsuario })
}

export async function reprovarInquerito(id: number, mensagem: string, idUsuario: number): Promise<void> {
  await apiPost<void>(`/reprovar/${id}`, { mensagem, idUsuario })
}

export async function finalizarReclamacao(id: number, mensagem: string, idUsuario: number): Promise<void> {
  await apiPost<void>(`/finalizar/${id}`, { mensagem, idUsuario })
}

export async function colocarEmAnalise(id: number, data: string, mensagem: string, idUsuario: number): Promise<void> {
  await apiPost<void>(`/analise/${id}`, { data, mensagem, idUsuario })
}

export async function colocarComoCriado(id: number, idUsuario: number): Promise<void> {
  await apiPost<void>(`/criado/${id}`, { idUsuario })
}

export async function getAllOcorrencias(telefone?: string, limit: number = 0, offset: number = 0): Promise<{ ocorrencias: Ocorrencia[]; total: number }> {
  const params = new URLSearchParams()
  if (telefone) params.set('telefone', telefone)
  params.set('limit', String(limit))
  params.set('offset', String(offset))
  const res = await apiGet<{ ocorrencias: Ocorrencia[]; total: number } | Ocorrencia[] | null>(`/ocorrencias?${params.toString()}`)
  if (Array.isArray(res)) return { ocorrencias: res, total: res.length }
  if (res && Array.isArray(res.ocorrencias)) return { ocorrencias: res.ocorrencias, total: res.total ?? res.ocorrencias.length }
  return { ocorrencias: [], total: 0 }
}

export async function getOcorrenciaById(id: number): Promise<Ocorrencia> {
  return apiGet<Ocorrencia>(`/ocorrencia/${id}`)
}

export async function createOcorrencia(data: OcorrenciaRequest): Promise<{ message: string; id: number }> {
  return apiPost<{ message: string; id: number }>('/ocorrencia', data)
}

export async function updateOcorrencia(id: number, data: OcorrenciaUpdateRequest): Promise<void> {
  await apiPut<void>(`/ocorrencia/${id}`, data)
}

export async function deleteOcorrencia(id: number): Promise<void> {
  await apiDelete<void>(`/ocorrencia/${id}`)
}
