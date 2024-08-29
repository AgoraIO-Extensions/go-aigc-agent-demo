## Type Notes

Please avoid `size_t` and `long` in API. They are not of the same length in different platforms / libclang bindings.

Special convertions:

`const char*` NUL-terminated string

In a callback function, if a field is `unsigned char*` pointer and the next field is `length`, it is considered a data arg and an array of bytes is constructed.

## Annotation Syntax

`@ANNOTATION` in comment adds additional information on specific elements.

Function annotations:

- `@ANNOTATION:GROUP:{group_name}`: group name, corresponding class name in C++
- `@ANNOTATION:CTOR:{group_name}`: this function is a constructor
- `@ANNOTATION:DTOR:{group_name}`: this function is a destructor
- `@ANNOTATION:RETURNS:{group_name}`: this function returns an agora object pointer
- `@ANNOTATION:OUT:{arg_name}`: an output argument (pointer)
- `@ANNOTATION:PAYLOAD:{arg}`: this field is a pointer to buffer, size is determined by context logic
- `@ANNOTATION:RAWHANDLE:{arg}`: this field contains a void pointer argument, which is a platform-native handle like view id, hwnd, etc

Struct annotations:

- `@ANNOTATION:TYPE:OBSERVER`: this struct is an observer. The struct content is copied so after register observer call you can safely free it.

Enum annotations:

- `@ANNOTATION:MESSAGE:foo bar baz`: human readable message for this enum entry

## Special cases

The following structs require manual mapping because buffer length is determined by other fields:

- audio_frame.buffer, size = samples_per_channel * channels * sizeof(int16_t)
- external_video_frame.buffer, user provided, no need to map from C to target language
- video_frame { y_buffer, u_buffer, v_buffer }, max location is x_stride * height
- agora_service_config.context

The following function(s) require manual binding because they are non-uniform:

- agora_device_info_get_device_name
- agora_rtm_message_get_raw_message_data

The following function(s) are not mapped due to failure of determining raw pointer length

- agora_video_frame_data
- agora_video_frame_mutable_data

The following function(s) should not be mapped due to API design problems:

- agora_parameter_get_string
- agora_parameter_get_array
- agora_parameter_convert_path
- agora_video_frame_fill_src
