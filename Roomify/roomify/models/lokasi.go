package models

type CreateLocationRequest struct {
	NamaLokasi string   `json:"nama_lokasi"`
	Capacity   int      `json:"capacity"`
	Ruangan    []Ruangan `json:"ruangan"`
}

type Ruangan struct {
	NamaRuangan string `json:"nama_ruangan"`
	Capacity    int    `json:"capacity"`
}

type Location struct {
	ID         int    `json:"id"`
	NamaLokasi string `json:"nama_lokasi"`
	Capacity   int    `json:"capacity"`
	CreatedAt  string `json:"created_at"`
}

type DetailLocation struct {
	ID          int    `json:"id"`
	IDLokasi    int    `json:"id_lokasi"`
	NamaRuangan string `json:"nama_ruangan"`
	Capacity    int    `json:"capacity"`
}

type LocationWithDetails struct {
	Location Location         `json:"lokasi"`
	Details  []DetailLocation `json:"details"`
}