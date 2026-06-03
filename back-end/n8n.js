/* ===== CONFIG ===== */
const GROUP_SIZE = 2;   // nº de frases por bloco no plain-text

/* ---------- util 0a: remove <think>...</think> e tags órfãs ---------- */
function stripThink(text) {
  if (!text || typeof text !== 'string') return '';
  return text
    .replace(/<think\b[^>]*>[\s\S]*?<\/think\s*>/gi, '')
    .replace(/<\/?\s*think\b[^>]*>/gi, '')
    .trim();
}

/* ---------- util 0b: torna o texto seguro para interpolação em JSON ---- */
function safeForJSON(text) {
  if (!text || typeof text !== 'string') return '';
  return text
    .replace(/\r\n?/g, '\n')
    .replace(/\n/g, '\\n')
    .replace(/(\\n){3,}/g, '\\n\\n')
    .trim();
}

/* ---------- util 1: normalizar texto (merge de vários {"messages":...}) --- */
function normalizeMessages(raw, groupSize = GROUP_SIZE) {
  if (!raw || typeof raw !== 'string') return [];

  /* 0) ■ VERIFICA SE EXISTEM VÁRIOS BLOCOS {"messages":...} ---------------- */
  const jsonBlocks = raw.match(/({\s*"messages"\s*:\s*\[[\s\S]*?\]\s*})/g);  // g = todos
  if (jsonBlocks && jsonBlocks.length) {
    let merged = [];
    for (const block of jsonBlocks) {
      try {
        const parsed = JSON.parse(block);
        if (parsed.messages && Array.isArray(parsed.messages))
          merged = merged.concat(parsed.messages);
      } catch (_) { /* ignora bloco mal-formado */ }
    }
    if (merged.length) return merged;
    // se achou blocos mas não conseguiu parsear, continua ↓
  }

  /* 1) ■ UM ÚNICO JSON {"messages":[...]} --------------------------------- */
  try {
    const parsed = JSON.parse(raw);
    if (parsed && typeof parsed === 'object' && parsed.messages) {
      return parsed.messages;
    }
    return [raw.trim()];
  } catch (_) { /* não era JSON único */ }

  /* 2) ■ Enumeração 1. 2. 3. ---------------------------------------------- */
  const bullets = raw.trim().split(/\s*\d+\.\s+/).filter(Boolean);
  if (bullets.length > 1) return bullets;

  /* 3) ■ Plain-text → divide em frases e agrupa --------------------------- */
  const sentences = raw
    .trim()
    .split(/(?<=[.!?])\s+/)           // após . ! ?
    .map(s => s.trim())
    .filter(Boolean);

  const grouped = [];
  for (let i = 0; i < sentences.length; i += groupSize) {
    grouped.push(sentences.slice(i, i + groupSize).join(' '));
  }
  return grouped.length ? grouped : [raw.trim()];
}


/* ---------- util 2: evita itens vazios ---------- */
const fillEmpty = (arr, fallback) =>
  arr.map(v => (v && v.trim()) ? v : fallback);

/* ---------- coleta streaming ---------- */
const data         = $input.all();
const payloadLines = data[0]?.json?.data?.split('\n') || [];

let conversationID  = '';
let finalMessageRaw = '';          // vai acumular!
let agentThoughts   = [];

for (const line of payloadLines) {
  const clean = line.replace(/^data:\s*/, '').trim();
  if (!clean || clean === '[DONE]') continue;

  try {
    const evt = JSON.parse(clean);

    if (evt.event === 'agent_thought') {
      if (evt.thought) agentThoughts.push(stripThink(evt.thought));
      conversationID = evt.conversation_id || conversationID;

    } else if (evt.event === 'agent_message') {
      // --- CONCATE N A  ---
      if (typeof evt.answer === 'string') {
        finalMessageRaw += evt.answer;
      }
    }

  } catch {
    console.log('Linha SSE não-JSON ignorada:', clean);
  }
}

/* ---------- blindagem ---------- */
const cleanRaw      = stripThink(finalMessageRaw);
const lastThought   = stripThink(agentThoughts.slice(-1)[0] || '') || 'Sem pensamento disponível';
const safeThoughts  = fillEmpty(agentThoughts, lastThought);

let normalized      = normalizeMessages(cleanRaw)
                        .map(safeForJSON)
                        .filter(Boolean);
if (!normalized.length) {
  normalized = ['Desculpe, não consegui processar sua mensagem agora. Por favor tente novamente em alguns minutos'];
}
const finalMessage  = JSON.stringify({ messages: normalized });

/* ---------- saída ---------- */
return [
  {
    json: {
      event:           'final_output',
      conversation_id: conversationID,
      agent_thoughts:  safeThoughts,
      final_message:   finalMessage      // string JSON segura
    }
  }
];