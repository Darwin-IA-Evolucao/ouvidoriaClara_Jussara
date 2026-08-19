import { apiGet, apiPost, apiPut, apiDelete } from '../utils/api'
import type { Ocorrencia, OcorrenciaRequest, OcorrenciaUpdateRequest } from '../types'

export async function aprovarInquerito(id: number, mensagem: string): Promise<void> {
  await apiPost<void>(`/aprovar/${id}`, { mensagem })
}

export async function aprovarRequerimento(id: number, mensagem: string): Promise<void> {
  await apiPost<void>(`/aprovar/requerimento/${id}`, { mensagem })
}

export async function aprovarCausaAnimal(id: number, mensagem: string): Promise<void> {
  await apiPost<void>(`/aprovar/causa-animal/${id}`, { mensagem })
}

export async function aprovarComoAmbos(id: number, mensagem: string): Promise<void> {
  await apiPost<void>(`/indicreq/${id}`, { mensagem })
}

export async function reprovarInquerito(id: number, mensagem: string): Promise<void> {
  await apiPost<void>(`/reprovar/${id}`, { mensagem })
}

export async function finalizarReclamacao(id: number, mensagem: string): Promise<void> {
  await apiPost<void>(`/finalizar/${id}`, { mensagem })
}

export async function colocarEmAnalise(id: number, data: string): Promise<void> {
  await apiPost<void>(`/analise/${id}`, { data })
}

export async function colocarComoCriado(id: number): Promise<void> {
  await apiPost<void>(`/criado/${id}`, {})
}

export async function getAllOcorrencias(telefone?: string): Promise<Ocorrencia[]> {
  const query = telefone ? `?telefone=${telefone}` : ''
  const data = await apiGet<Ocorrencia[] | null>(`/ocorrencias${query}`)
  return Array.isArray(data) ? data : []
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
