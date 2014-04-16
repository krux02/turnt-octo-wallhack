#pragma once
#include <stdlib.h>
#include <GL/glew.h>

//#undef GLEW_GET_FUN
//#define GLEW_GET_FUN(x) (*x)

extern void goDebugCallback(unsigned int source, unsigned int _type, unsigned int id, unsigned int severity, int length, char* message, void* userParam);

int initDebugContext();
