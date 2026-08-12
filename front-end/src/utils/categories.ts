import { apiGet } from './api'

export async function fetchCategorias(): Promise<string[]> {
  try {
    const data = await apiGet<string[]>('/categorias')
    return Array.isArray(data) ? data : []
  } catch {
    return []
  }
}

export function formatCategoryName(cat: string): string {
  if (categoryDisplayName[cat]) return categoryDisplayName[cat]
  return cat.split(' ').map((w) => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')
}

export const categoryDisplayName: Record<string, string> = {
  'geral': 'Geral',
  'maus tratos': 'Maus Tratos',
  'abandono presenciado': 'Abandono Presenciado',
  'animal apareceu na rua': 'Animal Apareceu Na Rua',
  'ajuda animal comunitario': 'Ajuda Animal Comunitário',
  'saude animal': 'Saúde Animal',
  'castracao eletiva': 'Castração Eletiva',
  'castracao emergencial': 'Castração Emergencial',
  'animais nao domiciliados': 'Animais Não Domiciliados',
  'animal desaparecido': 'Animal Desaparecido',
  'animal para ser adotado': 'Animal Para Ser Adotado',
  'adocao de animais': 'Adoção De Animais',
  'animal grande porte': 'Animal Grande Porte',
  'animal atropelado': 'Animal Atropelado',
  'cuidados animais': 'Cuidados Animais',
  'animais silvestres': 'Animais Silvestres',
  'equipamentos': 'Equipamentos',
  'asfalto': 'Asfalto',
  'transporte': 'Transporte',
  'saude': 'Saúde',
  'educacao': 'Educação',
  'seguranca': 'Segurança',
  'zeladoria': 'Zeladoria',
  'outros': 'Outros',
}
