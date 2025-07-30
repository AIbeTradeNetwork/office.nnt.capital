package team

import (
	"go.mongodb.org/mongo-driver/bson"
	"math/big"
	"server/internal/domain"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
	two  = big.NewInt(2)
)

func (s *Service) getUserTeamTypeByPlace(place *domain.UserPlace) domain.UserTeamType {
	teamType := domain.UserTeamTypeRight
	q := big.NewInt(0)
	r := big.NewInt(0)
	q.DivMod(place.Col, big.NewInt(2), r)
	if r.Int64() == 1 {
		teamType = domain.UserTeamTypeLeft
	}
	return teamType
}

func (s *Service) getOppositeUserTeamTypeByPlace(place *domain.UserPlace) domain.UserTeamType {
	teamType := domain.UserTeamTypeLeft
	q := big.NewInt(0)
	r := big.NewInt(0)
	q.DivMod(place.Col, big.NewInt(2), r)
	if r.Int64() == 1 {
		teamType = domain.UserTeamTypeRight
	}
	return teamType
}

// get start col for row in team
func (s *Service) getFromForRow(place *domain.UserPlace, row *big.Int, side domain.UserTeamType) *big.Int {
	a := new(big.Int)
	b := new(big.Int)
	c := new(big.Int)
	from := new(big.Int)
	a.Sub(place.Col, one)
	b.Sub(row, place.Row)
	c.Exp(lineCap, b, nil)
	from.Mul(a, c)
	if side == domain.UserTeamTypeRight {
		from.Add(from, s.getGapForRow(place, row, side))
	}
	return from
}

// get end col for row in team
func (s *Service) getToForRow(place *domain.UserPlace, row *big.Int, side domain.UserTeamType) *big.Int {
	a := new(big.Int)
	b := new(big.Int)
	to := new(big.Int)
	a.Sub(row, place.Row)
	b.Exp(lineCap, a, nil)
	to.Mul(place.Col, b)
	if side == domain.UserTeamTypeLeft {
		to.Sub(to, s.getGapForRow(place, row, side))
	}
	return to
}

// get gap between start and end col for row in team
func (s *Service) getGapForRow(place *domain.UserPlace, row *big.Int, side domain.UserTeamType) *big.Int {
	a := new(big.Int)
	gap := new(big.Int)
	a.Sub(row, place.Row)
	gap.Exp(lineCap, a, nil)
	if side == domain.UserTeamTypeLeft || side == domain.UserTeamTypeRight {
		gap.Div(gap, two)
	}
	return gap
}

// get all team position indexes for place
func (s *Service) getAllPositionsDown(place *domain.UserPlace, side domain.UserTeamType, rows *big.Int) [][2]*big.Int {
	// TODO: add capacity
	positions := make([][2]*big.Int, 0)
	start := new(big.Int).Add(place.Row, one)
	end := new(big.Int).Add(place.Row, rows)
	for row := new(big.Int).Set(start); row.Cmp(end) < 0; row.Add(row, one) {
		positions = append(positions, s.getRowPositions(place, row, side)...)
	}
	return positions
}

func (s *Service) getAllPositionsUp(place *domain.UserPlace) [][2]*big.Int {
	// TODO: add capacity
	positions := make([][2]*big.Int, 0)
	row := new(big.Int).Set(place.Row)
	col := new(big.Int).Set(place.Col)
	for row.Cmp(one) > 0 {
		row.Sub(row, one)
		var r *big.Int
		col, r = new(big.Int).DivMod(col, two, new(big.Int))
		if len(r.Bits()) > 0 {
			col.Add(col, one)
		}
		positions = append(positions, [2]*big.Int{new(big.Int).Set(row), new(big.Int).Set(col)})
	}
	return positions
}

func (s *Service) getPositionUp(place *domain.UserPlace) [2]*big.Int {
	row := new(big.Int).Set(place.Row)
	col := new(big.Int).Set(place.Col)
	row.Sub(row, big.NewInt(1))
	var r *big.Int
	col, r = new(big.Int).DivMod(col, two, new(big.Int))
	if len(r.Bits()) > 0 {
		col.Add(col, one)
	}
	return [2]*big.Int{new(big.Int).Set(row), new(big.Int).Set(col)}
}

// get team position indexes for row
func (s *Service) getRowPositions(place *domain.UserPlace, row *big.Int, side domain.UserTeamType) [][2]*big.Int {
	// TODO: add capacity
	positions := make([][2]*big.Int, 0)
	start := new(big.Int).Add(s.getFromForRow(place, row, side), one)
	end := s.getToForRow(place, row, side)
	for col := new(big.Int).Set(start); col.Cmp(end) < 0; col.Add(col, one) {
		positions = append(positions, [2]*big.Int{new(big.Int).Set(row), new(big.Int).Set(col)})
	}
	return positions
}

func (s *Service) getFilterForTeamDown(place *domain.UserPlace, rows *big.Int, side domain.UserTeamType) []bson.M {
	// TODO: add capacity
	ors := make([]bson.M, 0)
	start := new(big.Int).Add(place.Row, one)
	end := new(big.Int).Add(place.Row, rows)
	for row := new(big.Int).Set(start); row.Cmp(end) <= 0; row.Add(row, one) {
		ors = append(ors, bson.M{"$and": []bson.M{
			{"row": row.String()},
			{"col": bson.M{
				"$gt":  s.getFromForRow(place, row, side).String(),
				"$lte": s.getToForRow(place, row, side).String(),
			}},
		}})
	}
	if len(ors) == 0 {
		return nil
	}
	return ors
}

func (s *Service) getFilterForSideDown(place *domain.UserPlace, rows *big.Int, side domain.UserTeamType) []bson.M {
	// TODO: add capacity
	ors := make([]bson.M, 0)
	start := new(big.Int).Add(place.Row, one)
	end := new(big.Int).Add(place.Row, rows)
	for row := new(big.Int).Set(start); row.Cmp(end) < 0; row.Add(row, one) {
		var col *big.Int
		if side == domain.UserTeamTypeRight {
			col = s.getToForRow(place, row, side)
		} else {
			col = s.getFromForRow(place, row, side)
			col.Add(col, one)
		}
		ors = append(ors, bson.M{"$and": []bson.M{
			{"row": row.String()},
			{"col": col.String()},
		}})
	}
	return ors
}

func (s *Service) getUserTeamTypeByPlaceAndBuy(place *domain.UserPlace, buy *domain.Buy) domain.UserTeamType {
	midCol := s.getFromForRow(place, buy.Row, domain.UserTeamTypeRight)
	if buy.Col.Cmp(midCol) <= 0 {
		return domain.UserTeamTypeLeft
	} else {
		return domain.UserTeamTypeRight
	}
}
