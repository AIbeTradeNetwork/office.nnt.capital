package user_test

import (
	"context"
	"github.com/shopspring/decimal"
	"server/internal/domain"
	"server/internal/service/seed/data"
	"server/internal/service/user"
	"testing"
)

func TestService_RandSafeVariant(t *testing.T) {
	ctx := context.Background()

	//dbRepo, err := repository.NewDbRepo(ctx)
	//if err != nil {
	//	t.Fatal(err)
	//}

	results := make(map[int64]struct {
		Count int64
		Limit int64
		Type  domain.SafeVariantType
	})
	errors := 0
	sumAbt := decimal.Zero
	sumUsd := decimal.Zero

	userService := user.NewUserService(nil, nil, nil, nil, nil, nil)

	//seedService := seed.NewSeedService(dbRepo)
	//
	//seedService.SeedSafe(ctx)

	safe := data.Safes()[2]

	for i := 0; i < 10000; i++ {
		var err error
		variant, err := userService.RandSafeVariant(ctx, safe.Variants)
		if err != nil {
			errors++
			continue
		}
		_, ok := results[variant.Amount]
		if !ok {
			results[variant.Amount] = struct {
				Count int64
				Limit int64
				Type  domain.SafeVariantType
			}{Count: 0, Limit: variant.LimitCount, Type: variant.Type}
		}

		result := results[variant.Amount]
		result.Count += 1
		results[variant.Amount] = result

		if variant.Type == "abt" {
			sumAbt = sumAbt.Add(decimal.NewFromInt(variant.Amount).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(9)))))
		} else if variant.Type == "usd" {
			sumUsd = sumUsd.Add(decimal.NewFromInt(variant.Amount).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(2)))))
		}
	}

	t.Logf("errors: %v", errors)
	for amount, res := range results {
		if res.Type == "usd" {
			t.Logf("type: %v, amount: %v, count: %v, limit: %v", res.Type, amount, res.Count, res.Limit)
		}
	}
	for amount, res := range results {
		if res.Type == "abt" {
			t.Logf("type: %v, amount: %v, count: %v, limit: %v", res.Type, amount, res.Count, res.Limit)
		}
	}

	t.Logf("sum abt: %v", sumAbt)
	t.Logf("sum usd: %v", sumUsd)

	t.Logf("variants: %+v", safe.Variants)
}
