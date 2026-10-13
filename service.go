// Modified 2026-10-12: formatting and accurate placeholder documentation.

package gowin

// Service is a legacy placeholder. Windows service management is not implemented.
// Its fields and lifecycle methods remain private for compatibility.
type Service struct {
	name     string //
	dispName string //
	desc     string //
	exePath  string //
}

// TODO : make code for it
func (s *Service) install() (msj string, err error) {
	return
}

// TODO : make code for it
func (s *Service) remove() (msj string, err error) {
	return
}

// TODO : make code for it
func (s *Service) start() (msj string, err error) {
	return
}

// TODO : make code for it
func (s *Service) stop() (msj string, err error) {
	return
}

// TODO : make code for it
func (s *Service) restart() (msj string, err error) {
	return
}
