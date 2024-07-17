// Package searchseries response from /data/results/search_series
package searchseries

type SearchSeriesResults struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}
type ChunkInfo struct {
	ChunkSize       int      `json:"chunk_size"`
	NumChunks       int      `json:"num_chunks"`
	Rows            int      `json:"rows"`
	BaseDownloadURL string   `json:"base_download_url"`
	ChunkFileNames  []string `json:"chunk_file_names"`
}
type Params struct {
	CategoryIds   []int `json:"category_ids"`
	SeriesID      int   `json:"series_id"`
	SeasonYear    int   `json:"season_year"`
	SeasonQuarter int   `json:"season_quarter"`
	OfficialOnly  bool  `json:"official_only"`
	EventTypes    []int `json:"event_types"`
}
type Data struct {
	Success   bool      `json:"success"`
	ChunkInfo ChunkInfo `json:"chunk_info"`
	Params    Params    `json:"params"`
}
