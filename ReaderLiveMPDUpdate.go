package dashreader

import (
	"fmt"
	"log"
	"reflect"
)

//readerLiveMPDUpdate - Implement Reader of MPD
//  * Live
//  * MPD Updating
//  * SegmentTimeLine
//  * $Time$ based url
//  * $Number$ based url
type readerLiveMPDUpdate struct {
	readerBaseExtn
}

//MakeDASHReaderContext - Makes Reader Context
// Parameters:
//   1: Context received earlier... if first time pass nil
//   2: StreamSelector for the ContentType to select AdaptationSet
//   3: RepresentationSelector ... selector for Representation
// Return:
//   1: Context for current AdaptationSet,Representation
//   2: error
func (r *readerLiveMPDUpdate) MakeDASHReaderContext(rdrCtx ReaderContext, streamSelector StreamSelector, repSelector RepresentationSelector) (ReaderContext, error) {
	var curContext readerLiveMPDUpdateContext
	if rdrCtx != nil {
		v := rdrCtx.(*readerLiveMPDUpdateContext)
		curContext = *v
	} else {
		curContext = readerLiveMPDUpdateContext{
			readerBaseContext: readerBaseContext{
				ID:             r.ID,
				adaptSetID:     0,
				repID:          "",
				updCounter:     0,
				repSelector:    repSelector,
				streamSelector: streamSelector,
				StatzAgg:       r.StatzAgg,
			},
		}
	}
	curMpd, updCounter := r.readerBaseExtn.checkUpdate()
	if reflect.TypeOf(curContext.repSelector) != reflect.TypeOf(repSelector) {
		curContext.repSelector = repSelector
	}
	if reflect.TypeOf(curContext.streamSelector) != reflect.TypeOf(streamSelector) {
		curContext.streamSelector = streamSelector
	}
	if rdrCtx != nil {
		if updCounter == curContext.updCounter {
			//no update
			return &curContext, nil
		}
		err := curContext.adjustRepUpdate(r.readerBase, curMpd)
		if err == nil {
			return &curContext, nil
		}
		//Gaps TBD
		log.Printf("Adjust Rep Update Fail : %v", err)
	}
	//Incoming context is nil = new context
	//Locate the livePoint
	err := curContext.livePointLocate(r.readerBase, curMpd)
	if err != nil {
		return &curContext, fmt.Errorf("LivePoint Locate Failed: %w", err)
	}
	return &curContext, nil
}
