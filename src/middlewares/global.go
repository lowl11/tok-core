package middlewares

import "tok-core/src/events/client_session_event"

var (
	clientSession *client_session_event.Event
)

func SetClientSession(session *client_session_event.Event) {
	clientSession = session
}
