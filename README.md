# Google Photos Magic - Clean Architecture

This project has been refactored following Go clean architecture principles to improve maintainability, testability, and separation of concerns.

## 🏗️ Architecture Overview

The code follows the clean architecture pattern with the following layers:

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Core business objects (Album, AlbumsResponse)
- **Interfaces**: Contracts for repositories and use cases
- **Pure business logic** with no external dependencies

### 2. Use Case Layer (`internal/usecase/`)
- **Business Logic**: Orchestrates operations between entities
- **Application Rules**: Implements business workflows
- **No knowledge** of HTTP, databases, or external systems

### 3. Repository Layer (`internal/repository/`)
- **Data Access**: Implements domain interfaces
- **External API Integration**: Google Photos API, OAuth2
- **Data Persistence**: Token storage, API calls

### 4. Delivery Layer (`internal/delivery/`)
- **User Interface**: CLI handlers, HTTP handlers
- **Input/Output**: Formats data for user consumption
- **No business logic** - only presentation concerns

### 5. Main Application (`cmd/app/`)
- **Dependency Injection**: Wires all layers together
- **Configuration**: OAuth setup, API endpoints
- **Entry Point**: Application bootstrap

## 📁 Project Structure

```
google-photos-magic/
├── cmd/
│   └── app/
│       └── main.go              # Clean architecture entry point
├── internal/
│   ├── domain/                  # Business entities & interfaces
│   │   ├── album.go
│   │   └── oauth.go
│   ├── usecase/                 # Business logic
│   │   ├── album_usecase.go
│   │   └── oauth_usecase.go
│   ├── repository/              # Data access & external APIs
│   │   ├── google_photos_repository.go
│   │   └── oauth_repository.go
│   └── delivery/                # User interface
│       └── cli_handler.go
├── pkg/                         # Shared utilities (future use)
├── go.mod
└── README.md
```

## 🚀 Benefits of Clean Architecture

### ✅ **Separation of Concerns**
- Business logic is isolated from infrastructure
- Each layer has a single responsibility
- Easy to understand and maintain

### ✅ **Testability**
- Business logic can be tested without HTTP or databases
- Interfaces allow easy mocking
- Unit tests are fast and reliable

### ✅ **Flexibility**
- Easy to swap implementations (e.g., different databases)
- Add new features without changing existing code
- Support multiple delivery mechanisms (CLI, HTTP, gRPC)

### ✅ **Maintainability**
- Clear dependencies between layers
- Changes in one layer don't affect others
- Easy to onboard new developers

## 🔧 How to Use

```bash
go run cmd/app/main.go
```

## 🧪 Testing

The clean architecture makes testing much easier:

```go
// Mock the repository for testing business logic
mockRepo := &MockAlbumRepository{}
albumUseCase := usecase.NewAlbumUseCase(mockRepo)

// Test business logic without external dependencies
albums, err := albumUseCase.ListAlbums()
```

## 🔄 Migration Path

1. **Phase 1**: ✅ Complete - Core architecture implemented
2. **Phase 2**: Add comprehensive error handling
3. **Phase 3**: Implement logging and monitoring
4. **Phase 4**: Add configuration management
5. **Phase 5**: Implement caching layer

## 📚 Key Principles Applied

### Dependency Inversion
- High-level modules don't depend on low-level modules
- Both depend on abstractions (interfaces)

### Single Responsibility
- Each package has one clear purpose
- Functions do one thing well

### Interface Segregation
- Small, focused interfaces
- Clients only depend on methods they use

### Open/Closed Principle
- Open for extension (new implementations)
- Closed for modification (existing code)

## 🚨 Important Notes

- **credentials.json**: Required for OAuth2 authentication
- **token.json**: Automatically created after first OAuth flow
- **Dependencies**: Ensure all Go modules are properly installed

## 🔍 Code Examples

### Adding a New Feature

1. **Define interface** in domain layer
2. **Implement business logic** in use case layer
3. **Add data access** in repository layer
4. **Create user interface** in delivery layer
5. **Wire together** in main application

### Example: Adding Photo Management

```go
// 1. Domain interface
type PhotoRepository interface {
    GetPhotos(albumID string) ([]Photo, error)
}

// 2. Use case implementation
func (uc *PhotoUseCase) GetPhotos(albumID string) ([]Photo, error) {
    return uc.repo.GetPhotos(albumID)
}

// 3. Repository implementation
func (r *GooglePhotosRepository) GetPhotos(albumID string) ([]Photo, error) {
    // API call implementation
}

// 4. CLI handler
func (h *CLIHandler) HandleGetPhotos(albumID string) {
    photos, err := h.photoUseCase.GetPhotos(albumID)
    // Display logic
}
```

## 🤝 Contributing

When adding new features:
1. Follow the existing architecture patterns
2. Add interfaces in domain layer first
3. Implement business logic in use case layer
4. Add data access in repository layer
5. Create user interface in delivery layer
6. Update main application wiring

## 📖 Further Reading

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Clean Architecture Tutorial](https://medium.com/@letsCodeDevelopers/golang-clean-architecture-step-by-step-tutorial-b678c763c601)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
