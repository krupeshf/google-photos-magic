package delivery

import (
	"log"
	"time"

	"krupesh.faldu/internal/domain"
	"krupesh.faldu/internal/usecase"
)

// CLIHandler handles command-line interface interactions
type CLIHandler struct {
	albumUseCase *usecase.AlbumUseCase
	oauthUseCase *usecase.OAuthUseCase
}

// NewCLIHandler creates a new instance of CLIHandler
func NewCLIHandler(albumUseCase *usecase.AlbumUseCase, oauthUseCase *usecase.OAuthUseCase) *CLIHandler {
	return &CLIHandler{
		albumUseCase: albumUseCase,
		oauthUseCase: oauthUseCase,
	}
}

// HandleListAlbums handles the list albums command
func (h *CLIHandler) HandleListAlbums() {
	log.Printf("--- Listing Albums ---")

	response, err := h.albumUseCase.ListAlbums()
	if err != nil {
		log.Printf("Failed to list albums: %v", err)
		return
	}

	h.printAlbums(response.Albums)

	if response.NextPageToken != "" {
		log.Printf("Next page token: %s", response.NextPageToken)
		h.handleNextPage(response.NextPageToken)
	}
}

// HandleCreateAlbum handles the create album command
func (h *CLIHandler) HandleCreateAlbum() {
	log.Printf("--- Testing Album Creation ---")
	title := "test-album-" + time.Now().Format("2006-01-02-15-04-05")

	album, err := h.albumUseCase.CreateAlbum(title)
	if err != nil {
		log.Printf("Failed to create album: %v", err)
		return
	}

	log.Printf("Successfully created album: %s with ID: %s", album.Title, album.ID)
}

// HandleGetAlbum handles the get album by ID command
func (h *CLIHandler) HandleGetAlbum(albumID string) {
	log.Printf("--- Getting Album by ID ---")

	album, err := h.albumUseCase.GetAlbumByID(albumID)
	if err != nil {
		log.Printf("Failed to get album: %v", err)
		return
	}

	log.Printf("Album Info:")
	log.Printf("- ID: %s", album.ID)
	log.Printf("- Title: %s", album.Title)
}

// HandleNextPage handles fetching the next page of albums
func (h *CLIHandler) HandleNextPage(nextPageToken string) {
	log.Printf("--- Fetching Next Page ---")

	response, err := h.albumUseCase.FetchNextPage(nextPageToken)
	if err != nil {
		log.Printf("Failed to fetch next page: %v", err)
		return
	}

	if len(response.Albums) > 0 {
		log.Printf("Found %d albums on next page:", len(response.Albums))
		h.printAlbums(response.Albums)
	}
}

// printAlbums prints album information to the console
func (h *CLIHandler) printAlbums(albums []domain.Album) {
	if len(albums) == 0 {
		log.Printf("No albums found.")
		return
	}

	log.Printf("Albums:")
	for _, album := range albums {
		log.Printf("- %s (%s)", album.Title, album.ID)
	}
}

// handleNextPage is a helper method for handling next page requests
func (h *CLIHandler) handleNextPage(nextPageToken string) {
	h.HandleNextPage(nextPageToken)
}
