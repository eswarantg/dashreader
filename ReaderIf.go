package dashreader

import (
	"context"
	"net/url"
	"time"

	"github.com/eswarantg/statzagg"
)

//ChunkURL - URL extracted from MPD for playback
type ChunkURL struct {
	//ChunkURL - Actual URL
	ChunkURL url.URL
	//Range - Range Header
	Range string
	//FetchAt - WallClock Time when URL becomes available
	FetchAt time.Time
	//Duration - Duration of content available in this URL
	Duration time.Duration
}

//ChunkURLChannel - Channel of Chunk URLs
type ChunkURLChannel chan ChunkURL

//ReaderContext - Unique data for each Reader
type ReaderContext interface {
	//NextURL -
	//-- Once end is reached (io.EOF)
	//-- MakeDASHReaderContext has to be called again
	// Parameters;
	//   None
	// Return:
	//   1: Next URL
	//   2: error
	NextURL() (*ChunkURL, error)

	//GetURLs - Get URLs from Current MPD context
	//-- Once end of this list is reached
	//-- MakeDASHReaderContext has to be called again
	// Parameters;
	//   context for cancellation
	// Return:
	//   1: Channel of URLs, can be read till closed
	//   2: error
	NextURLs(context.Context) (<-chan ChunkURL, error)

	//GetFramerate - Framerate of content
	GetFramerate() float64
	//GetContentType - Content Type of content
	GetContentType() string
	//GetLang - Lang of content
	GetLang() string
	//GetCodecs - Codecs of content
	GetCodecs() string
}

//Reader - Read any DASH file and get Playback URLs
type Reader interface {
	//Update -
	// Parameters:
	//   MPD read
	// Return:
	//   1: MPD Updated - PublishTime Updated?
	//   2: error
	Update(*MPDtype) (bool, error)

	//MakeDASHReaderContext - Makes Reader Context
	// Parameters:
	//   1: Context received earlier... if first time pass nil
	//   2: StreamSelector for the ContentType to select AdaptationSet
	//   3: RepresentationSelector ... selector for Representation
	// Return:
	//   1: Context for current AdaptationSet,Representation
	//   2: error
	MakeDASHReaderContext(ReaderContext, StreamSelector, RepresentationSelector) (ReaderContext, error)

	//SetStatzAgg - Set StatzAgg for event forwarding
	// Parameters;
	//   StatzAgg
	// Return:
	//   NA
	SetStatzAgg(statzAgg statzagg.StatzAgg)
}
