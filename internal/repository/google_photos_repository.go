package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"krupesh.faldu/internal/domain"
)

const (
	albumsEndpoint = "https://photoslibrary.googleapis.com/v1/albums"
)

// GooglePhotosRepository implements the AlbumRepository interface
type GooglePhotosRepository struct {
	client *http.Client
}

// NewGooglePhotosRepository creates a new instance of GooglePhotosRepository
func NewGooglePhotosRepository(client *http.Client) domain.AlbumRepository {
	return &GooglePhotosRepository{
		client: client,
	}
}

// ListAlbums retrieves all albums from Google Photos API
func (r *GooglePhotosRepository) ListAlbums() (*domain.AlbumsResponse, error) {
	resp, err := r.makeAlbumsRequest(albumsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to make albums request: %v", err)
	}
	defer resp.Body.Close()

	return r.readAndParseResponse(resp)
}

// GetAlbumByID retrieves a specific album by ID
func (r *GooglePhotosRepository) GetAlbumByID(id string) (*domain.Album, error) {
	url := fmt.Sprintf("%s/%s", albumsEndpoint, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch album: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var album domain.Album
	if err := json.NewDecoder(resp.Body).Decode(&album); err != nil {
		return nil, fmt.Errorf("failed to decode album: %v", err)
	}

	return &album, nil
}

// CreateAlbum creates a new album
func (r *GooglePhotosRepository) CreateAlbum(title string) (*domain.Album, error) {
	body := map[string]interface{}{
		"album": map[string]string{
			"title": title,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", albumsEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create album failed: %v", err)
	}
	defer resp.Body.Close()

	var album domain.Album
	if err := json.NewDecoder(resp.Body).Decode(&album); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &album, nil
}

// FetchNextPage retrieves the next page of albums
func (r *GooglePhotosRepository) FetchNextPage(nextPageToken string) (*domain.AlbumsResponse, error) {
	nextPageURL := albumsEndpoint + "?pageToken=" + nextPageToken

	resp, err := r.makeAlbumsRequest(nextPageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch next page: %v", err)
	}
	defer resp.Body.Close()

	return r.readAndParseResponse(resp)
}

// makeAlbumsRequest creates and executes a request to the albums endpoint
func (r *GooglePhotosRepository) makeAlbumsRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return r.client.Do(req)
}

// readAndParseResponse reads and parses the HTTP response
func (r *GooglePhotosRepository) readAndParseResponse(resp *http.Response) (*domain.AlbumsResponse, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	log.Printf("Raw API Response: %s", string(body))

	var data domain.AlbumsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return &data, nil
}
