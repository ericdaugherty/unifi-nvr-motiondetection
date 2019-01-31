# UniFi NVR Motion Detector
[![Go Report Card](https://goreportcard.com/badge/github.com/ericdaugherty/imagefetcher)](https://goreportcard.com/report/github.com/ericdaugherty/imagefetcher)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/ericdaugherty/unifi-nvr-motiondetection/blob/master/LICENSE)

This package monitors the monitor.log file written by the Ubiquiti UniFi NVR, 
and triggers a callback method whenever motion is detected. You can trigger on 
either the start of the motion or the end of the motion.

This package was inspired by https://github.com/mzac/unifi-video-mqtt

