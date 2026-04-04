package app

import (
	"github.com/speech/fireworks-admin/internal/features/teltent"
)

// Registrars 内部结构体，用于收集所有 RouterRegistrar 实现。
type RegistrarIn struct {
	Teltent *teltent.Handler
	Health  *HealthRouter
}

// ProvideRegistrars 将所有 RouterRegistrar 实现收集为切片。
func ProvideRegistrars(r RegistrarIn) []RouterRegistrar {
	return []RouterRegistrar{
		r.Teltent,
		r.Health,
	}
}
