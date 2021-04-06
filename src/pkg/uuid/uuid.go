// FIXME: doc me
package uuid

import (
	"github.com/google/uuid"
)

func GenerateUuid() string {
	return uuid.NewString()
}
