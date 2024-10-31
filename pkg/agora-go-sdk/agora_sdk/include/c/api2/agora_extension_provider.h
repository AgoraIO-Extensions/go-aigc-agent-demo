//
//  Agora C SDK
//
//  Created by Hugo Chan in 2020.7
//  Copyright (c) 2020 Agora.io. All rights reserved.
//

#pragma once
#include "agora_base.h"
#include "agora_media_node_factory.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

typedef struct _extension_meta_info {
	int type;
	const char* extension_name;
} extension_meta_info;

/**
 * @ANNOTATION:GROUP:agora_extension_provider
 */
AGORA_API_C_INT agora_extension_provider(AGORA_HANDLE agora_extension_control, int level, const char* message);


#ifdef __cplusplus
}
#endif  // __cplusplus
