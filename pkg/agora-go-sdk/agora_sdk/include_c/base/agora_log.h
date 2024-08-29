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

typedef struct _log_config {
  /**The log file path, default is NULL for default log path
   */
  const char* file_path;
  /** The log file size, KB , set 1024KB to use default log size
   */
  uint32_t file_size_in_KB;
  /** The log level, set LOG_LEVEL_INFO to use default log level
   */
  int level;
} log_config;

/**
 * @ANNOTATION:GROUP:agora_log_writer
 */
AGORA_API_C_INT agora_log_writer_write_log(AGORA_HANDLE agora_log_writer, int level, const char* message, uint16_t length);

#ifdef __cplusplus
}
#endif  // __cplusplus