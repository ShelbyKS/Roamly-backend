package postgresql

import (
	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type UserConverter struct{}

func (UserConverter) ToDb(user model.User) orm.User {
	return orm.User{
		ID:        user.ID,
		Login:     user.Login,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}

func (UserConverter) ToDomain(user orm.User) model.User {
	return model.User{
		ID:        user.ID,
		Login:     user.Login,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
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

	var tripPlaces []*orm.Place
	for _, place := range trip.Places {
		placeDb := PlaceConverter{}.ToDb(*place)
		tripPlaces = append(tripPlaces, &placeDb)
	}

	var tripEvents []orm.Event
	for _, event := range trip.Events {
		eventDB := EventConverter{}.ToDb(event)
		tripEvents = append(tripEvents, eventDB)
	}

	var areaDb orm.Place
	if trip.Area != nil {
		areaDb = PlaceConverter{}.ToDb(*trip.Area)
	}

	return orm.Trip{
		ID:        trip.ID,
		Name:      trip.Name,
		Users:     users,
		StartTime: trip.StartTime,
		EndTime:   trip.EndTime,
		AreaID:    trip.AreaID,
		Places:    tripPlaces,
		Events:    tripEvents,
		Area:      areaDb,
	}
}

func (TripConverter) ToDomain(trip orm.Trip) model.Trip {
	var tripPlaces []*model.Place

	for _, place := range trip.Places {
		placeDomain := PlaceConverter{}.ToDomain(*place)
		tripPlaces = append(tripPlaces, &placeDomain)
	}

	// todo: тут не все поля юзера
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

	area := PlaceConverter{}.ToDomain(trip.Area)

	return model.Trip{
		ID:        trip.ID,
		Name:      trip.Name,
		Users:     users,
		StartTime: trip.StartTime,
		EndTime:   trip.EndTime,
		AreaID:    trip.AreaID,
		Area:      &area,
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
		ID: place.ID,
		// Photo:       place.Photo,
		// Name:        place.Name,
		// Closing:     place.Closing,
		// Opening:     place.Opening,
		// Rating:      place.Rating,
		Trips:       trips,
		GooglePlace: GooglePlaceConverter{}.ToDb(place.GooglePlace),
	}
}

func (PlaceConverter) ToDomain(place orm.Place) model.Place {
	var trips []*model.Trip
	for _, i := range place.Trips {
		trip := TripConverter{}.ToDomain(*i)
		trips = append(trips, &trip)
	}

	return model.Place{
		ID: place.ID,
		// Photo:       place.Photo,
		// Name:        place.Name,
		// Closing:     place.Closing,
		// Opening:     place.Opening,
		// Rating:      place.Rating,
		Trips:       trips,
		GooglePlace: GooglePlaceConverter{}.ToDomain(place.GooglePlace),
	}
}

type PhotoConverter struct{}

func (PhotoConverter) ToDb(photo model.Photo) orm.Photo {
	return orm.Photo{
		PhotoReference: photo.PhotoReference,
	}
}

func (PhotoConverter) ToDomain(photo orm.Photo) model.Photo {
	return model.Photo{
		PhotoReference: photo.PhotoReference,
	}
}

type GooglePlaceConverter struct{}

func (GooglePlaceConverter) ToDb(gp model.GooglePlace) orm.GooglePlace {
	photos := make([]orm.Photo, len(gp.Photos))
	for i, photo := range gp.Photos {
		photos[i] = PhotoConverter{}.ToDb(photo)
	}

	return orm.GooglePlace{
		FormattedAddress: gp.FormattedAddress,
		Geometry:         GeometryConverter{}.ToDb(gp.Geometry),
		Photos:           photos,
		Name:             gp.Name,
		PlaceID:          gp.PlaceID,
		Rating:           gp.Rating,
		Types:            gp.Types,
	}
}

func (GooglePlaceConverter) ToDomain(gp orm.GooglePlace) model.GooglePlace {
	photos := make([]model.Photo, len(gp.Photos))
	for i, photo := range gp.Photos {
		photos[i] = PhotoConverter{}.ToDomain(photo)
	}

	return model.GooglePlace{
		FormattedAddress: gp.FormattedAddress,
		Geometry:         GeometryConverter{}.ToDomain(gp.Geometry),
		Photos:           photos,
		Name:             gp.Name,
		PlaceID:          gp.PlaceID,
		Rating:           gp.Rating,
		Types:            gp.Types,
	}
}

type GeometryConverter struct{}

func (GeometryConverter) ToDb(geom model.Geometry) orm.Geometry {
	return orm.Geometry{
		Location: LocationConverter{}.ToDb(geom.Location),
	}
}

func (GeometryConverter) ToDomain(geom orm.Geometry) model.Geometry {
	return model.Geometry{
		Location: LocationConverter{}.ToDomain(geom.Location),
	}
}

type LocationConverter struct{}

func (LocationConverter) ToDb(loc model.Location) orm.Location {
	return orm.Location{
		Lat: loc.Lat,
		Lng: loc.Lng,
	}
}

func (LocationConverter) ToDomain(loc orm.Location) model.Location {
	return model.Location{
		Lat: loc.Lat,
		Lng: loc.Lng,
	}
}

// type PlaceConverter struct{}

// func (PlaceConverter) ToDb(place model.Place) orm.Place {
// 	var trips []*orm.Trip
// 	for _, i := range place.Trips {
// 		trip := TripConverter{}.ToDb(*i)
// 		trips = append(trips, &trip)
// 	}

// 	return orm.Place{
// 		ID:      place.ID,
// 		Photo:   place.Photo,
// 		Name:    place.Name,
// 		Closing: place.Closing,
// 		Opening: place.Opening,
// 		Rating:  place.Rating,
// 		Trips:   trips,
// 	}
// }

// func (PlaceConverter) ToDomain(place orm.Place) model.Place {
// 	var trips []*model.Trip
// 	for _, i := range place.Trips {
// 		trip := TripConverter{}.ToDomain(*i)
// 		trips = append(trips, &trip)
// 	}

// 	return model.Place{
// 		ID:     place.ID,
// 		Trips:  trips,
// 		Name:   place.Name,
// 		Photo:  place.Photo,
// 		Rating: place.Rating,
// 	}
// }

type EventConverter struct{}

func (EventConverter) ToDb(event model.Event) orm.Event {
	return orm.Event{
		ID:        event.ID,
		PlaceID:   event.PlaceID,
		TripID:    event.TripID,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		// Trip:      TripConverter{}.ToDb(event.Trip),
		// Place:     PlaceConverter{}.ToDb(event.Place),
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
