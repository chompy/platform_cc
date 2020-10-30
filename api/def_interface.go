package api

// DefInterface - definition interface
type DefInterface interface {
	Validate() []error
	//SetDefaults()
}
