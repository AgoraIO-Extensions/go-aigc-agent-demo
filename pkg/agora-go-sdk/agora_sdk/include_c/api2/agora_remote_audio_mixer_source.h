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
 * @ANNOTATION:GROUP:agora_remote_audio_mixer_source
 */
AGORA_API_C_INT agora_remote_audio_mixer_source_add_audio_track(AGORA_HANDLE agora_camera_capturer, AGORA_HANDLE track);
/**
 * @ANNOTATION:GROUP:agora_remote_audio_mixer_source
 */
AGORA_API_C_INT agora_remote_audio_mixer_source_remove_audio_track(AGORA_HANDLE agora_camera_capturer, AGORA_HANDLE track);
/**
 * @ANNOTATION:GROUP:agora_remote_audio_mixer_source
 */
AGORA_API_C_INT agora_remote_audio_mixer_source_get_mix_delay(AGORA_HANDLE agora_camera_capturer);

#ifdef __cplusplus
}
#endif  // __cplusplus
