package models

import "errors"

// ErrNotFound is used when document could not be found.
var ErrNotFound = errors.New("AdNetworkList does not exists")

// ErrUpsertFailed is used when document was not successfully upserted
var ErrUpsertFailed = errors.New("AdNetworkList failed to be upserted")
