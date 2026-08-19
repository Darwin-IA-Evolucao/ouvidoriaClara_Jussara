import * as React from 'react'
import { Box, Dialog, DialogContent, IconButton, Typography } from '@mui/material'
import { Brain, Play, Pause, X, FileText, Expand, Download } from 'lucide-react'
import type { MensagemChat } from '../types'

interface ChatBubbleProps {
  message: MensagemChat
}

const formatTime = (iso: string): string => {
  try {
    const d = new Date(iso)
    return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
  } catch {
    return ''
  }
}

const formatDuration = (s: number): string => {
  if (!Number.isFinite(s) || s < 0) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

// Waveform — mesmas alturas do chat-proprio-darwin (AudioMessage)
const WAVEFORM_HEIGHTS = [30, 52, 46, 74, 62, 82, 48, 58, 72, 88, 64, 54, 44, 60, 76, 56, 50, 66, 57, 42]

// Player de áudio customizado — replica AudioMessage do chat-proprio-darwin:
// botão circular + waveform (barras que trocam de cor no ponto do progresso) + tempo à direita da waveform
const AudioPlayer: React.FC<{ src: string; isCliente: boolean }> = ({ src, isCliente }) => {
  const audioRef = React.useRef<HTMLAudioElement | null>(null)
  const [playing, setPlaying] = React.useState(false)
  const [current, setCurrent] = React.useState(0)
  const [duration, setDuration] = React.useState(0)

  const togglePlay = () => {
    const a = audioRef.current
    if (!a) return
    if (playing) {
      a.pause()
      return
    }
    if (a.currentTime >= (a.duration || 0) - 0.1) a.currentTime = 0
    const p = a.play()
    if (p !== undefined) p.catch(() => {})
  }

  const progress = duration > 0 ? (current / duration) * 100 : 0
  // Cores: replicam Darwin — bolha escura(ia/agente)="bg-white/80" tocado, "bg-white/40" não tocado
  // bolha clara(cliente)="#128C7E" tocado, "gray-500" não tocado
  const playedColor = isCliente ? '#128C7E' : 'rgba(255,255,255,0.85)'
  const unplayedColor = isCliente ? 'hsl(var(--text-secondary))' : 'rgba(255,255,255,0.4)'
  const timeColor = isCliente ? 'hsl(var(--text-secondary))' : 'rgba(255,255,255,0.9)'

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, py: 0.75, px: 0.25, minWidth: 200, maxWidth: 240 }}>
      <audio
        ref={audioRef}
        src={src}
        preload="metadata"
        onLoadedMetadata={(e) => setDuration(e.currentTarget.duration || 0)}
        onTimeUpdate={(e) => setCurrent(e.currentTarget.currentTime)}
        onEnded={() => { setPlaying(false); setCurrent(0) }}
        onPlay={() => setPlaying(true)}
        onPause={() => setPlaying(false)}
        style={{ display: 'none' }}
      />

      {/* Botão play/pause — círculo verde WhatsApp, igual ao Darwin */}
      <Box
        component="button"
        onClick={togglePlay}
        sx={{
          flexShrink: 0,
          width: 36,
          height: 36,
          border: 'none',
          borderRadius: '50%',
          bgcolor: '#128C7E',
          cursor: 'pointer',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: '#fff',
          transition: 'background-color 0.15s',
          '&:hover': { bgcolor: '#0d7568' },
        }}
      >
        {playing ? (
          <Pause size={16} fill="currentColor" />
        ) : (
          <Play size={16} fill="currentColor" style={{ marginLeft: 2 }} />
        )}
      </Box>

      {/* Waveform + tempo — mesma linha, igual ao Darwin (bars flex-1, tempo min-w-[32px] à direita) */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flex: 1, minWidth: 0 }}>
        <Box
          onClick={(e: React.MouseEvent<HTMLDivElement>) => {
            const a = audioRef.current
            if (!a || !duration) return
            const rect = e.currentTarget.getBoundingClientRect()
            const ratio = Math.min(1, Math.max(0, (e.clientX - rect.left) / rect.width))
            a.currentTime = ratio * duration
          }}
          sx={{
            display: 'flex',
            alignItems: 'flex-end',
            gap: '2px',
            height: 12,
            flex: 1,
            minWidth: 0,
            cursor: 'pointer',
          }}
        >
          {WAVEFORM_HEIGHTS.map((h, i) => {
            const barProgress = (i / WAVEFORM_HEIGHTS.length) * 100
            const isPlayed = progress >= barProgress
            return (
              <Box
                key={i}
                sx={{
                  width: 3,
                  height: `${h}%`,
                  borderRadius: 999,
                  bgcolor: isPlayed ? playedColor : unplayedColor,
                  flexShrink: 0,
                }}
              />
            )
          })}
        </Box>
        <Typography sx={{ fontSize: 12, fontWeight: 500, color: timeColor, minWidth: 32, textAlign: 'right' }}>
          {formatDuration(playing ? current : duration || current)}
        </Typography>
      </Box>
    </Box>
  )
}

// Frases que indicam mensagens automáticas do sistema que devem ser escondidas
const FRASES_ESCONDIDAS = [
  'A url foi salva em:',
  'transcrição do audio:',
  'transcrição do áudio:',
]

const shouldHideMessage = (conteudo: string): boolean => {
  if (!conteudo) return false
  const lower = conteudo.toLowerCase()
  return FRASES_ESCONDIDAS.some((frase) => lower.includes(frase.toLowerCase()))
}

/* Formatação estilo WhatsApp: *negrito*, _itálico_, ~riscado~, `código` */
const escapeHtml = (s: string): string =>
  s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

const formatWhatsAppText = (text: string): string => {
  let html = escapeHtml(text)
  // Ordem importa: código primeiro (não aninha), depois negrito, itálico, riscado
  html = html.replace(/`([^`]+)`/g, '<span class="wa-code">$1</span>')
  html = html.replace(/\*([^*]+)\*/g, '<span class="wa-bold">$1</span>')
  html = html.replace(/_([^_]+)_/g, '<span class="wa-italic">$1</span>')
  html = html.replace(/~([^~]+)~/g, '<span class="wa-strike">$1</span>')
  return html
}

const ChatBubbleBase: React.FC<ChatBubbleProps> = ({ message }) => {
  const [imgOpen, setImgOpen] = React.useState(false)
  const [videoOpen, setVideoOpen] = React.useState(false)
  const [docOpen, setDocOpen] = React.useState(false)
  const isCliente = message.remetente === 'cliente'
  const isAgente = message.remetente === 'agente'
  const isImagem = message.tipo === 'imagem'
  const isAudio = message.tipo === 'audio'
  const isVideo = message.tipo === 'video'
  const isDocumento = message.tipo === 'documento'
  const hasMidia = !!message.linkMidia
  const isMedia = (isImagem || isAudio || isVideo) && hasMidia

  // Nome do documento: usa conteudo se tiver, senão extrai da URL (linkMidia)
  const docNome = React.useMemo(() => {
    if (message.conteudo && message.conteudo.trim()) return message.conteudo.trim()
    try {
      const url = new URL(message.linkMidia)
      const parts = url.pathname.split('/')
      return decodeURIComponent(parts[parts.length - 1] || 'Documento')
    } catch {
      const parts = message.linkMidia.split('/').slice(-1)[0]
      return parts || 'Documento'
    }
  }, [message.conteudo, message.linkMidia])

  // Extensão do documento para escolher ícone/preview
  const docExt = React.useMemo(() => {
    const m = docNome.toLowerCase().match(/\.([a-z0-9]+)$/)
    return m ? m[1] : ''
  }, [docNome])

  const isPdf = docExt === 'pdf'

  // Esconde mensagens com frases do sistema ("A url foi salva em:", "transcrição do audio:")
  if (shouldHideMessage(message.conteudo)) return null

  // Esconde mensagens completamente vazias (sem texto e sem mídia) — essas apareciam como uma
  // bolha "fantasma" só com o timestamp, sem fundo visível, logo abaixo de mensagens de mídia
  if (!message.conteudo?.trim() && !hasMidia) return null

  // Cor de fundo: cliente=surface, Ju(ia)=primary, Assessor(agente)=accent
  const bubbleBg = isCliente
    ? 'hsl(var(--surface-2))'
    : isAgente
      ? 'hsl(var(--accent))'
      : 'hsl(var(--primary))'
  const bubbleColor = isCliente ? 'hsl(var(--text-primary))' : '#fff'

  // Esconde o texto quando tem imagem ou vídeo (a descrição vem junto no conteudo)
  const shouldHideText = (isImagem || isVideo) && hasMidia

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: isCliente ? 'flex-start' : 'flex-end',
        my: 0.5,
        px: 1,
      }}
    >
      <Box
        sx={{
          position: 'relative',
          maxWidth: '75%',
          bgcolor: bubbleBg,
          color: bubbleColor,
          borderRadius: isCliente ? '12px 12px 12px 4px' : '12px 12px 4px 12px',
          px: isMedia ? 0.5 : 1.5,
          py: isMedia ? 0.5 : 1,
          // Espaço reservado no fim para o timestamp não sobrepor o conteúdo (igual pb-5 do Darwin p/ texto)
          pb: isAudio ? 0.5 : isMedia ? 0.5 : 2.25,
          pr: isAudio ? 4 : isMedia ? 0.5 : 1.5,
          boxShadow: '0 1px 2px rgba(0,0,0,0.08)',
        }}
      >
        {/* Tag de remetente (ia vs agente) */}
        {!isCliente && (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mb: message.conteudo ? 0.5 : 0 }}>
            {message.remetente === 'ia' && <Brain size={11} color="rgba(255,255,255,0.7)" />}
            <Typography sx={{ fontSize: 11, opacity: 0.7, fontWeight: 600 }}>
              {message.remetente === 'ia' ? 'Ju' : 'Assessor'}
            </Typography>
          </Box>
        )}

        {/* Mídia (imagem) — tamanho igual ao chat-proprio-darwin: max 250x180, object-cover */}
        {isImagem && hasMidia && (
          <Box
            sx={{
              maxWidth: 250,
              maxHeight: 180,
              overflow: 'hidden',
              borderRadius: 2,
              position: 'relative',
              mb: shouldHideText ? 0 : 1,
            }}
          >
            <img
              src={message.linkMidia}
              alt="imagem"
              style={{
                width: '100%',
                height: 'auto',
                maxHeight: 180,
                objectFit: 'cover',
                cursor: 'pointer',
                display: 'block',
              }}
              onClick={() => setImgOpen(true)}
            />
            {/* Timestamp sobre a imagem — pill escuro semi-transparente, estilo WhatsApp */}
            <Box
              sx={{
                position: 'absolute',
                bottom: 6,
                right: 6,
                bgcolor: 'rgba(0,0,0,0.65)',
                borderRadius: 1.5,
                px: 1,
                py: 0.5,
                pointerEvents: 'none',
                zIndex: 2,
              }}
            >
              <Typography sx={{ fontSize: 12, color: '#fff', whiteSpace: 'nowrap', lineHeight: 1, fontWeight: 500 }}>
                {formatTime(message.criadoEm)}
              </Typography>
            </Box>
          </Box>
        )}

        {/* Mídia (audio) — player customizado estilo WhatsApp/Darwin */}
        {isAudio && hasMidia && <AudioPlayer src={message.linkMidia} isCliente={isCliente} />}

        {/* Mídia (vídeo) — player com botão de expandir que abre modal de zoom */}
        {isVideo && hasMidia && (
          <Box
            sx={{
              maxWidth: 250,
              maxHeight: 180,
              overflow: 'hidden',
              borderRadius: 2,
              position: 'relative',
              mb: shouldHideText ? 0 : 1,
            }}
          >
            <video
              src={message.linkMidia}
              controls
              style={{
                width: '100%',
                maxHeight: 180,
                objectFit: 'cover',
                display: 'block',
                borderRadius: 8,
              }}
            />
            {/* Botão expandir — abre o vídeo em tela cheia no modal */}
            <IconButton
              onClick={() => setVideoOpen(true)}
              sx={{
                position: 'absolute',
                top: 4,
                right: 4,
                bgcolor: 'rgba(0,0,0,0.55)',
                color: '#fff',
                '&:hover': { bgcolor: 'rgba(0,0,0,0.75)' },
                width: 26,
                height: 26,
                zIndex: 2,
              }}
            >
              <Expand size={14} />
            </IconButton>
          </Box>
        )}

        {/* Documento — card estilo WhatsApp com ícone, nome e extensão; click abre modal de preview */}
        {isDocumento && hasMidia && (
          <Box
            onClick={() => setDocOpen(true)}
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 1.5,
              cursor: 'pointer',
              bgcolor: isCliente ? 'hsl(var(--surface))' : 'rgba(255,255,255,0.15)',
              borderRadius: 1.5,
              px: 1.5,
              py: 1.25,
              mb: 1,
              maxWidth: 260,
              '&:hover': { bgcolor: isCliente ? 'hsl(var(--border))' : 'rgba(255,255,255,0.25)' },
            }}
          >
            {/* Ícone colorido por tipo de documento */}
            <Box
              sx={{
                flexShrink: 0,
                width: 40,
                height: 40,
                borderRadius: 1,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                bgcolor: isPdf ? '#dc2626' : '#2563eb',
                color: '#fff',
              }}
            >
              <FileText size={22} />
            </Box>
            <Box sx={{ minWidth: 0, flex: 1 }}>
              <Typography
                sx={{
                  fontSize: 14,
                  fontWeight: 500,
                  color: bubbleColor,
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                }}
              >
                {docNome}
              </Typography>
              <Typography sx={{ fontSize: 12, color: isCliente ? 'hsl(var(--text-secondary))' : 'rgba(255,255,255,0.6)', textTransform: 'uppercase' }}>
                {docExt || 'arquivo'}
              </Typography>
            </Box>
          </Box>
        )}

        {/* Conteúdo textual — escondido quando tem imagem, vídeo ou áudio */}
        {message.conteudo && !shouldHideText && !isAudio && !isDocumento && (
          <Typography
            component="div"
            sx={{
              fontSize: 15,
              lineHeight: 1.45,
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-word',
              '& .wa-bold': { fontWeight: 700 },
              '& .wa-italic': { fontStyle: 'italic' },
              '& .wa-strike': { textDecoration: 'line-through' },
              '& .wa-code': { fontFamily: 'monospace', fontSize: 13.5, bgcolor: 'rgba(0,0,0,0.08)', px: 0.3, borderRadius: 0.5 },
            }}
            dangerouslySetInnerHTML={{ __html: formatWhatsAppText(message.conteudo) }}
          />
        )}

        {/* Timestamp — posicionado absoluto no canto inferior direito, igual ao MessageTime do Darwin
            (para imagem, o timestamp já está dentro do box da imagem com pill escuro) */}
        {!isImagem && (
          <Typography
            sx={{
              position: 'absolute',
              bottom: 4,
              right: 8,
              fontSize: 11,
              color: isCliente ? 'hsl(var(--text-secondary))' : 'rgba(255,255,255,0.7)',
              whiteSpace: 'nowrap',
              lineHeight: 1,
            }}
          >
            {formatTime(message.criadoEm)}
          </Typography>
        )}
      </Box>

      {/* Modal de imagem — igual ao chat-proprio-darwin, com botão X para fechar */}
      {isImagem && hasMidia && (
        <Dialog
          open={imgOpen}
          onClose={() => setImgOpen(false)}
          maxWidth={false}
          PaperProps={{
            sx: {
              maxWidth: '95vw',
              width: 'auto',
              p: 1,
              m: 0,
              bgcolor: 'transparent',
              boxShadow: 'none',
              overflow: 'visible',
            },
          }}
          BackdropProps={{
            sx: { bgcolor: 'rgba(0,0,0,0.85)' },
          }}
        >
          <DialogContent
            sx={{
              p: 0,
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              overflow: 'visible',
              position: 'relative',
            }}
          >
            <Box sx={{ position: 'relative', display: 'inline-block' }}>
              <img
                src={message.linkMidia}
                alt="imagem"
                style={{
                  maxWidth: '95vw',
                  maxHeight: '80vh',
                  borderRadius: 8,
                  objectFit: 'contain',
                  display: 'block',
                }}
                onClick={() => setImgOpen(false)}
              />
              {/* Timestamp sobre a imagem expandida — pill escuro semi-transparente */}
              <Box
                sx={{
                  position: 'absolute',
                  bottom: 12,
                  right: 12,
                  bgcolor: 'rgba(0,0,0,0.65)',
                  borderRadius: 1.5,
                  px: 1.5,
                  py: 0.75,
                  pointerEvents: 'none',
                  zIndex: 2,
                }}
              >
                <Typography sx={{ fontSize: 13, color: '#fff', whiteSpace: 'nowrap', lineHeight: 1, fontWeight: 500 }}>
                  {formatTime(message.criadoEm)}
                </Typography>
              </Box>
            </Box>
            {/* Botão X no canto superior direito para fechar o modal */}
            <IconButton
              onClick={() => setImgOpen(false)}
              sx={{
                position: 'absolute',
                top: -8,
                right: -8,
                bgcolor: 'rgba(0,0,0,0.6)',
                color: '#fff',
                '&:hover': { bgcolor: 'rgba(0,0,0,0.8)' },
                width: 32,
                height: 32,
                zIndex: 3,
              }}
            >
              <X size={18} />
            </IconButton>
          </DialogContent>
        </Dialog>
      )}

      {/* Modal de vídeo — com botão X para fechar, player expandido */}
      {isVideo && hasMidia && (
        <Dialog
          open={videoOpen}
          onClose={() => setVideoOpen(false)}
          maxWidth={false}
          PaperProps={{
            sx: {
              maxWidth: '95vw',
              width: 'auto',
              p: 1,
              m: 0,
              bgcolor: 'transparent',
              boxShadow: 'none',
              overflow: 'visible',
            },
          }}
          BackdropProps={{
            sx: { bgcolor: 'rgba(0,0,0,0.85)' },
          }}
        >
          <DialogContent
            sx={{
              p: 0,
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              overflow: 'visible',
              position: 'relative',
            }}
          >
            <Box sx={{ position: 'relative', display: 'inline-block' }}>
              <video
                src={message.linkMidia}
                controls
                autoPlay
                style={{
                  maxWidth: '95vw',
                  maxHeight: '80vh',
                  borderRadius: 8,
                  display: 'block',
                }}
              />
            </Box>
            <IconButton
              onClick={() => setVideoOpen(false)}
              sx={{
                position: 'absolute',
                top: -8,
                right: -8,
                bgcolor: 'rgba(0,0,0,0.6)',
                color: '#fff',
                '&:hover': { bgcolor: 'rgba(0,0,0,0.8)' },
                width: 32,
                height: 32,
                zIndex: 3,
              }}
            >
              <X size={18} />
            </IconButton>
          </DialogContent>
        </Dialog>
      )}

      {/* Modal de documento — preview do PDF em iframe + botão de download */}
      {isDocumento && hasMidia && (
        <Dialog
          open={docOpen}
          onClose={() => setDocOpen(false)}
          maxWidth={false}
          PaperProps={{
            sx: {
              maxWidth: '90vw',
              width: 'auto',
              p: 1,
              m: 0,
              bgcolor: 'transparent',
              boxShadow: 'none',
              overflow: 'visible',
            },
          }}
          BackdropProps={{
            sx: { bgcolor: 'rgba(0,0,0,0.85)' },
          }}
        >
          <DialogContent
            sx={{
              p: 0,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              overflow: 'visible',
              position: 'relative',
              gap: 1,
            }}
          >
            {/* Cabeçalho com nome do arquivo + botão de download */}
            <Box
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                gap: 1,
                width: '100%',
                bgcolor: 'rgba(0,0,0,0.6)',
                borderRadius: 1,
                px: 1.5,
                py: 1,
              }}
            >
              <Typography sx={{ fontSize: 13, color: '#fff', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', flex: 1 }}>
                {docNome}
              </Typography>
              <IconButton
                component="a"
                href={message.linkMidia}
                target="_blank"
                rel="noopener noreferrer"
                sx={{ color: '#fff', '&:hover': { bgcolor: 'rgba(255,255,255,0.15)' } }}
                title="Abrir / baixar"
              >
                <Download size={18} />
              </IconButton>
            </Box>
            {/* Preview: iframe para PDF, card com botão de download para outros tipos */}
            {isPdf ? (
              <iframe
                src={message.linkMidia}
                title={docNome}
                style={{
                  width: '80vw',
                  height: '75vh',
                  border: 'none',
                  borderRadius: 8,
                  background: '#fff',
                }}
              />
            ) : (
              <Box
                sx={{
                  width: '80vw',
                  maxWidth: 600,
                  height: 200,
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: 2,
                  bgcolor: 'hsl(var(--surface))',
                  borderRadius: 2,
                }}
              >
                <FileText size={48} color="hsl(var(--text-secondary))" />
                <Typography sx={{ fontSize: 14, color: 'hsl(var(--text-primary))' }}>
                  Preview não disponível para .{docExt}
                </Typography>
                <Box
                  component="a"
                  href={message.linkMidia}
                  target="_blank"
                  rel="noopener noreferrer"
                  sx={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: 1,
                    textDecoration: 'none',
                    bgcolor: 'hsl(var(--accent))',
                    color: '#fff',
                    borderRadius: 1.5,
                    px: 2,
                    py: 1,
                    fontSize: 13,
                    fontWeight: 500,
                    '&:hover': { bgcolor: 'hsl(var(--accent-hover))' },
                  }}
                >
                  <Download size={16} />
                  Baixar arquivo
                </Box>
              </Box>
            )}
            {/* Botão X para fechar */}
            <IconButton
              onClick={() => setDocOpen(false)}
              sx={{
                position: 'absolute',
                top: -8,
                right: -8,
                bgcolor: 'rgba(0,0,0,0.6)',
                color: '#fff',
                '&:hover': { bgcolor: 'rgba(0,0,0,0.8)' },
                width: 32,
                height: 32,
                zIndex: 3,
              }}
            >
              <X size={18} />
            </IconButton>
          </DialogContent>
        </Dialog>
      )}
    </Box>
  )
}

// Memo com comparação por conteúdo — o polling recria os objetos de mensagem a cada
// intervalo, mas isso não deve resetar o estado (playing/current) do AudioPlayer
export default React.memo(ChatBubbleBase, (prev, next) => {
  const m1 = prev.message
  const m2 = next.message
  return (
    m1.id === m2.id &&
    m1.remetente === m2.remetente &&
    m1.conteudo === m2.conteudo &&
    m1.tipo === m2.tipo &&
    m1.linkMidia === m2.linkMidia &&
    m1.criadoEm === m2.criadoEm
  )
})
