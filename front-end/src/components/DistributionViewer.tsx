import * as React from 'react'
import { useState, useMemo, useEffect } from 'react'
import { Box, Typography, Select, MenuItem, FormControl, Autocomplete, TextField, CircularProgress } from '@mui/material'
import InboxOutlined from '@mui/icons-material/InboxOutlined'
import {
  PieChart, Pie, Cell, Tooltip as ReTooltip, ResponsiveContainer,
} from 'recharts'
import GlassPanel from './GlassPanel'
import { formatCategoryName } from '../utils/categories'
import { isNA, displayField } from '../utils/display'
import { normalizeCidades } from '../utils/cityDistribution'
import { getOcorrenciasPorCidade } from '../services/reclamacaoService'
import type { Stat, Ocorrencia, OcorrenciaPorCidade } from '../types'

const STATUS_COLUMN_LABELS: Record<string, string> = {
  'criado': 'Sem Tratativa',
  'em análise': 'Em Análise',
  'em analise': 'Em Análise',
  'reprovado': 'Desqualificado',
  'finalizado': 'Finalizado',
}

const TYPE_COLUMN_LABELS: Record<string, string> = {
  'indicacao': 'Concluído como Indicação',
  'requerimento': 'Concluído como Requerimento',
  'causa animal': 'Concluído como Causa Animal',
  'outros': 'Concluído (Outros)',
}

interface DistributionViewerProps {
  stats: Stat
  kanbanData: { name: string; value: number }[]
  ocorrencias?: Ocorrencia[]
  dateStart?: string
  dateEnd?: string
}

type Tab = 'regiao' | 'categoria' | 'tipo' | 'cidade'
type CidadeModo = 'ranking' | 'ocorrencias'

const PIE_COLORS = [
  '#41669C', '#E89E70', '#66BB80', '#62A1D8', '#E2AF7A', '#A1A9B8', '#5282AE', '#53A16E',
  '#e63946', '#7b2cbf', '#3a86ff', '#ffbe0b', '#fb5607', '#8338ec', '#06d6a0', '#ef476f',
  '#118ab2', '#ffd166', '#073b4c', '#7209b7', '#560bad', '#f72585', '#4361ee', '#4cc9f0',
  '#80b918', '#55a630', '#2b9348', '#007f5f', '#aacc00', '#d4d700', '#e9c46a', '#e76f51',
  '#264653', '#2a9d8f', '#e07a5f', '#81b29a', '#f2cc8f', '#bc4749', '#6a994e', '#a7c957',
]

const TABS: { key: Tab; label: string }[] = [
  { key: 'regiao', label: 'Região' },
  { key: 'categoria', label: 'Categoria' },
  { key: 'cidade', label: 'Cidade' },
  { key: 'tipo', label: 'Status' },
]

const CIDADE_MODOS: { key: CidadeModo; label: string }[] = [
  { key: 'ranking', label: 'Ranking' },
  { key: 'ocorrencias', label: 'Ocorrências' },
]

interface TooltipProps {
  active?: boolean
  payload?: { name: string; value: number; payload: { percent?: number } }[]
  label?: string
}

const CustomTooltip: React.FC<TooltipProps> = ({ active, payload }) => {
  if (!active || !payload?.length) return null
  const { name, value } = payload[0]
  return (
    <Box sx={{ bgcolor: 'hsl(var(--surface))', border: '1px solid hsl(var(--border))', borderRadius: 2, px: 2, py: 1.5 }}>
      <Typography variant="caption" sx={{ color: 'hsl(var(--text-secondary))', display: 'block', mb: 0.3, fontSize: 11 }}>
        {name}
      </Typography>
      <Typography variant="body2" sx={{ color: 'hsl(var(--text-primary))', fontWeight: 700, fontSize: 15 }}>
        {value} <span style={{ color: 'hsl(var(--accent))', fontWeight: 400, fontSize: 12 }}>ocorrências</span>
      </Typography>
    </Box>
  )
}


const renderCustomLegend = (items: { name: string; value: number; color: string }[]) => (
  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1.5, justifyContent: 'center', mt: 2 }}>
    {items.map((item) => (
      <Box key={item.name} display="flex" alignItems="center" gap={0.8}>
        <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: item.color, flexShrink: 0 }} />
        <Typography variant="caption" sx={{ color: 'hsl(var(--text-secondary))', fontSize: 12 }}>
          {item.name}
        </Typography>
        <Typography variant="caption" sx={{ color: item.color, fontWeight: 700, fontSize: 12 }}>
          {item.value}
        </Typography>
      </Box>
    ))}
  </Box>
)

const DistributionViewer: React.FC<DistributionViewerProps> = ({ stats, ocorrencias = [], dateStart, dateEnd }) => {
  const [tab, setTab] = useState<Tab>('regiao')
  const [selectedRegiao, setSelectedRegiao] = useState('')
  const [cidadesData, setCidadesData] = useState<OcorrenciaPorCidade[]>([])
  const [cidadesLoading, setCidadesLoading] = useState(false)
  const [cidadesError, setCidadesError] = useState(false)
  const [cidadeModo, setCidadeModo] = useState<CidadeModo>('ranking')
  const [cidadeSelecionada, setCidadeSelecionada] = useState<string | null>(null)

  // Aba cidade: busca lazy ao entrar na aba ou ao mudar o filtro de datas
  useEffect(() => {
    if (tab !== 'cidade') return
    setCidadesLoading(true)
    setCidadesError(false)
    getOcorrenciasPorCidade(dateStart || undefined, dateEnd || undefined)
      .then(setCidadesData)
      .catch(() => { setCidadesData([]); setCidadesError(true) })
      .finally(() => setCidadesLoading(false))
  }, [tab, dateStart, dateEnd])

  const cidadesNorm = useMemo(() => normalizeCidades(cidadesData), [cidadesData])
  const cidadesTotalGeral = useMemo(() => cidadesNorm.reduce((sum, cidade) => sum + cidade.total, 0), [cidadesNorm])
  const maxCidadeTotal = cidadesNorm[0]?.total ?? 0

  const cidadePieData = useMemo(() => {
    const c = cidadesNorm.find((x) => x.cidade === cidadeSelecionada)
    if (!c) return []
    return Object.entries(c.categoria || {})
      .map(([cat, n], i) => ({ name: formatCategoryName(cat), value: n, color: PIE_COLORS[i % PIE_COLORS.length] }))
      .sort((a, b) => b.value - a.value)
  }, [cidadesNorm, cidadeSelecionada])

  const cidadeSelecionadaTotal = useMemo(
    () => cidadesNorm.find((x) => x.cidade === cidadeSelecionada)?.total ?? 0,
    [cidadesNorm, cidadeSelecionada],
  )

  const regiaoData = useMemo(() => {
    const map = new Map<string, number>()
    for (const r of stats.regioes || []) {
      const name = displayField(r.regiao)
      map.set(name, (map.get(name) || 0) + r.qtdRegiao)
    }
    return [...map.entries()]
      .sort((a, b) => b[1] - a[1])
      .map(([name, value], i) => ({ name, value, color: PIE_COLORS[i % PIE_COLORS.length] }))
  }, [stats.regioes])

  const regioesDisponiveis = useMemo(() => {
    const set = new Set<string>()
    ocorrencias.forEach((o) => {
      const r = o.detalhes?.regiao
      if (r) set.add(displayField(r))
    })
    return Array.from(set).sort()
  }, [ocorrencias])

  const categoriaData = useMemo(() => {
    const filtered = selectedRegiao
      ? ocorrencias.filter((o) => displayField(o.detalhes?.regiao) === selectedRegiao)
      : ocorrencias
    const map: Record<string, number> = {}
    for (const o of filtered) {
      const name = formatCategoryName(o.categoria || 'Sem Categoria')
      map[name] = (map[name] || 0) + 1
    }
    return Object.entries(map)
      .map(([name, value]) => ({ name, value, raw: name }))
      .sort((a, b) => b.value - a.value)
  }, [ocorrencias, selectedRegiao])

  const statusData = useMemo(() => {
    const filtered = selectedRegiao
      ? ocorrencias.filter((o) => displayField(o.detalhes?.regiao) === selectedRegiao)
      : ocorrencias
    const map: Record<string, number> = {}
    for (const o of filtered) {
      const st = (o.status || 'Sem Status').toLowerCase()
      const tipo = (o.tipo || '').toLowerCase()
      if (st === 'aprovado' && TYPE_COLUMN_LABELS[tipo]) {
        const name = TYPE_COLUMN_LABELS[tipo]
        map[name] = (map[name] || 0) + 1
      } else if (STATUS_COLUMN_LABELS[st]) {
        const name = STATUS_COLUMN_LABELS[st]
        map[name] = (map[name] || 0) + 1
      } else {
        const name = isNA(st) ? 'Não Informado' : st
        map[name] = (map[name] || 0) + 1
      }
    }
    return Object.entries(map)
      .map(([name, value]) => ({ name, value, raw: name }))
      .sort((a, b) => b.value - a.value)
  }, [ocorrencias, selectedRegiao])

  const activeData = tab === 'regiao' ? regiaoData : tab === 'categoria' ? categoriaData : statusData
  const hasData = tab === 'cidade'
    ? cidadesNorm.length > 0
    : activeData.length > 0 && activeData.some((d) => d.value > 0)

  return (
    <GlassPanel className="p-5" borderRadius={14}>
      <Box display="flex" alignItems="center" justifyContent="flex-end" mb={2}>
        <Box
          role="group"
          aria-label="Alternar distribuição"
          sx={{
            display: 'flex',
            flexWrap: 'nowrap',
            width: { xs: '100%', sm: 'auto' },
            justifyContent: 'flex-end',
            bgcolor: 'hsl(var(--surface-2))',
            borderRadius: 999,
            p: 0.4,
            gap: 0.4,
            border: '1px solid hsl(var(--border))',
          }}
        >
          {TABS.map((t) => (
            <Box
              component="button"
              key={t.key}
              type="button"
              aria-pressed={tab === t.key}
              onClick={() => setTab(t.key)}
              sx={{
                border: 'none',
                cursor: 'pointer',
                borderRadius: 999,
                py: 0.6,
                px: { xs: 0.5, sm: 2 },
                flex: { xs: 1, sm: 'initial' },
                minWidth: 0,
                whiteSpace: 'nowrap',
                fontSize: { xs: 10, sm: 11 },
                fontWeight: 600,
                fontFamily: 'inherit',
                letterSpacing: '0.03em',
                color: tab === t.key ? 'hsl(var(--text-primary))' : 'hsl(var(--text-secondary))',
                bgcolor: tab === t.key ? 'hsl(var(--accent) / 0.18)' : 'transparent',
                touchAction: 'manipulation',
                transition: 'background var(--dur-fast) var(--ease-brand), color var(--dur-fast) var(--ease-brand)',
                '&:hover': { color: 'hsl(var(--text-primary))' },
                '&:focus-visible': { outline: '2px solid hsl(var(--accent))', outlineOffset: 1 },
              }}
            >
              {t.label}
            </Box>
          ))}
        </Box>
      </Box>

      {tab === 'cidade' && cidadesLoading ? (
        <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height={300} gap={1.5}>
          <CircularProgress size={32} sx={{ color: 'hsl(var(--accent))' }} />
          <Typography variant="body2" sx={{ color: 'hsl(var(--text-secondary))', fontWeight: 600 }}>
            Carregando cidades…
          </Typography>
        </Box>
      ) : tab === 'cidade' && cidadesError ? (
        <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height={300} gap={1.5}>
          <InboxOutlined sx={{ fontSize: 48, color: 'hsl(var(--text-secondary) / 0.4)' }} />
          <Typography variant="body2" sx={{ color: 'hsl(var(--text-secondary))', fontWeight: 600 }}>
            Não foi possível carregar os dados de cidade
          </Typography>
          <Typography variant="caption" sx={{ color: 'hsl(var(--text-secondary) / 0.7)', textAlign: 'center', maxWidth: 320 }}>
            O endpoint /ocorrencias/cidades falhou. Verifique se o back-end está atualizado e tente novamente.
          </Typography>
        </Box>
      ) : !hasData ? (
        <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height={300} gap={1.5}>
          <InboxOutlined sx={{ fontSize: 48, color: 'hsl(var(--text-secondary) / 0.4)' }} />
          <Typography variant="body2" sx={{ color: 'hsl(var(--text-secondary))', fontWeight: 600 }}>
            Nenhuma ocorrência encontrada
          </Typography>
          <Typography variant="caption" sx={{ color: 'hsl(var(--text-secondary) / 0.7)', textAlign: 'center', maxWidth: 280 }}>
            Não há dados de {tab === 'regiao' ? 'região' : tab === 'categoria' ? 'categoria' : tab === 'cidade' ? 'cidade' : 'status'} para o período selecionado.
          </Typography>
        </Box>
      ) : (
        <>
          {tab === 'regiao' && (
            <Box>
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={regiaoData}
                    cx="50%"
                    cy="50%"
                    innerRadius="52%"
                    outerRadius="78%"
                    paddingAngle={3}
                    dataKey="value"
                    stroke="none"
                  >
                    {regiaoData.map((entry) => (
                      <Cell key={entry.name} fill={entry.color} />
                    ))}
                  </Pie>
                  <ReTooltip content={<CustomTooltip />} />
                </PieChart>
              </ResponsiveContainer>
              {renderCustomLegend(regiaoData)}
            </Box>
          )}

          {tab === 'categoria' && (
            <Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 2, flexWrap: 'wrap' }}>
                <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-secondary))', fontWeight: 500 }}>
                  Visualizar total de solicitações por Categoria × Região:
                </Typography>
                <FormControl size="small" sx={{ minWidth: 200 }}>
                  <Select
                    value={selectedRegiao}
                    onChange={(e) => setSelectedRegiao(e.target.value)}
                    displayEmpty
                    sx={{ fontSize: 12, height: 32, '& .MuiSelect-select': { py: 0.5 } }}
                  >
                    <MenuItem value="">Todas as Regiões</MenuItem>
                    {regioesDisponiveis.map((r) => (
                      <MenuItem key={r} value={r}>{displayField(r)}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Box>

              <Box sx={{ maxHeight: 220, overflowY: 'auto', mb: 2, pr: 0.5, '&::-webkit-scrollbar': { width: 10 }, '&::-webkit-scrollbar-thumb': { bgcolor: 'hsl(var(--scrollbar-thumb))', borderRadius: 3 } }}>
                {categoriaData.map((c, i) => {
                  const total = categoriaData.reduce((s, x) => s + x.value, 0)
                  const pct = total > 0 ? (c.value / total) * 100 : 0
                  const color = PIE_COLORS[i % PIE_COLORS.length]
                  return (
                    <Box key={c.raw} sx={{ display: 'flex', alignItems: 'center', gap: 1.5, px: 1.5, py: 1, mb: 0.5, borderRadius: 1.5, bgcolor: 'hsl(var(--surface-2) / 0.5)', border: '1px solid hsl(var(--border) / 0.5)', '&:last-child': { mb: 0 } }}>
                      <Typography sx={{ fontSize: 11, color: 'hsl(var(--text-secondary))', fontWeight: 600, minWidth: 18, textAlign: 'right' }}>#{i + 1}</Typography>
                      <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: color, flexShrink: 0 }} />
                      <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-primary))', fontWeight: 500, flex: 1 }}>{c.name}</Typography>
                      <Box sx={{ flex: 2, mx: 1 }}>
                        <Box sx={{ height: 5, borderRadius: 3, bgcolor: 'hsl(var(--border))', overflow: 'hidden' }}>
                          <Box sx={{ height: '100%', width: `${pct}%`, bgcolor: color, borderRadius: 3, transition: 'width 0.4s ease' }} />
                        </Box>
                      </Box>
                      <Typography sx={{ fontSize: 13, color: color, fontWeight: 700, minWidth: 28, textAlign: 'right' }}>{c.value}</Typography>
                      <Typography sx={{ fontSize: 10, color: 'hsl(var(--text-secondary))', minWidth: 36, textAlign: 'right' }}>{pct.toFixed(0)}%</Typography>
                    </Box>
                  )
                })}
              </Box>

              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={categoriaData}
                    cx="50%"
                    cy="50%"
                    innerRadius="52%"
                    outerRadius="78%"
                    paddingAngle={3}
                    dataKey="value"
                    stroke="none"
                  >
                    {categoriaData.map((entry, i) => (
                      <Cell key={entry.name} fill={PIE_COLORS[i % PIE_COLORS.length]} />
                    ))}
                  </Pie>
                  <ReTooltip content={<CustomTooltip />} />
                </PieChart>
              </ResponsiveContainer>
            </Box>
          )}

          {tab === 'tipo' && (
            <Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 2, flexWrap: 'wrap' }}>
                <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-secondary))', fontWeight: 500 }}>
                  Visualizar total de solicitações por Status × Região:
                </Typography>
                <FormControl size="small" sx={{ minWidth: 200 }}>
                  <Select
                    value={selectedRegiao}
                    onChange={(e) => setSelectedRegiao(e.target.value)}
                    displayEmpty
                    sx={{ fontSize: 12, height: 32, '& .MuiSelect-select': { py: 0.5 } }}
                  >
                    <MenuItem value="">Todas as Regiões</MenuItem>
                    {regioesDisponiveis.map((r) => (
                      <MenuItem key={r} value={r}>{displayField(r)}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Box>

              <Box sx={{ maxHeight: 220, overflowY: 'auto', mb: 2, pr: 0.5, '&::-webkit-scrollbar': { width: 10 }, '&::-webkit-scrollbar-thumb': { bgcolor: 'hsl(var(--scrollbar-thumb))', borderRadius: 3 } }}>
                {statusData.map((s, i) => {
                  const total = statusData.reduce((sum, x) => sum + x.value, 0)
                  const pct = total > 0 ? (s.value / total) * 100 : 0
                  const color = PIE_COLORS[i % PIE_COLORS.length]
                  return (
                    <Box key={s.raw} sx={{ display: 'flex', alignItems: 'center', gap: 1.5, px: 1.5, py: 1, mb: 0.5, borderRadius: 1.5, bgcolor: 'hsl(var(--surface-2) / 0.5)', border: '1px solid hsl(var(--border) / 0.5)', '&:last-child': { mb: 0 } }}>
                      <Typography sx={{ fontSize: 11, color: 'hsl(var(--text-secondary))', fontWeight: 600, minWidth: 18, textAlign: 'right' }}>#{i + 1}</Typography>
                      <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: color, flexShrink: 0 }} />
                      <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-primary))', fontWeight: 500, flex: 1 }}>{s.name}</Typography>
                      <Box sx={{ flex: 2, mx: 1 }}>
                        <Box sx={{ height: 5, borderRadius: 3, bgcolor: 'hsl(var(--border))', overflow: 'hidden' }}>
                          <Box sx={{ height: '100%', width: `${pct}%`, bgcolor: color, borderRadius: 3, transition: 'width 0.4s ease' }} />
                        </Box>
                      </Box>
                      <Typography sx={{ fontSize: 13, color: color, fontWeight: 700, minWidth: 28, textAlign: 'right' }}>{s.value}</Typography>
                      <Typography sx={{ fontSize: 10, color: 'hsl(var(--text-secondary))', minWidth: 36, textAlign: 'right' }}>{pct.toFixed(0)}%</Typography>
                    </Box>
                  )
                })}
              </Box>

              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={statusData}
                    cx="50%"
                    cy="50%"
                    innerRadius="52%"
                    outerRadius="78%"
                    paddingAngle={3}
                    dataKey="value"
                    stroke="none"
                  >
                    {statusData.map((entry, i) => (
                      <Cell key={entry.name} fill={PIE_COLORS[i % PIE_COLORS.length]} />
                    ))}
                  </Pie>
                  <ReTooltip content={<CustomTooltip />} />
                </PieChart>
              </ResponsiveContainer>
            </Box>
          )}

          {tab === 'cidade' && (
            <Box>
              {/* Toggle de modo: Ranking | Ocorrências */}
              <Box
                role="group"
                aria-label="Modo da visualização por cidade"
                sx={{
                  display: 'flex',
                  width: 'fit-content',
                  bgcolor: 'hsl(var(--surface-2))',
                  borderRadius: 999,
                  p: 0.4,
                  gap: 0.4,
                  border: '1px solid hsl(var(--border))',
                  mb: 2.5,
                }}
              >
                {CIDADE_MODOS.map((m) => (
                  <Box
                    component="button"
                    key={m.key}
                    type="button"
                    aria-pressed={cidadeModo === m.key}
                    onClick={() => setCidadeModo(m.key)}
                    sx={{
                      border: 'none',
                      cursor: 'pointer',
                      borderRadius: 999,
                      py: 0.6,
                      px: 2,
                      fontSize: 11,
                      fontWeight: 600,
                      fontFamily: 'inherit',
                      letterSpacing: '0.03em',
                      color: cidadeModo === m.key ? 'hsl(var(--text-primary))' : 'hsl(var(--text-secondary))',
                      bgcolor: cidadeModo === m.key ? 'hsl(var(--accent) / 0.18)' : 'transparent',
                      touchAction: 'manipulation',
                      transition: 'background var(--dur-fast) var(--ease-brand), color var(--dur-fast) var(--ease-brand)',
                      '&:hover': { color: 'hsl(var(--text-primary))' },
                      '&:focus-visible': { outline: '2px solid hsl(var(--accent))', outlineOffset: 1 },
                    }}
                  >
                    {m.label}
                  </Box>
                ))}
              </Box>

              {cidadeModo === 'ranking' && (
                <Box>
                  <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-secondary))', fontWeight: 500, mb: 1.5, px: 1.5 }}>
                    Cidades com mais solicitações — clique em uma cidade para ver suas ocorrências:
                  </Typography>
                  <Box sx={{ maxHeight: 320, overflowY: 'auto', pr: 0.5, '&::-webkit-scrollbar': { width: 10 }, '&::-webkit-scrollbar-thumb': { bgcolor: 'hsl(var(--scrollbar-thumb))', borderRadius: 3 } }}>
                    {cidadesNorm.map((c, i) => {
                      const pct = maxCidadeTotal > 0 ? (c.total / maxCidadeTotal) * 100 : 0
                      const share = cidadesTotalGeral > 0 ? Math.round((c.total / cidadesTotalGeral) * 100) : 0
                      const isTop = i === 0
                      return (
                        <Box
                          component="button"
                          type="button"
                          key={c.cidade}
                          title={`Ver ocorrências de ${c.cidade}`}
                          onClick={() => { setCidadeSelecionada(c.cidade); setCidadeModo('ocorrencias') }}
                          sx={{
                            display: 'flex', alignItems: 'center', gap: 1.5, width: '100%', fontFamily: 'inherit', textAlign: 'left',
                            px: 1.5, py: 1, mb: 0.5, borderRadius: 1.5,
                            bgcolor: isTop ? 'hsl(var(--accent) / 0.10)' : 'hsl(var(--surface-2) / 0.5)',
                            border: '1px solid', borderColor: isTop ? 'hsl(var(--accent) / 0.4)' : 'hsl(var(--border) / 0.5)',
                            cursor: 'pointer',
                            transition: 'background var(--dur-fast) var(--ease-brand), border-color var(--dur-fast) var(--ease-brand)',
                            '&:hover': { bgcolor: 'hsl(var(--surface-2))', borderColor: 'hsl(var(--accent) / 0.5)' },
                            '&:focus-visible': { outline: '2px solid hsl(var(--accent))', outlineOffset: 1 },
                            '&:last-child': { mb: 0 },
                          }}
                        >
                          <Typography component="span" sx={{ fontSize: 11, color: isTop ? 'hsl(var(--accent))' : 'hsl(var(--text-secondary))', fontWeight: 700, minWidth: 18, textAlign: 'right' }}>#{i + 1}</Typography>
                          <Typography component="span" noWrap sx={{ fontSize: 12, color: 'hsl(var(--text-primary))', fontWeight: isTop ? 700 : 500, flex: 1, minWidth: 0 }}>{c.cidade}</Typography>
                          <Box sx={{ flex: 2, mx: 1 }}>
                            <Box sx={{ height: 5, borderRadius: 3, bgcolor: 'hsl(var(--border))', overflow: 'hidden' }}>
                              <Box sx={{ height: '100%', width: `${pct}%`, background: 'linear-gradient(90deg, hsl(var(--accent)), hsl(var(--accent) / 0.55))', borderRadius: 3, transition: 'width 0.4s ease' }} />
                            </Box>
                          </Box>
                          <Typography component="span" sx={{ fontSize: 13, color: 'hsl(var(--accent))', fontWeight: 700, minWidth: 28, textAlign: 'right', fontVariantNumeric: 'tabular-nums' }}>{c.total}</Typography>
                          <Typography component="span" sx={{ fontSize: 10, color: 'hsl(var(--text-secondary))', minWidth: 36, textAlign: 'right', fontVariantNumeric: 'tabular-nums' }}>{share}%</Typography>
                        </Box>
                      )
                    })}
                  </Box>
                </Box>
              )}

              {cidadeModo === 'ocorrencias' && (
                <Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 2, flexWrap: 'wrap' }}>
                    <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-secondary))', fontWeight: 500 }}>
                      Buscar cidade:
                    </Typography>
                    <Autocomplete
                      size="small"
                      options={cidadesNorm.map((c) => c.cidade)}
                      value={cidadeSelecionada}
                      onChange={(_, v) => setCidadeSelecionada(v)}
                      sx={{ minWidth: 240 }}
                      slotProps={{
                        paper: { sx: { bgcolor: 'hsl(var(--surface))', color: 'hsl(var(--text-primary))', border: '1px solid hsl(var(--border))' } },
                      }}
                      renderInput={(params) => (
                        <TextField
                          {...params}
                          placeholder="Ex.: Sorocaba"
                          sx={{
                            '& .MuiOutlinedInput-root': {
                              color: 'hsl(var(--text-primary))',
                              borderRadius: '10px',
                              fontSize: 13,
                              '& fieldset': { borderColor: 'hsl(var(--border))' },
                              '&:hover fieldset': { borderColor: 'hsl(var(--primary-light))' },
                              '&.Mui-focused fieldset': { borderColor: 'hsl(var(--accent))', borderWidth: '1.5px' },
                              bgcolor: 'hsl(var(--surface-2))',
                            },
                          }}
                        />
                      )}
                    />
                    {cidadeSelecionada && (
                      <Typography sx={{ fontSize: 12, color: 'hsl(var(--accent))', fontWeight: 700 }}>
                        {cidadeSelecionadaTotal} solicitações
                      </Typography>
                    )}
                  </Box>

                  {cidadeSelecionada ? (
                    cidadePieData.length > 0 ? (
                      <>
                        <Box sx={{ maxHeight: 220, overflowY: 'auto', mb: 2, pr: 0.5, '&::-webkit-scrollbar': { width: 10 }, '&::-webkit-scrollbar-thumb': { bgcolor: 'hsl(var(--scrollbar-thumb))', borderRadius: 3 } }}>
                          {cidadePieData.map((c, i) => {
                            const total = cidadePieData.reduce((s, x) => s + x.value, 0)
                            const pctv = total > 0 ? (c.value / total) * 100 : 0
                            return (
                              <Box key={c.name} sx={{ display: 'flex', alignItems: 'center', gap: 1.5, px: 1.5, py: 1, mb: 0.5, borderRadius: 1.5, bgcolor: 'hsl(var(--surface-2) / 0.5)', border: '1px solid hsl(var(--border) / 0.5)', '&:last-child': { mb: 0 } }}>
                                <Typography sx={{ fontSize: 11, color: 'hsl(var(--text-secondary))', fontWeight: 600, minWidth: 18, textAlign: 'right' }}>#{i + 1}</Typography>
                                <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: c.color, flexShrink: 0 }} />
                                <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-primary))', fontWeight: 500, flex: 1 }}>{c.name}</Typography>
                                <Box sx={{ flex: 2, mx: 1 }}>
                                  <Box sx={{ height: 5, borderRadius: 3, bgcolor: 'hsl(var(--border))', overflow: 'hidden' }}>
                                    <Box sx={{ height: '100%', width: `${pctv}%`, bgcolor: c.color, borderRadius: 3, transition: 'width 0.4s ease' }} />
                                  </Box>
                                </Box>
                                <Typography sx={{ fontSize: 13, color: c.color, fontWeight: 700, minWidth: 28, textAlign: 'right' }}>{c.value}</Typography>
                                <Typography sx={{ fontSize: 10, color: 'hsl(var(--text-secondary))', minWidth: 36, textAlign: 'right' }}>{pctv.toFixed(0)}%</Typography>
                              </Box>
                            )
                          })}
                        </Box>

                        <ResponsiveContainer width="100%" height={300}>
                          <PieChart>
                            <Pie
                              data={cidadePieData}
                              cx="50%"
                              cy="50%"
                              innerRadius="52%"
                              outerRadius="78%"
                              paddingAngle={3}
                              dataKey="value"
                              stroke="none"
                            >
                              {cidadePieData.map((entry) => (
                                <Cell key={entry.name} fill={entry.color} />
                              ))}
                            </Pie>
                            <ReTooltip content={<CustomTooltip />} />
                          </PieChart>
                        </ResponsiveContainer>
                      </>
                    ) : (
                      <Typography variant="body2" sx={{ color: 'hsl(var(--text-secondary))', textAlign: 'center', py: 6 }}>
                        Essa cidade não tem solicitações no período.
                      </Typography>
                    )
                  ) : (
                    <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height={240} gap={1.5}>
                      <InboxOutlined sx={{ fontSize: 40, color: 'hsl(var(--text-secondary) / 0.4)' }} />
                      <Typography variant="body2" sx={{ color: 'hsl(var(--text-secondary))' }}>
                        Selecione uma cidade para ver a distribuição por categoria.
                      </Typography>
                    </Box>
                  )}
                </Box>
              )}
            </Box>
          )}
        </>
      )}
    </GlassPanel>
  )
}

export default DistributionViewer
