package service

import (
	"encoding/hex"
	"nano-shutter/dkg"
	"nano-shutter/internal/error"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shutter-network/shutter/shlib/shcrypto"
)

type EncryptRequest struct {
	CypherText string `json:"cypher_text" binding:"required"`
	Timestamp  int64  `json:"timestamp" binding:"required"`
}

type DecryptRequest struct {
	EncyptedMsg string `json:"encrypted_msg" binding:"required"`
	Timestamp   int64  `json:"timestamp" binding:"required"`
}

type Service struct {
	dkg.DKG
}

func NewService(dkg dkg.DKG) Service {
	return Service{
		DKG: dkg,
	}
}

func (srv *Service) Encrypt(c *gin.Context) {
	var requestBody EncryptRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		err := error.NewHttpError(
			"request body not found",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	epochStr := strconv.FormatInt(requestBody.Timestamp, 16)
	epochbyte, err := hex.DecodeString(epochStr)
	if err != nil {
		err := error.NewHttpError(
			"could not decode epoch id into bytes",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	epochId := shcrypto.ComputeEpochID(epochbyte)

	timestampStr := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	var sigma [32]byte
	copy(sigma[:], []byte(timestampStr))

	encryptedMsg := shcrypto.Encrypt([]byte(requestBody.CypherText), srv.PublicKey, epochId, sigma)

	c.JSON(http.StatusOK, gin.H{
		"message": hex.EncodeToString(encryptedMsg.Marshal()),
	})
}

func (srv *Service) Decrypt(c *gin.Context) {
	var requestBody DecryptRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		err := error.NewHttpError(
			"request body not found",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	if time.Now().UTC().Unix()-requestBody.Timestamp < 0 {
		err := error.NewHttpError(
			"time has not elapsed",
			"too early decryption",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	epochStr := strconv.FormatInt(requestBody.Timestamp, 16)
	iden, err := hex.DecodeString(epochStr)
	if err != nil {
		err := error.NewHttpError(
			"error parsing identifier",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	epochId := shcrypto.ComputeEpochID(iden)

	secretKeyShare := shcrypto.ComputeEpochSecretKeyShare(srv.SecretKeyShare, epochId)
	secretKey, err := shcrypto.ComputeEpochSecretKey([]int{0}, []*shcrypto.EpochSecretKeyShare{secretKeyShare}, 1)
	if err != nil {
		err := error.NewHttpError(
			"secretKey not found",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	var encMsg shcrypto.EncryptedMessage

	msg, err := hex.DecodeString(requestBody.EncyptedMsg)
	if err != nil {
		err := error.NewHttpError(
			"failed to decode encrypted msg",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	if err := encMsg.Unmarshal(msg); err != nil {
		println(err)
		err := error.NewHttpError(
			"failed to unmarshal encrypted msg",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	decryptedMsg, err := encMsg.Decrypt(secretKey)
	if err != nil {
		err := error.NewHttpError(
			"failed to decrypt",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": string(decryptedMsg),
	})
}
