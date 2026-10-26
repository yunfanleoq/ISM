package ISMScript

import (
	"regexp"
	"sync"
)

var (
	reDepBitGet        = regexp.MustCompile(`BitGet\s*\(\s*"([^"]+)"`)
	reDepGetDeviceData = regexp.MustCompile(`GetDeviceData\s*\(\s*"([^"]+)"`)
	reDepGetRealData   = regexp.MustCompile(`GetDeviceRealData\s*\(\s*"([^"]+)"`)
)

// ExtractScriptDeps returns distinct "device->point" literals referenced by the script.
func ExtractScriptDeps(content string) []string {
	seen := make(map[string]struct{})
	var deps []string
	add := func(m [][]string) {
		for _, g := range m {
			if len(g) < 2 {
				continue
			}
			key := g[1]
			if key == "" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			deps = append(deps, key)
		}
	}
	add(reDepBitGet.FindAllStringSubmatch(content, -1))
	add(reDepGetDeviceData.FindAllStringSubmatch(content, -1))
	add(reDepGetRealData.FindAllStringSubmatch(content, -1))
	return deps
}

var (
	wakeMu     sync.RWMutex
	wakeByKey  = make(map[string][]chan struct{})
	wakeByUUID = make(map[string]chan struct{})
)

func clearScriptWakes() {
	wakeMu.Lock()
	defer wakeMu.Unlock()
	wakeByKey = make(map[string][]chan struct{})
	wakeByUUID = make(map[string]chan struct{})
}

func registerScriptWake(scriptUUID string, deps []string, ch chan struct{}) {
	wakeMu.Lock()
	defer wakeMu.Unlock()
	wakeByUUID[scriptUUID] = ch
	for _, dep := range deps {
		wakeByKey[dep] = append(wakeByKey[dep], ch)
	}
}

// WakeScriptsForKey non-blocking wakes all anko-onchange scripts that depend on key.
func WakeScriptsForKey(key string) {
	wakeMu.RLock()
	chs := wakeByKey[key]
	wakeMu.RUnlock()
	for _, ch := range chs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
