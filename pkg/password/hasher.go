package password

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func NewHasher() *Hasher {
	return &Hasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *Hasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), h.cost)
}

func (h *Hasher) Compare(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
