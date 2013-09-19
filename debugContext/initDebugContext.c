#include "initDebugContext.h"

int initDebugContext() {
	if( GLEW_ARB_debug_output ) {
		glDebugMessageCallbackARB((GLDEBUGPROCARB)goDebugCallback, NULL);
		return 1;
	}
	return 0;
}