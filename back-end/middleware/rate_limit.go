package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitante struct {
	limiter  *rate.Limiter
	ultimoAt time.Time
}

type RateLimiter struct {
	mu         sync.Mutex
	visitantes map[string]*visitante
	taxa       rate.Limit
	burst      int
	ocioso     time.Duration
}

func NewRateLimiter(taxa rate.Limit, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitantes: make(map[string]*visitante),
		taxa:       taxa,
		burst:      burst,
		ocioso:     10 * time.Minute,
	}
	go rl.limpar()
	return rl
}

func (rl *RateLimiter) Permite(chave string) bool {
	rl.mu.Lock()
	v, ok := rl.visitantes[chave]
	if !ok {
		v = &visitante{limiter: rate.NewLimiter(rl.taxa, rl.burst)}
		rl.visitantes[chave] = v
	}
	v.ultimoAt = time.Now()
	rl.mu.Unlock()
	return v.limiter.Allow()
}

// limpar remove visitantes ociosos para o mapa nao crescer indefinidamente sob ataque.
func (rl *RateLimiter) limpar() {
	for {
		time.Sleep(rl.ocioso)
		rl.mu.Lock()
		for chave, v := range rl.visitantes {
			if time.Since(v.ultimoAt) > rl.ocioso {
				delete(rl.visitantes, chave)
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimitPorIP(rl *RateLimiter, retryAfter int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Permite(c.ClientIP()) {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "muitas requisicoes"})
			return
		}
		c.Next()
	}
}

// LimitarBodyPorPrefixo precisa ser registrado antes de qualquer middleware que leia
// o corpo da requisicao, senao o body ja foi carregado inteiro em memoria.
func LimitarBodyPorPrefixo(prefixo string, limite int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, prefixo) {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limite)
		}
		c.Next()
	}
}
