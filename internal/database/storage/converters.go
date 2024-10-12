package storage

import (
	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type UserConverter struct{}

func (UserConverter) ToDb(user model.User) orm.User {
	return orm.User{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}
}

func (UserConverter) ToDomain(user model.User) orm.User {
	return orm.User{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
	}
}

type TripConverter struct{}

func (TripConverter) ToDb(trip model.Trip) orm.Trip {
	users := make([]*orm.User, len(trip.Users))
	for i, user := range trip.Users {
		users[i] = &orm.User{
			ID:       user.ID,
			Login:    user.Login,
			Password: user.Password,
		}
	}
	return orm.Trip{
		ID:        trip.ID,
		Users:     users,
		StartTime: trip.StartTime,
		EndTime:   trip.EndTime,
		AreaID:    trip.AreaID,
	}
}

func (TripConverter) ToDomain(trip orm.Trip) model.Trip {
	var tripPlaces []model.Place

	for _, place := range trip.Places {
		tripPlaces = append(tripPlaces, model.Place{
			ID:     place.Payload.Data().PlaceID,
			Name:   place.Payload.Data().Name,
			Rating: place.Payload.Data().Rating,
		})
	}

	users := make([]*model.User, len(trip.Users))
	for i, user := range trip.Users {
		users[i] = &model.User{
			ID:       user.ID,
			Login:    user.Login,
			Password: user.Password,
		}
	}
	return model.Trip{
		ID:        trip.ID,
		Users:     users,
		StartTime: trip.StartTime,
		EndTime:   trip.EndTime,
		AreaID:    trip.AreaID,
	}
}
