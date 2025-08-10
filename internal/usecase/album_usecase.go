package usecase

import (
	"log"

	"krupesh.faldu/internal/domain"
)

// AlbumUseCase implements the business logic for album operations
type AlbumUseCase struct {
	repo domain.AlbumRepository
}

// NewAlbumUseCase creates a new instance of AlbumUseCase
func NewAlbumUseCase(repo domain.AlbumRepository) *AlbumUseCase {
	return &AlbumUseCase{
		repo: repo,
	}
}

// ListAlbums retrieves all albums with business logic
func (uc *AlbumUseCase) ListAlbums() (*domain.AlbumsResponse, error) {
	log.Printf("Fetching albums...")

	response, err := uc.repo.ListAlbums()
	if err != nil {
		log.Printf("Failed to fetch albums: %v", err)
		return nil, err
	}

	log.Printf("Successfully fetched %d albums", len(response.Albums))

	// Business logic: if there are more pages, log it
	if response.NextPageToken != "" {
		log.Printf("More albums available on next page")
	}

	return response, nil
}

// GetAlbumByID retrieves a specific album by ID
func (uc *AlbumUseCase) GetAlbumByID(id string) (*domain.Album, error) {
	log.Printf("Fetching album with ID: %s", id)

	album, err := uc.repo.GetAlbumByID(id)
	if err != nil {
		log.Printf("Failed to fetch album %s: %v", id, err)
		return nil, err
	}

	log.Printf("Successfully fetched album: %s", album.Title)
	return album, nil
}

// CreateAlbum creates a new album with business logic
func (uc *AlbumUseCase) CreateAlbum(title string) (*domain.Album, error) {
	log.Printf("Creating album with title: %s", title)

	album, err := uc.repo.CreateAlbum(title)
	if err != nil {
		log.Printf("Failed to create album %s: %v", title, err)
		return nil, err
	}

	log.Printf("Successfully created album: %s with ID: %s", album.Title, album.ID)
	return album, nil
}

// FetchNextPage retrieves the next page of albums
func (uc *AlbumUseCase) FetchNextPage(nextPageToken string) (*domain.AlbumsResponse, error) {
	log.Printf("Fetching next page of albums...")

	response, err := uc.repo.FetchNextPage(nextPageToken)
	if err != nil {
		log.Printf("Failed to fetch next page: %v", err)
		return nil, err
	}

	log.Printf("Successfully fetched %d albums from next page", len(response.Albums))
	return response, nil
}
