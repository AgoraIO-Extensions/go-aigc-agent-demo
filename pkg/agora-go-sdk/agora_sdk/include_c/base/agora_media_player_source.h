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

/*
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _agora_media_player_source_observer {

    void (*on_player_source_state_changed)(AGORA_HANDLE agora_media_player_source, int state,
                                           int ec);
    void (*on_position_changed)(AGORA_HANDLE agora_media_player_source, int64_t position);

    void (*on_player_event)(AGORA_HANDLE agora_media_player_source, int event, int64_t elapsedTime,
                            const char* message);

    void (*on_meta_data)(AGORA_HANDLE agora_media_player_source, const uint8_t* data, int length);

    void (*on_play_buffer_updated)(AGORA_HANDLE agora_media_player_source, int64_t play_cached_buffer);

    void (*on_preload_event)(AGORA_HANDLE agora_media_player_source, const char* src, int event);

    void (*on_completed)(AGORA_HANDLE agora_media_player_source);

    void (*on_agora_CDN_token_will_expire)(AGORA_HANDLE agora_media_player_source);

    void (*on_player_src_info_changed)(AGORA_HANDLE agora_media_player_source, const src_info* from, const src_info* to);

    void (*on_player_info_updated)(AGORA_HANDLE agora_media_player_source, const player_updated_info* info);

    void (*on_audio_volume_indication)(AGORA_HANDLE agora_media_player_source, int volume);
} agora_media_player_source_observer;

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_get_source_id(AGORA_HANDLE agora_media_player_source);  

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_open(AGORA_HANDLE agora_media_player_source, const char* url, int64_t start_pos);


/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_open_with_media_source(AGORA_HANDLE agora_media_player_source, int64_t start_pos, const media_source* source);


/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_play(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_pause(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_stop(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_resume(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_seek(AGORA_HANDLE agora_media_player_source, int64_t new_pos);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C int64_t AGORA_CALL_C agora_media_player_source_get_duration(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C int64_t AGORA_CALL_C agora_media_player_source_get_play_position(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C int64_t AGORA_CALL_C agora_media_player_source_get_stream_count(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C media_stream_info* AGORA_CALL_C agora_media_player_source_get_stream_info(AGORA_HANDLE agora_media_player_source, int64_t index);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_VOID agora_media_player_source_destroy_stream_info(AGORA_HANDLE agora_media_player_source, media_stream_info* info);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_set_loop_count(AGORA_HANDLE agora_media_player_source, int loop_count);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_mute_audio(AGORA_HANDLE agora_media_player_source, int audio_mute);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_is_audio_muted(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_mute_video(AGORA_HANDLE agora_media_player_source, int video_mute);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_is_video_muted(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_set_playback_speed(AGORA_HANDLE agora_media_player_source, int speed);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_select_audio_track(AGORA_HANDLE agora_media_player_source, int index);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_set_player_option(AGORA_HANDLE agora_media_player_source, const char* key, int value);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_take_screenshot(AGORA_HANDLE agora_media_player_source, const char* filename);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_select_internal_subtitle(AGORA_HANDLE agora_media_player_source, int index);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_set_external_subtitle(AGORA_HANDLE agora_media_player_source, const char* url);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_get_state(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_register_player_source_observer(AGORA_HANDLE agora_media_player_source, agora_media_player_source_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_unregister_player_source_observer(AGORA_HANDLE agora_media_player_source, agora_media_player_source_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_register_audio_frame_observer(AGORA_HANDLE agora_media_player_source, audio_frame_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_unregister_audio_frame_observer(AGORA_HANDLE agora_media_player_source, audio_frame_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_open_with_agora_CDN_src(AGORA_HANDLE agora_media_player_source, const char* src, int64_t start_pos);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_get_agora_CDN_line_count(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_switch_agora_CDN_line_by_index(AGORA_HANDLE agora_media_player_source, int index);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_get_current_agora_CDN_index(AGORA_HANDLE agora_media_player_source);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_enable_auto_switch_agora_CDN(AGORA_HANDLE agora_media_player_source, int enable);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_renew_agora_CDN_src_token(AGORA_HANDLE agora_media_player_source, const char* token, int64_t ts);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_switch_agora_CDN_src(AGORA_HANDLE agora_media_player_source, const char* src, int sync_pts);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_switch_src(AGORA_HANDLE agora_media_player_source, const char* src, int sync_pts);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_preload_src(AGORA_HANDLE agora_media_player_source, const char* src, int64_t start_pos);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_unload_src(AGORA_HANDLE agora_media_player_source, const char* src);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_play_preloaded_src(AGORA_HANDLE agora_media_player_source, const char* src);

/**
 * @ANNOTATION:GROUP:agora_media_player_source
 */
AGORA_API_C_INT agora_media_player_source_change_playback_speed(AGORA_HANDLE agora_media_player_source,int speed);

#ifdef __cplusplus
}
#endif // __cplusplus
