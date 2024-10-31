//
//  Agora C SDK
//
//  Created by Hugo Chan in 2020.7
//  Copyright (c) 2020 Agora.io. All rights reserved.
//

#pragma once
#include "agora_base.h"
#include "agora_media_base.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

#define kMaxCharBufferLength  50

/**
 * @brief The information of the media stream object.
 *
 */
typedef struct _player_stream_info {
  /** The index of the media stream. */
  int stream_index;

  /** The type of the media stream. See {@link MEDIA_STREAM_TYPE}. */
  int stream_type;

  /** The codec of the media stream. */
  char codec_name[kMaxCharBufferLength];

  /** The language of the media stream. */
  char language[kMaxCharBufferLength];

  /** The frame rate (fps) if the stream is video. */
  int video_frame_rate;

  /** The video bitrate (bps) if the stream is video. */
  int video_bit_rate;

  /** The video width (pixel) if the stream is video. */
  int video_width;

  /** The video height (pixel) if the stream is video. */
  int video_height;

  /** The rotation angle if the steam is video. */
  int video_rotation;

  /** The sample rate if the stream is audio. */
  int audio_sample_rate;

  /** The number of audio channels if the stream is audio. */
  int audio_channels;

  /** The number of bits per sample if the stream is audio. */
  int audio_bits_per_sample;

  /** The total duration (second) of the media stream. */
  int64_t duration;  
} player_stream_info;

/**
 * @brief The information of the media stream object.
 *
 */
typedef struct _src_info {
  /** The bitrate of the media stream. The unit of the number is kbps.
   *
   */
  int bitrate_in_kbps;

  /** The name of the media stream.
   *
  */
  const char* name;

} src_info;

typedef struct _cache_statistics {
  /**  total data size of uri
   */
  int64_t file_size;
  /**  data of uri has cached
   */
  int64_t cache_size;
  /**  data of uri has downloaded
   */
  int64_t download_size;
} cache_statistics;

typedef struct _player_updated_info {
  /** playerId has value when user trigger interface of opening
   */
  const char* player_id;

  /** deviceId has value when user trigger interface of opening
   */
  const char* device_id;

  /** cacheStatistics exist if you enable cache, triggered 1s at a time after openning url
   */
  cache_statistics cache_statistics;
} player_updated_info;



typedef struct _media_source {
  /**
   * The URL of the media file that you want to play.
   */
  const char* url;
  /**
   * The URI of the media file
   *
   * When caching is enabled, if the url cannot distinguish the cache file name,
   * the uri must be able to ensure that the cache file name corresponding to the url is unique.
   */
  const char* uri;
  /**
   * Set the starting position for playback, in ms.
   */
  int64_t start_pos;
  /**
  * Autoplay when media source is opened
  *
  */
  int auto_play;
  /**
   * Enable caching.
   */
  int enable_cache;
  /**
   * if the value is true, it means playing agora URL. 
   * The default value is false
   */
  int is_agora_source;
  /**
   * If it is set to true, it means that the live stream will be optimized for quick start. 
   * The default value is false
   */
  int is_live_source;
  
} media_source;

#ifdef __cplusplus
}
#endif  // __cplusplus