//
//  Agora C SDK
//
//  Created by Lyu Ge in 2020.10
//  Copyright (c) 2020 Agora.io. All rights reserved.
//
#pragma once

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

typedef enum _agora_error_code {

    // Catch-all library error
    // @ANNOTATION:MESSAGE:Generic error
    err_general = 1,

    // send attempted when endpoint write queue was full
    // @ANNOTATION:MESSAGE:Send queue full
    err_send_queue_full,

    // Attempted an operation using a payload that was improperly formatted
    // ex: invalid UTF8 encoding on a text message.
    // @ANNOTATION:MESSAGE:Payload violation
    err_payload_violation,

    // Attempted to open a secure connection with an insecure endpoint
    // @ANNOTATION:MESSAGE:Endpoint not secure
    err_endpoint_not_secure,

    // Attempted an operation that required an endpoint that is no longer
    // available. This is usually because the endpoint went out of scope
    // before a connection that it created.
    // @ANNOTATION:MESSAGE:Endpoint not available
    err_endpoint_unavailable,

    // An invalid uri was supplied
    // @ANNOTATION:MESSAGE:Invalid URI
    err_invalid_uri,

    // The endpoint is out of outgoing message buffers
    // @ANNOTATION:MESSAGE:No outgoing message buffers
    err_no_outgoing_buffers,

    // The endpoint is out of incoming message buffers
    // @ANNOTATION:MESSAGE:No incoming message buffers
    err_no_incoming_buffers,

    // The connection was in the wrong state for this operation
    // @ANNOTATION:MESSAGE:Invalid state
    err_invalid_state,

    // Unable to parse close code
    // @ANNOTATION:MESSAGE:Unable to extract close code
    err_bad_close_code,

    // Close code is in a reserved range
    // @ANNOTATION:MESSAGE:Extracted close code is in a reserved range
    err_reserved_close_code,

    // Close code is invalid
    // @ANNOTATION:MESSAGE:Extracted close code is in an invalid range
    err_invalid_close_code,

    // Invalid UTF-8
    // @ANNOTATION:MESSAGE:Invalid UTF-8
    err_invalid_utf8,

    // Invalid subprotocol
    // @ANNOTATION:MESSAGE:Invalid subprotocol
    err_invalid_subprotocol,

    // An operation was attempted on a connection that did not exist or was
    // already deleted.
    // @ANNOTATION:MESSAGE:Bad Connection
    err_bad_connection,

    // Unit testing utility error code
    // @ANNOTATION:MESSAGE:Test Error
    err_test,

    // Connection creation attempted failed
    // @ANNOTATION:MESSAGE:Connection creation attempt failed
    err_con_creation_failed,

    // Selected subprotocol was not requested by the client
    // @ANNOTATION:MESSAGE:Selected subprotocol was not requested by the client
    err_unrequested_subprotocol,

    // Attempted to use a client specific feature on a server endpoint
    // @ANNOTATION:MESSAGE:Feature not available on client endpoints
    err_client_only,

    // Attempted to use a server specific feature on a client endpoint
    // @ANNOTATION:MESSAGE:Feature not available on server endpoints
    err_server_only,

    // HTTP connection ended
    // @ANNOTATION:MESSAGE:HTTP connection ended
    err_http_connection_ended,

    // WebSocket opening handshake timed out
    // @ANNOTATION:MESSAGE:The opening handshake timed out
    err_open_handshake_timeout,

    // WebSocket close handshake timed out
    // @ANNOTATION:MESSAGE:The closing handshake timed out
    err_close_handshake_timeout,

    // Invalid port in URI
    // @ANNOTATION:MESSAGE:Invalid URI port
    err_invalid_port,

    // An async accept operation failed because the underlying transport has been
    // requested to not listen for new connections anymore.
    // @ANNOTATION:MESSAGE:Async Accept not listening
    err_async_accept_not_listening,

    // The requested operation was canceled
    // @ANNOTATION:MESSAGE:Operation canceled
    err_operation_canceled,

    // Connection rejected
    // @ANNOTATION:MESSAGE:Connection rejected
    err_rejected,

    // Upgrade Required. This happens if an HTTP request is made to a
    // WebSocket++ server that doesn't implement an http handler
    // @ANNOTATION:MESSAGE:Upgrade required
    err_upgrade_required,

    // Invalid WebSocket protocol version
    // @ANNOTATION:MESSAGE:Invalid version
    err_invalid_version,

    // Unsupported WebSocket protocol version
    // @ANNOTATION:MESSAGE:Unsupported version
    err_unsupported_version,

    // HTTP parse error
    // @ANNOTATION:MESSAGE:HTTP parse error
    err_http_parse_error,

    // Extension negotiation failed
    // @ANNOTATION:MESSAGE:Extension negotiation failed
    err_extension_neg_failed

} agora_error;

#ifdef __cplusplus
}
#endif  // __cplusplus
