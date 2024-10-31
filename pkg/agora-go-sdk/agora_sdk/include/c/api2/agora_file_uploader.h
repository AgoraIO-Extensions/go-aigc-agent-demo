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
#endif  // __cplusplus

typedef struct _image_payload_data {
  int seqid;
  int size;
  int width;
  int height;
  int64_t timestamp;
  uint8_t* buffer;
  void* privdata;
  int privsize;
} image_payload_data;

/**
 * @ANNOTATION:GROUP:agora_file_uploader_service
 */
AGORA_API_C_INT agora_file_uploader_service_start_image_upload(AGORA_HANDLE agora_file_uploader_service, const image_payload_data* img_data);
/**
 * @ANNOTATION:GROUP:agora_file_uploader_service
 */
AGORA_API_C_INT agora_file_uploader_service_stop_image_upload(AGORA_HANDLE agora_file_uploader_service);

#ifdef __cplusplus
}
#endif  // __cplusplus
