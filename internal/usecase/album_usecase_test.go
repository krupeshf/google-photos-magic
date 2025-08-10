package usecase

import (
	"testing"

	"krupesh.faldu/internal/domain"
)

// MockAlbumRepository is a mock implementation for testing
type MockAlbumRepository struct {
	albums []domain.Album
	err    error
}

func (m *MockAlbumRepository) ListAlbums() (*domain.AlbumsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.AlbumsResponse{
		Albums:        m.albums,
		NextPageToken: "",
	}, nil
}

func (m *MockAlbumRepository) GetAlbumByID(id string) (*domain.Album, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, album := range m.albums {
		if album.ID == id {
			return &album, nil
		}
	}
	return nil, nil
}

func (m *MockAlbumRepository) CreateAlbum(title string) (*domain.Album, error) {
	if m.err != nil {
		return nil, m.err
	}
	album := domain.Album{
		ID:    "test-id",
		Title: title,
	}
	return &album, nil
}

func (m *MockAlbumRepository) FetchNextPage(nextPageToken string) (*domain.AlbumsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.AlbumsResponse{
		Albums:        m.albums,
		NextPageToken: "",
	}, nil
}

func TestAlbumUseCase_ListAlbums(t *testing.T) {
	// Arrange
	mockRepo := &MockAlbumRepository{
		albums: []domain.Album{
			{ID: "1", Title: "Test Album 1"},
			{ID: "2", Title: "Test Album 2"},
		},
	}

	useCase := NewAlbumUseCase(mockRepo)

	// Act
	response, err := useCase.ListAlbums()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(response.Albums) != 2 {
		t.Errorf("Expected 2 albums, got %d", len(response.Albums))
	}

	if response.Albums[0].Title != "Test Album 1" {
		t.Errorf("Expected first album title 'Test Album 1', got '%s'", response.Albums[0].Title)
	}
}

func TestAlbumUseCase_CreateAlbum(t *testing.T) {
	// Arrange
	mockRepo := &MockAlbumRepository{}
	useCase := NewAlbumUseCase(mockRepo)
	title := "New Test Album"

	// Act
	album, err := useCase.CreateAlbum(title)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if album.Title != title {
		t.Errorf("Expected album title '%s', got '%s'", title, album.Title)
	}

	if album.ID != "test-id" {
		t.Errorf("Expected album ID 'test-id', got '%s'", album.ID)
	}
}
