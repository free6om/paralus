package service

import (
	"context"
	"time"

	"github.com/RafaySystems/rcloud-base/internal/models"
	"github.com/RafaySystems/rcloud-base/internal/persistence/provider/pg"
	commonv3 "github.com/RafaySystems/rcloud-base/proto/types/commonpb/v3"
	infrav3 "github.com/RafaySystems/rcloud-base/proto/types/infrapb/v3"
	"github.com/google/uuid"
	bun "github.com/uptrace/bun"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MetroService is the interface for metro operations
type MetroService interface {
	// create metro
	Create(ctx context.Context, metro *infrav3.Location) (*infrav3.Location, error)
	// get metro by id
	GetById(ctx context.Context, id uuid.UUID) (*infrav3.Location, error)
	// get metro by name
	GetByName(ctx context.Context, name string) (*infrav3.Location, error)
	// get metro id by name
	GetIDByName(ctx context.Context, name string) (uuid.UUID, error)
	// create or update metro
	Update(ctx context.Context, metro *infrav3.Location) (*infrav3.Location, error)
	// delete metro
	Delete(ctx context.Context, metro *infrav3.Location) (*infrav3.Location, error)
	// list metro
	List(ctx context.Context, partner string) (*infrav3.LocationList, error)
}

// metroService implements MetroService
type metroService struct {
	db *bun.DB
}

// NewProjectService return new project service
func NewMetroService(db *bun.DB) MetroService {
	return &metroService{db}
}

func (s *metroService) Create(ctx context.Context, metro *infrav3.Location) (*infrav3.Location, error) {

	var part models.Partner
	_, err := pg.GetByName(ctx, s.db, metro.Metadata.Partner, &part)
	if err != nil {
		return nil, err
	}

	//convert v3 spec to internal models
	metrodb := models.Metro{
		Name:           metro.Spec.Name,
		CreatedAt:      time.Now(),
		ModifiedAt:     time.Now(),
		Trash:          false,
		Latitude:       metro.Spec.Latitude,
		Longitude:      metro.Spec.Longitude,
		City:           metro.Spec.City,
		State:          metro.Spec.State,
		Country:        metro.Spec.Country,
		CountryCode:    metro.Spec.CountryCode,
		StateCode:      metro.Spec.StateCode,
		OrganizationId: uuid.Nil,
		PartnerId:      part.ID,
	}
	_, err = pg.Create(ctx, s.db, &metrodb)
	if err != nil {
		return nil, err
	}

	return metro, nil

}

func (s *metroService) GetByName(ctx context.Context, name string) (*infrav3.Location, error) {

	var metro infrav3.Location

	entity, err := pg.GetByName(ctx, s.db, name, &models.Metro{})
	if err != nil {
		return nil, err
	}

	if metrodb, ok := entity.(*models.Metro); ok {
		location := &infrav3.Location{
			Metadata: &commonv3.Metadata{
				Name:       metrodb.Name,
				ModifiedAt: timestamppb.New(metrodb.ModifiedAt),
			},
			Spec: &infrav3.Metro{
				Name:        metrodb.Name,
				Country:     metrodb.Country,
				City:        metrodb.City,
				State:       metrodb.State,
				Latitude:    metrodb.Latitude,
				Longitude:   metrodb.Longitude,
				StateCode:   metrodb.StateCode,
				CountryCode: metrodb.CountryCode,
			},
		}
		return location, nil

	}
	return &metro, nil
}

func (s *metroService) GetById(ctx context.Context, id uuid.UUID) (*infrav3.Location, error) {
	var location infrav3.Location

	entity, err := pg.GetByID(ctx, s.db, id, &models.Metro{})
	if err != nil {
		return nil, err
	}

	if metrodb, ok := entity.(*models.Metro); ok {

		location := &infrav3.Location{
			Metadata: &commonv3.Metadata{
				Name:       metrodb.Name,
				ModifiedAt: timestamppb.New(metrodb.ModifiedAt),
			},
			Spec: &infrav3.Metro{
				Name:        metrodb.Name,
				Country:     metrodb.Country,
				City:        metrodb.City,
				State:       metrodb.State,
				Latitude:    metrodb.Latitude,
				Longitude:   metrodb.Longitude,
				StateCode:   metrodb.StateCode,
				CountryCode: metrodb.CountryCode,
			},
		}

		return location, nil

	}
	return &location, nil
}

func (s *metroService) Update(ctx context.Context, metro *infrav3.Location) (*infrav3.Location, error) {

	entity, err := pg.GetByName(ctx, s.db, metro.Metadata.Name, &models.Metro{})
	if err != nil {
		return metro, err
	}

	if metrodb, ok := entity.(*models.Metro); ok {
		//update metro details
		metrodb.City = metro.Spec.City
		metrodb.Country = metro.Spec.Country
		metrodb.State = metro.Spec.State
		metrodb.StateCode = metro.Spec.StateCode
		metrodb.CountryCode = metro.Spec.CountryCode
		metrodb.Latitude = metro.Spec.Latitude
		metrodb.Longitude = metro.Spec.Longitude
		metrodb.ModifiedAt = time.Now()

		_, err = pg.Update(ctx, s.db, metrodb.ID, metrodb)
		if err != nil {
			return metro, err
		}
	}

	return metro, nil
}

func (s *metroService) Delete(ctx context.Context, metro *infrav3.Location) (*infrav3.Location, error) {

	entity, err := pg.GetByName(ctx, s.db, metro.Metadata.Name, &models.Metro{})
	if err != nil {
		return metro, err
	}
	if metrodb, ok := entity.(*models.Metro); ok {
		err = pg.Delete(ctx, s.db, metrodb.ID, metrodb)
		if err != nil {
			return metro, err
		}
	}

	return metro, nil
}

func (s *metroService) List(ctx context.Context, partner string) (*infrav3.LocationList, error) {
	var locations infrav3.LocationList
	var metros []*infrav3.Metro
	var metrodbs []models.Metro

	var part models.Partner
	_, err := pg.GetByName(ctx, s.db, partner, &part)
	if err != nil {
		return nil, err
	}

	entities, err := pg.List(ctx, s.db, uuid.NullUUID{UUID: part.ID, Valid: true}, uuid.NullUUID{UUID: uuid.Nil, Valid: false}, &metrodbs)
	if err != nil {
		return nil, err
	}

	if metrodbs, ok := entities.(*[]models.Metro); ok {
		for _, metrodb := range *metrodbs {

			metro := &infrav3.Metro{
				Name:        metrodb.Name,
				City:        metrodb.City,
				State:       metrodb.State,
				Country:     metrodb.Country,
				Latitude:    metrodb.Latitude,
				Longitude:   metrodb.Longitude,
				StateCode:   metrodb.StateCode,
				CountryCode: metrodb.CountryCode,
			}
			metros = append(metros, metro)
		}

		locations = infrav3.LocationList{
			Metadata: &commonv3.ListMetadata{
				Count: int64(len(metros)),
			},
			Items: metros,
		}
	}

	return &locations, nil
}

func (s *metroService) GetIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	entity, err := pg.GetByName(ctx, s.db, name, &models.Metro{})
	if err != nil {
		return uuid.Nil, err
	}

	if metrodb, ok := entity.(*models.Metro); ok {
		return metrodb.ID, nil
	}
	return uuid.Nil, nil
}
