package jwk

import (
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teamhanko/hanko/backend/persistence/models"
	"github.com/teamhanko/hanko/backend/test"
	"testing"
)

type mockJwkPersister struct {
	jwks []models.Jwk
}

func (m *mockJwkPersister) Get(i int) (*models.Jwk, error) {
	for _, v := range m.jwks {
		if v.ID == i {
			return &v, nil
		}
	}
	return nil, nil
}

func (m *mockJwkPersister) GetAll() ([]models.Jwk, error) {
	return m.jwks, nil
}

func (m *mockJwkPersister) GetLast() (*models.Jwk, error) {
	index := len(m.jwks)
	return &m.jwks[index-1], nil
}

func (m *mockJwkPersister) Create(jwk models.Jwk) error {
	//increment id
	index := len(m.jwks)
	jwk.ID = index

	m.jwks = append(m.jwks, jwk)
	return nil
}

func TestDefaultManager(t *testing.T) {
	keys := []string{"asfnoadnfoaegnq3094intoaegjnoadjgnoadng", "apdisfoaiegnoaiegnbouaebgn982"}
	//persister := mockJwkPersister{jwks: []models.Jwk{}}
	persister := test.NewJwkPersister(nil)

	dm, err := NewDefaultManager(keys, persister)
	require.NoError(t, err)
	all, err := persister.GetAll()

	require.NoError(t, err)
	assert.Equal(t, 2, len(all))

	js, err := dm.GetPublicKeys()
	require.NoError(t, err)
	assert.Equal(t, 2, js.Len())

	sk, err := dm.GetSigningKey()
	require.NoError(t, err)

	token := jwt.New()
	token.Set("Payload", "isJustFine")
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, sk))
	require.NoError(t, err)

	// Get Public Key of signing key
	pk, err := sk.PublicKey()
	require.NoError(t, err)

	// Parse and Verify
	tokenParsed, err := jwt.Parse(signed, jwt.WithKey(jwa.RS256, pk))
	assert.NoError(t, err)
	assert.Equal(t, token, tokenParsed)
}
