package dkg

import (
	"fmt"

	"github.com/shutter-network/shutter/shlib/puredkg"
)

type DKG struct {
	puredkg.Result
}

func StartDkg() DKG {

	pureDkg := puredkg.NewPureDKG(1, 1, 1, 0)

	polyCommitmentMsg, polyEvalMsg, err := pureDkg.StartPhase1Dealing()
	if err != nil {
		panic(fmt.Sprintf("error in phase 1 dealing: %w", err))
	}

	if err := pureDkg.HandlePolyCommitmentMsg(polyCommitmentMsg); err != nil {
		panic(fmt.Sprintf("error in HandlePolyCommitmentMsg: %w", err))
	}

	for _, msg := range polyEvalMsg {
		if err := pureDkg.HandlePolyEvalMsg(msg); err != nil {
			panic(fmt.Sprintf("error in HandlePolyEvalMsg: %w", err))
		}
	}

	accusationMsg := pureDkg.StartPhase2Accusing()

	for _, accusation := range accusationMsg {
		if err := pureDkg.HandleAccusationMsg(accusation); err != nil {
			panic(fmt.Sprintf("error in HandleAccusationMsg: %w", err))
		}
	}

	apologies := pureDkg.StartPhase3Apologizing()

	for _, apology := range apologies {
		if err := pureDkg.HandleApologyMsg(apology); err != nil {
			panic(fmt.Sprintf("error in HandleApologyMsg: %w", err))
		}
	}

	pureDkg.Finalize()

	result, err := pureDkg.ComputeResult()
	if err != nil {
		panic(fmt.Sprintf("error in computing result: %w", err))
	}
	return DKG{result}
}
