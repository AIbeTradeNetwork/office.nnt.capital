package ton

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tonkeeper/tongo/tonconnect"
)

func (r *Repo) GeneratePayload(ctx context.Context) (string, error) {
	payload, err := r.conTonGo.GeneratePayload()
	if err != nil {
		return "", err
	}

	return payload, nil
}

func (r *Repo) CheckProof(ctx context.Context, p string) (string, error) {
	var tp tonconnect.Proof

	payload := []byte(p)

	err := json.Unmarshal(payload, &tp)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	ok, _, err := r.conTonGo.CheckProof(ctx, &tp, func(p string) (bool, error) {
		return r.conTonGo.CheckPayload(p)
	}, func(d string) (bool, error) {
		return true, nil
	})
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("invalid proof")
	}

	return tp.Address, nil
}
