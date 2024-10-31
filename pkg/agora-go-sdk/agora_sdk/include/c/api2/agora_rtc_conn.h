//
//  Agora C SDK
//
//  Created by Ender Zheng in 2020.5
//  Copyright (c) 2020 Agora.io. All rights reserved.
//
#pragma once

#include "agora_base.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

 /**
 * The information on the RTC Connection.
 */
typedef struct _rtc_conn_info {
  /**
   * ID of the RTC Connection.
   */
  conn_id_t id;
  /**
   * ID of the target channel. NULL if you did not call the connect
   * method.
   */
  const char* channel_id;
  /**
   * The state of the current connection: #CONNECTION_STATE_TYPE.
   */
  int state;
  /**
   * ID of the local user.
   */
  const char* local_user_id;
  /**
   * Internal use only.
   */
  uid_t internal_uid;
} rtc_conn_info;

typedef struct _audio_subscription_options {
  /**
   * Determines whether to subscribe to audio packet only, i.e., RTP packet.
   * - 1: Subscribe to audio packet only, which means that the remote audio stream
   * is not be decoded at all. You can use this mode to receive audio packet and handle it
   * in applicatin.
   * Note: if set to true, other fileds in AudioSubscriptionOptions will be ignored.
   * - 0: Do not subscribe to audio packet only, which means that the remote audio stream
   * is decoded automatically.
  */
  int packet_only;
  /**
   * Determines whether to subscribe to PCM audio data only.
   * - 1: Subscribe to PCM audio data only, which means that the remote audio stream
   * is not be played by the built-in playback device automatically. You can use this
   * mode to pull PCM data and handle playback.
   * - 0: Do not subscribe to PCM audio only, which means that the remote audio stream
   * is played automatically.
   */
  int pcm_data_only;
  /**
   * The number of bytes that you expect for each audio sample.
   */
  uint32_t bytes_per_sample;
  /**
   * The number of audio channels that you expect.
   */
  uint32_t number_of_channels;
  /**
   * The audio sample rate (Hz) that you expect.
   */
  uint32_t sample_rate_hz;
} audio_subscription_options;

/**
 * Configurations for the RTC connection.
 */
typedef struct _rtc_conn_config {
  /**
   * Determines whether to subscribe to all audio streams automatically.
   * - 1: (Default) Subscribe to all audio streams automatically.
   * - 0: Do not subscribe to any audio stream automatically.
   */
  int auto_subscribe_audio;
  /**
   * Determines whether to subscribe to all video streams automatically.
   * - 1: (Default) Subscribe to all video streams automatically.
   * - 0: Do not subscribe to any video stream automatically.
   */
  int auto_subscribe_video;
  /**
   * Determines whether to enable audio recording or playout.
   * - true: It's used to publish audio and mix microphone, or subscribe audio and playout
   * - false: It's used to publish extenal audio frame only without mixing microphone, or no need audio device to playout audio either
   */
  int enable_audio_recording_or_playout;
  /**
   * The maximum sending bitrate.
   */
  int max_send_bitrate;
  /**
   * The minimum port.
   */
  int min_port;
  /**
   * The maximum port.
   */
  int max_port;
  /**
   * The options for audio subscription: AudioSubscriptionOptions.
   */
  audio_subscription_options audio_subs_options;
  /**
   * The role of the user: #CLIENT_ROLE_TYPE. The default user role is CLIENT_ROLE_AUDIENCE.
   */
  int client_role_type;

  int channel_profile;

  /**
   * Determines whether to receive audio media packet or not.
   */
  int audio_recv_media_packet;

  /**
   * Determines whether to receive audio encoded frame or not.
   */
  int audio_recv_encoded_frame;

    /**
   * Determines whether to receive video media packet or not.
   */
  int video_recv_media_packet;

} rtc_conn_config;

/**
 * The IRtcConnectionObserver class, which observes the connection state of the SDK.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _rtc_conn_observer {
  void (*on_connected)(AGORA_HANDLE agora_rtc_conn /* pointer to RefPtrHolder */, const rtc_conn_info* conn_info, int reason);
  void (*on_disconnected)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
  void (*on_connecting)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
  void (*on_reconnecting)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
  void (*on_reconnected)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason);
  void (*on_connection_lost)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info);

  void (*on_lastmile_quality)(AGORA_HANDLE agora_rtc_conn, int quality);
  void (*on_lastmile_probe_result)(AGORA_HANDLE agora_rtc_conn, const lastmile_probe_result* result);
  void (*on_token_privilege_will_expire)(AGORA_HANDLE agora_rtc_conn, const char* token);
  void (*on_token_privilege_did_expire)(AGORA_HANDLE agora_rtc_conn);
  void (*on_connection_license_validation_failure)(AGORA_HANDLE agora_rtc_conn, int reason);
  void (*on_connection_failure)(AGORA_HANDLE agora_rtc_conn, const rtc_conn_info* conn_info, int reason); 
  void (*on_user_joined)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id);
  void (*on_user_left)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int reason);
  void (*on_transport_stats)(AGORA_HANDLE agora_rtc_conn, const rtc_stats* stats);
  void (*on_change_role_success)(AGORA_HANDLE agora_rtc_conn, int old_role, int new_role, const client_role_options* new_role_options);
  void (*on_change_role_failure)(AGORA_HANDLE agora_rtc_conn, int reason, int current_role);
  void (*on_user_network_quality)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int tx_quality, int rx_quality);
  void (*on_network_type_changed)(AGORA_HANDLE agora_rtc_conn, int type);
  void (*on_api_call_executed)(AGORA_HANDLE agora_rtc_conn, int err, const char* api, const char* result);
  void (*on_content_inspect_result)(AGORA_HANDLE agora_rtc_conn, int result);
  void (*on_snapshot_taken)(AGORA_HANDLE agora_rtc_conn, const char* channel, uid_t uid, const char* file_path, int width, int height, int err_code);
  void (*on_error)(AGORA_HANDLE agora_rtc_conn, int error, const char* msg);
  void (*on_warning)(AGORA_HANDLE agora_rtc_conn, int warning, const char* msg);
  void (*on_channel_media_relay_state_changed)(AGORA_HANDLE agora_rtc_conn, int state, int code);
  void (*on_local_user_registered)(AGORA_HANDLE agora_rtc_conn, uid_t uid, const char* user_account);
  void (*on_user_account_updated)(AGORA_HANDLE agora_rtc_conn, uid_t uid, const char* user_account);
  void (*on_stream_message_error)(AGORA_HANDLE agora_rtc_conn, user_id_t user_id, int stream_id, int code, int missed, int cached);
  void (*on_encryption_error)(AGORA_HANDLE agora_rtc_conn, int error_type);
  void (*on_upload_log_result)(AGORA_HANDLE agora_rtc_conn, const char* request_id, int success, int reason);
} rtc_conn_observer;

/*
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _network_observer {
  void (*on_uplink_network_info_updated)(AGORA_HANDLE agora_rtc_conn /* pointer to RefPtrHolder */, const uplink_network_info* info);
  void (*on_downlink_network_info_updated)(AGORA_HANDLE agora_rtc_conn /* pointer to RefPtrHolder */, const downlink_network_info* info);
} network_observer;

/**
 * @ANNOTATION:GROUP:agora_service
 * @ANNOTATION:CTOR:agora_rtc_conn
 */
AGORA_API_C_HDL agora_rtc_conn_create(AGORA_HANDLE agora_svc, const rtc_conn_config* rtc_conn_config);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 * @ANNOTATION:DTOR:agora_rtc_conn
 */
AGORA_API_C_VOID agora_rtc_conn_destroy(AGORA_HANDLE agora_rtc_conn);


/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_connect(AGORA_HANDLE agora_rtc_conn, const char* token, const char* chan_id, user_id_t user_id);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_disconnect(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_start_lastmile_probe_test(AGORA_HANDLE agora_rtc_conn, const lastmile_probe_config* config);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_stop_lastmile_probe_test(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_renew_token(AGORA_HANDLE agora_rtc_conn, const char* token);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C rtc_conn_info* AGORA_CALL_C agora_rtc_conn_get_conn_info(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_VOID agora_rtc_conn_destroy_conn_info(AGORA_HANDLE agora_rtc_conn, rtc_conn_info* info);


/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 * @ANNOTATION:RETURNS:agora_local_user
 */
AGORA_API_C_HDL agora_rtc_conn_get_local_user(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_get_remote_users(AGORA_HANDLE agora_rtc_conn, AGORA_HANDLE users);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C user_info* agora_rtc_conn_get_user_info(AGORA_HANDLE agora_rtc_conn, user_id_t user_id);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_VOID agora_rtc_conn_destroy_user_info(AGORA_HANDLE agora_rtc_conn, user_info* info);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_register_observer(AGORA_HANDLE agora_rtc_conn, rtc_conn_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_unregister_observer(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_register_network_observer(AGORA_HANDLE agora_rtc_conn, network_observer* observer);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_unregister_network_observer(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C conn_id_t AGORA_CALL_C agora_rtc_conn_get_conn_id(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C rtc_stats* AGORA_CALL_C agora_rtc_conn_get_transport_stats(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_VOID agora_rtc_conn_destroy_transport_stats(AGORA_HANDLE agora_rtc_conn, rtc_stats* stats);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_HDL agora_rtc_conn_get_agora_parameter(AGORA_HANDLE agora_rtc_conn);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 * @ANNOTATION:OUT:stream_id
 */
AGORA_API_C_INT agora_rtc_conn_create_data_stream(AGORA_HANDLE agora_rtc_conn, int* stream_id, int reliable, int ordered);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_send_stream_message(AGORA_HANDLE agora_rtc_conn, int stream_id, const char* data, uint32_t length);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_enable_encryption(AGORA_HANDLE agora_rtc_conn, int enabled, const encryption_config* config);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_INT agora_rtc_conn_send_custom_report_message(AGORA_HANDLE agora_rtc_conn, const char* id, const char* category, const char* event, const char* label, int value);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_HDL agora_rtc_conn_get_user_info_by_user_account(AGORA_HANDLE agora_rtc_conn, const char* user_account, user_info* user_info);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_HDL agora_rtc_conn_get_user_info_by_uid(AGORA_HANDLE agora_rtc_conn, uid_t uid, user_info* user_info);

/**
 * @ANNOTATION:GROUP:agora_rtc_conn
 */
AGORA_API_C_HDL agora_rtc_conn_get_ntp_time(AGORA_HANDLE agora_rtc_conn);


#ifdef __cplusplus
}
#endif  // __cplusplus
