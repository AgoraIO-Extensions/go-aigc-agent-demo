//
//  Agora C SDK
//
//  Created by Hugo Chan in 2020.7
//  Copyright (c) 2020 Agora.io. All rights reserved.
//

#pragma once

#include "agora_base.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * The ICameraCaptureObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _camera_capture_observer {
  void (*on_camera_focus_area_changed)(AGORA_HANDLE agora_camera_capture_observer, int image_width, int image_height, int x, int y);
  void (*on_face_position_changed)(AGORA_HANDLE agora_camera_capture_observer, int image_width, const rectangle* vc_rectangle, const int* vec_distance, int num_faces);
  void (*on_camera_exposure_area_changed)(AGORA_HANDLE agora_camera_capture_observer, int x, int y, int width, int height);
  void (*on_camera_state_changed)(AGORA_HANDLE agora_camera_capture_observer, int state, int source);
} camera_capture_observer;

/**
 * @ANNOTATION:GROUP:agora_device_info
 * @ANNOTATION:DTOR:agora_device_info
 */
AGORA_API_C_VOID agora_camera_capturer_release_device_info(AGORA_HANDLE agora_device_info);

/**
 * @ANNOTATION:GROUP:agora_device_info
 */
AGORA_API_C uint32_t AGORA_CALL_C agora_device_info_number_of_devices(AGORA_HANDLE agora_device_info);

/**
 * @ANNOTATION:GROUP:agora_device_info
 * @ANNOTATION:OUT:device_name_utf8
 * @ANNOTATION:OUT:device_unique_id_utf8
 * @ANNOTATION:OUT:product_unique_id_utf8
 */
AGORA_API_C_INT agora_device_info_get_device_name(AGORA_HANDLE agora_device_info,
                                                uint32_t device_number, char* device_name_utf8,
                                                uint32_t device_name_length, char* device_unique_id_utf8,
                                                uint32_t device_unique_id_length, char* product_unique_id_utf8,
                                                uint32_t product_unique_id_length);

/**
 * @ANNOTATION:GROUP:agora_device_info
 */
AGORA_API_C_INT agora_device_info_number_of_capabilities(AGORA_HANDLE agora_device_info, const char* device_unique_id_utf8);

/**
 * @ANNOTATION:GROUP:agora_device_info
 * @ANNOTATION:OUT:capability
 */
AGORA_API_C_INT agora_device_info_get_capability(AGORA_HANDLE agora_device_info, const char* device_unique_id_utf8,
                                                const uint32_t device_capability_number,
                                                video_format* capability);


/**
 * Camera capturer
 */
#if defined(__ANDROID__) || (defined(__APPLE__) && TARGET_OS_IPHONE)

 /**
  * @ANNOTATION:GROUP:agora_camera_capturer
  */
AGORA_API_C_INT agora_camera_capturer_set_camera_source(AGORA_HANDLE agora_camera_capturer, int source);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_switch_camera(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_is_zoom_supported(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_set_camera_zoom(AGORA_HANDLE agora_camera_capturer, float zoom_value);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_get_camera_max_zoom(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_is_focus_supported(AGORA_HANDLE agora_camera_capturer, float x, float y);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_is_auto_face_focus_supported(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_set_camera_focus(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_set_camera_auto_face_focus(AGORA_HANDLE agora_camera_capturer, int enable);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_enable_face_detection(AGORA_HANDLE agora_camera_capturer, int enable);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_is_camera_face_detect_supported(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_is_camera_torch_supported(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_set_camera_torchOn(AGORA_HANDLE agora_camera_capturer, int is_on);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_get_camera_source(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_is_camera_exposure_position_supported(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_set_camera_exposure_position(AGORA_HANDLE agora_camera_capturer, float position_xin_view, float position_yin_view);

#if (defined(__APPLE__) && TARGET_OS_IOS)

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_is_camera_auto_exposure_face_mode_supported(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_set_camera_auto_exposure_face_mode_enabled(AGORA_HANDLE agora_camera_capturer, int enabled);

#endif

#elif defined(_WIN32) || (defined(__linux__) && !defined(__ANDROID__)) || \
    (defined(__APPLE__) && TARGET_OS_MAC && !TARGET_OS_IPHONE)

 /**
  * @ANNOTATION:GROUP:agora_camera_capturer
  * @ANNOTATION:CTOR:agora_device_info
  */
AGORA_API_C_HDL agora_camera_capturer_create_device_info(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_init_with_device_id(AGORA_HANDLE agora_camera_capturer, const char* device_id);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_INT agora_camera_capturer_init_with_device_name(AGORA_HANDLE agora_camera_capturer, const char* device_name);

#endif  // __ANDROID__ || (__APPLE__ && TARGET_OS_IPHONE)

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_VOID agora_camera_capturer_set_device_orientation(AGORA_HANDLE agora_camera_capturer, int orientation);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_VOID agora_camera_capturer_set_capture_format(AGORA_HANDLE agora_camera_capturer, const video_format* capture_format);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C video_format* AGORA_CALL_C agora_camera_capturer_get_capture_format(AGORA_HANDLE agora_camera_capturer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_VOID agora_camera_capturer_destroy_capture_format(AGORA_HANDLE agora_camera_capturer, video_format* stats);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_VOID agora_camera_capturer_register_camera_observer(AGORA_HANDLE agora_camera_capturer, camera_capture_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_camera_capturer
 */
AGORA_API_C_VOID agora_camera_capturer_unregister_camera_observer(AGORA_HANDLE agora_camera_capturer, camera_capture_observer* observer);

#ifdef __cplusplus
}
#endif //__cpusplus
