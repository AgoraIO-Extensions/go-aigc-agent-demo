//
//  Agora C SDK
//
//  Created by Tommy Miao in 2020.5
//  Copyright (c) 2020 Agora.io. All rights reserved.
//
#pragma once

#include <stdint.h>
#include <stdlib.h>

#include "agora_api.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

typedef unsigned int uid_t;
typedef unsigned int track_id_t;
typedef unsigned int conn_id_t;

typedef void* view_t;
typedef const char* user_id_t;

typedef struct _encoded_video_frame_info encoded_video_frame_info;

#define k_max_codec_name_len 50

typedef struct _audio_parameters {
  int sample_rate;
  size_t channels;
  size_t frames_per_buffer;
} audio_parameters;

/** The packet information, i.e., infomation contained in RTP heaader.
*/
typedef struct _packet_options {
  // Rtp timestamp
  uint32_t timestamp;
  // Audio level indication.
  uint8_t audio_level_indication;
} packet_options;

typedef struct _external_video_frame {
  /** The buffer type. See #VIDEO_BUFFER_TYPE
   */
  int type;
  /** The pixel format. See #VIDEO_PIXEL_FORMAT
   */
  int format;
  /** The video buffer.
   */
  void* buffer;
  /** Line spacing of the incoming video frame, which must be in pixels instead of bytes. For
   * textures, it is the width of the texture.
   */
  int stride;
  /** Height of the incoming video frame.
   */
  int height;
  /** [Raw data related parameter] The number of pixels trimmed from the left. The default value is
   * 0.
   */
  int crop_left;
  /** [Raw data related parameter] The number of pixels trimmed from the top. The default value is
   * 0.
   */
  int crop_top;
  /** [Raw data related parameter] The number of pixels trimmed from the right. The default value is
   * 0.
   */
  int crop_right;
  /** [Raw data related parameter] The number of pixels trimmed from the bottom. The default value
   * is 0.
   */
  int crop_bottom;
  /** [Raw data related parameter] The clockwise rotation of the video frame. You can set the
   * rotation angle as 0, 90, 180, or 270. The default value is 0.
   */
  int rotation;
  /** Timestamp of the incoming video frame (ms). An incorrect timestamp results in frame loss or
   * unsynchronized audio and video.
   */
  long long timestamp;
  /**
   * [Texture-related parameter]
   * When using the OpenGL interface (javax.microedition.khronos.egl.*) defined by Khronos, set EGLContext to this field.
   * When using the OpenGL interface (android.opengl.*) defined by Android, set EGLContext to this field.
   */
  void *egl_context;
  /**
   * [Texture related parameter] Texture ID used by the video frame.
   */
  int egl_type;
  /**
   * [Texture related parameter] Incoming 4 &times; 4 transformational matrix. The typical value is a unit matrix.
   */
  int texture_id;
  /**
   * [Texture related parameter] Incoming 4 &times; 4 transformational matrix. The typical value is a unit matrix.
   */
  float matrix[16];
  /**
   * [Texture related parameter] The MetaData buffer.
   *  The default value is NULL
   */
  uint8_t* metadata_buffer;
  /**
   * [Texture related parameter] The MetaData size.
   *  The default value is 0
   */
  int metadata_size;
  /** The alpha buffer.
   */
  void* alpha_buffer;
} external_video_frame;

/** Definition of VideoFrame.

The video data format is in YUV420. The buffer provides a pointer to a pointer. However, the
interface cannot modify the pointer of the buffer, but can only modify the content of the buffer.

*/
typedef struct _video_frame {
  int type;
  /** Video pixel width.
   */
  int width;  // width of video frame
  /** Video pixel height.
   */
  int height;  // height of video frame
  /** Line span of Y buffer in YUV data.
   */
  int y_stride;  // stride of Y data buffer
  /** Line span of U buffer in YUV data.
   */
  int u_stride;  // stride of U data buffer
  /** Line span of V buffer in YUV data.
   */
  int v_stride;  // stride of V data buffer
  /** Pointer to the Y buffer pointer in the YUV data.
   */
  uint8_t* y_buffer;  // Y data buffer
  /** Pointer to the U buffer pointer in the YUV data.
   */
  uint8_t* u_buffer;  // U data buffer
  /** Pointer to the V buffer pointer in the YUV data
   */
  uint8_t* v_buffer;  // V data buffer
  /** Set the rotation of this frame before rendering the video, and it supports 0, 90, 180, 270
   * degrees clockwise.
   */
  int rotation;  // rotation of this frame (0, 90, 180, 270)
  /** Timestamp to render the video stream. It instructs the users to use this timestamp to
  synchronize the video stream render while rendering the video streams.

  Note: This timestamp is for rendering the video stream, not for capturing the video stream.
  */
  int64_t render_time_ms;
  int avsync_type;
  /**
   * [Texture related parameter] The MetaData buffer.
   *  The default value is NULL
   */
  uint8_t* metadata_buffer;
  /**
   * [Texture related parameter] The MetaData size.
   *  The default value is 0
   */
  int metadata_size;
  /**
   * [Texture related parameter], egl context.
   */
  void* shared_context;
  /**
   * [Texture related parameter], Texture ID used by the video frame.
   */
  int texture_id;
  /**
   * [Texture related parameter], Incoming 4 &times; 4 transformational matrix.
   */
  float matrix[16];
  /**
   *  Portrait Segmentation meta buffer, dimension of which is the same as VideoFrame.
   *  Pixl value is between 0-255, 0 represents totally background, 255 represents totally foreground.
   *  The default value is NULL
   */
  uint8_t* alpha_buffer;
} video_frame;

// Stereo, 32 kHz, 60 ms (2 * 32 * 60)
#define k_max_data_size_samples 3840

typedef struct _audio_pcm_frame {
  uint32_t capture_timestamp;
  uint32_t samples_per_channel;
  int sample_rate_hz;
  uint32_t num_channels;
  uint32_t bytes_per_sample;
  int16_t data[k_max_data_size_samples];
} audio_pcm_frame;

// '2' is the size of above audio_pcm_frame.data[0]
#define k_max_data_size_bytes (k_max_data_size_samples * 2)

/**
 * Sync Video Frame Observer (the other agora::media::base::IVideoFrameObserver will only be used by Media Player)
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _video_frame_observer {
  /* return value stands for a 'bool' in C++: 1 for success, 0 for failure */
  int (*on_capture_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_pre_encode_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_secondary_camera_capture_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_secondary_pre_encode_camera_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_screen_capture_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_pre_encode_screen_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_media_player_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame, int media_player_id);
  int (*on_secondary_screen_sapture_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_secondary_pre_encode_screen_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
  int (*on_render_video_frame)(AGORA_HANDLE agora_video_frame_observer, const char* channel_id, uid_t uid, video_frame* frame);
  int (*on_transcoded_video_frame)(AGORA_HANDLE agora_video_frame_observer, video_frame* frame);
} video_frame_observer;

/* 'observer' could be the address of a local variable since we will save it by value */
/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 * @ANNOTATION:CTOR:agora_video_frame_observer
 */
AGORA_API_C_HDL agora_video_frame_observer_create(video_frame_observer* observer);

/* destroy before exiting SDK */
/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 * @ANNOTATION:DTOR:agora_video_frame_observer
 */
AGORA_API_C_VOID agora_video_frame_observer_destroy(AGORA_HANDLE agora_video_frame_observer);

/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 */
AGORA_API_C_INT agora_video_frame_observer_get_video_frame_process_mode(AGORA_HANDLE agora_video_frame_observer);

/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 */
AGORA_API_C_INT agora_video_frame_observer_get_video_pixel_format_preference(AGORA_HANDLE agora_video_frame_observer);

/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 */
AGORA_API_C_INT agora_video_frame_observer_get_rotation_applied(AGORA_HANDLE agora_video_frame_observer);

/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 */
AGORA_API_C_INT agora_video_frame_observer_get_mirror_applied(AGORA_HANDLE agora_video_frame_observer);

/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 */
AGORA_API_C_SIZE_T agora_video_frame_observer_get_observed_frame_position(AGORA_HANDLE agora_video_frame_observer);

/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 */
AGORA_API_C_INT agora_video_frame_observer_is_external(AGORA_HANDLE agora_video_frame_observer);

/**
 * @ANNOTATION:GROUP:agora_video_frame_observer
 */
AGORA_API_C_INT agora_video_frame_observer_through_the_interface(AGORA_HANDLE agora_video_frame_observer);

typedef struct _audio_frame {
  int type;
  /**
   * The number of samples per channel in this frame.
   */
  int samples_per_channel;
  /**
   * The number of bytes per sample: Two for PCM 16.
   */
  int bytes_per_sample;
  /**
   * The number of channels (data is interleaved, if stereo).
   */
  int channels;
  /**
   * The Sample rate.
   */
  int samples_per_sec;
  /**
   * The pointer to the data buffer.
   */
  void* buffer;
  /**
   * The timestamp to render the audio data. Use this member to synchronize the audio renderer while
   * rendering the audio streams.
   *
   * @note
   * This timestamp is for audio stream rendering. Set it as 0.
   */
  int64_t render_time_ms;
  int avsync_type;
} audio_frame;

typedef enum _RAW_AUDIO_FRAME_OP_MODE_TYPE {
  RAW_AUDIO_FRAME_OP_MODE_READ_ONLY = 0,
  RAW_AUDIO_FRAME_OP_MODE_READ_WRITE = 2
} RAW_AUDIO_FRAME_OP_MODE_TYPE;

typedef enum _AUDIO_FRAME_POSITION {
  AUDIO_FRAME_POSITION_PLAYBACK = 0x0001,
  AUDIO_FRAME_POSITION_RECORD = 0x0002,
  AUDIO_FRAME_POSITION_MIXED = 0x0004,
  AUDIO_FRAME_POSITION_BEFORE_MIXING = 0x0008,
} AUDIO_FRAME_POSITION;

typedef struct _audio_params {
  int sample_rate;
  int channels;
  RAW_AUDIO_FRAME_OP_MODE_TYPE mode;
  int samples_per_call;
} audio_params;

/**
 * Async Audio Frame Observer (the other agora::media::base::IAudioFrameObserver will only be used by Media Player)
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _audio_frame_observer {
  /* return value stands for a 'bool' in C++: 1 for success, 0 for failure */
  int (*on_record_audio_frame)(AGORA_HANDLE agora_local_user /* raw pointer */,const char* channelId, const audio_frame* frame);
  int (*on_playback_audio_frame)(AGORA_HANDLE agora_local_user, const char* channelId, const audio_frame* frame);
  int (*on_mixed_audio_frame)(AGORA_HANDLE agora_local_user, const char* channelId, const audio_frame* frame);
  int (*on_ear_monitoring_audio_frame)(AGORA_HANDLE agora_local_user, const audio_frame* frame);
  int (*on_playback_audio_frame_before_mixing)(AGORA_HANDLE agora_local_user, const char* channelId, user_id_t uid, const audio_frame* frame);
  int (*on_get_audio_frame_position)(AGORA_HANDLE agora_local_user);
  audio_params (*on_get_playback_audio_frame_param)(AGORA_HANDLE agora_local_user);
  audio_params (*on_get_record_audio_frame_param)(AGORA_HANDLE agora_local_user);
  audio_params (*on_get_mixed_audio_frame_param)(AGORA_HANDLE agora_local_user);
  audio_params (*on_get_ear_monitoring_audio_frame_param)(AGORA_HANDLE agora_local_user);
} audio_frame_observer;

typedef struct _audio_spectrum_data {
  /**
   * The audio spectrum data of audio.
   */
  const float *audio_spectrum_data;
  /**
   * The data length of audio spectrum data.
   */
  int data_length;
} audio_spectrum_data;

typedef struct _user_audio_spectrum_info  {
  /**
   * User ID of the speaker.
   */
  uid_t uid;
  /**
   * The audio spectrum data of audio.
   */
  audio_spectrum_data spectrum_data;
} user_audio_spectrum_info;

/**
 * The IAudioSpectrumObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _audio_spectrum_observer {
  int (*on_local_audio_spectrum)(AGORA_HANDLE agora_local_user, const audio_spectrum_data* data);
  int (*on_remote_audio_spectrum)(AGORA_HANDLE agora_local_user, const user_audio_spectrum_info* spectrums, unsigned int spectrum_number);
} audio_spectrum_observer;

/**
 * The IVideoEncodedFrameObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _video_encoded_frame_observer {
  int (*on_encoded_video_frame)(AGORA_HANDLE agora_video_encoded_frame_observer, uid_t uid, const uint8_t* image_buffer, size_t length,
                                const encoded_video_frame_info* video_encoded_frame_info);
} video_encoded_frame_observer;

/**
 * @ANNOTATION:GROUP:agora_video_encoded_frame_observer
 * @ANNOTATION:CTOR:agora_video_encoded_frame_observer
 */
AGORA_API_C_HDL agora_video_encoded_frame_observer_create(video_encoded_frame_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_video_encoded_frame_observer
 * @ANNOTATION:DTOR:agora_video_encoded_frame_observer
 */
AGORA_API_C_VOID agora_video_encoded_frame_observer_destroy(AGORA_HANDLE agora_video_encoded_frame_observer);

typedef struct _media_stream_info { /* the index of the stream in the media file */
  int stream_index;

  /* stream type */
  int stream_type;

  /* stream encoding name */
  char codec_name[k_max_codec_name_len];

  /* streaming language */
  char language[k_max_codec_name_len];

  /* If it is a video stream, video frames rate */
  int video_frame_rate;

  /* If it is a video stream, video bit rate */
  int video_bit_rate;

  /* If it is a video stream, video width */
  int video_width;

  /* If it is a video stream, video height */
  int video_height;

  /* If it is a video stream, video rotation */
  int video_rotation;

  /* If it is an audio stream, audio bit rate */
  int audio_sample_rate;

  /* If it is an audio stream, the number of audio channels */
  int audio_channels;

  /* If it is an audio stream, bits per sample */
  int audio_bits_per_sample;

  /* stream duration in second */
  int64_t duration;
} media_stream_info;

/** Definition of contentinspect
 */
#define MAX_CONTENT_INSPECT_MODULE_COUNT 32
typedef struct _content_inspect_module {
  /**
   * The content inspect module type.
   */
  int type;
  /**The content inspect frequency, default is 0 second.
   * the frequency <= 0 is invalid.
   */
  unsigned int frequency;
} content_inspect_module;
/** Definition of ContentInspectConfig.
 */
typedef struct _content_inspect_config {
  /** jh on device.*/
  int device_work;

/** jh on cloud.*/
  int cloud_work;

  /**the type of jh on device.*/
  int devicework_type;
  const char* extraInfo;

  /**The content inspect modules, max length of modules is 32.
   * the content(snapshot of send video stream, image) can be used to max of 32 types functions.
   */
  content_inspect_module modules[MAX_CONTENT_INSPECT_MODULE_COUNT];
  /**The content inspect module count.
   */
  int moduleCount;
} content_inspect_config;

#ifdef __cplusplus
}
#endif  // __cplusplus
