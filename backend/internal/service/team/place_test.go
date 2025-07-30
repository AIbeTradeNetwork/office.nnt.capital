package team

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"server/internal/domain"
	"server/internal/service/team/mocks"
	"testing"
)

func TestService_getUserTeamTypeByPlaceAndBuy(t *testing.T) {
	repo := mocks.NewDbRepository(t)

	s := NewTeamService(repo)

	teamType := s.getUserTeamTypeByPlaceAndBuy(&domain.UserPlace{
		Row: big.NewInt(3),
		Col: big.NewInt(8),
	}, &domain.Buy{
		Row: big.NewInt(7),
		Col: big.NewInt(16),
	})
	assert.Equal(t, domain.UserTeamTypeLeft, teamType)

	teamType = s.getUserTeamTypeByPlaceAndBuy(&domain.UserPlace{
		Row: big.NewInt(3),
		Col: big.NewInt(8),
	}, &domain.Buy{
		Row: big.NewInt(7),
		Col: big.NewInt(150),
	})
	assert.Equal(t, domain.UserTeamTypeRight, teamType)
}

func TestService_getUserTeamTypeByPlace(t *testing.T) {
	repo := mocks.NewDbRepository(t)

	s := NewTeamService(repo)

	teamType := s.getUserTeamTypeByPlace(&domain.UserPlace{
		Row: big.NewInt(8),
		Col: big.NewInt(128),
	})
	assert.Equal(t, domain.UserTeamTypeRight, teamType)

	teamType = s.getUserTeamTypeByPlace(&domain.UserPlace{
		Row: big.NewInt(8),
		Col: big.NewInt(1),
	})
	assert.Equal(t, domain.UserTeamTypeLeft, teamType)
}

func TestService_getOppositeUserTeamTypeByPlace(t *testing.T) {
	repo := mocks.NewDbRepository(t)

	s := NewTeamService(repo)

	teamType := s.getOppositeUserTeamTypeByPlace(&domain.UserPlace{
		Row: big.NewInt(8),
		Col: big.NewInt(128),
	})
	assert.Equal(t, domain.UserTeamTypeLeft, teamType)

	teamType = s.getOppositeUserTeamTypeByPlace(&domain.UserPlace{
		Row: big.NewInt(8),
		Col: big.NewInt(1),
	})
	assert.Equal(t, domain.UserTeamTypeRight, teamType)
}
