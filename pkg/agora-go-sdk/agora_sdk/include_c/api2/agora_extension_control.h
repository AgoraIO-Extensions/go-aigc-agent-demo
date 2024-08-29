//
//  Agora C SDK
//
//  Created by Hugo Chan in 2020.7
//  Copyright (c) 2020 Agora.io. All rights reserved.
//

#pragma once
#include "agora_base.h"
#include "agora_video_frame.h"
#include "agora_log.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

/**
 * Agora Extension Capabilities.
 */
typedef struct _capabilities {
	/**
	 * Whether to support audio extensions.
	 */
	int audio;
	/**
	 * Whether to support video extensions.
	 */
	int video;
} capabilities;

/**
 * @ANNOTATION:GROUP:agora_extension_control
 */
AGORA_API_C_INT agora_extension_control_get_capabilities(AGORA_HANDLE agora_extension_control, capabilities* capabilities);

/**
 * @ANNOTATION:GROUP:agora_extension_control
 */
AGORA_API_C_INT agora_extension_control_recycle_video_cache(AGORA_HANDLE agora_extension_control);

/**
 * @ANNOTATION:GROUP:agora_extension_control
 */
AGORA_API_C_INT agora_extension_control_dump_video_frame(AGORA_HANDLE agora_extension_control, video_frame* frame, const char* file);

/**
 * @ANNOTATION:GROUP:agora_extension_control
 */
AGORA_API_C_INT agora_extension_control_log(AGORA_HANDLE agora_extension_control, int level, const char* message);

/**
 * @ANNOTATION:GROUP:agora_extension_control
 */
AGORA_API_C_INT agora_extension_control_fire_event(AGORA_HANDLE agora_extension_control, const char* provider, const char* extension, const char* event_key, const char* value);

/**
 * @ANNOTATION:GROUP:agora_extension_control
 */
AGORA_API_C_INT agora_extension_control_register_provider(AGORA_HANDLE agora_extension_control, const char* provider, AGORA_HANDLE instance);

#ifdef __cplusplus
}
#endif  // __cplusplus
