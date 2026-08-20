import * as React from 'react'
import { useState, useEffect, useMemo, useCallback } from 'react'
import {
  Box, Typography, TextField, Button, MenuItem, Select, FormControl,
  Dialog, DialogTitle, DialogContent, DialogActions, InputAdornment,
  IconButton, Table, TableBody, TableCell, TableContainer, TableHead,
  TableRow, Chip,
} from '@mui/material'
import { Search, Plus, Edit as EditIcon, Trash2, MapPin, ChevronFirst, ChevronLast, ChevronLeft, ChevronRight } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import ConfirmDialog from '../components/ConfirmDialog'
import PageLoader from '../components/PageLoader'
import { inputSx, dialogSx } from '../utils/inputSx'
import { getAllEnderecos, createEndereco, updateEndereco, deleteEndereco } from '../services/enderecoService'
import type { Logradouro } from '../types'

const REGIOES = ['Centro', 'Zona Sul', 'Zona Norte', 'Zona Leste', 'Zona Oeste', 'Zona Industrial']

const ITEMS_PER_PAGE = 12

const EnderecosPage: React.FC = () => {
  const [enderecos, setEnderecos] = useState<Logradouro[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [filtroRegiao, setFiltroRegiao] = useState('')
  const [page, setPage] = useState(0)
  const [pageInput, setPageInput] = useState('')

  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Logradouro | null>(null)
  const [form, setForm] = useState({ logradouro: '', bairro: '', regiao: '' })
  const [formError, setFormError] = useState('')

  const [confirmOpen, setConfirmOpen] = useState(false)
  const [toDelete, setToDelete] = useState<Logradouro | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getAllEnderecos()
      setEnderecos(data)
    } catch (err: any) {
      setError(err.message || 'Erro ao carregar endereços')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const regioesDisponiveis = useMemo(() => {
    const set = new Set<string>()
    enderecos.forEach((e) => { if (e.regiao) set.add(e.regiao) })
    return Array.from(set).sort()
  }, [enderecos])

  const filtered = useMemo(() => {
    const term = search.trim().toLowerCase()
    return enderecos.filter((e) => {
      if (filtroRegiao && e.regiao !== filtroRegiao) return false
      if (term) {
        const matchLog = e.logradouro?.toLowerCase().includes(term)
        const matchBairro = e.bairro?.toLowerCase().includes(term)
        const matchRegiao = e.regiao?.toLowerCase().includes(term)
        return matchLog || matchBairro || matchRegiao
      }
      return true
    })
  }, [enderecos, search, filtroRegiao])

  const paged = useMemo(() => {
    const start = page * ITEMS_PER_PAGE
    return filtered.slice(start, start + ITEMS_PER_PAGE)
  }, [filtered, page])

  const totalPages = Math.max(1, Math.ceil(filtered.length / ITEMS_PER_PAGE))

  const openCreate = () => {
    setEditing(null)
    setForm({ logradouro: '', bairro: '', regiao: '' })
    setFormError('')
    setModalOpen(true)
  }

  const openEdit = (e: Logradouro) => {
    setEditing(e)
    setForm({ logradouro: e.logradouro, bairro: e.bairro, regiao: e.regiao })
    setFormError('')
    setModalOpen(true)
  }

  const handleSave = async () => {
    if (!form.logradouro.trim() || !form.bairro.trim() || !form.regiao.trim()) {
      setFormError('Preencha todos os campos')
      return
    }
    try {
      if (editing && editing.id != null) {
        await updateEndereco(editing.id, form)
      } else {
        await createEndereco(form)
      }
      setModalOpen(false)
      await load()
    } catch (err: any) {
      setFormError(err.message || 'Erro ao salvar endereço')
    }
  }

  const handleDelete = async () => {
    if (!toDelete || toDelete.id == null) return
    try {
      await deleteEndereco(toDelete.id)
      setConfirmOpen(false)
      setToDelete(null)
      await load()
    } catch (err: any) {
      setError(err.message || 'Erro ao excluir endereço')
      setConfirmOpen(false)
    }
  }

  if (loading) return <PageLoader message="Carregando endereços..." />

  if (error) {
    return (
      <Box className="page-root">
        <PageHeader title="Endereços" />
        <Typography color="error" sx={{ textAlign: 'center' }}>{error}</Typography>
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
          <Button onClick={load} variant="contained" sx={{ textTransform: 'none' }}>Tentar novamente</Button>
        </Box>
      </Box>
    )
  }

  return (
    <Box className="page-root">
      <PageHeader title="Endereços" />

      {/* Toolbar */}
      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1.5, mb: 2, alignItems: 'center' }}>
        <TextField
          size="small"
          placeholder="Buscar por logradouro, bairro ou região..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(0) }}
          sx={{ ...inputSx, flex: 1, minWidth: 200 }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search size={16} style={{ color: 'hsl(var(--text-secondary))' }} />
              </InputAdornment>
            ),
          }}
        />
        <FormControl size="small" sx={{ minWidth: 150 }}>
          <Select
            value={filtroRegiao}
            onChange={(e) => { setFiltroRegiao(e.target.value); setPage(0) }}
            displayEmpty
            sx={{ fontSize: 12, ...inputSx, '& .MuiSelect-select': { py: 0.7 } }}
          >
            <MenuItem value="">Todas as regiões</MenuItem>
            {regioesDisponiveis.map((r) => (
              <MenuItem key={r} value={r}>{r}</MenuItem>
            ))}
          </Select>
        </FormControl>
        <Button
          variant="contained"
          startIcon={<Plus size={16} />}
          onClick={openCreate}
          sx={{
            textTransform: 'none', borderRadius: 2, fontWeight: 600, fontSize: 13,
            background: 'hsl(var(--primary))', whiteSpace: 'nowrap',
          }}
        >
          Novo Endereço
        </Button>
      </Box>

      {/* Table */}
      {filtered.length === 0 ? (
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', py: 6, gap: 1.5 }}>
          <MapPin size={40} style={{ color: 'hsl(var(--text-secondary) / 0.4)' }} />
          <Typography sx={{ color: 'hsl(var(--text-secondary))', fontWeight: 600, fontSize: 14 }}>
            Nenhum endereço encontrado
          </Typography>
          <Typography sx={{ color: 'hsl(var(--text-secondary) / 0.7)', fontSize: 12 }}>
            {search || filtroRegiao ? 'Tente ajustar os filtros' : 'Clique em "Novo Endereço" para começar'}
          </Typography>
        </Box>
      ) : (
        <>
          <TableContainer sx={{
            borderRadius: 2,
            border: '1px solid hsl(var(--border))',
            bgcolor: 'hsl(var(--surface-2) / 0.3)',
            overflow: 'auto',
          }}>
            <Table size="small">
              <TableHead>
                <TableRow sx={{ '& .MuiTableCell-head': { bgcolor: 'hsl(var(--surface-2))', fontWeight: 700, fontSize: 11, color: 'hsl(var(--text-secondary))', textTransform: 'uppercase', letterSpacing: '0.04em' } }}>
                  <TableCell>Logradouro</TableCell>
                  <TableCell>Bairro</TableCell>
                  <TableCell>Região</TableCell>
                  <TableCell align="right" sx={{ width: 100 }}>Ações</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {paged.map((e) => (
                  <TableRow
                    key={e.id}
                    sx={{
                      '&:hover': { bgcolor: 'hsl(var(--surface-2) / 0.5)' },
                      '& .MuiTableCell-body': { fontSize: 13, color: 'hsl(var(--text-primary))', py: 1.2 },
                    }}
                  >
                    <TableCell>{e.logradouro}</TableCell>
                    <TableCell>{e.bairro}</TableCell>
                    <TableCell>
                      <Chip
                        label={e.regiao}
                        size="small"
                        sx={{
                          fontSize: 11, fontWeight: 600, height: 22,
                          bgcolor: 'hsl(var(--accent) / 0.12)',
                          color: 'hsl(var(--accent))',
                          border: '1px solid hsl(var(--accent) / 0.2)',
                        }}
                      />
                    </TableCell>
                    <TableCell align="right">
                      <IconButton size="small" onClick={() => openEdit(e)} sx={{ color: 'hsl(var(--text-secondary))', '&:hover': { color: 'hsl(var(--accent))' } }}>
                        <EditIcon size={16} />
                      </IconButton>
                      <IconButton size="small" onClick={() => { setToDelete(e); setConfirmOpen(true) }} sx={{ color: 'hsl(var(--text-secondary))', '&:hover': { color: 'hsl(var(--error))' } }}>
                        <Trash2 size={16} />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>

          {/* Pagination */}
          <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', mt: 1.5, flexWrap: 'wrap', gap: 1 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <IconButton
                size="small"
                disabled={page === 0}
                onClick={() => setPage(0)}
                sx={{ color: 'hsl(var(--text-secondary))', '&:hover': { color: 'hsl(var(--accent))' }, '&.Mui-disabled': { color: 'hsl(var(--text-secondary) / 0.3)' } }}
              >
                <ChevronFirst size={18} />
              </IconButton>
              <IconButton
                size="small"
                disabled={page === 0}
                onClick={() => setPage((p) => Math.max(0, p - 1))}
                sx={{ color: 'hsl(var(--text-secondary))', '&:hover': { color: 'hsl(var(--accent))' }, '&.Mui-disabled': { color: 'hsl(var(--text-secondary) / 0.3)' } }}
              >
                <ChevronLeft size={18} />
              </IconButton>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mx: 0.5 }}>
                <TextField
                  size="small"
                  value={pageInput}
                  onChange={(e) => setPageInput(e.target.value.replace(/\D/g, ''))}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      const p = parseInt(pageInput, 10)
                      if (!isNaN(p) && p >= 1 && p <= totalPages) {
                        setPage(p - 1)
                      }
                      setPageInput('')
                    }
                  }}
                  placeholder={(page + 1).toString()}
                  sx={{ ...inputSx, width: 90, '& .MuiFilledInput-input': { fontSize: 12, textAlign: 'center', py: 0.5, px: 0.5 } }}
                  inputProps={{ autoComplete: 'off' }}
                />
                <Button
                  size="small"
                  onClick={() => {
                    const p = parseInt(pageInput, 10)
                    if (!isNaN(p) && p >= 1 && p <= totalPages) {
                      setPage(p - 1)
                    }
                    setPageInput('')
                  }}
                  sx={{ textTransform: 'none', fontSize: 12, minWidth: 40, px: 1.5, color: 'hsl(var(--accent))' }}
                >
                  Ir
                </Button>
                <Typography sx={{ fontSize: 12, color: 'hsl(var(--text-secondary))', mx: 0.5 }}>
                  de {totalPages}
                </Typography>
              </Box>
              <IconButton
                size="small"
                disabled={page >= totalPages - 1}
                onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
                sx={{ color: 'hsl(var(--text-secondary))', '&:hover': { color: 'hsl(var(--accent))' }, '&.Mui-disabled': { color: 'hsl(var(--text-secondary) / 0.3)' } }}
              >
                <ChevronRight size={18} />
              </IconButton>
              <IconButton
                size="small"
                disabled={page >= totalPages - 1}
                onClick={() => setPage(totalPages - 1)}
                sx={{ color: 'hsl(var(--text-secondary))', '&:hover': { color: 'hsl(var(--accent))' }, '&.Mui-disabled': { color: 'hsl(var(--text-secondary) / 0.3)' } }}
              >
                <ChevronLast size={18} />
              </IconButton>
            </Box>
          </Box>
        </>
      )}

      {/* Create/Edit Modal */}
      <Dialog open={modalOpen} onClose={() => setModalOpen(false)} maxWidth="sm" fullWidth PaperProps={{ sx: dialogSx }}>
        <DialogTitle sx={{ color: 'hsl(var(--accent))', fontWeight: 700, fontSize: 16 }}>
          {editing ? 'Editar Endereço' : 'Novo Endereço'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            <TextField
              label="Logradouro"
              fullWidth
              value={form.logradouro}
              onChange={(e) => setForm({ ...form, logradouro: e.target.value })}
              variant="filled"
              sx={inputSx}
              inputProps={{ autoComplete: 'off' }}
            />
            <TextField
              label="Bairro"
              fullWidth
              value={form.bairro}
              onChange={(e) => setForm({ ...form, bairro: e.target.value })}
              variant="filled"
              sx={inputSx}
              inputProps={{ autoComplete: 'off' }}
            />
            <TextField
              select
              label="Região"
              fullWidth
              value={form.regiao}
              onChange={(e) => setForm({ ...form, regiao: e.target.value })}
              variant="filled"
              sx={inputSx}
            >
              <MenuItem value="">Selecione...</MenuItem>
              {REGIOES.map((r) => (
                <MenuItem key={r} value={r}>{r}</MenuItem>
              ))}
            </TextField>
            {formError && (
              <Typography sx={{ fontSize: 12, color: 'hsl(var(--error))', fontWeight: 500 }}>
                {formError}
              </Typography>
            )}
          </Box>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={() => setModalOpen(false)} sx={{ color: 'hsl(var(--text-secondary))', textTransform: 'none' }}>Cancelar</Button>
          <Button onClick={handleSave} variant="contained" sx={{ background: 'hsl(var(--primary))', textTransform: 'none', borderRadius: 2 }}>
            {editing ? 'Salvar' : 'Criar'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation */}
      <ConfirmDialog
        open={confirmOpen}
        title="Excluir Endereço"
        message={`Tem certeza que deseja excluir "${toDelete?.logradouro}" (${toDelete?.bairro})?`}
        confirmLabel="Excluir"
        destructive
        onConfirm={handleDelete}
        onCancel={() => { setConfirmOpen(false); setToDelete(null) }}
      />
    </Box>
  )
}

export default EnderecosPage
