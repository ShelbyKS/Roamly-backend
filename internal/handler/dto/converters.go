package dto

import (
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type UserConverter struct{}

func (UserConverter) ToDto(user model.User) GetUser {
	return GetUser{
		ID:        user.ID,
		Login:     user.Login,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}

type TripConverter struct{}

func (TripConverter) ToDto(trip model.Trip) TripResponse {
	places := make([]GooglePlace, len(trip.Places))
	for i, place := range trip.Places {
		placeDto := GooglePlaceConverter{}.ToDto(place.GooglePlace)
		places[i] = placeDto
	}

	// todo: тут не все поля юзера
	users := make([]GetUser, len(trip.Users))
	for i, user := range trip.Users {
		users[i] = GetUser{
			ID:    user.ID,
			Login: user.Login,
			Email: user.Email,
		}
	}

	events := make([]GetEvent, len(trip.Events))
	for i, event := range trip.Events {
		events[i] = GetEvent{
			PlaceID:   event.PlaceID,
			TripID:    event.TripID,
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
		}
	}

	area := GooglePlaceConverter{}.ToDto(trip.Area.GooglePlace)

	return TripResponse{
		ID:        trip.ID,
		Name:      trip.Name,
		Users:     users,
		StartTime: trip.StartTime,
		EndTime:   trip.EndTime,
		AreaID:    trip.AreaID,
		Area:      area,
		Places:    places,
		Events:    events,
	}
}

type PhotoConverter struct{}

func (PhotoConverter) ToDb(photo model.Photo) Photo {
	return Photo{
		PhotoReference: photo.PhotoReference,
	}
}

type GooglePlaceConverter struct{}

func (GooglePlaceConverter) ToDto(gp model.GooglePlace) GooglePlace {
	photos := make([]Photo, len(gp.Photos))
	for i, photo := range gp.Photos {
		photos[i] = PhotoConverter{}.ToDb(photo)
	}

	return GooglePlace{
		FormattedAddress: gp.FormattedAddress,
		Geometry:         GeometryConverter{}.ToDto(gp.Geometry),
		Name:             gp.Name,
		PlaceID:          gp.PlaceID,
		Rating:           gp.Rating,
		Types:            gp.Types,
		Photos:           photos,
		Vicinity:         gp.Vicinity,
	}
}

type GeometryConverter struct{}

func (GeometryConverter) ToDto(geom model.Geometry) Geometry {
	return Geometry{
		Location: LocationConverter{}.ToDto(geom.Location),
	}
}

type LocationConverter struct{}

func (LocationConverter) ToDto(loc model.Location) Location {
	return Location{
		Lat: loc.Lat,
		Lng: loc.Lng,
	}
}

type EventConverter struct{}

func (EventConverter) ToDto(event model.Event) GetEvent {
	return GetEvent{
		PlaceID:   event.PlaceID,
		TripID:    event.TripID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
	}
}
