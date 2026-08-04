import * as React from 'react'
import { useState, useEffect, useCallback } from 'react'
import {
  Box, Typography, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, Alert, IconButton,
} from '@mui/material'
import { Plus, Pencil, Trash2, Megaphone } from 'lucide-react'
import GlassPanel from '../components/GlassPanel'
import PageHeader from '../components/PageHeader'
import PageLoader from '../components/PageLoader'
import ConfirmDialog from '../components/ConfirmDialog'
import { inputSx, dialogSx } from '../utils/inputSx'
import {
  getAllCampanhas, createCampanha, updateCampanha, deleteCampanha,
} from '../services/campanhaService'
import type { Campanha } from '../types'

const CampanhasPage: React.FC = () => {
  const [campanhas, setCampanhas] = useState<Campanha[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [openModal, setOpenModal] = useState(false)
  const [editing, setEditing] = useState<Campanha | null>(null)
  const [palavraChave, setPalavraChave] = useState('')
  const [saving, setSaving] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [toDelete, setToDelete] = useState<Campanha | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getAllCampanhas()
      setCampanhas(Array.isArray(data) ? data : [])
    } catch (err: any) {
      setError(err.message || 'Erro ao carregar campanhas')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const openCreate = () => {
    setEditing(null)
    setPalavraChave('')
    setOpenModal(true)
  }

  const openEdit = (c: Campanha) => {
    setEditing(c)
    setPalavraChave(c.palavraChave)
    setOpenModal(true)
  }

  const save = async () => {
    const valor = palavraChave.trim()
    if (!valor) {
      setError('Informe a palavra-chave')
      return
    }
    setSaving(true)
    setError(null)
    try {
      if (editing) await updateCampanha(editing.id, { palavraChave: valor })
      else await createCampanha({ palavraChave: valor })
      setOpenModal(false)
      await load()
    } catch (err: any) {
      setError(err.message || 'Erro ao salvar campanha')
    } finally {
      setSaving(false)
    }
  }

  const confirmDelete = async () => {
    if (!toDelete) return
    setError(null)
    try {
      await deleteCampanha(toDelete.id)
      setConfirmOpen(false)
      setToDelete(null)
      await load()
    } catch (err: any) {
      setError(err.message || 'Erro ao excluir campanha')
      setConfirmOpen(false)
    }
  }

  const filtered = campanhas.filter((c) =>
    c.palavraChave.toLowerCase().includes(search.trim().toLowerCase())
  )

  if (loading) return <PageLoader message="Carregando campanhas..." />

  return (
    <Box className="page-root">
      <PageHeader title="Campanhas" />

      {error && (
        <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: 1.5,
          mb: 2,
          alignItems: 'center',
          justifyContent: 'center',
          p: 2,
          borderRadius: 2,
          bgcolor: 'hsl(var(--surface-2))',
          border: '1px solid hsl(var(--border))',
        }}
      >
        <TextField
          size="small"
          label="Buscar"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          variant="filled"
          sx={{ ...inputSx, minWidth: 220, flex: 1 }}
          inputProps={{ autoComplete: 'off' }}
        />
        <Button
          variant="contained"
          startIcon={<Plus size={16} />}
          onClick={openCreate}
          sx={{
            background: 'hsl(var(--primary))',
            textTransform: 'none',
            borderRadius: 2,
            boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
            whiteSpace: 'nowrap',
          }}
        >
          Nova campanha
        </Button>
      </Box>

      {filtered.length === 0 ? (
        <GlassPanel style={{ padding: 32, textAlign: 'center' }}>
          <Megaphone size={32} style={{ color: 'hsl(var(--text-secondary))', margin: '0 auto 12px' }} />
          <Typography sx={{ color: 'hsl(var(--text-secondary))', fontSize: 14 }}>
            {search ? 'Nenhuma campanha encontrada' : 'Nenhuma campanha cadastrada'}
          </Typography>
        </GlassPanel>
      ) : (
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
          {filtered.map((c) => (
            <GlassPanel key={c.id} style={{ padding: '14px 18px' }}>
              <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, minWidth: 0 }}>
                  <Box
                    sx={{
                      width: 36,
                      height: 36,
                      borderRadius: 2,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      bgcolor: 'hsl(var(--primary) / 0.15)',
                      color: 'hsl(var(--accent))',
                      flexShrink: 0,
                    }}
                  >
                    <Megaphone size={16} />
                  </Box>
                  <Box sx={{ minWidth: 0 }}>
                    <Typography
                      sx={{
                        color: 'hsl(var(--text-primary))',
                        fontWeight: 600,
                        fontSize: 14,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      }}
                    >
                      {c.palavraChave}
                    </Typography>
                    <Typography sx={{ color: 'hsl(var(--text-secondary))', fontSize: 12 }}>
                      ID {c.id}
                    </Typography>
                  </Box>
                </Box>
                <Box sx={{ display: 'flex', gap: 0.5, flexShrink: 0 }}>
                  <IconButton
                    size="small"
                    onClick={() => openEdit(c)}
                    sx={{ color: 'hsl(var(--accent))' }}
                    aria-label="Editar"
                  >
                    <Pencil size={16} />
                  </IconButton>
                  <IconButton
                    size="small"
                    onClick={() => { setToDelete(c); setConfirmOpen(true) }}
                    sx={{ color: 'hsl(var(--error))' }}
                    aria-label="Excluir"
                  >
                    <Trash2 size={16} />
                  </IconButton>
                </Box>
              </Box>
            </GlassPanel>
          ))}
        </Box>
      )}

      <Dialog open={openModal} onClose={() => setOpenModal(false)} maxWidth="sm" fullWidth PaperProps={{ sx: dialogSx }}>
        <DialogTitle sx={{ color: 'hsl(var(--accent))', fontWeight: 700 }}>
          {editing ? 'Editar campanha' : 'Nova campanha'}
        </DialogTitle>
        <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
          <TextField
            label="Palavra-chave"
            value={palavraChave}
            onChange={(e) => setPalavraChave(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter') save() }}
            variant="filled"
            sx={inputSx}
            autoFocus
            inputProps={{ autoComplete: 'off' }}
            helperText="Usada para identificar a origem do contato na primeira mensagem"
          />
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={() => setOpenModal(false)} sx={{ color: 'hsl(var(--text-secondary))', textTransform: 'none' }}>
            Cancelar
          </Button>
          <Button
            onClick={save}
            disabled={saving}
            variant="contained"
            sx={{ background: 'hsl(var(--primary))', textTransform: 'none', borderRadius: 2 }}
          >
            Salvar
          </Button>
        </DialogActions>
      </Dialog>

      <ConfirmDialog
        open={confirmOpen}
        title="Excluir campanha"
        message={`Confirmar exclusão de "${toDelete?.palavraChave}"?`}
        onConfirm={confirmDelete}
        onCancel={() => { setConfirmOpen(false); setToDelete(null) }}
        destructive
      />
    </Box>
  )
}

export default CampanhasPage
