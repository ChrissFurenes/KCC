package RAMDisk

import "sync"

type RAMDisk struct {
	mu    sync.Mutex
	files map[string][]byte
}

func NewRAMDisk() *RAMDisk {
	return &RAMDisk{
		files: make(map[string][]byte),
	}
}

func (r *RAMDisk) WriteFile(name string, data []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.files[name] = append(r.files[name], data...)
}
func (r *RAMDisk) ReadFile(name string) ([]byte, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, exists := r.files[name]
	if !exists {
		return nil, false
	}
	return append([]byte(nil), data...), true
}
func (r *RAMDisk) DeleteFile(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.files, name)
}
