//
//  Agora C SDK
//
//  Created by Tommy Miao in 2020.5
//  Copyright (c) 2020 Agora.io. All rights reserved.
//
#pragma once

#include "agora_audio_track.h"
#include "agora_video_track.h"
#include "agora_rtmp_conn.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

/**
 * The IRtmpLocalUserObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _rtmp_local_user_observer {
	void (*on_audio_track_publish_success)(AGORA_HANDLE rtmp_local_user, AGORA_HANDLE audio_track);
	void (*on_audio_track_publication_failure)(AGORA_HANDLE rtmp_local_user, AGORA_HANDLE audio_track, int error);
	void (*on_video_track_publish_success)(AGORA_HANDLE rtmp_local_user, AGORA_HANDLE video_track);
	void (*on_video_track_publication_failure)(AGORA_HANDLE rtmp_local_user, AGORA_HANDLE video_track, int error);
} rtmp_local_user_observer;

/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_set_audio_stream_conf(AGORA_HANDLE agora_rtmp_local_user, const rtmp_streaming_audio_conf* conf);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_set_video_stream_conf(AGORA_HANDLE agora_rtmp_local_user, const rtmp_streaming_video_conf* conf);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_adjust_recording_signal_volume(AGORA_HANDLE agora_rtmp_local_user, int volume);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_get_recording_signal_volume(AGORA_HANDLE agora_rtmp_local_user, int32_t* volume);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_set_audio_enabled(AGORA_HANDLE agora_rtmp_local_user, int enabled);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_adjust_video_bitrate(AGORA_HANDLE agora_rtmp_local_user, int type);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_set_video_enabled(AGORA_HANDLE agora_rtmp_local_user, int enabled);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_publish_audio(AGORA_HANDLE agora_rtmp_local_user, AGORA_HANDLE audio_track);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_unpublish_audio(AGORA_HANDLE agora_rtmp_local_user, AGORA_HANDLE audio_track);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_publish_media_player_audio(AGORA_HANDLE agora_rtmp_local_user, AGORA_HANDLE audio_track, int32_t player_id);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_unpublish_media_player_audio(AGORA_HANDLE agora_rtmp_local_user, AGORA_HANDLE audio_track, int32_t player_id);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_publish_video(AGORA_HANDLE agora_rtmp_local_user, AGORA_HANDLE video_track);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_unpublish_video(AGORA_HANDLE agora_rtmp_local_user, AGORA_HANDLE video_track);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_register_rtmp_user_observer(AGORA_HANDLE agora_rtmp_local_user, rtmp_local_user_observer* observe, void(*safe_deleter)(rtmp_local_user_observer*));
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_unregister_rtmp_user_observer(AGORA_HANDLE agora_rtmp_local_user, rtmp_local_user_observer* observe);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_register_audio_frame_observer(AGORA_HANDLE agora_rtmp_local_user, audio_frame_observer* observe);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_unregister_audio_frame_observer(AGORA_HANDLE agora_rtmp_local_user, audio_frame_observer* observe);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_register_video_frame_observer(AGORA_HANDLE agora_rtmp_local_user, video_frame_observer* observe);
/**
 * @ANNOTATION:GROUP:agora_rtmp_local_user
 */
AGORA_API_C_INT agora_rtmp_local_user_unregister_video_frame_observer(AGORA_HANDLE agora_rtmp_local_user, video_frame_observer* observe);

#ifdef __cplusplus
}
#endif  // __cplusplus
