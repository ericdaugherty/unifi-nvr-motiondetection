package nvr

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGetCameraID(t *testing.T) {

	line1 := "1548877926.436 2019-01-30 12:52:06.436/MST: INFO   Camera[F09FC22F4D1D] type:start event:8078 clock:10895263924 (Front Door) in ApplicationEvtBus-0"

	id, event := getCameraID(line1)

	if id != "F09FC22F4D1D" || event != "start" {
		t.Errorf("Expected start event for cameraId F09FC22F4D1D but got event,start: %v, %v", event, id)
	}

	line2 := "1548877946.785 2019-01-30 12:52:26.785/MST: INFO   Camera[F09FC22F4D1D] type:stop event:8078 clock:10895284314 (Front Door) in ApplicationEvtBus-9"

	id, event = getCameraID(line2)

	if id != "F09FC22F4D1D" || event != "stop" {
		t.Errorf("Expected stop event for cameraId F09FC22F4D1D but got event,start: %v, %v", event, id)
	}
}

func TestDetectMotionNoFile(t *testing.T) {

	_, err := DetectMotion("/Non/existent/file/should/fail")
	if err == nil {
		t.Error("DetectMotion should return an error immediately for a non-existent file.")
	}
}

func TestDetectMotionOffset(t *testing.T) {

	fn, err := createTempFile(10)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(fn)

	c := make(chan string)
	md, err := DetectMotion(fn)
	if err != nil {
		t.Error(err)
		return
	}
	md.AddStartMotionCallback("NEWCAM", func(cameraId string, event string) {
		c <- cameraId
	})

	time.Sleep(1000)

	err = writeToFile(fn, "1548877926.436 2019-01-30 12:52:06.436/MST: INFO   Camera[NEWCAM] type:start event:8078 clock:10895263924 (Front Door) in ApplicationEvtBus-0\n")
	if err != nil {
		t.Error(err)
		return
	}

	id := <-c

	if id != "NEWCAM" {
		t.Errorf("Expected new motion line 'NEWCAM' to be read but got %v", id)
	}
}

func TestDetectMotionMultiIds(t *testing.T) {

	fn, err := createTempFile(10)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(fn)

	c := make(chan string)
	md, err := DetectMotion(fn)
	if err != nil {
		t.Error(err)
		return
	}
	md.AddStartMotionCallback("F09FC22F4D1D", func(cameraId string, event string) {
		c <- cameraId
	})
	md.AddStartMotionCallback("NEWCAM", func(cameraId string, event string) {
		c <- cameraId
	})

	time.Sleep(1000)

	err = writeToFile(fn, "1548877926.436 2019-01-30 12:52:06.436/MST: INFO   Camera[NEWCAM] type:start event:8078 clock:10895263924 (Front Door) in ApplicationEvtBus-0\n")
	if err != nil {
		t.Error(err)
		return
	}

	id := <-c

	if id != "NEWCAM" {
		t.Errorf("Expected new motion line 'NEWCAM' to be read but got %v", id)
	}
}

func TestDetectMotionStop(t *testing.T) {

	fn, err := createTempFile(10)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(fn)

	c := make(chan string)
	md, err := DetectMotion(fn)
	if err != nil {
		t.Error(err)
		return
	}
	md.AddStopMotionCallback("F09FC22F4D1D", func(cameraId string, event string) {
		c <- cameraId
	})

	err = writeToFile(fn, "1548877946.785 2019-01-30 12:52:26.785/MST: INFO   Camera[F09FC22F4D1D] type:stop event:8078 clock:10895284314 (Front Door) in ApplicationEvtBus-9\n")
	if err != nil {
		t.Error(err)
		return
	}

	id := <-c

	if id != "F09FC22F4D1D" {
		t.Errorf("Expected new motion line 'F09FC22F4D1D' to be read but got %v", id)
	}
}

func createTempFile(lines int) (string, error) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		return "", err
	}

	for i := 0; i < lines; i++ {
		_, err := f.WriteString(fmt.Sprintf("1548877926.436 2019-01-30 12:52:06.436/MST: INFO   Camera[F09FC22F4D1D] type:start event:8078 clock:10895263924 (Front Door) in ApplicationEvtBus-0 Line: %v\n", i))
		if err != nil {
			return "", err
		}
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func writeToFile(fn string, data string) error {

	f, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	_, err = f.WriteString(data)
	f.Close()
	return err
}
