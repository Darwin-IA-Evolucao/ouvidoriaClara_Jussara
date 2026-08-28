import * as React from 'react'
import { useState, useMemo } from 'react'
import { Box, Typography, Select, MenuItem, FormControl } from '@mui/material'
import InboxOutlined from '@mui/icons-material/InboxOutlined'
import {
  PieChart, Pie, Cell, Tooltip as ReTooltip, ResponsiveContainer,
} from 'recharts'
import GlassPanel from './GlassPanel'
import { formatCategoryName } from '../utils/categories'
import type { Stat, Ocorrencia } from '../types'

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
}

interface DistributionViewerProps {
  stats: Stat
  kanbanData: { name: string; value: number }[]
  ocorrencias?: Ocorrencia[]
}

type Tab = 'regiao' | 'categoria' | 'tipo'

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
  { key: 'tipo', label: 'Colunas' },
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

const DistributionViewer: React.FC<DistributionViewerProps> = ({ stats, ocorrencias = [] }) => {
  const [tab, setTab] = useState<Tab>('regiao')
  const [selectedRegiao, setSelectedRegiao] = useState('')

  const regiaoData = [...(stats.regioes || [])]
    .sort((a, b) => b.qtdRegiao - a.qtdRegiao)
    .map((r, i) => ({ name: r.regiao, value: r.qtdRegiao, color: PIE_COLORS[i % PIE_COLORS.length] }))

  const regioesDisponiveis = useMemo(() => {
    const set = new Set<string>()
    ocorrencias.forEach((o) => {
      const r = o.detalhes?.regiao
      if (r) set.add(r)
    })
    return Array.from(set).sort()
  }, [ocorrencias])

  const categoriaData = useMemo(() => {
    const filtered = selectedRegiao
      ? ocorrencias.filter((o) => o.detalhes?.regiao === selectedRegiao)
      : ocorrencias
    const map: Record<string, number> = {}
    for (const o of filtered) {
      const cat = o.categoria || 'Sem Categoria'
      map[cat] = (map[cat] || 0) + 1
    }
    return Object.entries(map)
      .map(([categoria, value]) => ({ name: formatCategoryName(categoria), value, raw: categoria }))
      .sort((a, b) => b.value - a.value)
  }, [ocorrencias, selectedRegiao])

  const statusData = useMemo(() => {
    const filtered = selectedRegiao
      ? ocorrencias.filter((o) => o.detalhes?.regiao === selectedRegiao)
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
        map[st] = (map[st] || 0) + 1
      }
    }
    return Object.entries(map)
      .map(([name, value]) => ({ name, value, raw: name }))
      .sort((a, b) => b.value - a.value)
  }, [ocorrencias, selectedRegiao])

  const activeData = tab === 'regiao' ? regiaoData : tab === 'categoria' ? categoriaData : statusData
  const hasData = activeData.length > 0 && activeData.some((d) => d.value > 0)

  return (
    <GlassPanel className="p-5" borderRadius={14}>
      <Box display="flex" alignItems="center" flexWrap="wrap" gap={1.5} mb={3}>
        <Typography variant="h6" sx={{ fontWeight: 600, fontSize: 14, letterSpacing: '0.04em', textTransform: 'uppercase', color: 'hsl(var(--accent))', mr: 1.5 }}>
          Distribuição por:
        </Typography>
        <Box display="flex" gap={1}>
          {TABS.map((t) => (
            <Box
              key={t.key}
              onClick={() => setTab(t.key)}
              sx={{
                px: 2, py: 0.7, borderRadius: '20px', cursor: 'pointer',
                fontSize: 12, fontWeight: 600, letterSpacing: '0.03em',
                border: `1.5px solid ${tab === t.key ? 'hsl(var(--accent))' : 'hsl(var(--border))'}`,
                bgcolor: tab === t.key ? 'hsl(var(--accent) / 0.12)' : 'transparent',
                color: tab === t.key ? 'hsl(var(--accent))' : 'hsl(var(--text-secondary))',
                transition: 'all 0.2s ease',
                '&:hover': { borderColor: 'hsl(var(--accent))', color: 'hsl(var(--accent))' },
              }}
            >
              {t.label}
            </Box>
          ))}
        </Box>
      </Box>

      {!hasData ? (
        <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height={300} gap={1.5}>
          <InboxOutlined sx={{ fontSize: 48, color: 'hsl(var(--text-secondary) / 0.4)' }} />
          <Typography variant="body2" sx={{ color: 'hsl(var(--text-secondary))', fontWeight: 600 }}>
            Nenhuma ocorrência encontrada
          </Typography>
          <Typography variant="caption" sx={{ color: 'hsl(var(--text-secondary) / 0.7)', textAlign: 'center', maxWidth: 280 }}>
            Não há dados de {tab === 'regiao' ? 'região' : tab === 'categoria' ? 'categoria' : 'status'} para o período selecionado.
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
                      <MenuItem key={r} value={r}>{r}</MenuItem>
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
                      <MenuItem key={r} value={r}>{r}</MenuItem>
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
        </>
      )}
    </GlassPanel>
  )
}

export default DistributionViewer
