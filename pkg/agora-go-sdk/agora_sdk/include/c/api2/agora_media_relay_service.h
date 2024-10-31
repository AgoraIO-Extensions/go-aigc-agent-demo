//
//  Agora C SDK
//
//  Created by Hugo Chan in 2020.7
//  Copyright (c) 2020 Agora.io. All rights reserved.
//

#pragma once
#include "agora_base.h"
#include "agora_service.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

/**
 * The IMediaRelayObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _media_relay_observer {
  void (*on_channel_media_relay_state_changed)(AGORA_HANDLE agora_media_relay_observer, int state, int code);
  void (*on_channel_media_relay_event)(AGORA_HANDLE agora_media_relay_observer, int code);
} media_relay_observer;

/**
 * @ANNOTATION:GROUP:agora_media_relay_service
 */
AGORA_API_C_INT agora_media_relay_service_start_channel_media_relay(AGORA_HANDLE agora_media_relay_sev, const channel_media_relay_config* config);
/**
 * @ANNOTATION:GROUP:agora_media_relay_service
 */
AGORA_API_C_INT agora_media_relay_service_update_channel_media_relay(AGORA_HANDLE agora_media_relay_sev, const channel_media_relay_config* config);
/**
 * @ANNOTATION:GROUP:agora_media_relay_service
 */
AGORA_API_C_INT agora_media_relay_service_stop_channel_media_relay(AGORA_HANDLE agora_media_relay_sev);
/**
 * @ANNOTATION:GROUP:agora_media_relay_service
 */
AGORA_API_C_INT agora_media_relay_service_pause_all_channel_media_relay(AGORA_HANDLE agora_media_relay_sev);
/**
 * @ANNOTATION:GROUP:agora_media_relay_service
 */
AGORA_API_C_INT agora_media_relay_service_resume_all_channel_media_relay(AGORA_HANDLE agora_media_relay_sev);
/**
 * @ANNOTATION:GROUP:agora_media_relay_service
 */
AGORA_API_C_INT agora_media_relay_service_register_event_handler(AGORA_HANDLE agora_media_relay_sev, media_relay_observer* event_observer, void(*safeDeleter)(media_relay_observer*));
/**
 * @ANNOTATION:GROUP:agora_media_relay_service
 */
AGORA_API_C_INT agora_media_relay_service_unregister_event_handler(AGORA_HANDLE agora_media_relay_sev, media_relay_observer* event_observer);

#ifdef __cplusplus
}
#endif  // __cplusplus
