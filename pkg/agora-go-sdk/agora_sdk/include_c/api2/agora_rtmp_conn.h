//
//  Agora C SDK
//
//  Created by Tommy Miao in 2020.5
//  Copyright (c) 2020 Agora.io. All rights reserved.
//
#pragma once

#include "agora_base.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

/**
 * Configurations for the RTMP audio stream.
 */
typedef struct _rtmp_streaming_audio_conf {
  /**
   * Sampling rate (Hz). The default value is 44100 (i.e. 44.1kHz).
   */
  int sample_rate_hz;

  /**
   * Number of bytes per sample. The default value is 2 (i.e. 16-bit PCM).
   */
  int bytes_per_sample;

  /**
   * Number of channels. The default value is 1 (i.e. mono).
   */
  int number_of_channels;

  /**
   * The target bitrate (Kbps) of the output audio stream to be published.
   * The default value is 48 Kbps.
  */
  int bitrate;
} rtmp_streaming_audio_conf;

/**
 * Configurations for the RTMP video stream.
 */
typedef struct _rtmp_streaming_video_conf {
	/**
   * The width (in pixels) of the video. The default value is 640.
   *
   * @note
   * - The value of the dimension (with the |height| below) does not indicate the orientation mode
   * of the output ratio. For how to set the video orientation,
   * see {@link OrientationMode OrientationMode}.
   */
  int width;

  /**
   * The height (in pixels) of the video. The default value is 360.
   *
   * @note
   * - The value of the dimension (with the |width| above) does not indicate the orientation mode
   * of the output ratio. For how to set the video orientation,
   * see {@link OrientationMode OrientationMode}.
   */
  int height;

  /**
   * Frame rate (fps) of the output video stream to be published. The default
   * value is 15 fps.
   */
  int framerate;

	/**
   * The target bitrate (Kbps) of the output video stream to be published.
   * The default value is 800 Kbps.
   */
  int bitrate;

  /**
   *  (For future use) The maximum bitrate (Kbps) for video.
   *  The default value is 960 Kbps.
   */
  int max_bitrate;

  /**
   *  (For future use) The minimum bitrate (Kbps) for video.
   *  The default value is 600 Kbps.
   */
  int min_bitrate;

  /**
   *  The interval between two keyframes.
   *  The default value is 2000ms.
   */
  unsigned int gop_in_ms;

	/**
   * The orientation mode.
   * See {@link ORIENTATION_MODE ORIENTATION_MODE}.
   */
  int orientation_mode;
} rtmp_streaming_video_conf;

/**
 * Configurations for the RTMP connection.
 */
typedef struct _rtmp_conn_conf {
  rtmp_streaming_audio_conf audio_conf;
  rtmp_streaming_video_conf video_conf;
  int enable_write_flv_file;
} rtmp_conn_conf;

/**
 * The information on the RTMP Connection.
 */
typedef struct _rtmp_conn_info {
  /**
   * The state of the current connection: #RTMP_CONNECTION_STATE.
   */
  int state;
} rtmp_conn_info;

/**
 * The IRtmpLocalUserObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _rtmp_conn_observer {
	void (*on_connected)(AGORA_HANDLE rtmp_local_user, const rtmp_conn_info* conn_info);
	void (*on_disconnected)(AGORA_HANDLE rtmp_local_user, const rtmp_conn_info* conn_info);
	void (*on_reconnecting)(AGORA_HANDLE rtmp_local_user, const rtmp_conn_info* conn_info);
	void (*on_reconnected)(AGORA_HANDLE rtmp_local_user, const rtmp_conn_info* conn_info);
	void (*on_connection_failure)(AGORA_HANDLE rtmp_local_user, const rtmp_conn_info* conn_info, int err_code);
	void (*on_transfer_statistics)(AGORA_HANDLE rtmp_local_user, uint64_t video_bitrate, uint64_t audio_bitrate, uint64_t video_frame_rate, uint64_t push_video_frame_cnt, uint64_t pop_video_frame_cnt);
} rtmp_conn_observer;

/**
 * @ANNOTATION:GROUP:agora_rtmp_conn
 */
AGORA_API_C_INT agora_rtmp_conn_connect(AGORA_HANDLE agora_rtmp_conn, const char* url);
/**
 * @ANNOTATION:GROUP:agora_rtmp_conn
 */
AGORA_API_C_INT agora_rtmp_conn_disconnect(AGORA_HANDLE agora_rtmp_conn);
/**
 * @ANNOTATION:GROUP:agora_rtmp_conn
 */
AGORA_API_C_INT agora_rtmp_conn_get_connection_info(AGORA_HANDLE agora_rtmp_conn);
/**
 * @ANNOTATION:GROUP:agora_rtmp_conn
 */
AGORA_API_C_INT agora_rtmp_conn_get_rtmp_local_user(AGORA_HANDLE agora_rtmp_conn);
/**
 * @ANNOTATION:GROUP:agora_rtmp_conn
 */
AGORA_API_C_INT agora_rtmp_conn_register_observer(AGORA_HANDLE agora_rtmp_conn, rtmp_conn_observer* observer, void(*safe_deleter)(rtmp_conn_observer*));
/**
 * @ANNOTATION:GROUP:agora_rtmp_conn
 */
AGORA_API_C_INT agora_rtmp_conn_unregister_observer(AGORA_HANDLE agora_rtmp_conn, rtmp_conn_observer* observer);


#ifdef __cplusplus
}
#endif  // __cplusplus
