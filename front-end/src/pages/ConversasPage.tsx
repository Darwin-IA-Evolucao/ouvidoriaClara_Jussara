import * as React from 'react'
import { useState, useEffect, useCallback, useMemo, useRef } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Box, useMediaQuery, useTheme } from '@mui/material'
import PageLoader from '../components/PageLoader'
import ConversationList from '../components/ConversationList'
import ConversationChat from '../components/ConversationChat'
import { listConversas } from '../services/difyService'
import { getAllClientes } from '../services/clienteService'
import type { ConversaResumo, Cliente } from '../types'

const LIST_POLL_INTERVAL = 15000 // 15s — atualiza a lista silenciosamente

const ConversasPage: React.FC = () => {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))
  const [searchParams, setSearchParams] = useSearchParams()

  const [conversas, setConversas] = useState<ConversaResumo[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [mobileView, setMobileView] = useState<'list' | 'chat'>('list')
  const listAbortRef = useRef<AbortController | null>(null)
  const clientesMapRef = useRef<Map<string, string>>(new Map())

  const telFromUrl = searchParams.get('tel') || ''

  const load = useCallback(async (silent = false) => {
    if (!silent) setLoading(true)
    if (listAbortRef.current) listAbortRef.current.abort()
    const controller = new AbortController()
    listAbortRef.current = controller
    try {
      // Busca conversas e clientes em paralelo
      const [data, clientesRes] = await Promise.all([
        listConversas(),
        getAllClientes(0, 0).catch(() => ({ clientes: [] as Cliente[], total: 0 })),
      ])
      if (controller.signal.aborted) return

      // Mapa telefone → nome do cliente
      const clientesMap = new Map<string, string>()
      const clientesList = Array.isArray(clientesRes?.clientes) ? clientesRes.clientes : []
      clientesList.forEach((c) => {
        if (c.nome) clientesMap.set(c.telefone, c.nome)
      })
      clientesMapRef.current = clientesMap

      // Enriquece o nome: usa nome da tabela cliente se disponível,
      // senão mantém o nome que veio do backend (tabela contatos),
      // senão string vazia (o frontend mostra só o telefone).
      // O backend pode retornar "?" quando o nome não existe — tratamos como vazio.
      const isNomeValido = (n: string | undefined | null): boolean =>
        !!n && n.trim() !== '' && n.trim() !== '?' && n.trim() !== 'null'

      const enriched = (Array.isArray(data) ? data : []).map((c) => ({
        ...c,
        nome: isNomeValido(clientesMap.get(c.telefone))
          ? clientesMap.get(c.telefone)!
          : isNomeValido(c.nome)
            ? c.nome
            : '',
      }))

      setConversas(enriched)
    } catch (err) {
      if (controller.signal.aborted) return
      console.error('Erro ao carregar conversas:', err)
      if (!silent) setConversas([])
    } finally {
      if (!controller.signal.aborted && !silent) setLoading(false)
    }
  }, [])

  // Carga inicial
  useEffect(() => {
    load()
  }, [load])

  // Polling silencioso da lista a cada 15s (sem loading spinner)
  useEffect(() => {
    const interval = setInterval(() => {
      load(true)
    }, LIST_POLL_INTERVAL)
    return () => clearInterval(interval)
  }, [load])

  // Restaura a conversa selecionada a partir da URL (?tel=...)
  // depois que a lista carrega. Sobrevive ao F5.
  const selected = useMemo<ConversaResumo | null>(() => {
    if (!telFromUrl || conversas.length === 0) return null
    return conversas.find((c) => c.telefone === telFromUrl) ?? null
  }, [telFromUrl, conversas])

  // Se a URL tem ?tel=... e estamos no mobile, começa na view de chat
  useEffect(() => {
    if (isMobile && telFromUrl && conversas.length > 0) {
      setMobileView('chat')
    }
  }, [isMobile, telFromUrl, conversas.length])

  const handleSelect = useCallback(
    (conversa: ConversaResumo) => {
      setSearchParams({ tel: conversa.telefone }, { replace: true })
      if (isMobile) setMobileView('chat')
    },
    [isMobile, setSearchParams]
  )

  const handleBack = useCallback(() => {
    setMobileView('list')
  }, [])

  // Recarrega a lista silenciosamente quando uma mensagem é enviada
  // para que a última mensagem da conversa atualize no painel esquerdo
  const handleMessageSent = useCallback(() => {
    load(true)
  }, [load])

  if (loading) {
    return <PageLoader />
  }

  // Mobile: mostra lista OU chat (não os dois ao mesmo tempo)
  if (isMobile) {
    return (
      <Box
        sx={{
          height: 'calc(100vh - 64px)',
          display: 'flex',
          flexDirection: 'column',
          bgcolor: 'hsl(var(--background))',
        }}
      >
        {mobileView === 'list' ? (
          <ConversationList
            conversas={conversas}
            selected={selected}
            onSelect={handleSelect}
            search={search}
            onSearchChange={setSearch}
            loading={loading}
          />
        ) : (
          <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', minHeight: 0 }}>
            <Box
              onClick={handleBack}
              sx={{
                px: 2,
                py: 1,
                bgcolor: 'hsl(var(--surface))',
                borderBottom: '1px solid hsl(var(--border))',
                cursor: 'pointer',
                fontSize: 13,
                color: 'hsl(var(--accent))',
                fontWeight: 600,
              }}
            >
              ← Voltar para lista
            </Box>
            <ConversationChat
              conversa={selected}
              onMessageSent={handleMessageSent}
              clientesMap={clientesMapRef.current}
            />
          </Box>
        )}
      </Box>
    )
  }

  // Desktop: lista + chat lado a lado
  return (
    <Box
      sx={{
        height: 'calc(100vh - 64px)',
        display: 'flex',
        bgcolor: 'hsl(var(--background))',
        overflow: 'hidden',
      }}
    >
      <ConversationList
        conversas={conversas}
        selected={selected}
        onSelect={handleSelect}
        search={search}
        onSearchChange={setSearch}
        loading={loading}
      />
      <ConversationChat
        conversa={selected}
        onMessageSent={handleMessageSent}
        clientesMap={clientesMapRef.current}
      />
    </Box>
  )
}

export default ConversasPage
