//
//  Agora C SDK
//
//  Created by Tommy Miao in 2020.5
//  Copyright (c) 2020 Agora.io. All rights reserved.
//
#pragma once

#include <stdint.h>
#include <stdlib.h>

#include "agora_api.h"
#include "agora_media_base.h"

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

/**
 * for Agora handles that are allocated by SDK (got by API call) and need to be deallocated by SDK
 */
AGORA_API_C_HDL agora_alloc(uint32_t size);
AGORA_API_C_VOID agora_free(AGORA_HANDLE* buf);

typedef struct _user_info {
  /**
   * ID of the user.
   */
  const char* user_id;
  /**
   * Whether the user has enabled audio:
   * - 1: The user has enabled audio.
   * - 0: The user has disabled audio.
   */
  int has_audio;
  /**
   * Whether the user has enabled video:
   * - 1: The user has enabled video.
   * - 0: The user has disabled video.
   */
  int has_video;
} user_info;


typedef struct _account_info {
  /**
   * string account of the user.
   */
  const char* user_account;
  /**
   * ID of the user
   */
  int uid;
} account_info;


// TODO(tomiao): typedef util::AList<UserInfo> UserList;
//typedef user_info user_info_list;

typedef struct _video_dimensions {
  /**
   * The width of the video in number of pixels.
   */
  int width;
  /**
   * The height of the video in number of pixels.
   */
  int height;
} video_dimensions;

typedef struct _encoded_audio_frame_info {
  /**
   * Determines whether the audio frame source is a speech.
   * - 1: (Default) The audio frame source is a speech.
   * - 0: The audio frame source is not a speech.
   */
  int speech;
  /**
   * The audio codec: AUDIO_CODEC_TYPE.
   */
  int codec;
  /**
   * The sample rate (Hz) of the audio frame.
   */
  int sample_rate_hz;
  /**
   * The number of samples per audio channel.
   *
   * If this value is not set, it is 1024 for AAC, 960 for OPUS default.
   */
  int samples_per_channel;
  /**
   * Determines whether to sent the audio frame even when it is empty.
   * - 1: (Default) Send the audio frame even when it is empty.
   * - 0: Do not send the audio frame when it is empty.
   */
  int send_even_if_empty;
  /**
   * The number of channels of the audio frame.
   */
  int number_of_channels;

} encoded_audio_frame_info;

typedef struct _audio_pcm_data_info {
  // The samples count you expect.
  /**
   * The sample count of the Pcm data that you expect.
   */
  uint32_t samples_per_channel;
  int16_t channel_num;

  // Output
  /**
   * The number of output samples.
   */
  uint32_t samples_out;
  int64_t elapsed_time_ms;
  int64_t ntp_time_ms;
} audio_pcm_data_info;

typedef struct _encoded_video_frame_info {
  /**
   * The video codec: #VIDEO_CODEC_TYPE.
   */
  int codec_type;
  /**
   * The width (px) of the video.
   */
  int width;
  /**
   * The height (px) of the video.
   */
  int height;
  /**
   * The number of video frames per second.
   * This value will be used for calculating timestamps of the encoded image.
   * If framesPerSecond equals zero, then real timestamp will be used.
   * Otherwise, timestamp will be adjusted to the value of framesPerSecond set.
   */
  int frames_per_second;
  /**
   * The frame type of the encoded video frame: #VIDEO_FRAME_TYPE.
   */
  int frame_type;
  /**
   * The rotation information of the encoded video frame: #VIDEO_ORIENTATION.
   */
  int rotation;
  int track_id;  // This can be reserved for multiple video tracks, we need to create different ssrc
                 // and additional payload for later implementation.
  /**
   * This is a input parameter which means the timestamp for capturing the video.
   */
  int64_t capture_time_ms;
  /**
   * The timestamp for rendering the video.
   */
  int64_t decode_time_ms;
  /**
   * ID of the user.
   */
  uid_t uid;
  /**
   * The stream type of video frame.
   */
  int stream_type;
} encoded_video_frame_info;

typedef struct _video_encoder_config {
  /**
   * The video encoder code type: #VIDEO_CODEC_TYPE.
   */
  int codec_type;
  /**
   * The video frame dimension used to specify the video quality and measured
   * by the total number of pixels along a frame's width and height: VideoDimensions.
   */
  video_dimensions dimensions;
  /**
   * Frame rate of the video: int type, but can accept #FRAME_RATE for backward
   * compatibility.
   */
  int frame_rate;
  /**
   * Video encoding target bitrate (Kbps).
   *
   * Choose one of the following options:
   *
   * - #STANDARD_BITRATE: (Recommended) Standard bitrate.
   *   - Communication profile: The encoding bitrate equals the base bitrate.
   *   - Live broadcast profile: The encoding bitrate is twice the base bitrate.
   * - #COMPATIBLE_BITRATE: Compatible bitrate: The bitrate stays the same
   * regardless of the profile.
   *
   * The Communication profile prioritizes smoothness, while the Live Broadcast
   * profile prioritizes video quality (requiring a higher bitrate). Agora
   * recommends setting the bitrate mode a #STANDARD_BITRATE or simply to
   * address this difference.
   *
   * The following table lists the recommended video encoder configurations,
   * where the base bitrate applies to the communication profile. Set your
   * bitrate based on this table. If the bitrate you set is beyond the proper
   * range, the SDK automatically sets it to within the range.

   | Resolution             | Frame Rate (fps) | Base Bitrate (Kbps, for Communication) | Live
   Bitrate (Kbps, for Live Broadcast)|
   |------------------------|------------------|----------------------------------------|----------------------------------------|
   | 160 &times; 120        | 15               | 65                                     | 130 |
   | 120 &times; 120        | 15               | 50                                     | 100 |
   | 320 &times; 180        | 15               | 140                                    | 280 |
   | 180 &times; 180        | 15               | 100                                    | 200 |
   | 240 &times; 180        | 15               | 120                                    | 240 |
   | 320 &times; 240        | 15               | 200                                    | 400 |
   | 240 &times; 240        | 15               | 140                                    | 280 |
   | 424 &times; 240        | 15               | 220                                    | 440 |
   | 640 &times; 360        | 15               | 400                                    | 800 |
   | 360 &times; 360        | 15               | 260                                    | 520 |
   | 640 &times; 360        | 30               | 600                                    | 1200 |
   | 360 &times; 360        | 30               | 400                                    | 800 |
   | 480 &times; 360        | 15               | 320                                    | 640 |
   | 480 &times; 360        | 30               | 490                                    | 980 |
   | 640 &times; 480        | 15               | 500                                    | 1000 |
   | 480 &times; 480        | 15               | 400                                    | 800 |
   | 640 &times; 480        | 30               | 750                                    | 1500 |
   | 480 &times; 480        | 30               | 600                                    | 1200 |
   | 848 &times; 480        | 15               | 610                                    | 1220 |
   | 848 &times; 480        | 30               | 930                                    | 1860 |
   | 640 &times; 480        | 10               | 400                                    | 800 |
   | 1280 &times; 720       | 15               | 1130                                   | 2260 |
   | 1280 &times; 720       | 30               | 1710                                   | 3420 |
   | 960 &times; 720        | 15               | 910                                    | 1820 |
   | 960 &times; 720        | 30               | 1380                                   | 2760 |
   | 1920 &times; 1080      | 15               | 2080                                   | 4160 |
   | 1920 &times; 1080      | 30               | 3150                                   | 6300 |
   | 1920 &times; 1080      | 60               | 4780                                   | 6500 |
   | 2560 &times; 1440      | 30               | 4850                                   | 6500 |
   | 2560 &times; 1440      | 60               | 6500                                   | 6500 |
   | 3840 &times; 2160      | 30               | 6500                                   | 6500 |
   | 3840 &times; 2160      | 60               | 6500                                   | 6500 |
   */
  int bitrate;

  /**
   * (For future use) The minimum encoding bitrate (Kbps).
   *
   * The Agora SDK automatically adjusts the encoding bitrate to adapt to the
   * network conditions.
   *
   * Using a value greater than the default value forces the video encoder to
   * output high-quality images but may cause more packet loss and hence
   * sacrifice the smoothness of the video transmission. That said, unless you
   * have special requirements for image quality, Agora does not recommend
   * changing this value.
   *
   * @note
   * This parameter applies to the Live Broadcast profile only.
   */
  int min_bitrate;
  /**
   * (For future use) The video orientation mode of the video: #ORIENTATION_MODE.
   */
  int orientation_mode;
  /**
   *
   * The video degradation preference when the bandwidth is a constraint:
   * #DEGRADATION_PREFERENCE. Currently, this member supports `MAINTAIN_QUALITY`(0)
   * only.
   */
  int degradation_preference;
  /**
   * If mirror_type is set to VIDEO_MIRROR_MODE_ENABLED, then the video frame would be mirrored before encoding.
   */
  int mirror_mode;
} video_encoder_config;

/** Data stream config
*/
typedef struct _data_stream_config {
   /** syncWithAudio Sets whether or not the recipients receive the data stream sync with current audio stream.
    */
   int sync_with_audio;
   /** ordered Sets whether or not the recipients receive the data stream in the sent order:
    */
   int ordered;
} data_stream_config;

typedef struct _simulcast_stream_config {
  /**
   * The video frame dimension, which is used to specify the video quality and measured by the total number
   * of pixels along a frame's width and height: VideoDimensions.
   */
  video_dimensions dimensions;
  /**
   * The video bitrate in Kbps.
   */
  int bitrate;
  /**
   * The video framerate.
   */
  int framerate;
} simulcast_stream_config;

/** The relative location of the region to the screen or window.
 */
typedef struct _rectangle {
  /** The horizontal offset from the top-left corner.
   */
  int x;
  /** The vertical offset from the top-left corner.
   */
  int y;
  /** The width of the region.
   */
  int width;
  /** The height of the region.
   */
  int height;
} rectangle;

/** The options of the watermark image to be added. */
typedef struct _watermark_ratio {
  /**
   * The ratio of the width of the video, see #WATERMARK_FIT_MODE::FIT_MODE_USE_IMAGE_RATIO
   */
  float x_ratio;
  /**
   * The ratio of the height of the video, see #WATERMARK_FIT_MODE::FIT_MODE_USE_IMAGE_RATIO
   */
  float y_ratio;
  /**
   * The ratio of the width of the video, see #WATERMARK_FIT_MODE::FIT_MODE_USE_IMAGE_RATIO
   */
  float width_ratio;
} watermark_ratio;

/** The options of the watermark image to be added. */
typedef struct _watermark_dptions {
  /** Sets whether or not the watermark image is visible in the local video preview:
   * - true: (Not support) The watermark image is visible in preview.
   * - false: (Default) The watermark image is not visible in preview.
   */
  int visible_in_preview;
  /**
   * The watermark position in the landscape mode. See Rectangle.
   * For detailed information on the landscape mode, see the advanced guide *Video Rotation*.
   */
  rectangle position_in_landscape_mode;
  /**
   * The watermark position in the portrait mode. See #WATERMARK_FIT_MODE::FIT_MODE_USE_IMAGE_RATIO.
   */
  rectangle position_in_portrait_mode;
  /**
   * The watermark position in the ratio mode. See Rectangle.
   * For detailed information on the portrait mode, see the advanced guide *Video Rotation*.
   */
  watermark_ratio watermark_ratio;
  /**
   * The fit mode of watermark.
   */
  int mode;
} watermark_dptions;

/** Statistics of a channel. */
typedef struct _rtc_stats {
  /**
   * The connection ID.
   */
  unsigned int connection_id;
  /**
   * The call duration (s), represented by an aggregate value.
   */
  unsigned int duration;
  /**
   * The total number of bytes transmitted, represented by an aggregate value.
   */
  unsigned int tx_bytes;
  /**
   * The total number of bytes received, represented by an aggregate value.
   */
  unsigned int rx_bytes;
  /**
   * Total number of audio bytes sent (bytes), represented by an aggregate value.
   */
  unsigned int tx_audio_bytes;
  /**
   * Total number of video bytes sent (bytes), represented by an aggregate value.
   */
  unsigned int tx_video_bytes;
  /**
   * Total number of audio bytes received (bytes) before network countermeasures, represented by an aggregate value.
   */
  unsigned int rx_audio_bytes;
  /**
   * Total number of video bytes received (bytes), represented by an aggregate value.
   */
  unsigned int rx_video_bytes;
  /**
   * The transmission bitrate (Kbps), represented by an instantaneous value.
   */
  unsigned short tx_k_bit_rate;
  /**
   * The receive bitrate (Kbps), represented by an instantaneous value.
   */
  unsigned short rx_k_bit_rate;
  /**
   * Audio receive bitrate (Kbps), represented by an instantaneous value.
   */
  unsigned short rx_audio_k_bit_rate;
  /**
   * The audio transmission bitrate (Kbps), represented by an instantaneous value.
   */
  unsigned short tx_audio_k_bit_rate;
  /**
   * The video receive bitrate (Kbps), represented by an instantaneous value.
   */
  unsigned short rx_video_k_bit_rate;
  /**
   * The video transmission bitrate (Kbps), represented by an instantaneous value.
   */
  unsigned short tx_video_k_bit_rate;
  /**
   * The VOS client-server latency (ms).
   */
  unsigned short lastmile_delay;
  /**
   * The number of users in the channel.
   */
  unsigned int user_count;
  /**
   * The application CPU usage (%).
   */
  double cpu_app_usage;
  /**
   * The system CPU usage (%).
   */
  double cpu_total_usage;
  /** 
   * gateway Rtt
  */
  int gateway_rtt;
  /**
   * The memory usage ratio of the app (%).
   */
  double memory_app_usage_ratio;
  /**
   * The memory usage ratio of the system (%).
   */
  double memory_total_usage_ratio;
  /**
   * The memory usage of the app (KB).
   */
  int memory_app_usage_in_kbytes;
  /**
   * The duration (ms) between connection establish and connect start , 0 if not valid
   */
  int connect_time_ms;
  /**
   * The duration(ms) between first audio packet received and connection start, 0 if not valid
   */
  int first_audio_packet_duration;
  /**
   * The duration(ms) between first video packet received and connection start, 0 if not valid
   */
  int first_video_packet_duration;
  /**
   * The duration(ms) between first video key frame received and connection start, 0 if not valid
   */
  int first_video_key_frame_packet_duration;
  /**
   * Video packet number before first video key frame received, 0 if not valid
   */
  int packets_before_first_key_frame_packet;
  /**
   * The duration (ms) between the last time unmute audio and the first audio packet is received.
   * 0 indicates that this member does not apply.
   */
  int first_audio_packet_duration_after_unmute;
  /**
   * The duration (ms) between the last time unmute video and the first video packet is received.
   * 0 indicates that this member does not apply.
   */
  int first_video_packet_duration_after_unmute;
  /**
   * The duration (ms) between the last time unmute video and the first video key frame is received.
   * 0 indicates that this member does not apply.
   */
  int first_video_key_frame_packet_duration_after_unmute;
  /**
   * The duration (ms) between the last time unmute video and the first video key frame is decoded.
   * 0 indicates that this member does not apply.
   */
  int first_video_key_frame_decoded_duration_after_unmute;
  /**
   * The duration (ms) between the last time unmute video and the first video key frame is rendered.
   * 0 indicates that this member does not apply.
   */
  int first_video_key_frame_rendered_duration_after_unmute;
  /**
   * The packet loss rate of sender(broadcaster).
   */
  int tx_packet_loss_rate;
  /**
   * The packet loss rate of receiver(audience).
   */
  int rx_packet_loss_rate;
} rtc_stats;

/** Client role options, contains audience latency level.
 */
typedef struct _client_role_options {
  /**
  Audience latency level, support 0.5s and 1.2s.
  */
  int audience_latency_level;
} client_role_options;

/** Audio statistics of a remote user */
typedef struct _remote_audio_stats {
  /**
   * User ID of the remote user sending the audio streams.
   */
  uid_t uid;
  /**
   * Audio quality received by the user: #QUALITY_TYPE.
   */
  int quality;
  /**
   * Network delay (ms) from the sender to the receiver.
   */
  int network_transport_delay;
  /**
   * Network delay (ms) from the receiver to the jitter buffer.
   */
  int jitter_buffer_delay;
  /**
   * The audio frame loss rate in the reported interval.
   */
  int audio_loss_rate;
  /**
   * The number of channels.
   */
  int num_channels;
  /**
   * The sample rate (Hz) of the received audio stream in the reported interval.
   */
  int received_sample_rate;
  /**
   * The average bitrate (Kbps) of the received audio stream in the reported
   * interval.
   */
  int received_bitrate;
  /**
   * The total freeze time (ms) of the remote audio stream after the remote
   * user joins the channel.
   *
   * In a session, audio freeze occurs when the audio frame loss rate reaches 4%.
   */
  int total_frozen_time;
  /**
   * The total audio freeze time as a percentage (%) of the total time when the
   * audio is available.
   */
  int frozen_rate;
  /**
   * The quality of the remote audio stream as determined by the Agora
   * real-time audio MOS (Mean Opinion Score) measurement method in the
   * reported interval. The return value ranges from 0 to 500. Dividing the
   * return value by 100 gets the MOS score, which ranges from 0 to 5. The
   * higher the score, the better the audio quality.
   *
   * | MOS score       | Perception of audio quality                                                                                                                                 |
   * |-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
   * | Greater than 4  | Excellent. The audio sounds clear and smooth.                                                                                                               |
   * | From 3.5 to 4   | Good. The audio has some perceptible impairment, but still sounds clear.                                                                                    |
   * | From 3 to 3.5   | Fair. The audio freezes occasionally and requires attentive listening.                                                                                      |
   * | From 2.5 to 3   | Poor. The audio sounds choppy and requires considerable effort to understand.                                                                               |
   * | From 2 to 2.5   | Bad. The audio has occasional noise. Consecutive audio dropouts occur, resulting in some information loss. The users can communicate only with difficulty.  |
   * | Less than 2     | Very bad. The audio has persistent noise. Consecutive audio dropouts are frequent, resulting in severe information loss. Communication is nearly impossible. |
   */
  int mos_value;
  /**
   * The total time (ms) when the remote user neither stops sending the audio
   * stream nor disables the audio module after joining the channel.
   */
  int total_active_time;
  /**
   * The total publish duration (ms) of the remote audio stream.
   */
  int publish_duration;
} remote_audio_stats;

/** The max value (px) of the width. */
#define k_max_width_in_pixels 3840
/** The max value (px) of the height. */
#define k_max_height_in_pixels 2160
/** The max value (fps) of the frame rate. */
#define k_max_fps 60

typedef struct _video_format {
  /**
   * The width (px) of the video.
   */
  int width;
  /**
   * The height (px) of the video.
   */
  int height;
  /**
   * The video frame rate (fps).
   */
  int fps;
} video_format;

/**
 * The struct of RemoteVideoTrackInfo.
 */
typedef struct _video_track_info {
  /**
   * Whether the video track is local or remote.
   * - true: The video track is local.
   * - false: The video track is remote.
   */
  int is_local;
  /**
   * ID of the user who owns the video track.
   */
  uid_t owner_uid;
  /**
   * ID of the video track.
   */
  track_id_t track_id;
  /**
   * The channel ID of the video track.
   */
  const char* channel_id;
  /**
   * The video stream type: #VIDEO_STREAM_TYPE.
   */
  int stream_type;
  /**
   * The video codec type: #VIDEO_CODEC_TYPE.
   */
  int codec_type;
  /**
   * Whether the video track contains encoded video frame only.
   * - 1: The video track contains encoded video frame only.
   * - 0: The video track does not contain encoded video frame only.
   */
  int encoded_frame_only;
  /**
   * The video source type: #VIDEO_SOURCE_TYPE
   */
  int source_type;
  /**
   * the frame position for the video observer: #VIDEO_MODULE_POSITION
   */
  uint32_t observation_position;
} video_track_info;

/**
 * The struct of RemoteVideoStreamInfo.
 */
typedef struct _remote_video_stream_info {
  /**
   * ID of the user who owns the remote video track.
   */
  uid_t uid;
  /**
   * The remote video stream type: #VIDEO_STREAM_TYPE.
   */
  uint8_t stream_type;
  /**
   * The remote video downscale type: #REMOTE_VIDEO_DOWNSCALE_LEVEL.
   */
  uint8_t current_downscale_level;
  /**
   * Total downscale level counts
   */
  uint8_t total_downscale_level_counts;
} remote_video_stream_info;

/** Properties of the audio volume information.
 */
typedef struct _audio_volume_info {
  /** User ID of the speaker.
   */
  user_id_t user_id;
  /** The volume of the speaker that ranges from 0 (lowest volume) to 255 (highest volume).
   */
  unsigned int volume;  // [0,255]
  /*
   * The activity status of remote users
   */
  unsigned int vad;
  /**
   * Voice pitch frequency in Hz
   */
  double voicePitch;
} audio_volume_info;

/**
 * The definition of the DeviceInfo struct
 */
typedef struct _device_info {
  /*
   * Whether the device support low latency audio. Not support by default
   */
  int is_low_latency_audio_supported;
} device_info;

// TODO(tomiao): rename to packet_observer_packet if meet any naming conflict
typedef struct _packet {
  const unsigned char* buffer;
  unsigned int size;
} packet;

/**
 * Related to IRtcConnectionEx::setPacketObserver() and unused in RTC Context currently
 * ANNOTATION:TYPE:OBSERVER
 */
// typedef struct _packet_observer {
//   int (*on_send_audio_packet)(packet* packet);
//   int (*on_send_video_packet)(packet* packet);
//   int (*on_receive_audio_packet)(packet* packet);
//   int (*on_receive_video_packet)(packet* packet);
// } packet_observer;

typedef struct _local_audio_stats {
  /**
   * The number of channels.
   */
  int num_channels;
  /**
   * The sample rate (Hz).
   */
  int sent_sample_rate;
  /**
   * The average sending bitrate (Kbps).
   */
  int sent_bitrate;
  /**
   * The internal payload type
   */
  int internal_codec;
  
  /**
   * Voice pitch frequency in Hz
   */
  double voice_pitch;
} local_audio_stats;

typedef struct _rtc_image {
  /** URL address of the watermark on the live broadcast video.
   */
  const char* url;
  /** Position of the watermark on the upper left of the live broadcast video on
  the horizontal axis.
  */
  int x;
  /** Position of the watermark on the upper left of the live broadcast video on
  the vertical axis.
  */
  int y;
  /** Width of the watermark on the live broadcast video.
   */
  int width;
  /** Height of the watermark on the live broadcast video.
   */
  int height;
  /**
   * Order attribute for an ordering of overlapping two-dimensional objects.
   */
  int z_order;
  /** The transparency level of the image. The value ranges between 0 and 1.0:
   * - 0: Completely transparent
   * - 1.0: (Default) Opaque
   */
  double alpha;
} rtc_image;

typedef struct _transcoding_user {
  /** User ID of the CDN live.
   */
  uid_t uid;
  user_id_t user_id;
  /** Horizontal position of the top left corner of the video frame.
   */
  int x;
  /** Vertical position of the top left corner of the video frame.
   */
  int y;
  /** Width of the video frame.
   */
  int width;
  /** Height of the video frame.
   */
  int height;

  /** The layer of the video frame. Between 1 and 100:
  - 1: Default, lowest
  - 100: Highest
  */
  int z_order;
  /** The transparency of the video frame.
   */
  double alpha;
  /**
   * The audio channel of the sound. The default value is 0:
   *  - 0: (default) Supports dual channels at most, depending on the upstream of
   *  the broadcaster.
   *  - 1: The audio stream of the broadcaster is in the FL audio channel. If the
   *  upstream of the broadcaster uses dual sound channel, only the left sound
   *  channel will be used for streaming.
   *  - 2: The audio stream of the broadcaster is in the FC audio channel. If the
   *  upstream of the broadcaster uses dual sound channel, only the left sound
   *  channel will be used for streaming.
   *  - 3: The audio stream of the broadcaster is in the FR audio channel. If the
   *  upstream of the broadcaster uses dual sound channel, only the left sound
   *  channel will be used for streaming.
   *  - 4: The audio stream of the broadcaster is in the BL audio channel. If the
   *  upstream of the broadcaster uses dual sound channel, only the left sound
   *  channel will be used for streaming.
   *  - 5: The audio stream of the broadcaster is in the BR audio channel.
   *  If the upstream of the broadcaster uses dual sound channel, only the left
   *  sound channel will be used for streaming.
  */
  int audio_channel;
} transcoding_user;

 typedef struct _live_stream_advanced_feature {
  /** The advanced feature for high-quality video with a lower bitrate. */
  // static const char* LBHQ = "lbhq";
  /** The advanced feature for the optimized video encoder. */
  // static const char* VEO = "veo";

  /** The name of the advanced feature. It contains LBHQ and VEO.
   * "lbhq"
   * "veo"
   */
  const char* feature_name;

  /** Whether to enable the advanced feature:
   * - true: Enable the advanced feature.
   * - false: (Default) Disable the advanced feature.
   */
  unsigned int opened;
} live_stream_advanced_feature;

typedef struct _live_transcoding {
  /**
   * The width of the video.
   */
  int width;
  /**
   * The height of the video.
   */
  int height;
  /**
   * The bitrate (Kbps) of the output data stream set for CDN live. The default value is 400.
   */
  int video_bitrate;
  /**
   * The frame rate (fps) of the output data stream set for CDN live. The default value is 15.
   */
  int video_frame_rate;
  /**
   * Determines whether to enable low latency.
   * - 1: Low latency with unassured quality.
   * - 0: (Default) High latency with assured quality.
   */
  int low_latency;
  /**
   * The time interval (s) between two consecutive I frames. The default value is 2.
   */
  int video_gop;
  /**
   * Self-defined video codec profiles: VIDEO_CODEC_TYPE_FOR_STREAM.
   */
  int video_codec_type_for_stream;
  /**
   * Self-defined video codec profiles: VIDEO_CODEC_PROFILE_TYPE.
   */
  int video_codec_profile;
  /**
   * The background color to set in RGB hex value. Value only, do not include a #.
   * For example, 0xFFB6C1 (light pink). The default value is 0x000000 (black).
   */
  unsigned int background_color;
  /**
   * The number of users in the live broadcast.
   */
  unsigned int user_count;
  /**
   * The TranscodingUser class.
   */
  transcoding_user* transcoding_users;
  /**
   * Extra user-defined information to be sent to the CDN client. The extra
   * infomation will be transmitted by SEI packets.
   */
  const char* transcoding_extra_info;
  /**
   * The pointer to the metadata to be sent to CDN client defined by rtmp or FLV metadata.
   */
  const char* metadata;
  /**
   * The number of watermark in the live broadcast.
   */
  unsigned int watermark_count;
  /**
   * The HTTP/HTTPS URL address of the watermark image added to the CDN
   * publishing stream. The audience of the CDN publishing stream can see the
   * watermark.
   */
  rtc_image* watermark;
  /**
   * The number of background image in the live broadcast.
   */
  unsigned int background_image_count;
  /**
   * The HTTP/HTTPS URL address of the background image added to the CDN
   * publishing stream. The audience of the CDN publishing stream can see the
   * background image.
   */
  rtc_image* background_image;
  /**
   * The self-defined audio-sampling rates: #AUDIO_SAMPLE_RATE_TYPE
   */
  int audio_sample_rate;
  /**
   * The bitrate (Kbps) of the audio-output stream set for CDN live. The highest
   * value is 128.
   */
  int audio_bitrate;
  /**
   * Agora's self-defined audio-channel types. Agora recommends choosing 1 or
   * 2:
   * - 1: Mono (default)
   * - 2: Dual-sound channels
   * - 3: Three-sound channels
   * - 4: Four-sound channels
   * - 5: Five-sound channels
   */
  int audio_channels;
  /**
   * The audio codec profile type: #AUDIO_CODEC_PROFILE_TYPE.
   */
  int audio_codec_profile;
  /**
   * The number of advanced feature count in the live broadcast.
   */
  int advanced_feature_count;
  /** 
   * Advanced features of the RTMP or RTMPS streaming with transcoding. See LiveStreamAdvancedFeature.
   */
  live_stream_advanced_feature* advancedFeatures;
} live_transcoding;

/**
 * The definition of the LocalTranscodingVideoStream struct.
 */
typedef struct _transcoding_video_stream {
  /**
   * Source type of video stream.
   */
  int source_type;
  /**
   * Remote user uid if sourceType is VIDEO_SOURCE_REMOTE.
   */
  uid_t remote_user_uid;
  /**
   * RTC image if sourceType is VIDEO_SOURCE_RTC_IMAGE.
   */
  const char* image_url;
  /**
   * The horizontal position of the top left corner of the video frame.
   */
  int x;
  /**
   * The vertical position of the top left corner of the video frame.
   */
  int y;
  /**
   * The width of the video frame.
   */
  int width;
  /**
   * The height of the video frame.
   */
  int height;
  /**
   * The layer of the video frame that ranges from 1 to 100:
  * - 1: (Default) The lowest layer.
  * - 100: The highest layer.
  */
  int zOrder;
  /**
   * The transparency of the video frame.
   */
  double alpha;
  /**
   * mirror of the source video frame (only valid for camera streams)
   */
  int mirror;
} transcoding_video_stream;

/**
 * The definition of the LocalTranscodingConfiguration struct.
 */
typedef struct _local_transcoder_config {
  /**
   * The number of VideoInputStreams in the transcoder.
   */
  unsigned int stream_count;
  /**
   * The video stream layout configuration in the transcoder.
   */
  transcoding_video_stream* video_input_streams;
  /**
   * The video encoder configuration of transcoded video.
   */
  video_encoder_config video_output_config;
} local_transcoder_config;

/**
 * The last mile network probe configuration.
 */
typedef struct _lastmile_probe_config {
  /**
   * Sets whether or not to test the uplink network. Some users, for example,
   * the audience in a Live-broadcast channel, do not need such a test:
   * - 1: Test.
   * - 0: Do not test.
   */
  int probe_uplink;
  /**
   * Sets whether or not to test the downlink network:
   * - 1: Test.
   * - 0: Do not test.
   */
  int probe_downlink;
  /**
   * The expected maximum sending bitrate (bps) of the local user. The value
   * ranges between 100000 and 5000000. We recommend setting this parameter
   * according to the bitrate value set by
   * \ref IRtcEngine::setVideoEncoderConfiguration "setVideoEncoderConfiguration".
   */
  unsigned int expected_uplink_bitrate;
  /**
   * The expected maximum receiving bitrate (bps) of the local user. The value
   * ranges between 100000 and 5000000.
   */
  unsigned int expected_downlink_bitrate;
} lastmile_probe_config;

/**
 * The uplink or downlink last-mile network probe test result.
 */
typedef struct _lastmile_probe_one_way_result {
  /**
   * The packet loss rate (%).
   */
  unsigned int packet_loss_rate;
  /**
   * The network jitter (ms).
   */
  unsigned int jitter;
  /**
   * The estimated available bandwidth (bps).
   */
  unsigned int available_bandwidth;
} lastmile_probe_one_way_result;

/**
 * The uplink and downlink last-mile network probe test result.
 */
typedef struct _lastmile_probe_result {
  /**
   * The state of last mile network probe test: LASTMILE_PROBE_RESULT_STATE.
   */
  int state;
  /**
   * The uplink last-mile network probe test result: LastmileProbeOneWayResult.
   */
  lastmile_probe_one_way_result uplink_report;
  /**
   * The downlink last-mile network probe test result: LastmileProbeOneWayResult.
   */
  lastmile_probe_one_way_result downlink_report;
  /**
   * The round-trip delay time (ms).
   */
  unsigned int rtt;
} lastmile_probe_result;

/**
 * The VideoCanvas class, which defines the video display window.
 */
typedef struct _video_canvas {
  /**
   * The video display window.
   */
  view_t view;
  /**
   * The video display mode: #RENDER_MODE_TYPE.
   */
  int render_mode;
  /**
   * The video mirror mode:
   */
  int mirror_mode;
  /**
   * The user ID.
   */
  uid_t uid;
  void* priv;  // private data (underlying video engine denotes it)
  int is_screen_view;
  size_t priv_size;
  int source_type;
} video_canvas;

typedef struct _virtual_background_source {
  /** The type of the custom background image. See #BACKGROUND_SOURCE_TYPE.
   */
  int background_source_type;

  /**
   * The color of the custom background image. The format is a hexadecimal integer defined by RGB, without the # sign,
   * such as 0xFFB6C1 for light pink. The default value is 0xFFFFFF, which signifies white. The value range
   * is [0x000000,0xFFFFFF]. If the value is invalid, the SDK replaces the original background image with a white
   * background image.
   *
   * @note This parameter takes effect only when the type of the custom background image is `BACKGROUND_COLOR`.
   */
  unsigned int color;

  /**
   * The local absolute path of the custom background image. PNG and JPG formats are supported. If the path is invalid,
   * the SDK replaces the original background image with a white background image.
   *
   * @note This parameter takes effect only when the type of the custom background image is `BACKGROUND_IMG`.
   */
  const char* source;

  /** blur degree */
  int blur_degree;
} virtual_background_source;

typedef struct _segmentation_property {
  int model_type;
  float green_capacity;
} segmentation_property;

/**
 * The definition of the screen sharing encoding parameters.
 */
typedef struct _screen_capture_parameters {
  /**
   * The dimensions of the shared region in terms of width &times; height. The default value is 0, which means
   * the original dimensions of the shared screen.
   */
  video_dimensions dimensions;
  /**
   * The frame rate (fps) of the shared region. The default value is 5. We do not recommend setting
   * this to a value greater than 15.
   */
  int frame_rate;
  /**
   * The bitrate (Kbps) of the shared region. The default value is 0, which means the SDK works out a bitrate
   * according to the dimensions of the current screen.
   */
  int bitrate;
  /** Sets whether or not to capture the mouse for screen sharing:
  * - true: (Default) Capture the mouse.
  * - false: Do not capture the mouse.
  */
  int capture_mouse_cursor;
  /** Whether to bring the window to the front when calling
  * \ref IRtcEngine::startScreenCaptureByWindowId "startScreenCaptureByWindowId" to share the window:
  * - true: Bring the window to the front.
  * - false: (Default) Do not bring the window to the front.
  */
  int window_focus;
  /**
   * A list of IDs of windows to be blocked.
   */
  view_t *exclude_window_list;
  /**
   * The number of windows to be blocked.
   */
  int exclude_window_count;
} screen_capture_parameters;

/**
 * The Audio file recording options.
 */
typedef struct _audio_recording_config {
  /**
   * The path of recording file.
   * The string of the file path is in UTF-8 code.
   */
  const char* file_path;
  /**
   * Determines whether to encode audio data.
   * - true: Encode the audio data with AAC Encoder.
   * - false: (Default) Do not encode the audio data. Save audio data as a wav file.
   */
  int encode;
  /**
   * The sample rate of audio data. Default is 32000.
   * The optional value is 16000, 32000, 44100, or 48000.
   */
  int sample_rate;
  /**
   * The recording type of audio data.
   */
  int file_recording_type;
  /**
   * The recording quality of audio data.
   */
  int quality;
} audio_recording_config;

/**
 * The Audio encoded frame receiver options.
 *
 */
typedef struct _audio_encoded_frame_observer_config {
  /**
   * The position where SDK record the audio, and callback to encoded audio frame receiver.
   */
  int postion_type;
  /**
   * The audio encoding type of encoded frame.
   */
  int encoding_type;
} audio_encoded_frame_observer_config;

/**
 * The IAudioEncodedFrameObserver class.
 * @ANNOTATION:TYPE:OBSERVER
 */
typedef struct _audio_encoded_frame_observer {
  void (*on_record_audio_encoded_frame)(AGORA_HANDLE agora_audio_encoded_frame, 
        const uint8_t* frame_buffer,  int length, const encoded_audio_frame_info* audio_encoded_frame_info);
  void (*on_playback_audio_encoded_frame)(AGORA_HANDLE agora_audio_encoded_frame, 
        const uint8_t* frame_buffer,  int length, const encoded_audio_frame_info* audio_encoded_frame_info);
  void (*on_mixed_audio_encoded_frame)(AGORA_HANDLE agora_audio_encoded_frame, 
        const uint8_t* frame_buffer,  int length, const encoded_audio_frame_info* audio_encoded_frame_info);
} audio_encoded_frame_observer;

/** The definition of ChannelMediaInfo.
 */
typedef struct _channel_media_info {
    /** The channel name. The default value is NULL, which means that the SDK
     * applies the current channel name.
     */
  const char* channel_name;
    /** The token that enables the user to join the channel. The default value
     * is NULL, which means that the SDK applies the current token.
     */
  const char* token;
    /** The user ID.
     */
  uid_t uid;
} channel_media_info;

/** The definition of ChannelMediaRelayConfiguration.
 */
typedef struct _channel_media_relay_config {
  /** Pointer to the source channel: ChannelMediaInfo.
    *
    * @note
    * - `uid`: ID of the user whose media stream you want to relay. We
    * recommend setting it as 0, which means that the SDK relays the media
    * stream of the current broadcaster.
    * - If you do not use a token, we recommend using the default values of
    * the parameters in ChannelMediaInfo.
    * - If you use a token, set uid as 0, and ensure that the token is
    * generated with the uid set as 0.
    */
  channel_media_info *src_info;
  /** Pointer to the destination channel: ChannelMediaInfo. If you want to
    * relay the media stream to multiple channels, define as many
    * ChannelMediaInfo structs (at most four).
    *
    * @note `uid`: ID of the user who is in the source channel.
    */
  channel_media_info *dest_infos;
  /** The number of destination channels. The default value is 0, and the
    * value range is [0,4). Ensure that the value of this parameter
    * corresponds to the number of ChannelMediaInfo structs you define in
    * `destInfos`.
    */
  int dest_count;
} channel_media_relay_config;

typedef struct _uplink_network_info {
  /**
   * The target video encoder bitrate bps.
   */
  int video_encoder_target_bitrate_bps;
} uplink_network_info;

typedef struct _peer_downlink_info {
  /**
   * The ID of the user who owns the remote video stream.
   */
  const char* uid;
  /**
   * The remote video stream type: #VIDEO_STREAM_TYPE.
   */
  int stream_type;
  /**
   * The remote video downscale type: #REMOTE_VIDEO_DOWNSCALE_LEVEL.
   */
  int current_downscale_level;
  /**
   * The expected bitrate in bps.
   */
  int expected_bitrate_bps;
} peer_downlink_info;

typedef struct _downlink_network_info {
  /**
   * The queue time in ms.
   */
  int lastmile_buffer_delay_time_ms;

  /**
   * The bandwidth estimation bps.
   */
  int bandwidth_estimation_bps;
  /**
   * The total video downscale level count.
   */
  int total_downscale_level_count;
  /**
   * The peer video downlink info array.
   */
  peer_downlink_info* peer_downlink_info;
  /**
   * The total video received count.
   */
  int total_received_video_count;
} downlink_network_info;

/** Configurations of the built-in encryption schemas. */
typedef struct _encryption_config {
  /**
   * The encryption mode. The default encryption mode is `AES_128_GCM2`. See #ENCRYPTION_MODE.
   */
  int encryption_mode;
  /**
   * The encryption key in the string format.
   *
   * @note If you do not set an encryption key or set it as NULL, you cannot use the built-in encryption, and the SDK returns #ERR_INVALID_ARGUMENT (-2).
   */
  const char* encryption_key;
  uint8_t encryption_kdf_salt[32];
} encryption_config;



/** 
 * Spatial audio parameters
 */
typedef struct _spatial_audio_params {
  /**
   * optional azimuth: speaker azimuth in a spherical coordinate system centered on the listener
   */
  double speaker_azimuth;
  /**
   * optional azimuth: speaker elevation in a spherical coordinate system centered on the listener
   */
  double speaker_elevation;
  /**
   * distance between speaker and listener
   */
  double speaker_distance;
  /**
   * speaker orientation [0-180]: 0 degree is the same with listener orientation
   */
  int speaker_orientation;
  /**
   * enable blur or not for the speaker
   */
  int enable_blur;
  /**
   * enable air absorb or not for the speaker
   */
  int enable_air_absorb;
} spatial_audio_params;

/**
 * The configuration for creating encoded video image track.
 */
typedef struct _sender_options {
  /**
   * Determins whether to enable the CC mode: #TCcMode.
   */
  int cc_mode;

  /**
   * Determins which codec type is used for encoded image: #VIDEO_CODEC_TYPE
   */
  int codec_type;

  /**
   * Video encoding target bitrate (Kbps).
   *
   * Choose one of the following options:
   *
   * - #STANDARD_BITRATE: (Recommended) Standard bitrate.
   *   - Communication profile: The encoding bitrate equals the base bitrate.
   *   - Live broadcast profile: The encoding bitrate is twice the base bitrate.
   * - #COMPATIBLE_BITRATE: Compatible bitrate: The bitrate stays the same
   * regardless of the profile.
   *
   * The Communication profile prioritizes smoothness, while the Live Broadcast
   * profile prioritizes video quality (requiring a higher bitrate). Agora
   * recommends setting the bitrate mode a #STANDARD_BITRATE or simply to
   * address this difference.
   *
   * The following table lists the recommended video encoder configurations,
   * where the base bitrate applies to the communication profile. Set your
   * bitrate based on this table. If the bitrate you set is beyond the proper
   * range, the SDK automatically sets it to within the range.

   | Resolution             | Frame Rate (fps) | Base Bitrate (Kbps, for Communication) | Live Bitrate (Kbps, for Live Broadcast)|
   |------------------------|------------------|----------------------------------------|----------------------------------------|
   | 160  &times; 120       | 15               | 65                                     | 130  |
   | 120  &times; 120       | 15               | 50                                     | 100  |
   | 320  &times; 180       | 15               | 140                                    | 280  |
   | 180  &times; 180       | 15               | 100                                    | 200  |
   | 240  &times; 180       | 15               | 120                                    | 240  |
   | 320  &times; 240       | 15               | 200                                    | 400  |
   | 240  &times; 240       | 15               | 140                                    | 280  |
   | 424  &times; 240       | 15               | 220                                    | 440  |
   | 640  &times; 360       | 15               | 400                                    | 800  |
   | 360  &times; 360       | 15               | 260                                    | 520  |
   | 640  &times; 360       | 30               | 600                                    | 1200 |
   | 360  &times; 360       | 30               | 400                                    | 800  |
   | 480  &times; 360       | 15               | 320                                    | 640  |
   | 480  &times; 360       | 30               | 490                                    | 980  |
   | 640  &times; 480       | 15               | 500                                    | 1000 |
   | 480  &times; 480       | 15               | 400                                    | 800  |
   | 640  &times; 480       | 30               | 750                                    | 1500 |
   | 480  &times; 480       | 30               | 600                                    | 1200 |
   | 848  &times; 480       | 15               | 610                                    | 1220 |
   | 848  &times; 480       | 30               | 930                                    | 1860 |
   | 640  &times; 480       | 10               | 400                                    | 800  |
   | 1280 &times; 720       | 15               | 1130                                   | 2260 |
   | 1280 &times; 720       | 30               | 1710                                   | 3420 |
   | 960  &times; 720       | 15               | 910                                    | 1820 |
   | 960  &times; 720       | 30               | 1380                                   | 2760 |
   | 1920 &times; 1080      | 15               | 2080                                   | 4160 |
   | 1920 &times; 1080      | 30               | 3150                                   | 6300 |
   | 1920 &times; 1080      | 60               | 4780                                   | 6500 |
   | 2560 &times; 1440      | 30               | 4850                                   | 6500 |
   | 2560 &times; 1440      | 60               | 6500                                   | 6500 |
   | 3840 &times; 2160      | 30               | 6500                                   | 6500 |
   | 3840 &times; 2160      | 60               | 6500                                   | 6500 |
   */
  int target_bitrate;
} sender_options;
#ifdef __cplusplus
}
#endif  // __cplusplus
