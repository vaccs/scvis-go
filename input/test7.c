//
//-- Illustrate layout of activation record, overflows (using argv[1])
//
//
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int simpleFunction(char *cpparm1){
   char  lbuf1[3];

   strcpy(lbuf1,cpparm1);
   return(0);
}
int main(int argc, char **argv){
  simpleFunction(argv[1]);
}
