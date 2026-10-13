// O agente de IA grava "N/A" quando o lead não informa um campo.
// Nunca exibir o placeholder cru na UI — sempre "Não Informado".
const NA_PATTERN = /^n\s*[./\\]\s*a\.?$/i

export function isNA(v: unknown): boolean {
  return typeof v === 'string' && NA_PATTERN.test(v.trim())
}

// Para renderização de campos: null/undefined/''/'N/A' → 'Não Informado'.
export function displayField(v: string | null | undefined): string {
  return v == null || v.trim() === '' || isNA(v) ? 'Não Informado' : v
}

// Para inputs controlados: 'N/A' aparece como campo vazio (placeholder "Não informado"),
// mas o valor cru continua no estado e é persistido sem alteração.
export function editValue(v: string | undefined | null): string {
  return isNA(v) ? '' : (v ?? '')
}
