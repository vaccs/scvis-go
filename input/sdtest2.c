#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/resource.h>

#include "sd.h"

int main() {
  struct rlimit rlpDefine;
  rlpDefine.rlim_max = 0;
  setrlimit(RLIMIT_CORE,&rlpDefine);

  char *track = (char*)malloc(15);
  VACCSSensitiveAnnotate(track);
  mlock(track,15);

  strcpy(track,"HELLO");
  for (int i = 0; i < 15; i++)
     track[i] = 0;

  munlock(track, 15);

  int insecure;
  VACCSSensitiveAnnotate(insecure);

  insecure = 15;

  printf("Done\n");
}
