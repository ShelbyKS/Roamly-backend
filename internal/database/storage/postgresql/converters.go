package postgresql

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
	var tripPlaces []*model.Place

	for _, place := range trip.Places {
		tripPlaces = append(tripPlaces, &model.Place{
			ID:     place.ID,
			Photo:  place.Photo,
			Name:   place.Name,
			Rating: place.Rating,
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

	events := make([]model.Event, len(trip.Events))
	for i, event := range trip.Events {
		events[i] = model.Event{
			PlaceID:   event.PlaceID,
			TripID:    event.TripID,
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
		}
	}

	tripArea := &model.Place{
		ID:     trip.Area.ID,
		Photo:  trip.Area.Photo,
		Name:   trip.Area.Name,
		Rating: trip.Area.Rating,
	}

	return model.Trip{
		ID:        trip.ID,
		Users:     users,
		StartTime: trip.StartTime,
		EndTime:   trip.EndTime,
		AreaID:    trip.AreaID,
		Area:      tripArea,
		Places:    tripPlaces,
		Events:    events,
	}
}

type PlaceConverter struct{}

func (PlaceConverter) ToDb(place model.Place) orm.Place {
	var trips []*orm.Trip
	for _, i := range place.Trips {
		trip := TripConverter{}.ToDb(*i)
		trips = append(trips, &trip)
	}

	return orm.Place{
		ID:      place.ID,
		Photo:   place.Photo,
		Name:    place.Name,
		Closing: place.Closing,
		Opening: place.Opening,
		Rating:  place.Rating,
		Trips:   trips,
	}
}

func (PlaceConverter) ToDomain(place orm.Place) model.Place {
	var trips []*model.Trip
	for _, i := range place.Trips {
		trip := TripConverter{}.ToDomain(*i)
		trips = append(trips, &trip)
	}

	return model.Place{
		ID:     place.ID,
		Trips:  trips,
		Name:   place.Name,
		Photo:  place.Photo,
		Rating: place.Rating,
	}
}

type EventConverter struct{}

func (EventConverter) ToDb(event model.Event) orm.Event {
	return orm.Event{
		PlaceID:   event.PlaceID,
		TripID:    event.TripID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		Trip:      TripConverter{}.ToDb(event.Trip),
		Place:     PlaceConverter{}.ToDb(event.Place),
	}
}

func (EventConverter) ToDomain(event orm.Event) model.Event {
	return model.Event{
		PlaceID:   event.PlaceID,
		TripID:    event.TripID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		Trip:      TripConverter{}.ToDomain(event.Trip),
		Place:     PlaceConverter{}.ToDomain(event.Place),
	}
}
