package gdit

import "sync"

type ProviderMap struct {
	m sync.Map
}

func (pm *ProviderMap) Store() {
}

func (pm *ProviderMap) Load() {

}
