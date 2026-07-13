package cache

import (
	"context"
	"fmt"
	"proj/internal/entity"
)

func (c RedisVerificationRepository) Delete2FA(ctx context.Context, email entity.VerificationEmail) error {
	err := c.Client.Del(ctx, fmt.Sprintf("2fa:email:%s", email)).Err()
	if err != nil {
		return err
	}
	return nil
}
