package transport

type URL struct {
	host               string
	port               int
	parameter          map[string]string
	isSharable         bool
	needAuthorityCheck bool
	serialType         string
	providerApp        string
	loadBalance        string
}
