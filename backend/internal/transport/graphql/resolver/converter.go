package resolver

import (
	"server/internal/domain"
	"server/internal/transport/graphql/generated"
)

func userRankConvertToTransport(r *domain.UserRank) *generated.UserRank {
	return &generated.UserRank{
		Code:     r.RankCode,
		StartAt:  r.StartAt.UnixMilli(),
		EndAt:    r.EndAt.UnixMilli(),
		Priority: r.Priority,
	}
}
func userRankConvertToTransportList(ranks []*domain.UserRank) []*generated.UserRank {
	result := make([]*generated.UserRank, len(ranks))
	for i, r := range ranks {
		result[i] = userRankConvertToTransport(r)
	}
	return result
}
func userPlanConvertToTransport(p *domain.UserPlan) *generated.UserPlan {
	return &generated.UserPlan{
		Code:     p.PlanCode,
		StartAt:  p.StartAt.UnixMilli(),
		EndAt:    p.EndAt.UnixMilli(),
		Priority: p.Priority,
	}
}
func userPlanConvertToTransportList(plans []*domain.UserPlan) []*generated.UserPlan {
	result := make([]*generated.UserPlan, len(plans))
	for i, p := range plans {
		result[i] = userPlanConvertToTransport(p)
	}
	return result
}

func userProductConvertToTransport(p *domain.UserProduct) *generated.UserProduct {
	if p == nil {
		return nil
	}
	return &generated.UserProduct{
		Code:       p.ProductCode,
		Category:   p.ProductCategory,
		StartAt:    p.StartAt.UnixMilli(),
		EndAt:      p.EndAt.UnixMilli(),
		Priority:   p.Priority,
		Multiplier: p.Multiplier,
	}
}
func userProductConvertToTransportList(products []*domain.UserProduct) []*generated.UserProduct {
	result := make([]*generated.UserProduct, len(products))
	for i, p := range products {
		result[i] = userProductConvertToTransport(p)
	}
	return result
}

func userPlaceConvertToTransport(p *domain.UserPlace) *generated.UserPlace {
	return &generated.UserPlace{
		Row:       p.Row.String(),
		Col:       p.Col.String(),
		CreatedAt: p.CreatedAt.UnixMilli(),
		MatchUID:  p.MatchUID,
	}
}

func userActivityConvertToTransport(a *domain.UserActivity) *generated.UserActivity {
	return &generated.UserActivity{
		StartAt: a.StartAt.UnixMilli(),
		EndAt:   a.EndAt.UnixMilli(),
		Cv:      a.CvAmount,
	}
}

func teamTypeConvertToTransport(t domain.UserTeamType) generated.TeamType {
	switch t {
	case domain.UserTeamTypeRight:
		return generated.TeamTypeRight
	case domain.UserTeamTypeLeft:
		return generated.TeamTypeLeft
	}
	return generated.TeamTypeAuto
}

func userConfigConvertToTransport(c *domain.UserConfig) *generated.UserConfig {
	return &generated.UserConfig{
		TeamType: teamTypeConvertToTransport(c.TeamType),
	}
}

func teamUserConvertToTransport(u *domain.TeamUser) *generated.TeamUser {
	role := generated.UserRoleClient

	var plan *generated.UserPlan = nil
	var plans []*generated.UserPlan = nil
	if len(u.Plans) > 0 {
		plan = userPlanConvertToTransport(u.Plans[0])
		plans = userPlanConvertToTransportList(u.Plans)
	}

	var rank *generated.UserRank = nil
	var ranks []*generated.UserRank = nil
	if len(u.Ranks) > 0 {
		rank = userRankConvertToTransport(u.Ranks[0])
		ranks = userRankConvertToTransportList(u.Ranks)
	}

	var place *generated.UserPlace = nil
	if u.Place != nil {
		place = userPlaceConvertToTransport(u.Place)
		role = generated.UserRoleDistributor
	}

	var activity *generated.UserActivity = nil
	if u.Activity != nil {
		activity = userActivityConvertToTransport(u.Activity)
	}

	return &generated.TeamUser{
		UID:       u.User.UID,
		RefUID:    u.User.RefUID,
		Email:     u.User.Email,
		Nickname:  u.User.Nickname,
		CreatedAt: u.User.CreatedAt.UnixMilli(),
		Role:      role,
		Plan:      plan,
		Plans:     plans,
		Rank:      rank,
		Ranks:     ranks,
		Place:     place,
		Activity:  activity,
		TeamCount: u.TeamCount,
	}
}

func teamUserConvertToTransportList(users []*domain.TeamUser) []*generated.TeamUser {
	result := make([]*generated.TeamUser, len(users))
	for i, u := range users {
		result[i] = teamUserConvertToTransport(u)
	}
	return result
}

func friendUserConvertToTransport(u *domain.User) *generated.FriendUser {
	return &generated.FriendUser{
		UID:       u.UID,
		RefUID:    u.RefUID,
		Email:     u.Email,
		Nickname:  u.Nickname,
		CreatedAt: u.CreatedAt.UnixMilli(),
		TeamCount: u.TeamCount,
	}
}

func friendUserConvertToTransportList(users []*domain.User) []*generated.FriendUser {
	result := make([]*generated.FriendUser, len(users))
	for i, u := range users {
		result[i] = friendUserConvertToTransport(u)
	}
	return result
}

func buyTypeConvertToTransport(t domain.BuyType) generated.BuyType {
	return generated.BuyType(t)
}

func buyConvertToTransport(b *domain.Buy) *generated.Buy {
	return &generated.Buy{
		UID:          b.UID,
		UserUID:      b.UserUID,
		RefUID:       b.RefUID,
		MatchUID:     b.MatchUID,
		Row:          b.Row.String(),
		Col:          b.Col.String(),
		Type:         buyTypeConvertToTransport(b.Type),
		CreatedAt:    b.CreatedAt.UnixMilli(),
		PaidAt:       b.PaidAt.UnixMilli(),
		ApprovedAt:   b.ApprovedAt.UnixMilli(),
		ChargedAt:    b.ChargedAt.UnixMilli(),
		RefundedAt:   b.RefundedAt.UnixMilli(),
		PlanCode:     b.PlanCode,
		ProductCode:  b.ProductCode,
		CurrencyCode: b.CurrencyCode,
		Amount:       b.Amount,
		Cv:           b.Cv,
	}
}

func buyConvertToTransportList(buys []*domain.Buy) []*generated.Buy {
	result := make([]*generated.Buy, len(buys))
	for i, b := range buys {
		result[i] = buyConvertToTransport(b)
	}
	return result
}

func transactionTypeConvertToTransport(t domain.TransactionType) generated.TransactionType {
	return generated.TransactionType(t)
}

func transportMsgCodesToTransport(msgCodes []domain.TransactionMsgCode) []string {
	result := make([]string, len(msgCodes))
	for i, msgCode := range msgCodes {
		result[i] = string(msgCode)
	}
	return result
}

func transactionConvertToTransport(t *domain.Transaction) *generated.Transaction {
	return &generated.Transaction{
		UserUID:    t.UserUID,
		FromUID:    t.FromUID,
		Percent:    t.Percent,
		Level:      int64(t.Level),
		Type:       transactionTypeConvertToTransport(t.Type),
		RankCode:   t.RankCode,
		Amount:     t.Amount,
		PosAmount:  t.PosAmount,
		FullAmount: t.FullAmount,
		BuyUID:     t.BuyUID,
		PayoutUID:  t.PayoutUID,
		CreatedAt:  t.CreatedAt.UnixMilli(),
		ChargedAt:  t.ChargedAt.UnixMilli(),
		MsgCodes:   transportMsgCodesToTransport(t.MsgCodes),
	}
}

func transactionConvertToTransportList(transactions []*domain.Transaction) []*generated.Transaction {
	result := make([]*generated.Transaction, len(transactions))
	for i, t := range transactions {
		result[i] = transactionConvertToTransport(t)
	}
	return result
}

func payoutConvertToTransport(p *domain.Payout) *generated.Payout {
	return &generated.Payout{
		UID:           p.UID,
		Amount:        p.Amount,
		Commission:    p.Fee,
		CurrencyCode:  p.CurrencyCode,
		MethodCode:    p.MethodCode,
		AccountNumber: p.AccountNumber,
		AccountName:   p.AccountName,
		CreatedAt:     p.CreatedAt.UnixMilli(),
		ApprovedAt:    p.ApprovedAt.UnixMilli(),
		ChargedAt:     p.ChargedAt.UnixMilli(),
		CancelledAt:   p.CancelledAt.UnixMilli(),
		Reason:        p.Reason,
	}
}

func payoutConvertToTransportList(payouts []*domain.Payout) []*generated.Payout {
	result := make([]*generated.Payout, len(payouts))
	for i, p := range payouts {
		result[i] = payoutConvertToTransport(p)
	}
	return result
}

func configConvertToTransport(conf *domain.Config) *generated.Cfg {
	return &generated.Cfg{
		PayoutAmountMin:     conf.PayoutAmountMin,
		PayoutFeeMin:        conf.PayoutFeeMin,
		PayoutFeePercent:    conf.PayoutFeePercent,
		DistributorPrice:    conf.DistributorPrice,
		DefaultCurrencyCode: conf.DefaultCurrencyCode,
		DefaultRefUID:       conf.DefaultRefUid,
		CoinRefBonus:        conf.CoinRefBonus,
		CoinCode:            conf.CoinCode,
		UnlimInvite:         conf.UnlimInvite,
		TonWallet:           conf.TonWallet,
	}
}

func claimConvertToTransport(claim *domain.Claim) *generated.Claim {
	return &generated.Claim{
		Code:         claim.Code,
		MinPeriod:    int64(claim.MinPeriod),
		MaxPeriod:    int64(claim.MaxPeriod),
		Amount:       claim.Amount,
		CurrencyCode: claim.CurrencyCode,
		Precision:    int64(claim.Precision),
	}
}

func notificationConvertToTransportList(notifications []*domain.Notification) []*generated.Notification {
	result := make([]*generated.Notification, len(notifications))
	for i, n := range notifications {
		result[i] = notificationConvertToTransport(n)
	}
	return result
}

func notificationConvertToTransport(n *domain.Notification) *generated.Notification {
	return &generated.Notification{
		UID:       n.UID,
		ToUserUID: n.ToUserUID,
		Texts:     localeConvertToTransport(n.Texts),
		CreatedAt: n.CreatedAt.UnixMilli(),
	}
}

func localeConvertToTransport(l map[string]string) *generated.Locale {
	ru, _ := l["ru"]
	en, _ := l["en"]

	return &generated.Locale{
		Ru: ru,
		En: en,
	}
}

func userLevelConvertToTransport(l *domain.UserLevel) *generated.UserLevel {
	return &generated.UserLevel{
		Code:        l.Code,
		Balance:     l.Balance,
		InviteLimit: l.InviteLimit,
	}
}

func userLevelConvertToTransportList(levels []*domain.UserLevel) []*generated.UserLevel {
	result := make([]*generated.UserLevel, len(levels))
	for i, l := range levels {
		result[i] = userLevelConvertToTransport(l)
	}
	return result
}

func taskConvertToTransportList(tasks []*domain.Task) []*generated.Task {
	result := make([]*generated.Task, len(tasks))
	for i, t := range tasks {
		result[i] = taskConvertToTransport(t)
	}
	return result
}

func taskConvertToTransport(t *domain.Task) *generated.Task {
	return &generated.Task{
		Code:         t.Code,
		Texts:        localeConvertToTransport(t.Texts),
		CurrencyCode: t.CurrencyCode,
		Precision:    int64(t.Precision),
		Amount:       t.Amount,
		Link:         t.Link,
		Completed:    t.Completed,
		IsApprove:    t.IsApprove,
	}
}

func comboConvertToTransport(c *domain.Combo) *generated.Combo {
	return &generated.Combo{
		UID:          c.UID,
		CurrencyCode: c.CurrencyCode,
		Precision:    int64(c.Precision),
		Amount:       c.Amount,
		PriseCode:    c.PriseCode,
		Limit:        int64(c.Limit),
		Count:        int64(c.Count),
	}
}

func comboConvertToTransportList(combos []*domain.Combo) []*generated.Combo {
	result := make([]*generated.Combo, len(combos))
	for i, c := range combos {
		result[i] = comboConvertToTransport(c)
	}
	return result
}
