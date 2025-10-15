package entity

type Event struct {
	EventName               string `json:"event_name"`
	EventBannerPresignedURL string `json:"event_banner_presigned_url"`
}
