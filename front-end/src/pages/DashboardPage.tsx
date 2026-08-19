import * as React from 'react'
import { useState, useEffect, useMemo } from 'react'
import { Typography, Box, TextField, Button, useMediaQuery, useTheme } from '@mui/material'
import { FileText, Users, TrendingUp, TrendingDown, BarChart3, XCircle, Inbox, Heart, CheckCircle } from 'lucide-react'
import GlassPanel from '../components/GlassPanel'
import PageHeader from '../components/PageHeader'
import { inputSx } from '../utils/inputSx'
import { getStats } from '../services/statsService'
import { getAllOcorrencias } from '../services/reclamacaoService'
import { getAllContatos } from '../services/contatoService'
import { toUTCDate } from '../utils/date'
import DistributionViewer from '../components/DistributionViewer'
import PageLoader from '../components/PageLoader'
import type { Stat, Ocorrencia, Contato } from '../types'

interface TrendInfo {
  direction: 'up' | 'down' | 'neutral'
  percent: number
  prevLabel: string
}

interface StatCardProps {
  label: string
  value: string | number
  icon: React.ElementType
  color?: string
  subtitle?: string
  trend?: TrendInfo | null
  sx?: React.CSSProperties
}

const StatCard: React.FC<StatCardProps> = ({ label, value, icon: Icon, color = 'hsl(var(--accent))', subtitle, trend, sx }) => (
  <GlassPanel
    className="p-5 flex flex-col gap-3"
    borderRadius={14}
    style={{ transition: 'transform 0.2s ease, box-shadow 0.2s ease', cursor: 'default', ...sx }}
  >
    <Box display="flex" alignItems="center" justifyContent="space-between">
      <Box>
        <Typography variant="body2" sx={{ color: 'hsl(var(--text-secondary))', fontWeight: 500, fontSize: 12, letterSpacing: '0.04em', textTransform: 'uppercase' }}>
          {label}
        </Typography>
        {subtitle && (
          <Typography variant="caption" sx={{ color: 'hsl(var(--text-secondary) / 0.6)', fontSize: 10, display: 'block', mt: 0.3 }}>
            {subtitle}
          </Typography>
        )}
      </Box>
      <Box sx={{ width: 36, height: 36, borderRadius: '10px', display: 'flex', alignItems: 'center', justifyContent: 'center', background: `${color}1a` }}>
        <Icon size={18} style={{ color }} />
      </Box>
    </Box>
    <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 1 }}>
      <Typography variant="h4" sx={{ color: 'hsl(var(--text-primary))', fontWeight: 700, letterSpacing: '-0.02em', lineHeight: 1 }}>
        {value}
      </Typography>
      {trend && (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.3 }}>
          {trend.direction === 'up' ? (
            <TrendingUp size={16} color="#66BB80" />
          ) : trend.direction === 'down' ? (
            <TrendingDown size={16} color="#D16670" />
          ) : null}
          <Typography sx={{ fontSize: 11, fontWeight: 600, color: trend.direction === 'up' ? '#66BB80' : trend.direction === 'down' ? '#D16670' : 'hsl(var(--text-secondary))' }}>
            {Math.abs(trend.percent).toFixed(0)}% comparado a {trend.prevLabel}
          </Typography>
        </Box>
      )}
    </Box>
    <Box sx={{ height: 3, borderRadius: 2, background: `${color}` }} />
  </GlassPanel>
)

const DashboardPage: React.FC = () => {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'))
  const [stats, setStats] = useState<Stat | null>(null)
  const [ocorrencias, setOcorrencias] = useState<Ocorrencia[]>([])
  const [contatos, setContatos] = useState<Contato[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')

  useEffect(() => {
    Promise.all([getStats(), getAllOcorrencias(), getAllContatos()])
      .then(([s, o, c]) => {
        setStats(s)
        setOcorrencias(o || [])
        setContatos(c || [])
      })
      .catch(() => setError('Erro ao carregar estatísticas'))
      .finally(() => setLoading(false))
  }, [])

  const dateRange = useMemo(() => {
    const start = startDate ? new Date(`${startDate}T00:00:00`) : null
    const end = endDate ? new Date(`${endDate}T23:59:59`) : null
    return { start, end }
  }, [startDate, endDate])

  const filteredOcorrencias = useMemo(() => {
    if (!dateRange.start && !dateRange.end) return ocorrencias
    return ocorrencias.filter((o) => {
      if (!o.dataCriacao) return false
      const d = toUTCDate(o.dataCriacao)
      if (isNaN(d.getTime())) return false
      if (dateRange.start && d < dateRange.start) return false
      if (dateRange.end && d > dateRange.end) return false
      return true
    })
  }, [ocorrencias, dateRange])

  const filteredContatos = useMemo(() => {
    if (!dateRange.start && !dateRange.end) return contatos
    return contatos.filter((c) => {
      if (!c.data_criacao) return false
      const d = toUTCDate(c.data_criacao)
      if (isNaN(d.getTime())) return false
      if (dateRange.start && d < dateRange.start) return false
      if (dateRange.end && d > dateRange.end) return false
      return true
    })
  }, [contatos, dateRange])

  const computedStats = useMemo<Stat | null>(() => {
    if (!stats) return null
    const hasFilter = Boolean(startDate || endDate)
    if (!hasFilter) return stats

    const totalReclamacoes = filteredOcorrencias.length
    const totalIndicacoes = filteredOcorrencias.filter((o) => o.tipo === 'indicacao').length
    const totalRequerimentos = filteredOcorrencias.filter((o) => o.tipo === 'requerimento').length
    const percIndicacao = totalReclamacoes > 0 ? (totalIndicacoes / totalReclamacoes) * 100 : 0
    const percRequerimento = totalReclamacoes > 0 ? (totalRequerimentos / totalReclamacoes) * 100 : 0
    const reclamacoesTelefones = new Set(filteredOcorrencias.map((o) => o.telefone))
    const contatosTelefones = new Set(filteredContatos.map((c) => c.telefone))
    const allTelefones = new Set([...reclamacoesTelefones, ...contatosTelefones])
    const numPessoas = allTelefones.size

    const regiaoMap: Record<string, number> = {}
    for (const o of filteredOcorrencias) {
      const regiao = o.detalhes?.regiao || 'Sem Região Definida'
      regiaoMap[regiao] = (regiaoMap[regiao] || 0) + 1
    }
    const regioes = Object.entries(regiaoMap).map(([regiao, qtdRegiao]) => ({ regiao, qtdRegiao }))

    const categoriaMap: Record<string, number> = {}
    for (const o of filteredOcorrencias) {
      const cat = o.categoria || 'Sem Categoria'
      categoriaMap[cat] = (categoriaMap[cat] || 0) + 1
    }
    const categorias = Object.entries(categoriaMap).map(([categoria, qtdCategoria]) => ({ categoria, qtdCategoria }))

    return {
      ...stats,
      numReclamacoes: totalReclamacoes,
      totalIndicacoes,
      totalRequerimentos,
      percIndicacao,
      percRequerimento,
      numPessoas,
      regioes,
      categorias,
      indicacoesAprovadas: totalIndicacoes,
      requerimentosAprovados: totalRequerimentos,
    }
  }, [stats, filteredOcorrencias, filteredContatos, startDate, endDate])

  const trends = useMemo<{ reclamacoes: TrendInfo | null; indicacoes: TrendInfo | null; requerimentos: TrendInfo | null; pessoas: TrendInfo | null; reprovadas: TrendInfo | null }>(() => {
    const hasFilter = Boolean(startDate && endDate)
    if (!hasFilter) return { reclamacoes: null, indicacoes: null, requerimentos: null, pessoas: null, reprovadas: null }

    const start = new Date(`${startDate}T00:00:00`)
    const end = new Date(`${endDate}T23:59:59`)
    const durationMs = end.getTime() - start.getTime()
    const prevEnd = new Date(start.getTime() - 1)
    const prevStart = new Date(prevEnd.getTime() - durationMs)

    const prevOcorrencias = ocorrencias.filter((o) => {
      if (!o.dataCriacao) return false
      const d = toUTCDate(o.dataCriacao)
      if (isNaN(d.getTime())) return false
      return d >= prevStart && d <= prevEnd
    })

    const prevContatos = contatos.filter((c) => {
      if (!c.data_criacao) return false
      const d = toUTCDate(c.data_criacao)
      if (isNaN(d.getTime())) return false
      return d >= prevStart && d <= prevEnd
    })

    const fmtShort = (d: Date) => d.toLocaleDateString('pt-BR', { day: '2-digit', month: 'long', timeZone: 'UTC' })
    const prevLabel = `${fmtShort(prevStart)}–${fmtShort(prevEnd)}`

    const calcTrend = (current: number, previous: number): TrendInfo => {
      if (previous === 0) return { direction: current > 0 ? 'up' : 'neutral', percent: current > 0 ? 100 : 0, prevLabel }
      const pct = ((current - previous) / previous) * 100
      if (Math.abs(pct) < 1) return { direction: 'neutral', percent: 0, prevLabel }
      return { direction: pct > 0 ? 'up' : 'down', percent: Math.abs(pct), prevLabel }
    }

    const currReclamacoes = filteredOcorrencias.length
    const currIndicacoes = filteredOcorrencias.filter((o) => o.tipo === 'indicacao').length
    const currRequerimentos = filteredOcorrencias.filter((o) => o.tipo === 'requerimento').length
    const currPessoas = new Set([...filteredOcorrencias.map((o) => o.telefone), ...filteredContatos.map((c) => c.telefone)]).size
    const currReprovadas = filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'reprovado').length

    const prevReclamacoes = prevOcorrencias.length
    const prevIndicacoes = prevOcorrencias.filter((o) => o.tipo === 'indicacao').length
    const prevRequerimentos = prevOcorrencias.filter((o) => o.tipo === 'requerimento').length
    const prevPessoas = new Set([...prevOcorrencias.map((o) => o.telefone), ...prevContatos.map((c) => c.telefone)]).size
    const prevReprovadas = prevOcorrencias.filter((o) => o.status.toLowerCase() === 'reprovado').length

    return {
      reclamacoes: calcTrend(currReclamacoes, prevReclamacoes),
      indicacoes: calcTrend(currIndicacoes, prevIndicacoes),
      requerimentos: calcTrend(currRequerimentos, prevRequerimentos),
      pessoas: calcTrend(currPessoas, prevPessoas),
      reprovadas: calcTrend(currReprovadas, prevReprovadas),
    }
  }, [filteredOcorrencias, filteredContatos, ocorrencias, contatos, startDate, endDate])

  const kanbanData = (() => {
    const counts: Record<string, number> = {
      'Sem Tratativa': 0,
      'Em Análise': 0,
      'Aprovar como Requerimento': 0,
      'Aprovar como Indicação': 0,
      'Aprovar como Causa Animal': 0,
      'Desqualificado': 0,
      'Finalizado': 0,
    }
    for (const o of filteredOcorrencias) {
      const status = o.status.toLowerCase()
      if (status === 'criado') counts['Sem Tratativa']++
      else if (status === 'em análise' || status === 'em analise') counts['Em Análise']++
      else if (status === 'aprovado' && o.tipo === 'requerimento') counts['Aprovar como Requerimento']++
      else if (status === 'aprovado' && o.tipo === 'indicacao') counts['Aprovar como Indicação']++
      else if (status === 'aprovado' && o.tipo === 'causa animal') counts['Aprovar como Causa Animal']++
      else if (status === 'reprovado') counts['Desqualificado']++
      else if (status === 'finalizado') counts['Finalizado']++
    }
    return Object.entries(counts)
      .map(([name, value]) => ({ name, value }))
      .sort((a, b) => b.value - a.value)
  })()

  if (loading) {
    return <PageLoader message="Carregando estatísticas..." />
  }

  if (error) {
    return (
      <Box className="page-root"><Typography color="error">{error}</Typography></Box>
    )
  }

  return (
    <Box className="page-root">
      <PageHeader title="Estatísticas" />
      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          justifyContent: 'center',
          gap: 1.5,
          mb: 3,
          alignItems: 'center',
          p: 2,
          borderRadius: 2,
          bgcolor: 'hsl(var(--surface-2))',
          border: '1px solid hsl(var(--border))',
        }}
      >
        <TextField
          type="date"
          size="small"
          label="De"
          InputLabelProps={{ shrink: true }}
          inputProps={{ style: { textAlign: isMobile ? 'center' : 'left' }, autoComplete: 'off' }}
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
          sx={{ ...inputSx, minWidth: 150 }}
        />
        <TextField
          type="date"
          size="small"
          label="Até"
          InputLabelProps={{ shrink: true }}
          inputProps={{ style: { textAlign: isMobile ? 'center' : 'left' }, autoComplete: 'off' }}
          value={endDate}
          onChange={(e) => setEndDate(e.target.value)}
          sx={{ ...inputSx, minWidth: 150 }}
        />
        <Button
          size="small"
          onClick={() => { setStartDate(''); setEndDate('') }}
          sx={{ color: 'hsl(var(--text-secondary))', textTransform: 'none', borderRadius: 2, width: isMobile ? '100%' : undefined }}
        >
          Limpar Filtros
        </Button>
      </Box>
      {/* Cards por coluna do Kanban (7 colunas, agrupados em 3 colunas visuais) */}
      <Box sx={{ display: 'grid', gridTemplateColumns: isMobile ? '1fr' : 'repeat(3, 1fr)', gap: 2, mb: 2.5, alignItems: 'stretch' }}>
        {/* Coluna 1 */}
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <StatCard sx={{ flex: 1 }} label="Sem Tratativa" value={filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'criado').length} icon={Inbox} color="#A1A9B8" />
          <StatCard sx={{ flex: 1 }} label="Em Análise" value={filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'em análise' || o.status.toLowerCase() === 'em analise').length} icon={BarChart3} color="#62A1D8" />
        </Box>

        {/* Coluna 2 */}
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <StatCard sx={{ flex: 1 }} label="Aprovar como Requerimento" value={filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'aprovado' && o.tipo === 'requerimento').length} icon={FileText} color="#E89E70" />
          <StatCard sx={{ flex: 1 }} label="Aprovar como Indicação" value={filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'aprovado' && o.tipo === 'indicacao').length} icon={TrendingUp} color="#66BB80" />
          <StatCard sx={{ flex: 1 }} label="Aprovar como Causa Animal" value={filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'aprovado' && o.tipo === 'causa animal').length} icon={Heart} color="#B98CE8" />
        </Box>

        {/* Coluna 3 */}
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          <StatCard sx={{ flex: 1 }} label="Desqualificado" value={filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'reprovado').length} icon={XCircle} color="#D16670" />
          <StatCard sx={{ flex: 1 }} label="Finalizado" value={filteredOcorrencias.filter((o) => o.status.toLowerCase() === 'finalizado').length} icon={CheckCircle} color="#5B8FE0" />
        </Box>
      </Box>

      {/* Cards gerais */}
      <Box sx={{ display: 'flex', flexDirection: isMobile ? 'column' : 'row', gap: 2.5, alignItems: 'stretch', mb: 2.5 }}>
        <StatCard sx={{ flex: 1 }} label="Total de Solicitações" value={computedStats?.numReclamacoes ?? 0} icon={FileText} color="#41669C" trend={trends.reclamacoes} />
        <StatCard sx={{ flex: 1 }} label="Total de Pessoas Atendidas" value={computedStats?.numPessoas ?? 0} icon={Users} color="#62A1D8" trend={trends.pessoas} />
      </Box>

      {computedStats && (
        <Box mt={4}>
          <DistributionViewer stats={computedStats} kanbanData={kanbanData} ocorrencias={filteredOcorrencias} />
        </Box>
      )}
    </Box>
  )
}

export default DashboardPage
