package facilserver

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	x402 "github.com/coinbase/x402/go"
	"github.com/gin-gonic/gin"
)

// Server wraps a Facilitator interface with HTTP handlers.
type Server struct {
	facilitator Facilitator
	logger      *slog.Logger
}

// New creates a new facilitator HTTP server.
// Accepts the Facilitator interface — works with *x402.X402Facilitator or any mock.
func New(facilitator Facilitator, logger *slog.Logger) *Server {
	return &Server{
		facilitator: facilitator,
		logger:      logger,
	}
}

type verifySettleRequest struct {
	X402Version         int             `json:"x402Version"`
	PaymentPayload      json.RawMessage `json:"paymentPayload"`
	PaymentRequirements json.RawMessage `json:"paymentRequirements"`
}

// HandleVerify handles POST /verify.
func (s *Server) HandleVerify(c *gin.Context) {
	req, err := s.parseRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log decoded payload for debugging
	var payloadMap, reqMap map[string]any
	json.Unmarshal(req.PaymentPayload, &payloadMap)
	json.Unmarshal(req.PaymentRequirements, &reqMap)

	s.logger.Info("verify request",
		"x402Version", req.X402Version,
		"paymentPayload", payloadMap,
		"paymentRequirements", reqMap,
	)

	resp, err := s.facilitator.Verify(c.Request.Context(), req.PaymentPayload, req.PaymentRequirements)
	if err != nil {
		s.logger.Error("verify failed", "error", err)
		if verifyErr, ok := err.(*x402.VerifyError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"isValid":        false,
				"invalidReason":  verifyErr.InvalidReason,
				"invalidMessage": verifyErr.InvalidMessage,
				"payer":          verifyErr.Payer,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("verify success", "isValid", resp.IsValid, "payer", resp.Payer)
	c.JSON(http.StatusOK, resp)
}

// HandleSettle handles POST /settle.
func (s *Server) HandleSettle(c *gin.Context) {
	req, err := s.parseRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payloadMap map[string]any
	json.Unmarshal(req.PaymentPayload, &payloadMap)

	s.logger.Info("settle request",
		"x402Version", req.X402Version,
		"paymentPayload", payloadMap,
	)

	resp, err := s.facilitator.Settle(c.Request.Context(), req.PaymentPayload, req.PaymentRequirements)
	if err != nil {
		s.logger.Error("settle failed", "error", err)
		if settleErr, ok := err.(*x402.SettleError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":      false,
				"errorReason":  settleErr.ErrorReason,
				"errorMessage": settleErr.ErrorMessage,
				"transaction":  settleErr.Transaction,
				"network":      settleErr.Network,
				"payer":        settleErr.Payer,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("settle success",
		"success", resp.Success,
		"transaction", resp.Transaction,
		"network", resp.Network,
		"payer", resp.Payer,
	)
	c.JSON(http.StatusOK, resp)
}

// HandleSupported handles GET /supported.
func (s *Server) HandleSupported(c *gin.Context) {
	resp := s.facilitator.GetSupported()
	if len(resp.Kinds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no supported payment kinds"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

const maxRequestBodySize = 1 << 20 // 1 MB

// parseRequest reads and decodes the JSON request body with a size limit.
func (s *Server) parseRequest(c *gin.Context) (*verifySettleRequest, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodySize)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, ErrReadBody
	}

	var req verifySettleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, ErrInvalidJSON
	}
	return &req, nil
}
