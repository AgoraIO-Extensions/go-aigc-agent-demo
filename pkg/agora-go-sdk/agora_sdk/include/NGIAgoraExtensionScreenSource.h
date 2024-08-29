//
//  Agora SDK
//
//  Copyright (c) 2021 Agora.io. All rights reserved.
//

#pragma once  // NOLINT(build/header_guard)

#include "AgoraBase.h"
#include "IAgoraLog.h"
#include "AgoraRefPtr.h"
#include "NGIAgoraVideoFrame.h"
#include "AgoraMediaBase.h"

namespace agora {
namespace rtc {

struct ScreenCaptureProfilingStatistics {
  int capture_type;
  uint32_t captured_frame_width;
  uint32_t captured_frame_height;
  uint32_t total_captured_frames;
  uint64_t per_frame_cap_time_ms;
  uint64_t per_capture_cpu_cycles;
  bool capture_mouse_cursor;

  ScreenCaptureProfilingStatistics()
    : capture_type(-1), captured_frame_width(0), captured_frame_height(0),
      total_captured_frames(0), per_frame_cap_time_ms(0),
      per_capture_cpu_cycles(0), capture_mouse_cursor(true) {}
};

class IScreenCaptureSource : public RefCountInterface {
 public:
  class Control : public RefCountInterface {
  public:
    virtual int postEvent(const char* key, const char* value) = 0;
    virtual void printLog(commons::LOG_LEVEL level, const char* format, ...) = 0;
    virtual int pushAudioFrame(const media::IAudioFrameObserver::AudioFrame& captured_frame) = 0;
    virtual bool timeToPushVideo() = 0;
    virtual int pushVideoFrame(const agora::agora_refptr<IVideoFrame>& captured_frame) = 0;
    virtual agora::agora_refptr<IVideoFrameMemoryPool> getMemoryPool() = 0;
  };

  struct AudioCaptureConfig {
    uint32_t volume;
    int sample_rate_hz;
    int num_channels;
    AudioCaptureConfig() : volume(0), sample_rate_hz(0), num_channels(0) {}
  };

#if defined (__ANDROID__) || (defined(TARGET_OS_IPHONE) && TARGET_OS_IPHONE)
  struct VideoCaptureConfig {
    agora::rtc::VideoDimensions dimensions;
    VideoCaptureConfig()
      : dimensions(640, 360) {}
  };
#else
  struct VideoCaptureConfig {
    enum CaptureType {
      CaptureWindow,
      CaptureScreen,
    };
    CaptureType type;
    Rectangle screen_rect;
    Rectangle region_offset;
    uint32_t display_id;
    view_t window_id;

    VideoCaptureConfig()
      : type(CaptureScreen), screen_rect(), region_offset(), display_id(0), window_id(NULL) { }
  };
#endif
  enum CaptureMode {
    kPull, // SDK needs to poll the captured frame actively
    kPush // Capture source pushes the captured frame to sdk
  };
  #if defined(_WIN32)
  enum VideoContentSubType {
  UNSPECIFIED = 0,  // if camera, trade as camera; if share trade as document
  SHARE_DOCUMENT = 1,
  SHARE_GAMING = 2,
  SHARE_VIDEO = 3,
  SHARE_RDC = 4,   // remote desktop control
  SHARE_HFHD = 5,  // high frame-rate high definition screen share
  MAX = 16,
};
#endif 


  virtual ~IScreenCaptureSource() {}

  virtual int initializeCapture(const agora_refptr<Control>& control) = 0;

  // Start video capture interface for desktop capturing
  virtual int startVideoCapture(const VideoCaptureConfig& config) = 0;
  virtual int stopVideoCapture() = 0;

  virtual CaptureMode getVideoCaptureMode() = 0;

  // Implementation of the following interfaces are not mandatory
  virtual int startAudioCapture(const AudioCaptureConfig& config) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int stopAudioCapture() {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int setAudioVolume(uint32_t volume) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int setFrameRate(int fps) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int setScreenCaptureDimensions(const agora::rtc::VideoDimensions& dimensions) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int updateCaptureRegion(const agora::rtc::Rectangle& captureRegion) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int setExcludeWindowList(void* const * handles, int count) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int captureMouseCursor(bool capture) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int capture(agora::agora_refptr<IVideoFrame>& frame) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int getProfilingStats(ScreenCaptureProfilingStatistics& stats) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int getScreenDimensions(int& width, int& height) {
    return ERR_NOT_SUPPORTED;
  }
  virtual int setProperty(const char* key, const char* json_value) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int setCustomContext(const char* key, const void* context) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual int getProperty(const char* key, char* json_value, int& length) {
    return -ERR_NOT_SUPPORTED;
  }
  virtual void* getCustomContext(const char* key) {
    return NULL;
  }
  virtual void* getScreenCaptureSources(int thumb_cx, int thumb_cy, int icon_cx, int icon_cy,
                                                            bool include_screen) {
    return NULL;
  }
#if defined(_WIN32)
  virtual int SetContentType(VideoContentSubType type) {
    return -ERR_NOT_SUPPORTED;
  }
#endif 
};

}  // namespace rtc
}  // namespace agora
