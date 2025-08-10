package domain

// Album represents a Google Photos album
type Album struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// AlbumsResponse represents the API response for listing albums
type AlbumsResponse struct {
	Albums        []Album `json:"albums"`
	NextPageToken string  `json:"nextPageToken"`
}

// AlbumRepository defines the interface for album operations
type AlbumRepository interface {
	ListAlbums() (*AlbumsResponse, error)
	GetAlbumByID(id string) (*Album, error)
	CreateAlbum(title string) (*Album, error)
	FetchNextPage(nextPageToken string) (*AlbumsResponse, error)
}

// AlbumUseCase defines the business logic for album operations
type AlbumUseCase interface {
	ListAlbums() (*AlbumsResponse, error)
	GetAlbumByID(id string) (*Album, error)
	CreateAlbum(title string) (*Album, error)
	FetchNextPage(nextPageToken string) (*AlbumsResponse, error)
}
