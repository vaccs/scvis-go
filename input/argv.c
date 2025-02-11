#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int main(int argc, char* argv[]) {

   char *buf;

   buf = malloc(strlen(argv[1]+1));
   strncpy(buf,argv[1],strlen(argv[1])); 
   printf("buf = %s\n",buf);
}
