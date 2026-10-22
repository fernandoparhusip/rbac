package password

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(value string) (string, error)
	Compare(hash, value string) error
}

type Bcrypt struct {
	cost int
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{
		cost: bcrypt.DefaultCost,
	}
}

func (b *Bcrypt) Hash(value string) (string, error) {
	out, err := bcrypt.GenerateFromPassword(
		[]byte(value),
		b.cost,
	)

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (b *Bcrypt) Compare(hash, value string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(value),
	)
}
