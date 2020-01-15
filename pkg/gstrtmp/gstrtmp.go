package gstrtmp

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0

#include "gst.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func init() {
	go C.gst_rtmp_start_mainloop()
}

type Pipeline struct {
	Pipeline *C.GstElement
}

func CreatePipeline(rtmpUrl string) *Pipeline {
	pipelineStr := fmt.Sprintf(`
appsrc is-live=true do-timestamp=true name=videosrc ! 
h264parse config-interval=-1 ! 
flvmux streamable=true name=mux ! 
queue ! 
rtmpsink sync=true location='%s live=1' 
appsrc is-live=true do-timestamp=true name=audiosrc ! 
queue ! 
opusparse ! 
opusdec ! 
audioconvert ! 
audioresample ! 
audio/x-raw, rate=48000 ! 
queue !  
voaacenc ! 
audio/mpeg, rate=48000, channels=2, mpegversion=4 ! 
aacparse ! 
queue ! 
mux. 
`, rtmpUrl)
	fmt.Println(pipelineStr)
	pipelineStrUnsafe := C.CString(pipelineStr)
	defer C.free(unsafe.Pointer(pipelineStrUnsafe))
	return &Pipeline{Pipeline: C.gst_rtmp_create_pipeline(pipelineStrUnsafe)}
}

func CreateVideoPipeline(rtmpUrl string) *Pipeline {
	pipelineStr := fmt.Sprintf(`
appsrc is-live=true do-timestamp=true name=videosrc ! 
h264parse config-interval=-1 ! 
flvmux streamable=true name=mux ! 
queue ! 
rtmpsink sync=true location='%s live=1'
`, rtmpUrl)
	fmt.Println(pipelineStr)
	pipelineStrUnsafe := C.CString(pipelineStr)
	defer C.free(unsafe.Pointer(pipelineStrUnsafe))
	return &Pipeline{Pipeline: C.gst_rtmp_create_pipeline(pipelineStrUnsafe)}
}

func (p *Pipeline) Start() {
	C.gst_rtmp_start_pipeline(p.Pipeline)
}

func (p *Pipeline) Stop() {
	C.gst_rtmp_stop_pipeline(p.Pipeline)
}

func (p *Pipeline) PushVideo(buffer []byte) {
	b := C.CBytes(buffer)
	defer C.free(unsafe.Pointer(b))
	C.gst_rtmp_push_video_buffer(p.Pipeline, b, C.int(1), C.int(len(buffer)))
}

func (p *Pipeline) PushAudio(buffer []byte) {
	b := C.CBytes(buffer)
	defer C.free(unsafe.Pointer(b))
	C.gst_rtmp_push_audio_buffer(p.Pipeline, b, C.int(1), C.int(len(buffer)))
}
