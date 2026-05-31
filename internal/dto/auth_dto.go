// internal/dto/auth_dto.go
package dto

import "time"

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ProfileResponse dengan data detail (Guru/Siswa)
type ProfileResponse struct {
	User    UserResponse `json:"user"`
	Profile interface{}  `json:"profile"` // Bisa GuruResponse atau SiswaResponse
	Role    string       `json:"role"`
}

type GuruProfileResponse struct {
	GuruID          uint       `json:"guru_id"`
	GuruNama        string     `json:"guru_nama"`
	GuruNIP         *string    `json:"guru_nip,omitempty"`
	GuruNUPTK       *string    `json:"guru_nuptk,omitempty"`
	GuruJenkel      string     `json:"guru_jenkel"`
	GuruTempatLahir *string    `json:"guru_tempat_lahir,omitempty"`
	GuruTglLahir    *time.Time `json:"guru_tgl_lahir,omitempty"`
	GuruEmail       *string    `json:"guru_email,omitempty"`
	GuruNoTelp      *string    `json:"guru_no_telp,omitempty"`
	GuruPhoto       *string    `json:"guru_photo,omitempty"`
	StatusAktif     bool       `json:"status_aktif"`
}
