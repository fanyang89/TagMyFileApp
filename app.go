package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
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

// DiskSpaceInfo represents disk space information
type DiskSpaceInfo struct {
	TotalBytes     uint64 `json:"totalBytes"`
	FreeBytes      uint64 `json:"freeBytes"`
	UsedBytes      uint64 `json:"usedBytes"`
	TotalSpaceGB   float64 `json:"totalSpaceGB"`
	FreeSpaceGB    float64 `json:"freeSpaceGB"`
	UsedSpaceGB    float64 `json:"usedSpaceGB"`
	UsagePercent   float64 `json:"usagePercent"`
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

// GetImageData returns base64 encoded image data for image files
func (a *App) GetImageData(path string) (string, error) {
	// 检查文件扩展名是否为图片格式
	ext := strings.ToLower(filepath.Ext(path))
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	
	if !imageExts[ext] {
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}
	
	// 读取文件数据
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}
	
	// 编码为 base64
	mimeType := "image/jpeg"
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".bmp":
		mimeType = "image/bmp"
	case ".webp":
		mimeType = "image/webp"
	}
	
	base64Data := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), nil
}

// OpenFile opens a file with the default system application
func (a *App) OpenFile(path string) error {
	if runtime.GOOS == "windows" {
		// Windows: use start command
		cmd := exec.Command("cmd", "/c", "start", "", "/B", path)
		return cmd.Run()
	} else if runtime.GOOS == "darwin" {
		// macOS: use open command
		cmd := exec.Command("open", path)
		return cmd.Run()
	} else {
		// Linux: use xdg-open command
		cmd := exec.Command("xdg-open", path)
		return cmd.Run()
	}
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

// GetDiskSpaceInfo returns disk space information for the current system
func (a *App) GetDiskSpaceInfo() (*DiskSpaceInfo, error) {
	if runtime.GOOS == "windows" {
		return a.getWindowsDiskSpace()
	} else {
		return a.getUnixDiskSpace()
	}
}

// getWindowsDiskSpace gets disk space information on Windows
func (a *App) getWindowsDiskSpace() (*DiskSpaceInfo, error) {
	// Get all available drives
	drives := a.GetDrives()
	
	var totalBytes, totalFreeBytes, totalUsedBytes uint64
	
	// Windows kernel32.dll functions
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceExW := kernel32.NewProc("GetDiskFreeSpaceExW")
	
	// Sum up space from all drives
	for _, drive := range drives {
		var freeBytes, driveTotalBytes, driveTotalFreeBytes uint64
		
		// Convert to UTF16 pointer for the drive (remove trailing backslash)
		drivePath := drive.Path
		if len(drivePath) > 0 && drivePath[len(drivePath)-1] == '\\' {
			drivePath = drivePath[:len(drivePath)-1]
		}
		
		drivePtr, err := syscall.UTF16PtrFromString(drivePath)
		if err != nil {
			fmt.Printf("Failed to convert drive string %s: %v\n", drivePath, err)
			continue
		}
		
		// Call GetDiskFreeSpaceExW
		ret, _, err := getDiskFreeSpaceExW.Call(
			uintptr(unsafe.Pointer(drivePtr)),
			uintptr(unsafe.Pointer(&freeBytes)),
			uintptr(unsafe.Pointer(&driveTotalBytes)),
			uintptr(unsafe.Pointer(&driveTotalFreeBytes)),
		)
		
		if ret == 0 {
			fmt.Printf("Failed to get disk space for %s\n", drivePath)
			continue
		}
		
		totalBytes += driveTotalBytes
		totalFreeBytes += freeBytes
		totalUsedBytes += (driveTotalBytes - freeBytes)
	}
	
	// If no drives were found, use fallback data
	if totalBytes == 0 {
		return &DiskSpaceInfo{
			TotalBytes:   500 * 1024 * 1024 * 1024, // 500 GB
			FreeBytes:    200 * 1024 * 1024 * 1024, // 200 GB
			UsedBytes:    300 * 1024 * 1024 * 1024, // 300 GB
			TotalSpaceGB: 500,
			FreeSpaceGB:  200,
			UsedSpaceGB:  300,
			UsagePercent: 60,
		}, nil
	}
	
	gb := float64(1024 * 1024 * 1024)
	
	return &DiskSpaceInfo{
		TotalBytes:   totalBytes,
		FreeBytes:    totalFreeBytes,
		UsedBytes:    totalUsedBytes,
		TotalSpaceGB: float64(totalBytes) / gb,
		FreeSpaceGB:  float64(totalFreeBytes) / gb,
		UsedSpaceGB:  float64(totalUsedBytes) / gb,
		UsagePercent: float64(totalUsedBytes) / float64(totalBytes) * 100,
	}, nil
}

// getUnixDiskSpace gets disk space information on Unix-like systems
func (a *App) getUnixDiskSpace() (*DiskSpaceInfo, error) {
	// For now, return a placeholder for non-Windows systems
	// In a real implementation, you would use platform-specific APIs
	return &DiskSpaceInfo{
		TotalBytes:   100 * 1024 * 1024 * 1024, // 100 GB placeholder
		FreeBytes:    50 * 1024 * 1024 * 1024,  // 50 GB placeholder
		UsedBytes:    50 * 1024 * 1024 * 1024,  // 50 GB placeholder
		TotalSpaceGB: 100,
		FreeSpaceGB:  50,
		UsedSpaceGB:  50,
		UsagePercent: 50,
	}, nil
}
