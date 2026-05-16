package dto

import (
	"time"

	"github.com/entertrans/go-base-project.git/internal/model"
)

/* =========================
   Responses
========================= */

type LampiranResponse struct {
	IDLampiran      uint64    `json:"id_lampiran"`
	SiswaNIS        string    `json:"siswa_nis"`
	DokumenJenis    string    `json:"dokumen_jenis"`
	StorageProvider string    `json:"storage_provider"`
	ObjectKey       string    `json:"object_key"` // untuk gdrive: file_id
	FileName        string    `json:"file_name"`
	MimeType        string    `json:"mime_type"`
	SizeBytes       int64     `json:"size_bytes"`
	ETag            string    `json:"etag"` // gdrive: kosong
	UploadedAt      time.Time `json:"uploaded_at"`
	ViewURL         string    `json:"view_url"`
}

type SiswaLampiranResponse struct {
	SiswaID   uint   `json:"siswa_id"`
	SiswaNIS  string `json:"siswa_nis"`
	SiswaNama string `json:"siswa_nama"`
	Lampiran  []LampiranResponse `json:"lampiran"`
}

type DeleteLampiranResponse struct {
	Message      string `json:"message"`
	SiswaNIS     string `json:"siswa_nis"`
	DokumenJenis string `json:"dokumen_jenis"`
	ObjectKey    string `json:"object_key"`
}

/* =========================
   Mappers
========================= */

func MapLampiranToResponse(m model.Lampiran) LampiranResponse {
	return LampiranResponse{
		IDLampiran:      m.IDLampiran,
		SiswaNIS:        m.SiswaNIS,
		DokumenJenis:    m.DokumenJenis,
		StorageProvider: m.StorageProvider,
		ObjectKey:       m.ObjectKey,
		FileName:        m.FileName,
		MimeType:        m.MimeType,
		SizeBytes:       m.SizeBytes,
		ETag:            m.ETag,
		UploadedAt:      m.UploadedAt,
		ViewURL:         buildViewURL(m.IDLampiran),
	}
}

func MapSiswaLampiranToResponse(s model.Siswa, list []model.Lampiran) SiswaLampiranResponse {
	out := SiswaLampiranResponse{
		SiswaID: s.SiswaID,
	}

	if s.SiswaNIS != nil {
		out.SiswaNIS = *s.SiswaNIS
	}
	if s.SiswaNama != nil {
		out.SiswaNama = *s.SiswaNama
	}

	out.Lampiran = make([]LampiranResponse, 0, len(list))
	for _, l := range list {
		out.Lampiran = append(out.Lampiran, MapLampiranToResponse(l))
	}
	return out
}

func buildViewURL(id uint64) string {
	// route kamu: /api/v1/lampiran/file/:id/view
	return "/lampiran/file/" + uintToString(id) + "/view"
}

func uintToString(v uint64) string {
	if v == 0 {
		return "0"
	}
	var b [32]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + (v % 10))
		v /= 10
	}
	return string(b[i:])
}
