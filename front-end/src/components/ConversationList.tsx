import * as React from 'react'
import { Box, Typography, TextField, InputAdornment, CircularProgress } from '@mui/material'
import { Search, MessageCircle } from 'lucide-react'
import { formatPhoneDisplay } from '../utils/phone'
import { formatRelativeTime, formatDate, toUTCDate } from '../utils/date'
import type { ConversaResumo } from '../types'

const formatChatDate = (value: string): string => {
  if (!value) return ''
  const d = toUTCDate(value)
  if (isNaN(d.getTime())) return ''
  const now = new Date()
  const diffMs = now.getTime() - d.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  if (diffDays < 7) return formatRelativeTime(value)
  if (d.getUTCFullYear() === now.getUTCFullYear()) {
    return d.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short', timeZone: 'UTC' })
  }
  return formatDate(value)
}

interface ConversationListProps {
  conversas: ConversaResumo[]
  selected: ConversaResumo | null
  onSelect: (conversa: ConversaResumo) => void
  search: string
  onSearchChange: (value: string) => void
  loading: boolean
}

const ConversationList: React.FC<ConversationListProps> = ({
  conversas,
  selected,
  onSelect,
  search,
  onSearchChange,
  loading,
}) => {
  const filtered = React.useMemo(() => {
    const term = search.trim().toLowerCase()
    const digits = search.replace(/\D/g, '')
    if (!term) return conversas
    return conversas.filter((c) => {
      const matchName = c.nome.toLowerCase().includes(term)
      const matchPhone = digits.length > 0 && c.telefone.includes(digits)
      return matchName || matchPhone
    })
  }, [conversas, search])

  return (
    <Box
      sx={{
        width: 340,
        flexShrink: 0,
        borderRight: '1px solid hsl(var(--border))',
        bgcolor: 'hsl(var(--surface))',
        display: 'flex',
        flexDirection: 'column',
        minHeight: 0,
      }}
    >
      {/* Header */}
      <Box sx={{ p: 2, borderBottom: '1px solid hsl(var(--border))' }}>
        <Typography sx={{ fontSize: 16, fontWeight: 700, color: 'hsl(var(--text-primary))', mb: 1.5 }}>
          Conversas
        </Typography>
        <TextField
          fullWidth
          size="small"
          placeholder="Buscar contato..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          sx={{
            '& .MuiOutlinedInput-root': {
              bgcolor: 'hsl(var(--surface-2))',
              borderRadius: 2,
              '& fieldset': { borderColor: 'hsl(var(--border))' },
              '&:hover fieldset': { borderColor: 'hsl(var(--accent))' },
              '&.Mui-focused fieldset': { borderColor: 'hsl(var(--accent))' },
            },
            '& input': { color: 'hsl(var(--text-primary))', fontSize: 13 },
          }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search size={16} color="hsl(var(--text-secondary))" />
              </InputAdornment>
            ),
          }}
        />
      </Box>

      {/* Lista */}
      <Box sx={{ flex: 1, overflowY: 'auto', minHeight: 0 }}>
        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress size={24} sx={{ color: 'hsl(var(--accent))' }} />
          </Box>
        ) : filtered.length === 0 ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 6, px: 2 }}>
            <MessageCircle size={32} color="hsl(var(--text-secondary))" style={{ opacity: 0.4, marginBottom: 12 }} />
            <Typography sx={{ fontSize: 13, color: 'hsl(var(--text-secondary))', textAlign: 'center' }}>
              {conversas.length === 0
                ? 'Nenhuma conversa encontrada.'
                : 'Nenhum contato encontrado para a busca.'}
            </Typography>
          </Box>
        ) : (
          filtered.map((conversa) => {
            const isActive = selected?.telefone === conversa.telefone
            return (
              <Box
                key={conversa.telefone}
                onClick={() => onSelect(conversa)}
                sx={{
                  px: 2,
                  py: 1.5,
                  cursor: 'pointer',
                  borderBottom: '1px solid hsl(var(--border) / 0.5)',
                  bgcolor: isActive ? 'hsl(var(--primary) / 0.15)' : 'transparent',
                  borderLeft: isActive
                    ? '3px solid hsl(var(--accent))'
                    : '3px solid transparent',
                  transition: 'all 0.15s ease',
                  '&:hover': {
                    bgcolor: isActive ? 'hsl(var(--primary) / 0.15)' : 'hsl(var(--surface-2))',
                  },
                }}
              >
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', gap: 1 }}>
                  <Typography
                    sx={{
                      fontSize: 13.5,
                      fontWeight: isActive ? 600 : 500,
                      color: 'hsl(var(--text-primary))',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                      flex: 1,
                    }}
                  >
                    {conversa.nome || formatPhoneDisplay(conversa.telefone)}
                  </Typography>
                  {conversa.nome && (
                    <Typography
                      sx={{
                        fontSize: 11,
                        color: 'hsl(var(--text-secondary))',
                        whiteSpace: 'nowrap',
                        flexShrink: 0,
                      }}
                    >
                      {formatPhoneDisplay(conversa.telefone)}
                    </Typography>
                  )}
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline', gap: 1, mt: 0.25 }}>
                  <Typography
                    sx={{
                      fontSize: 11.5,
                      color: 'hsl(var(--text-secondary))',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                      flex: 1,
                    }}
                  >
                    {conversa.ultimaMensagem || formatPhoneDisplay(conversa.telefone)}
                  </Typography>
                  {conversa.criadoEm && (
                    <Typography
                      sx={{
                        fontSize: 10.5,
                        color: 'hsl(var(--text-secondary) / 0.7)',
                        whiteSpace: 'nowrap',
                        flexShrink: 0,
                      }}
                    >
                      {formatChatDate(conversa.criadoEm)}
                    </Typography>
                  )}
                </Box>
              </Box>
            )
          })
        )}
      </Box>
    </Box>
  )
}

export default ConversationList
