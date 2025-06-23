package ethereum

import v2 "github.com/blockopsnetwork/telescope/internal/static/integrations/v2"

func init() {
	v2.Register(&Config{}, v2.TypeMultiplex)
}
