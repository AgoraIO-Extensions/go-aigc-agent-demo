//
//  Agora C SDK
//
//  Created by Hugo Chan in 2020.7
//  Copyright (c) 2020 Agora.io. All rights reserved.
//

#pragma once
#include "agora_base.h"
#include "agora_media_base.h"
#include "agora_media_player_types.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

/**
 * @brief The input SEI data
 *
 */
typedef struct _input_sei_data {
  int32_t         type;           ///< SEI type
  int64_t         timestamp;      ///< the frame timestamp which be attached
  int64_t         frame_index;    ///< the frame index which be attached
  uint8_t*        private_data;   ///< SEI really data
  int32_t         data_size;      ///< size of really data
} input_sei_data;

/**
 * The IMediaPlayerSourceObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _media_streaming_source_observer {
  void (*on_state_changed)(AGORA_HANDLE agora_media_streaming_source, int state, int err_code);
  void (*on_open_done)(AGORA_HANDLE agora_media_streaming_source, int err_code);
  void (*on_seek_done)(AGORA_HANDLE agora_media_streaming_source, int err_code);
  void (*on_eof_once)(AGORA_HANDLE agora_media_streaming_source, int64_t progress_ms, int64_t repeat_count);
  void (*on_progress)(AGORA_HANDLE agora_media_streaming_source, int64_t position_ms);
  void (*on_meta_data)(AGORA_HANDLE agora_media_streaming_source, const void* data, int length);
} media_streaming_source_observer;

/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_open(AGORA_HANDLE agora_media_streaming_source, const char* url, int64_t start_pos, int auto_play);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_close(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_get_source_id(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_is_video_valid(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_is_audio_valid(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_get_duration(AGORA_HANDLE agora_media_streaming_source, int64_t* duration);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_get_stream_count(AGORA_HANDLE agora_media_streaming_source, int64_t* count);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_get_stream_info(AGORA_HANDLE agora_media_streaming_source, int64_t index, player_stream_info* out_info);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_set_loop_count(AGORA_HANDLE agora_media_streaming_source, int64_t loop_count);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_play(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_pause(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_stop(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_seek(AGORA_HANDLE agora_media_streaming_source, int64_t new_pos);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_get_curr_position(AGORA_HANDLE agora_media_streaming_source, int64_t* pos);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_get_curr_state(AGORA_HANDLE agora_media_streaming_source);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_append_sei_data(AGORA_HANDLE agora_media_streaming_source, const input_sei_data* in_sei_date);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_register_observer(AGORA_HANDLE agora_media_streaming_source, media_streaming_source_observer* observer);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_unregister_observer(AGORA_HANDLE agora_media_streaming_source, media_streaming_source_observer* observer);
/**
 * @ANNOTATION:GROUP:agora_media_streaming_source
 */
AGORA_API_C_INT agora_media_streaming_source_parse_media_info(AGORA_HANDLE agora_media_streaming_source, const char* url, 
                                                              player_stream_info* video_info, 
                                                              player_stream_info* audio_info);
#ifdef __cplusplus
}
#endif  // __cplusplus