package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// FileSystemItem represents a file or directory
type FileSystemItem struct {
	Key         string  `json:"key"`
	Label       string  `json:"label"`
	Path        string  `json:"path"`
	Name        string  `json:"name"`
	IsDirectory bool    `json:"isDirectory"`
	Size        int64   `json:"size"`
	Modified    int64   `json:"modified"`
	IsLeaf      bool    `json:"isLeaf"`
}

// GetDrives returns available drives on Windows
func (a *App) GetDrives() []FileSystemItem {
	var drives []FileSystemItem
	
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		if _, err := os.Stat(drivePath); err == nil {
			drives = append(drives, FileSystemItem{
				Key:         drivePath,
				Label:       string(drive) + ":",
				Path:        drivePath,
				Name:        string(drive) + ":",
				IsDirectory: true,
				IsLeaf:      false,
			})
		}
	}
	
	return drives
}

// GetDirectoryContents returns the contents of a directory
func (a *App) GetDirectoryContents(path string) ([]FileSystemItem, error) {
	var items []FileSystemItem
	
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}
	
	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		item := FileSystemItem{
			Key:         fullPath,
			Label:       entry.Name(),
			Path:        fullPath,
			Name:        entry.Name(),
			IsDirectory: entry.IsDir(),
			Size:        info.Size(),
			Modified:    info.ModTime().Unix(),
			IsLeaf:      !entry.IsDir(),
		}
		
		items = append(items, item)
	}
	
	return items, nil
}

// GetFileSystemItem returns detailed information about a specific file or directory
func (a *App) GetFileSystemItem(path string) (*FileSystemItem, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}
	
	return &FileSystemItem{
		Key:         path,
		Label:       filepath.Base(path),
		Path:        path,
		Name:        filepath.Base(path),
		IsDirectory: info.IsDir(),
		Size:        info.Size(),
		Modified:    info.ModTime().Unix(),
		IsLeaf:      !info.IsDir(),
	}, nil
}

// GetRootDirectories returns root level directories (for cross-platform compatibility)
func (a *App) GetRootDirectories() ([]FileSystemItem, error) {
	// For Windows, return drives
	if strings.HasPrefix(os.Getenv("OS"), "Windows") {
		return a.GetDrives(), nil
	}
	
	// For Unix-like systems, return root directory
	rootItems, err := a.GetDirectoryContents("/")
	if err != nil {
		return nil, err
	}
	
	return rootItems, nil
}
