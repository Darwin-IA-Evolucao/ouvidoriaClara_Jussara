import { apiGet, apiPost, apiPut, apiDelete } from '../utils/api'
import type { Campanha } from '../types'

export async function getAllCampanhas(): Promise<Campanha[]> {
  return apiGet<Campanha[]>('/campanhas')
}

export async function createCampanha(campanha: Pick<Campanha, 'palavraChave'>): Promise<Campanha> {
  return apiPost<Campanha>('/campanhas', campanha)
}

export async function updateCampanha(id: number, campanha: Pick<Campanha, 'palavraChave'>): Promise<Campanha> {
  return apiPut<Campanha>(`/campanhas/${id}`, campanha)
}

export async function deleteCampanha(id: number): Promise<{ message: string }> {
  return apiDelete<{ message: string }>(`/campanhas/${id}`)
}
