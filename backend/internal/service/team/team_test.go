package team

import (
	"testing"
)

func TestService_PlaceCreateForRef(t *testing.T) {
	//ctx := context.Background()

	//rootUser := &domain.User{
	//	UID:       "ADMIN",
	//	RefUID:    "",
	//	Nickname:  "admin",
	//	Email:     "admin@mail.com",
	//	CreatedAt: time.Now().UTC(),
	//}
	//dummyUser := &domain.User{
	//	UID:       "DUMMY",
	//	RefUID:    "ADMIN",
	//	Nickname:  "dummy",
	//	Email:     "dummy@mail.com",
	//	CreatedAt: time.Now().UTC(),
	//}
	//users := []*domain.User{
	//	{
	//		UID:       "NEWBEE",
	//		RefUID:    "DUMMY",
	//		Nickname:  "newbee",
	//		Email:     "newbee@mail.com",
	//		CreatedAt: time.Now().UTC(),
	//	},
	//	{
	//		UID:       "NOOBEE",
	//		RefUID:    "DUMMY",
	//		Nickname:  "noobee",
	//		Email:     "noobee@mail.com",
	//		CreatedAt: time.Now().UTC(),
	//	},
	//	{
	//		UID:       "NOOBI",
	//		RefUID:    "DUMMY",
	//		Nickname:  "noobi",
	//		Email:     "noobi@mail.com",
	//		CreatedAt: time.Now().UTC(),
	//	},
	//	{
	//		UID:       "MJPAIN",
	//		RefUID:    "DUMMY",
	//		Nickname:  "mjpain",
	//		Email:     "mjpain@mail.com",
	//		CreatedAt: time.Now().UTC(),
	//	},
	//	{
	//		UID:       "MJFROST",
	//		RefUID:    "MJPAIN",
	//		Nickname:  "mjfrost",
	//		Email:     "mjfrost@mail.com",
	//		CreatedAt: time.Now().UTC(),
	//	},
	//}

	//root := &domain.UserPlace{
	//	UserUID:   "ADMIN",
	//	RefUID:    "",
	//	Row:       1,
	//	Col:       1,
	//	CreatedAt: time.Now().UTC(),
	//}
	//place := &domain.UserPlace{
	//	UserUID:   "DUMMY",
	//	RefUID:    "ADMIN",
	//	Row:       3,
	//	Col:       3,
	//	CreatedAt: time.Now().UTC(),
	//}
	//
	//places := []*domain.UserPlace{
	//	{
	//		UserUID:   "NEWBEE",
	//		RefUID:    "DUMMY",
	//		Row:       4,
	//		Col:       5,
	//		CreatedAt: time.Now().UTC(),
	//	},
	//	{
	//		UserUID:   "NOOBEE",
	//		RefUID:    "DUMMY",
	//		Row:       4,
	//		Col:       6,
	//		CreatedAt: time.Now().UTC(),
	//	},
	//	{
	//		UserUID:   "NOOBI",
	//		RefUID:    "DUMMY",
	//		Row:       5,
	//		Col:       11,
	//		CreatedAt: time.Now().UTC(),
	//	},
	//}
	//
	//repo := mocks.NewDbRepository(t)

	//repo, err := repository.NewDbRepo(ctx)
	//err = repo.DropTest(ctx)

	//err = repo.UserCreate(ctx, rootUser)
	//err = repo.UserCreate(ctx, dummyUser)
	//for _, u := range users {
	//	err = repo.UserCreate(ctx, u)
	//}
	//
	//err = repo.UserPlaceCreate(ctx, root)
	//err = repo.UserPlaceCreate(ctx, place)
	//for _, p := range places {
	//	err = repo.UserPlaceCreate(ctx, p)
	//}

	//service := NewTeamService(repo)
	//
	//from := service.getFromForRow(place, 5, domain.UserTeamTypeUndefined)
	//assert.Equal(t, from, uint64(8))
	//
	//from = service.getFromForRow(place, 7, domain.UserTeamTypeLeft)
	//assert.Equal(t, from, uint64(32))
	//
	//from = service.getFromForRow(place, 10, domain.UserTeamTypeRight)
	//assert.Equal(t, from, uint64(320))
	//
	//placesLeft := places
	//placesRight := places[0:0]

	//repo.On("UserPlaceGet", ctx, mock.Anything, mock.Anything).Return(places, nil)
	//repo.On("UserPlaceGetAll", ctx, mock.Anything, mock.Anything).Return(places, nil)
	//newPlace, err := service.PlaceGetNew(ctx, place, domain.UserTeamTypeUndefined)
	//assert.NoError(t, err)
	//assert.Equal(t, newPlace.Row, uint64(5))
	//assert.Equal(t, newPlace.Col, uint64(9))

	//repo.ExpectedCalls = nil
	//repo.On("UserPlaceGet", ctx, mock.Anything, mock.Anything).Return(places[2:3], nil)
	//repo.On("UserPlaceGetAll", ctx, mock.Anything, mock.Anything).Return(placesLeft, nil)
	//newPlace, err := service.PlaceGetNew(ctx, place, domain.UserTeamTypeLeft)
	//assert.NoError(t, err)
	//assert.Equal(t, newPlace.Row, uint64(5))
	//assert.Equal(t, newPlace.Col, uint64(9))
	//
	//repo.ExpectedCalls = nil
	//repo.On("UserPlaceGet", ctx, mock.Anything, mock.Anything).Return(places[1:2], nil)
	//repo.On("UserPlaceGetAll", ctx, mock.Anything, mock.Anything).Return(placesRight, nil)
	//newPlace, err = service.PlaceGetNew(ctx, place, domain.UserTeamTypeRight)
	//assert.NoError(t, err)
	//assert.Equal(t, newPlace.Row, uint64(4))
	//assert.Equal(t, newPlace.Col, uint64(6))

	//repo.ExpectedCalls = nil
	//repo.On("UserPlaceGetAll", ctx, mock.Anything, mock.Anything).Return(placesRight, nil)
	//newPlace, err = service.PlaceCreateForRef(ctx, "MJFROST")
	//assert.NoError(t, err)
	//assert.Equal(t, newPlace.Row, uint64(5))
	//assert.Equal(t, newPlace.Col, uint64(9))
	//assert.Equal(t, newPlace.RefUID, "DUMMY")
	//conf, err := repo.UserConfigGetByUserUID(ctx, "DUMMY")
	//assert.Equal(t, conf.LastTeamType, domain.UserTeamTypeLeft)
	//
	//repo.ExpectedCalls = nil
	//repo.On("UserPlaceGetAll", ctx, mock.Anything, mock.Anything).Return(placesRight, nil)
	//newPlace, err = service.PlaceCreateForRef(ctx, "MJPAIN")
	//assert.NoError(t, err)
	//assert.Equal(t, newPlace.Row, uint64(5))
	//assert.Equal(t, newPlace.Col, uint64(12))
	//assert.Equal(t, newPlace.RefUID, "DUMMY")
	//conf, err = repo.UserConfigGetByUserUID(ctx, "DUMMY")
	//assert.Equal(t, conf.LastTeamType, domain.UserTeamTypeRight)
}
