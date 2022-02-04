package scraper

type YtInitialData struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer struct {
							Contents Contents `json:"contents"`
						} `json:"ItemSectionRenderer"`
					} `json:"contents"`
				} `json:"SectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
	EstimatedResults string `json:"estimatedResults"`
}

type Contents []struct {
	VideoRenderer VideoRenderer `json:"videoRenderer"`
}

type VideoRenderer struct {
	Badges []struct {
		MetadataBadgeRenderer struct {
			Style          string `json:"style"`
			Label          string `json:"label"`
			TrackingParams string `json:"trackingParams"`
		} `json:"metadataBadgeRenderer"`
	} `json:"badges"`
	ChannelThumbnailSupportedRenderers struct {
		ChannelThumbnailWithLinkRenderer struct {
			Accessibility struct {
				AccessibilityData struct {
					Label string `json:"label"`
				} `json:"accessibilityData"`
			} `json:"accessibility"`
			NavigationEndpoint struct{} `json:"navigationEndpoint"`
			Thumbnail          struct {
				Thumbnails []struct {
					Url    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
		} `json:"ChannelThumbnailWithLinkRenderer"`
	} `json:"ChannelThumbnailSupportedRenderers"`
	DetailedMetadataSnippets []struct{} `json:"detailedMetadataSnippets"`
	IsWatched                bool       `json:"isWatched"`
	LengthText               struct {
		SimpleText string `json:"simpleText"`
	} `json:"lengthText"`
	LongBylineText     struct{} `json:"longBylineText"`
	Menu               struct{} `json:"menu"`
	NavigationEndpoint struct{} `json:"navigationEndpoint"`
	OwnerText          struct {
		Runs []struct {
			Text string
		}
	} `json:"ownerText"`
	PublishedTimeText  struct{} `json:"publishedTimeText"`
	RichThumbnail      struct{} `json:"richThumbnail"`
	ShortBylineText    struct{} `json:"shortBylineText"`
	ShortViewCountText struct {
		SimpleText string `json:"simpleText"`
	} `json:"shortViewCountText"`
	ShowActionMenu bool `json:"showActionMenu"`
	Thumbnail      struct {
		Thumbnails []struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"thumbnails"`
	} `json:"thumbnail"`
	ThumbnailOverlays []struct{} `json:"thumbnailOverlays"`
	Title             struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	VideoId       string `json:"videoId"`
	ViewCountText struct {
		SimpleText string `json:"simpleText"`
	} `json:"viewCountText"`
}