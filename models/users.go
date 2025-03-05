package models

type User struct {
	Model
	Name,
	Email string
	Pass []byte `json:"-" c3po:"-"`
}

func (u *User) ToMAP() map[string]any {
	return map[string]any{
		"uuid":  u.UUID,
		"name":  u.Name,
		"email": u.Email,
	}
}
