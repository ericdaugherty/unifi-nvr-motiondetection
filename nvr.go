package nvr

import (
	"os"
	"regexp"
	"sync"

	"github.com/hpcloud/tail"
)

var filter *regexp.Regexp

func init() {
	filter = regexp.MustCompile(`Camera\[(.+)\] type:(\w+) `)
}

// MotionDetector manages the callbacks for start and stop events.
type MotionDetector struct {
	startCallbacks sync.Map // map[string]func(string, string)
	stopCallbacks  sync.Map // map[string]func(string, string)
}

// DetectMotion starts monitoring the specified motion.log file for motion start and stop events.
// DetectMotion will return an error if the specified file does not exist.
// Call AddStartMotionCallback or AddStopMotionCallback to setup callback methods for each camera ID to monitor.
func DetectMotion(motionlog string) (*MotionDetector, error) {

	md := &MotionDetector{}

	fi, err := os.Stat(motionlog)
	if err != nil {
		return md, err
	}

	tc := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Location: &tail.SeekInfo{
			Offset: fi.Size(),
		},
	}

	t, err := tail.TailFile(motionlog, tc)
	if err != nil {
		return nil, err
	}

	go func() {
		for line := range t.Lines {
			matchedID, event := getCameraID(line.Text)
			var f interface{}
			var ok bool

			switch event {
			case "start":
				f, ok = md.startCallbacks.Load(matchedID)
			case "stop":
				f, ok = md.stopCallbacks.Load(matchedID)
			default:
				ok = false
			}

			if ok {
				f1 := f.(func(string, string))
				f1(matchedID, event)
			}
		}
	}()

	return md, nil
}

// AddStartMotionCallback adds the specified callback for 'start' motion events.
func (md *MotionDetector) AddStartMotionCallback(cameraID string, callback func(string, string)) {
	md.startCallbacks.Store(cameraID, callback)

}

// AddStopMotionCallback adds the specified callback for 'stop' motion events.
func (md *MotionDetector) AddStopMotionCallback(cameraID string, callback func(string, string)) {
	md.stopCallbacks.Store(cameraID, callback)
}

// RemoveStartMotionCallback removes the specified callback for 'start' motion events.
func (md *MotionDetector) RemoveStartMotionCallback(cameraID string) {
	md.startCallbacks.Delete(cameraID)
}

// RemoveStopMotionCallback removes the specified callback for 'start' motion events.
func (md *MotionDetector) RemoveStopMotionCallback(cameraID string) {
	md.stopCallbacks.Delete(cameraID)
}

func getCameraID(logline string) (string, string) {

	s := filter.FindStringSubmatch(logline)
	if len(s) > 2 {
		return s[1], s[2]
	}
	return "", ""
}
