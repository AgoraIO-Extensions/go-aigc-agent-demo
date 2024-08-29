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

typedef struct _texture_id {
  uintptr_t id;
} texture_id;

/**
 * This structure defines the video frame of texture type on Android
 * @note For technical preview, not supported for the moment. Use RawPixelBuffer instead.
 * 
 */
typedef struct _texture_info {
  int texture_type;
  int context_type;
  void* shared_context;
  int texture_id;
  float transform_matrix[16];
} texture_info;

/**
 * This structure defines the raw video frame data in memory
 * 
 */
typedef struct _raw_pixel_buffer {
  int format;
  uint8_t* data;
  int size;
} raw_pixel_buffer;

typedef struct _padded_raw_pixel_buffer {
  int format;
  uint8_t* data;
  int size;
  int stride;
} padded_raw_pixel_buffer;

typedef struct _color_space {
  int primaries;
  int transfer;
  int matrix;
  int range;
} color_space;

typedef struct _video_frame_data {
  int type;
  union {
    texture_info texture; // Android (To be supported)
    raw_pixel_buffer pixels; // All platform
    void* cvpixelbuffer; // iOS (To be supported)
  };
  int width;
  int height;
  int rotation;
  color_space color_space;
  int64_t timestamp_ms; // Capture time in milli-seconds

  padded_raw_pixel_buffer padded_pixels;
} video_frame_data;

typedef struct _alpha_channel {
  uint8_t* data;
  int size;
} alpha_channel;

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_get_video_frame_data(AGORA_HANDLE agora_video_frame, video_frame_data* data);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_fill_video_frame_data(AGORA_HANDLE agora_video_frame, const video_frame_data* data);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_get_video_frame_meta_data(AGORA_HANDLE agora_video_frame, void* data);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_fill_video_frame_meta_data(AGORA_HANDLE agora_video_frame, int type, const void* data);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_type(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_format(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_width(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_height(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_size(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_rotation(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_VOID agora_video_frame_set_rotation(AGORA_HANDLE agora_video_frame, int rotation);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C int64_t AGORA_CALL_C agora_video_frame_timestamp_us(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_VOID agora_video_frame_set_timestamp_us(AGORA_HANDLE agora_video_frame, int64_t timestamp_us);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C const uint8_t* AGORA_CALL_C agora_video_frame_data(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C uint8_t* AGORA_CALL_C agora_video_frame_mutable_data(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_resize(AGORA_HANDLE agora_video_frame, int width, int height);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C texture_id* AGORA_CALL_C agora_video_frame_texture_id(AGORA_HANDLE agora_video_frame);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_VOID agora_video_frame_destroy_texture_id(AGORA_HANDLE agora_video_frame, texture_id* id);

/**
 * TODO: unclear API, is src input or output?
 * @ANNOTATION:GROUP:agora_video_frame
 * @ANNOTATION:PAYLOAD:src
 */
AGORA_API_C_INT agora_video_frame_fill_src(AGORA_HANDLE agora_video_frame, int format, int width, int height, int rotation, const uint8_t* src);

/**
 * @ANNOTATION:GROUP:agora_video_frame
 */
AGORA_API_C_INT agora_video_frame_fill_texture(AGORA_HANDLE agora_video_frame, int format, int width, int height, int rotation, texture_id* id);

/**
 * @ANNOTATION:GROUP:agora_video_frame_memory_pool
 */
AGORA_API_C void* AGORA_CALL_C agora_video_frame_memory_pool_create_video_frame(AGORA_HANDLE agora_video_frame_memory_pool, const video_frame_data* data, const int* metatypes, int count);

#ifdef __cplusplus
}
#endif  // __cplusplus
