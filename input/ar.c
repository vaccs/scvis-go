//
//-- Illustrate layout of activation record, overflows (using argv[1])
// 
//
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int simpleFunction(int iparm, char *cpparm1, char *cpparm2, char *cpparm3){
   char  lbuf1[3];
   int   localInt;
   char  lbuf2[8];
   char  *lbuf3;

   strcpy(lbuf1,cpparm1);
   localInt=iparm;
   strcpy(lbuf2,cpparm2);
   lbuf3=malloc(8);
   strcpy(lbuf3,cpparm3);

   free(lbuf3);

   return(0);
}
int main(int argc, char **argv){
  int x;
  simpleFunction(argc,argv[1],argv[2],argv[3]);
  x=1;
 printf("Done\n");
}
