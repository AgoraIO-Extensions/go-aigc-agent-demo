//
//  vad.h
//
//  Created by ZhouRui on 2024/2/8.
//

#ifndef vad_h
#define vad_h

#include <stdint.h>
#include <stdlib.h>

#define AGORA_UAP_VAD_VERSION (20240708)

enum VAD_STATE {
    VAD_STATE_NONE_SPEAKING = 0,
    VAD_STATE_START_SPEAKING = 1,
    VAD_STATE_SPEAKING = 2,
    VAD_STATE_STOP_SPEAKING = 3,
};

typedef struct Vad_Config_ {
  int fftSz;  // fft-size, only support: 128, 256, 512, 1024, default value is 1024
  int hopSz;  // fft-Hop Size, will be used to check, default value is 160
  int anaWindowSz;  // fft-window Size, will be used to calc rms, default value is 768
  int frqInputAvailableFlag;  // whether Aed_InputData will contain external freq. power-sepctra, default value is 0
  int useCVersionAIModule; // whether to use the C version of AI submodules, default value is 0
  float voiceProbThr;  // voice probability threshold 0.0f ~ 1.0f, default value is 0.8
  float rmsThr; // rms threshold in dB, default value is -40.0
  float jointThr; // joint threshold in dB, default value is 0.0
  float aggressive; // aggressive factor, greater value means more aggressive, default value is 5.0
  int startRecognizeCount; // start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
  int stopRecognizeCount; // max recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 6
  int preStartRecognizeCount; // pre start recognize count, buffer size for 10ms 16KHz 16bit 1channel PCM, default value is 10
  float activePercent; // active percent, if over this percent, will be recognized as speaking, default value is 0.6
  float inactivePercent; // inactive percent, if below this percent, will be recognized as non-speaking, default value is 0.2
} Vad_Config;

typedef struct Vad_AudioData_ {
  void* audioData; // this frame's input signal
  int size;
} Vad_AudioData;

#ifdef __cplusplus
extern "C" {
#endif

/****************************************************************************
 * Agora_UAP_VAD_Create(...)
 *
 * This function creats a state handler from nothing, which is NOT ready for
 * processing
 *
 * Input:
 *
 * Output:
 *      - stPtr         : buffer to store the returned state handler
 *
 * Return value         :  0 - Ok
 *                        -1 - Error
 */
int Agora_UAP_VAD_Create(void** stPtr, const Vad_Config* config);

/****************************************************************************
 * Agora_UAP_VAD_Destroy(...)
 *
 * destroy VAD instance, and releasing all the dynamically allocated memory
 *
 * Input:
 *      - stPtr         : buffer of State Handler, after this method, this
 *                        handler won't be usable anymore
 *
 * Output:
 *
 * Return value         :  0 - Ok
 *                        -1 - Error
 */
int Agora_UAP_VAD_Destroy(void** stPtr);

/****************************************************************************
 * Agora_UAP_VAD_Proc(...)
 *
 * process a single frame
 *
 * Input:
 *      - stPtr         : State Handler which has gone through create and
 *                        memAllocate and reset
 *      - pIn           : input data stream
 *
 * Output:
 *      - pOut          : output data
 *
 * Return value         :  0 - Ok
 *                        -1 - Error
 */
int Agora_UAP_VAD_Proc(void* stPtr, const Vad_AudioData* pIn, Vad_AudioData* pOut, enum VAD_STATE* state);

#ifdef __cplusplus
}
#endif
#endif /* vad_h */
