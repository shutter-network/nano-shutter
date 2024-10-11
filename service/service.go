package service

import (
	"encoding/hex"
	"fmt"
	"nano-shutter/dkg"
	"nano-shutter/internal/error"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shutter-network/shutter/shlib/shcrypto"
)

type EncryptRequest struct {
	CypherText string `json:"cypher_text" binding:"required"`
}

type DecryptRequest struct {
	EncyptedMsg string `json:"encrypted_msg" binding:"required"`
	Identifier  string `json:"identifier" binding:"required"`
}

type Service struct {
	dkg.DKG
	CurrentEpochTimestamp int64
	EpochDelay            int64
}

func NewService(dkg dkg.DKG) Service {
	epochDelayStr := os.Getenv("EPOCH_DELAY")
	epochDelay, err := strconv.ParseInt(epochDelayStr, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("failed to get epoch delay: %w", err))
	}
	return Service{
		DKG:        dkg,
		EpochDelay: epochDelay,
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

	currentTimestamp := time.Now().UTC().Unix()

	if srv.CurrentEpochTimestamp == 0 || currentTimestamp-srv.CurrentEpochTimestamp >= srv.EpochDelay {
		srv.CurrentEpochTimestamp = currentTimestamp
	}
	epochStr := strconv.FormatInt(srv.CurrentEpochTimestamp, 16)
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
		"message":    hex.EncodeToString(encryptedMsg.Marshal()),
		"identifier": epochStr,
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

	timestamp, err := strconv.ParseInt(requestBody.Identifier, 16, 0)
	if err != nil {
		err := error.NewHttpError(
			"failed to marshal identifier",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}
	if time.Now().UTC().Unix()-timestamp <= srv.EpochDelay {
		err := error.NewHttpError(
			"epoch did not end",
			"too early decryption",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	iden, err := hex.DecodeString(requestBody.Identifier)
	if err != nil {
		err := error.NewHttpError(
			"error parsing identifier",
			"request body unmarshalling error",
			http.StatusBadRequest,
		)
		c.Error(err)
		return
	}

	//TODO: check for time elapsed
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
