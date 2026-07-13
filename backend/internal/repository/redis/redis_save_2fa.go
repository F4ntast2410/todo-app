package cache

import (
	"context"
	"fmt"
	"proj/internal/entity"
	"time"
)

func (c RedisVerificationRepository) Save2FA(ctx context.Context, email entity.VerificationEmail, unique_code string, time_duration time.Duration) error {
	err := c.Client.Set(ctx, fmt.Sprintf("2fa:email:%s", email), unique_code, time_duration).Err()
	if err != nil {
		return err
	}
	return nil
}
